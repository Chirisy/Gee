## 1 http basis

### gee的雏形  
gee/  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|--gee.go  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|--go.mod  
main.go  
go.mod  
### 主要目标是实现Handler接口
```go
package http
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
func ListenAndServe(address string, h Handler) error
```
* main.go  
使用New()创建 gee 的实例，使用 GET()方法添加路由，最后使用Run()启动Web服务。
* gee.go  
1.当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中，(*Engine).Run()方法，是 ListenAndServe 的包装。  
2.Engine实现的 ServeHTTP 方法的作用就是，解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND 。  
3.首先定义了类型HandlerFunc，这是提供给框架用户的，用来定义路由映射的处理方法。我们在Engine中，添加了一张路由映射表router，key 由请求方法和静态路由地址构成，例如GET-/、GET-/hello、POST-/hello，这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法。

