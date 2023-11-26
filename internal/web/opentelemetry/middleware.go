package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"template/internal/web"
	"template/internal/web/demo"
)

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next demo.HandlerFunc) demo.HandlerFunc {
		return func(ctx *demo.Context) {
			_, span := m.tracer.Start(ctx.Req.Context(), "my-span")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			// 请求路径
			// /ussr/123 /user/456
			// /assssssssdddddddddddddfffffffffdggggggggggggggggggggggtttttttttttttthhhhhhhhhyyyyyyyyyjjjjjj
			// [:1024]防止攻击者攻击时使用很长的URL
			span.SetAttributes(attribute.String("htto.path", ctx.Req.URL.Path[:1024]))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			next(ctx)

		}
	}
}
