package demo

import (
	"net/http"
)

type HandlerFunc func(context *Context)

type Server interface {
	// Start 启动监听器
	Start(add string) error

	addRoute(method string, path string, handler HandlerFunc)
}

type HTTPServer struct {
	router
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.Serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Serve(ctx *Context) {
	mi, found := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	if !found || mi.n.handle == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("not found"))
		return
	}
	ctx.Params = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handle(ctx)
}

func (h *HTTPServer) Group(prefix string) *Group {
	return &Group{
		prefix: prefix,
		s:      h,
	}
}

type Group struct {
	prefix string
	s      Server
}

func (g *Group) AddRoute(method, path string, handler HandlerFunc) {
	g.s.addRoute(method, g.prefix+path, handler)
}

func (s *HTTPServer) Get(path string, handle HandlerFunc) {
	s.addRoute(http.MethodGet, path, handle)
}

func (s *HTTPServer) Post(path string, handle HandlerFunc) {
	s.addRoute(http.MethodPost, path, handle)
}

type MyGroup struct {
	prefix string
	s      Server
}

func NewMyGroup(server Server, prefix string) *MyGroup {
	return &MyGroup{
		prefix: prefix,
		s:      server,
	}
}

func (m *MyGroup) addRoute(method, path string, hanler HandlerFunc) {
	m.s.addRoute(method, m.prefix+path, hanler)
}
