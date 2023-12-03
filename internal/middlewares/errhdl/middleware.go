package errhdl

import v2 "github.com/ExerciseCoding/template/internal/web_server/v2"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m MiddlewareBuilder) build() v2.Middleware {
	return func(next v2.HandlerFunc) v2.HandlerFunc {
		return func(ctx *v2.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				// 篡改结果
				ctx.RespData = resp
			}
		}
	}
}
