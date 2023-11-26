package accesslog

import (
	"encoding/json"
	v2 "template/internal/web_server/v2"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() v2.Middleware {
	return func(next v2.HandlerFunc) v2.HandlerFunc {
		return func(ctx *v2.Context) {
			// 记录请求
			defer func() {
				l := accesslog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()

			next(ctx)
		}
	}
}

type accesslog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
