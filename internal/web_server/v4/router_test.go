package v4

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// 验证思路
	// 第一步：构建路由树
	// 第二步：验证路由树
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc/*",
		//},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/order/query/:ordername",
		},
	}
	var mockHandler HandlerFunc = func(context *Context) {}
	r := newRouter()

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
							},
							"query": &node{
								path: "query",
								paramChild: &node{
									path:    ":ordername",
									handler: mockHandler,
								},
							},
						},
						startChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
				handler: mockHandler,
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"login": &node{
						path:    "login",
						handler: mockHandler,
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:    "create",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	msg, equal := wantRouter.equal(&r)
	assert.True(t, equal, msg)

	r = newRouter()
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web: 路径必须以 / 开头")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c////", mockHandler)
	}, "web: 路径必须以 / 结尾")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b///c", mockHandler)
	}, "web: 不能出现连续的 /")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "web: 路由冲突[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "web: 路由冲突[/a/b/c]")
}

// 返回一个错误信息，帮助排查问题
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径path不同"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点个数不同"), false
	}

	if n.startChild != nil {
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}
	// 比较handler
	nhandler := reflect.ValueOf(n.handler)
	yhandler := reflect.ValueOf(y.handler)
	if nhandler != yhandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, v := range y.children {
		dst, ok := n.children[path]
		if !ok {
			return fmt.Sprintf("节点不存在"), false
		}
		msg, ok := v.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/order/query/:ordername",
		},
	}
	var mockHandler HandlerFunc = func(context *Context) {}
	r := newRouter()
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name string

		method string
		path   string

		wantFound bool
		wantInfo  *matchInfo
	}{
		{
			name: "method not found",

			method: http.MethodOptions,
			path:   "/order/detail",
		},
		{
			// 完全命中
			name: "order detail",

			method: http.MethodGet,
			path:   "/order/detail",

			wantFound: true,
			wantInfo: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
			},
		},
		// 命中了但是没有handler
		{
			name: "order",

			method: http.MethodGet,
			path:   "/order",

			wantFound: true,
			wantInfo: &matchInfo{
				n: &node{
					path: "order",
					children: map[string]*node{
						"detail": &node{
							path:    "detail",
							handler: mockHandler,
						},
						"query": &node{
							path: "query",
							paramChild: &node{
								path:    ":ordername",
								handler: mockHandler,
							},
						},
					},
					startChild: &node{
						path:    "*",
						handler: mockHandler,
					},
				},
			},
		},
		{
			name: "path not found",

			method: http.MethodGet,
			path:   "/aaaaa",
		},
		// 根节点
		{
			name: "root",

			method: http.MethodGet,
			path:   "/",

			wantFound: true,
			wantInfo: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
					children: map[string]*node{
						"user": &node{
							path:    "user",
							handler: mockHandler,
							children: map[string]*node{
								"home": &node{
									path:    "home",
									handler: mockHandler,
								},
							},
						},
						"order": &node{
							path: "order",
							children: map[string]*node{
								"detail": &node{
									path:    "detail",
									handler: mockHandler,
								},
								"query": &node{
									path: "query",
									paramChild: &node{
										path:    ":ordername",
										handler: mockHandler,
									},
								},
							},
							startChild: &node{
								path:    "*",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},

		// 通配符匹配
		{
			name: "/order/start",

			method: http.MethodGet,
			path:   "/order/start",

			wantFound: true,
			wantInfo: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},

		// 命中/order/query/:ordername
		{
			name:      ":ordername",
			method:    http.MethodGet,
			path:      "/order/query/123",
			wantFound: true,
			wantInfo: &matchInfo{
				n: &node{
					path:    ":ordername",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"ordername": "123",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.wantInfo.pathParams, n.pathParams)
			msg, ok := n.n.equal(tc.wantInfo.n)

			assert.True(t, ok, msg)
		})
	}
}

func TestRouter_findRouter_Middleware(t *testing.T) {
	mdlBuilder := func(i byte) Middleware {
		return func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				ctx.RespData = append(ctx.RespData, i)
				next(ctx)
			}
		}
	}
	testRouters := []struct {
		method string
		path   string

		mdls []Middleware
	}{
		{
			method: http.MethodGet,
			path:   "/a/b",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b')},
		},
		{
			method: http.MethodGet,
			path:   "/a/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('*')},
		},
		{
			method: http.MethodGet,
			path:   "/a/b/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('*')},
		},
		{
			method: http.MethodPost,
			path:   "/a/b/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('*')},
		},
		{
			method: http.MethodPost,
			path:   "/a/*/c",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('*'), mdlBuilder('c')},
		},
		{
			method: http.MethodPost,
			path:   "/a/b/c",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('c')},
		},
		{
			method: http.MethodDelete,
			path:   "/*",
			mdls:   []Middleware{mdlBuilder('*')},
		},
		{
			method: http.MethodDelete,
			path:   "/",
			mdls:   []Middleware{mdlBuilder('/')},
		},
	}
	//var mockHandler HandlerFunc = func(context *Context) {}
	r := newRouter()
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, nil, route.mdls...)
	}

	testCases := []struct {
		name string

		method string
		path   string

		wantResp string
	}{
		{
			name:   "static, not match",
			method: http.MethodGet,
			path:   "/a",
		},
		{
			name:     "static, match",
			method:   http.MethodGet,
			path:     "/a/c",
			wantResp: "a*",
		},
		{
			name:     "static and star",
			method:   http.MethodGet,
			path:     "/a/b",
			wantResp: "a*ab",
		},
		{
			name:     "static and star",
			method:   http.MethodGet,
			path:     "/a/b/c",
			wantResp: "a*abab*",
		},
		{
			name:     "abc",
			method:   http.MethodPost,
			path:     "/a/b/c",
			wantResp: "a*cab*abc",
		},
		{
			name:     "root",
			method:   http.MethodDelete,
			path:     "/",
			wantResp: "/",
		},
		{
			name:     "root star",
			method:   http.MethodDelete,
			path:     "/a",
			wantResp: "/*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, _ := r.findRoute(tc.method, tc.path)
			mdls := mi.mdls
			var root HandlerFunc = func(ctx *Context) {
				assert.Equal(t, tc.wantResp, string(ctx.RespData))
			}
			for i := len(mdls) - 1; i >= 0; i-- {
				root = mdls[i](root)
			}

			root(&Context{
				RespData: make([]byte, 0, len(tc.wantResp)),
			})
		})
	}
}
