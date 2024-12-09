package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 定义了一个函数类型, 用于代表处理HTTP请求的方法
type HandlerFunc func(c *Context)

type RouteGroup struct {
	prefix      string        //分组的前缀
	middlewares []HandlerFunc //作用在这个分组上的中间件
	parent      *RouteGroup   //父分组,用于支持分组嵌套
	engine      *Engine       //保存engine,赋予分组访问router的能力
}

type Engine struct {
	*RouteGroup
	router *router
	groups []*RouteGroup
}

// ServeHTTP 实现了 http.Handler 接口, 作为 Engine 的实例方法.
// 受到请求时, 会根据请求路径查找路由映射表 router, 如果查到, 就执行注册的处理方法; 查不到, 就返回 404 NOT FOUND
// 实现后, Engine 就可以作为一个 HTTP 服务端被启动
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	var c = newContext(writer, request)
	c.Handlers = middlewares
	e.router.handle(c)
}

// New 是 Engine 的构造函数, 返回一个实例
func New() *Engine {
	e := &Engine{router: newRouter()}
	e.RouteGroup = &RouteGroup{engine: e}  //初始化顶层分组
	e.groups = []*RouteGroup{e.RouteGroup} //将顶层分组加入分组数组
	return e
}

// Group 创建子分组
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	newGroup := &RouteGroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      group.engine,
	}
	//加入分组数组
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouteGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp //该分组下的路由都有相同前缀
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run 定义了启动 http 服务的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
