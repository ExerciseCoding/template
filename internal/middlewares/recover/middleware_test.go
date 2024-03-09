package recover

import (
	"fmt"
	v2 "github.com/ExerciseCoding/template/internal/web_server/v2"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("内部错误"),
		Log: func(ctx *v2.Context) {
			fmt.Printf("错误路径: %s", ctx.Req.URL.String())
		},
	}

	server := v2.NewHTTPServer(v2.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *v2.Context) {
		panic("errs")
	})
	server.Start(":8081")
}
