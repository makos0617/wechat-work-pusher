package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Context HTTP 请求上下文
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   map[string]string
}

// Handler 处理函数类型
type Handler func(*Context)

// Middleware 中间件类型
type Middleware func(Handler) Handler

// Router 路由器
type Router struct {
	routes      map[string]map[string]Handler // method -> path -> handler
	middlewares []Middleware
	basePath    string
}

// NewRouter 创建新的路由器
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]Handler),
	}
}

// Use 添加中间件
func (r *Router) Use(middleware Middleware) *Router {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// Group 创建路由组
func (r *Router) Group(basePath string) *Router {
	// 确保分组路径以 '/' 开头
	if basePath != "" && basePath[0] != '/' {
		basePath = "/" + basePath
	}
	return &Router{
		routes:      r.routes,
		middlewares: r.middlewares,
		basePath:    r.basePath + basePath,
	}
}

// Post 添加 POST 路由
func (r *Router) Post(path string, handler Handler) {
	r.addRoute("POST", path, handler)
}

// Get 添加 GET 路由
func (r *Router) Get(path string, handler Handler) {
	r.addRoute("GET", path, handler)
}

// addRoute 添加路由
func (r *Router) addRoute(method, path string, handler Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}
	fullPath := r.basePath + path
	r.routes[method][fullPath] = handler
}

// ServeHTTP 实现 http.Handler 接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:  req,
		Response: w,
		Params:   make(map[string]string),
	}

	// 查找路由
	if methodRoutes, exists := r.routes[req.Method]; exists {
		if handler, exists := methodRoutes[req.URL.Path]; exists {
			// 应用中间件
			finalHandler := handler
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				finalHandler = r.middlewares[i](finalHandler)
			}
			finalHandler(ctx)
			return
		}
	}

	// 404
	http.NotFound(w, req)
}

// BindForm 绑定表单数据
func (c *Context) BindForm(v interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}
	
	// 简单的表单绑定实现
	// 根据结构体字段名从表单中获取值
	// 这里为了简化，直接在控制器中手动获取参数
	return nil
}

// Json 返回 JSON 响应
func (c *Context) Json(data interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Response).Encode(data)
}

// JsonSuccess 返回成功响应
func (c *Context) JsonSuccess(message string) {
	c.Json(map[string]interface{}{
		"result":  0,
		"message": message,
	})
}

// RestRet 响应结构
type RestRet struct {
	Result  int         `json:"result"`
	Message interface{} `json:"message"`
}

// 响应结果常量
const (
	ResultOK      = 0
	ResultErr     = 1
	ResultAuthErr = 2
)

// Server HTTP 服务器
type Server struct {
	router *Router
	port   int
}

// NewServer 创建新的服务器
func NewServer(port int) *Server {
	return &Server{
		router: NewRouter(),
		port:   port,
	}
}

// AddRoutes 添加路由
func (s *Server) AddRoutes(setupFunc func(*Router)) {
	setupFunc(s.router)
}

// Run 启动服务器
func (s *Server) Run() error {
	addr := ":" + strconv.Itoa(s.port)
	log.Printf("服务器启动在端口 %d", s.port)
	return http.ListenAndServe(addr, s.router)
}

// String 简单的字符串类型，用于兼容原有的 class.String
type String struct {
	String string
	Valid  bool
}

// 中间件辅助函数

// StopExecution 停止执行（通过 panic 实现）
func (c *Context) StopExecution() {
	panic("stop_execution")
}

// Next 继续执行下一个中间件
func (c *Context) Next() {
	// 在实际的中间件链中，这个方法会被中间件系统处理
}

// AuthMiddleware 创建认证中间件
func AuthMiddleware(authFunc func(*Context) bool) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) {
			defer func() {
				if r := recover(); r != nil {
					if r == "stop_execution" {
						return
					}
					panic(r)
				}
			}()
			
			if !authFunc(ctx) {
				return
			}
			next(ctx)
		}
	}
}