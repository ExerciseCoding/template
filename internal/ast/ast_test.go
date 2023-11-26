package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAst (t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", `package ast
	import (
		"fmt"
		"go/ast"
		"reflect"
	)
	
	type printVisitor struct {}

	func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
		return p
	}
	`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ast.Walk(&printVisitor{}, f)
}