package demo

import "strings"

type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handleFunc HandlerFunc) {
	root, ok := r.trees[method]
	if !ok {
		//根节点
		root = &node{path: "/"}
		r.trees[method] = root
	}

	if path == "/" {
		root.handle = handleFunc
		return
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")

	cur := root
	for _, seg := range segs {
		cur = cur.childOrCreate(seg)
	}
	cur.handle = handleFunc
	cur.route = path
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}
	segs := strings.Split(strings.Trim(path, "/"), "/")

	cur := root
	for _, seg := range segs {
		if cur.children == nil {
			if cur.paramsChild != nil {
				mi := &matchInfo{
					n: cur.paramsChild,
					pathParams: map[string]string{
						cur.paramsChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.startChild}, cur.startChild != nil
		}
		child, ok := cur.children[seg]
		if !ok {
			if cur.paramsChild != nil {
				mi := &matchInfo{
					n: cur.paramsChild,
					pathParams: map[string]string{
						cur.paramsChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.startChild}, cur.startChild != nil
		}
		cur = child
	}

	return &matchInfo{n: cur}, true
}

// childOrCreate 查找子节点，如果子节点不存在就创建一个
// 并且将子节点放回去了 children 中
func (n *node) childOrCreate(path string) *node {
	// /a/*/c
	if path == "*" {
		if n.startChild == nil {
			n.startChild = &node{
				path: path,
			}
		}
		return n.startChild
	}

	if path[0] == ':' {
		if n.paramsChild == nil {
			n.paramsChild = &node{
				path: path,
			}
		}
		return n.paramsChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
	}
	n.children[path] = child
	n = child

	return child
}

type node struct {
	// /a/b/c 中的b这一段
	handle HandlerFunc
	path   string

	// path => 到子节点的映射
	children map[string]*node
	// children []*node

	//通配符匹配
	startChild *node

	paramsChild *node

	// route到达该节点的完整路由路径
	route string
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
