package ast

import "go/ast"

type printVisitor struct {

}

// node的子节点才会用到w
func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	return p
}