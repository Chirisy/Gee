package gee

import (
	"net/http"
	"strings"
)

// HandlerFunc 定义了一个函数类型, 用于代表处理HTTP请求的方法
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

// ServeHTTP 实现了 http.Handler 接口, 作为 Engine 的实例方法.
// 受到请求时, 会根据请求路径查找路由映射表 router, 如果查到, 就执行注册的处理方法; 查不到, 就返回 404 NOT FOUND
// 实现后, Engine 就可以作为一个 HTTP 服务端被启动
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var key strings.Builder
	key.WriteString(request.Method)
	key.WriteString("-")
	key.WriteString(request.URL.Path)
	if handler, ok := e.router[key.String()]; ok {
		handler(writer, request)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("404 NOT FOUND: " + request.URL.Path + "\n"))
	}
}

// New 是 Engine 的构造函数, 返回一个实例
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// method 是请求的方法, 比如 GET、POST, pattern 是请求的路径, handler 是处理请求的方法
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	var key strings.Builder // strings.Builder 是 Go 1.10 引入的, 用于高效地构建字符串
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(pattern)         //优点是可以避免大量内存拷贝, 因为字符串是只读的, 无法直接修改
	e.router[key.String()] = handler // 将处理方法和路由注册到映射表 router 中
}

// GET 定义了添加 GET 请求的方法
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST 定义了添加 POST 请求的方法
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run 定义了启动 http 服务的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
