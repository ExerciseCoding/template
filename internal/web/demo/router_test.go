package demo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
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
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path: "/user/detail/:user_sn",
		},
	}

	testCases := []struct {
		name   string
		method string
		path   string

		found    bool
		wantNode *node
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
			path:   "/",

			found:    false,
			wantNode: nil,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",

			found:    false,
			wantNode: nil,
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",

			found: true,
			wantNode: &node{
				path:   "/",
				handle: handle,
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",

			found: true,
			wantNode: &node{
				path:   "user",
				handle: handle,
			},
		},
		{
			name:   "no handler",
			method: http.MethodGet,
			path:   "/user",

			found: true,
			wantNode: &node{
				path: "user",
			},
		},
		{
			name:   "/order/abc",
			method: http.MethodGet,
			path:   "/order/abc",

			found: true,
			wantNode: &node{
				path:   "*",
				handle: handle,
			},
		},
		{
			name:   "/order/create",
			method: http.MethodGet,
			path:   "/order/create",

			found: true,
			wantNode: &node{
				path:   "create",
				handle: handle,
			},
		},
		{
			name: "/user/detail/:user_sn",
			method: http.MethodGet,
			path: "/user/detail/:user_sn",

			found: true,
			wantNode: &node{
				path: ":user_sn",
				handle: handle,
			},
		},
	}
	r := newRouter()
	for _, rt := range testRoutes {
		r.addRoute(rt.method, rt.path, handle)
	}

	for _, tc := range testCases {
		n, found := r.findRoute(tc.method, tc.path)
		assert.Equal(t, tc.found, found)
		if !found {
			continue
		}
		wantVal := reflect.ValueOf(tc.wantNode.handle)
		nVal := reflect.ValueOf(n.n.handle)
		assert.Equal(t, wantVal, nVal)
	}
}

func Test_router_addRoute(t *testing.T) {
	tests := []struct {
		name string

		// 输入
		method string
		path   string
	}{
		// 静态匹配
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
			path:   "/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/order/create/add",
		},
		{
			method: http.MethodPost,
			path:   "/order/cancel",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: "乱写方法",
			path:   "/",
		},
		{
			method: http.MethodGet,
			path: "/user/detail/:detail_sn",
		},
	}
	var handle HandlerFunc = func(coontext *Context) {}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:   "user",
						handle: handle,
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								paramsChild: &node{
									path: ":detail_sn",
									handle: handle,
								},
							},
						},
					},
					"home": &node{
						path:   "home",
						handle: handle,
					},
					"order": &node{
						path: "order",
						startChild: &node{
							path:   "*",
							handle: handle,
						},
						children: map[string]*node{
							"create": &node{
								path: "create",
								children: map[string]*node{
									"add": &node{
										path:   "add",
										handle: handle,
									},
								},
								handle: handle,
							},
						},
					},
				},
				handle: handle,
			},
			"乱写方法": &node{
				path:   "/",
				handle: handle,
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"cancel": &node{
								path:   "cancel",
								handle: handle,
							},
						},
					},
				},
			},
		},
	}

	res := &router{
		trees: map[string]*node{},
	}
	for _, tc := range tests {
		res.addRoute(tc.method, tc.path, handle)
	}
	errStr, ok := wantRouter.equal(res)
	assert.True(t, ok, errStr)
}

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标router里面没有方法 %s 的路由数", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("#{n.path} 节点 path 不相等 x #{n.path}, y {y.path}"), false
	}
	nhv := reflect.ValueOf(n.handle)
	yhv := reflect.ValueOf(y.handle)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点hander 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}

	if len(n.children) == 0 {
		return "", true
	}
	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

var handle HandlerFunc = func(coontext *Context) {}

func TestFunc(t *testing.T) {
	handle1 := func(ctx *context.Context) {
		fmt.Println("handle1")
	}
	val1 := reflect.ValueOf(handle1)
	val2 := reflect.ValueOf(handle1)
	assert.Equal(t, val1, val2)
}
