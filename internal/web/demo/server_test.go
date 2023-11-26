package demo

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(context *Context) {
		context.Resp.Write([]byte("hello, world!"))
	})
	s.Get("/user", func(context *Context) {
		context.Resp.Write([]byte("hello, user!"))
	})

	s.Get("/user/*", func(context *Context) {
		context.Resp.Write([]byte("hello, user *!"))
	})

	s.Get("/user/home/:id", func(context *Context) {
		context.Resp.Write([]byte("hello, user id: !"))
	})
	s.Get("/user/home/:id", func(context *Context) {
		context.Resp.Write([]byte(fmt.Sprintf("%s", context.Params["id"])))
	})

	g := s.Group("/order")
	g.AddRoute(http.MethodGet, "/detail", func(context *Context) {
		context.Resp.Write([]byte("hello, order detail"))
	})
	s.Start(":8081")
}
