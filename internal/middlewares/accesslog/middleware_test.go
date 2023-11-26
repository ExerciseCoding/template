package accesslog

import (
	"fmt"
	"net/http"
	v2 "template/internal/web_server/v2"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdls := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	server := v2.NewHTTPServer(v2.ServerWithMiddleware(mdls))
	server.Post("/a/b/*", func(ctx *v2.Context) {
		fmt.Println("hello, it's me")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
