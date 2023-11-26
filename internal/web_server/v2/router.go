package v1

import (
	"fmt"
	"strings"
)

// 用来支持对路由树的操作
// 代表路由树(森林)
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// GET也有一棵树，POST也有一棵树

	// http method => 路由树根节点
	trees map[string]*node
}

type node struct {
	route string

	path string

	// 静态节点
	// 子 path 到子节点的映射
	children map[string]*node

	// 动态节点
	startChild *node

	// 参数
	paramChild *node

	// 缺一个代表用户注册的业务逻辑
	handler HandlerFunc
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// AddRoute 加一些限制
// path 必须以 / 开头，不能以 / 结尾， 中间也不能有连续的 //
func (r *router) addRoute(method string, path string, handleFunc HandlerFunc) {
	if path == "" {
		panic("web: 路径不能为空字符串")
	}

	// 开头不能没有/
	if path[0] != '/' {
		panic("web: 路径必须以 / 开头")
	}

	//结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路径不能以 / 结尾")
	}

	// 首先找到树
	root, ok := r.trees[method]
	if !ok {
		// 说明还没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}
	//切割path, /use/home切割会被分成三段
	segs := strings.Split(path[1:], "/")

	for _, seg := range segs {
		if seg == "" {
			panic("web: 不能出现连续的 /")
		}
		// 递归下去，找准位置
		// 如果中途有节点不存在，就要创建出来
		children := root.childOrCreate(seg)
		root = children
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	root.handler = handleFunc
	root.route = path
}

func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {
		n.startChild = &node{
			path: seg,
		}
		return n.startChild
	}
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {

	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	// 这里把前置和后置的 / 都去掉
	segs := strings.Split(strings.Trim(path, "/"), "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}

		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}

	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

// childOf 优先考虑静态匹配，匹配不上，再考虑通配符匹配
// 第一个返回值是子节点
// 第二个是标记是否是路径参数
// 第三个标记命中了没有
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}

	return child, false, ok
}
