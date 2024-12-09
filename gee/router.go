package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //每种请求方式的前缀树
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路由,进行切分,遇到*停止
func parsePattern(pattern string) []string {
	splits := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range splits {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 封装路由方法,添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	//检查是否已经有请求方法的前缀树,没有就需要新建树根节点
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	//将路由插入对应方法的前缀树,并建立路由和处理函数的映射
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) //传入路径
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0) //搜索匹配的路由节点
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { //例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				///static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}
				params[part[1:]] = strings.Join(searchParts[index:], "/") //重新将后面的part组合成一个字符串
			}
		}
		return n, params //返回匹配到的路由规则节点和路由参数
	}
	return nil, nil
}

// 封装路由方法,处理路由
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern                //为路径匹配路由用于获取处理方法
		c.Handlers = append(c.Handlers, r.handlers[key]) //将匹配到的核心业务逻辑加入到handlers数组
	} else {
		c.Handlers = append(c.Handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next() //开始调用中间件链
}
