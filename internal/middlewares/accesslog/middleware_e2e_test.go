//go:build e2e

// go:build 标签，你可以在源代码文件中指定条件，以控制编译器在构建时是否包含该文件或特定部分的代码。
// 这样可以根据不同的构建环境或目标平台，选择性地包含或排除代码。
package accesslog

import (
	"fmt"
	//v2 "template/internal/web_server/v2"
	"testing"
)

func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdls := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	server := v2.NewHTTPServer(v2.ServerWithMiddleware(mdls))
	server.Get("/a/b/*", func(ctx *v2.Context) {
		ctx.Resp.Write([]byte("hello, it's me"))
	})
	server.Start(":8081")
}
