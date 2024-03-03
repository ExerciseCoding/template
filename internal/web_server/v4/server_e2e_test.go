package v4

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()
	h.addRoute(http.MethodGet, "/user", func(context *Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	//handle1 := func(ctx *Context) {
	//	fmt.Println("处理第一件事")
	//}

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})

	h.Get("/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.Get("/check/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不正确"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", id)))
	})

	type User struct {
		Name string `json:"name"`
	}
	h.Get("/user/:id", func(ctx *Context) {
		ctx.RespJson(User{
			Name: "Tom",
		})
	})
	h.Get("/user/student/:id", func(ctx *Context) {
		s := SafeContext{
			ctx: ctx,
		}
		s.RespJSONOK(User{
			Name: "Tom",
		})
	})

	h.Start("127.0.0.1:19999")
}
