package v1

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(context *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	addRoute(method string, path string, handlerFunc HandlerFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	router

	mdls []Middleware

	log func(msg string, any ...any)
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
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
	root := s.serve
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	var m Middleware = func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			// 就设置好了RespData 和 RespStatusCode
			next(ctx)
			s.flushResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
	//s.serve(ctx)
}

func (h *HTTPServer) flushResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("写入响应失败 %v", err)
	}

}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) serve(ctx *Context) {
	n, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		// 路由没有命中，就是404
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("NOT FOUND")
		//ctx.Resp.WriteHeader(404)
		//ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.PathParams = n.pathParams
	ctx.MatchRoute = n.n.route
	n.n.handler(ctx)
}

// 第一个问题：相对路径还是绝对路径
// 你的配置文件格式. json, yaml, xml
//func NewHTTPServerV2(cfgFilePath string) *HTTPServer {
// 你在这里加载配置，解析，然后初始化 HTTPServer
//}
