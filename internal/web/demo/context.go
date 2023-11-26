package demo

import "net/http"

type Context struct {
	Req    *http.Request
	Resp   http.ResponseWriter
	Params map[string]string

	// 命中路由
	MatchedRoute string
}
