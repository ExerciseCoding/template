package opentelemetry

import (
	"github.com/ExerciseCoding/template/internal/web_server/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// 包名
const instrumentationName = "template/middlewares/internal/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() v2.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next v2.HandlerFunc) v2.HandlerFunc {
		return func(ctx *v2.Context) {
			reqCtx := ctx.Req.Context()
			// 尝试和客户端的trace结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			reqCtx, span := m.Tracer.Start(reqCtx, "unkown")
			//defer span.End()
			defer func() {
				span.SetName(ctx.MatchRoute)
				// 把响应码加上去
				span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
				span.End()
			}()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.schema", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			// 直接调用下一步

			ctx.Req = ctx.Req.WithContext(reqCtx)
			next(ctx)

			// 这个是只有执行完 next 才可能有值
			//span.SetName(ctx.MatchRoute)
			//span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
