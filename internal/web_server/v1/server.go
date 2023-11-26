package v1

import "net/http"

type HandlerFunc func(context *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	addRoute(method string, path string, handlerFunc HandlerFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (s *HTTPServer) Get(path string, handler HandlerFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) Post(path string, handler HandlerFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Options(path string, handler HandlerFunc) {
	s.addRoute(http.MethodOptions, path, handler)
}

func (s *HTTPServer) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: write,
	}
	s.serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) serve(ctx *Context) {
	n, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		// 路由没有命中，就是404
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.PathParams = n.pathParams
	n.n.handler(ctx)
}
