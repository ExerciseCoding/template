package recover

import v2 "github.com/ExerciseCoding/template/internal/web_server/v2"

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	Log        func(ctx *v2.Context)
}

func (m MiddlewareBuilder) Build() v2.Middleware {
	return func(next v2.HandlerFunc) v2.HandlerFunc {
		return func(ctx *v2.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.Data
				}
				m.Log(ctx)
			}()
			next(ctx)
		}
	}
}
