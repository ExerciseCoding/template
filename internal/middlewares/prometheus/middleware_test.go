package prometheus

import (
	v2 "github.com/ExerciseCoding/template/internal/web_server/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		Namespace: "template",
		Subsystem: "web",
		Name:      "http_response",
	}

	server := v2.NewHTTPServer(v2.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *v2.Context) {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)
		ctx.RespJsonAndStatus(202, User{
			Name: "Tom",
		})
	})
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
	server.Start(":8081")
}

type User struct {
	Name string
}
