package test

import (
	"github.com/ExerciseCoding/template/internal/session"
	"github.com/ExerciseCoding/template/internal/session/cookie"
	"github.com/ExerciseCoding/template/internal/session/memory"
	v4 "github.com/ExerciseCoding/template/internal/web_server/v4"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	// 非常简单的登录检验
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store: memory.NewStore(time.Minute * 15),
		CtxSessKey: "sessKey",
	}
	server := v4.NewHTTPServer(v4.ServerWithMiddleware(func(next v4.HandlerFunc) v4.HandlerFunc {
		return func(ctx *v4.Context) {
			if ctx.Req.URL.Path == "/login" {
				//放过去， 用户准备登录
				next(ctx)
				return
			}

			//sessId, err := p.Extract(ctx.Req)
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登陆")
				return
			}
			// 刷新session的过期时间
			_ = m.RefreshSession(ctx)
			// 登录成功调用next
			next(ctx)
		}
	}))

	server.Post("/login", func(ctx *v4.Context) {
		// 要在这之前校验用户名和密码
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败了")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "xiaoming")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败了")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
		return
	})


	// 退出登录
	server.Post("/logout", func(ctx *v4.Context) {
		//  清理各种数据
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出登录")
	})

	server.Get("/user", func(ctx *v4.Context) {
		sess, _ := m.GetSession(ctx)

		// 假如说我要把昵称从session里面拿出来
		val, _ := sess.Get(ctx.Req.Context(), "nickname")

		ctx.RespData=[]byte(val.(string))

	})
	server.Start(":8081")
}