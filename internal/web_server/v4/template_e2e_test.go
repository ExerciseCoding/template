//go:build e2e

package v4

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"testing"
)

func TestLoginPage(t *testing.T) {

	tpl, err := template.ParseGlob("/Users/halisliu/go/src/template/internal/testdata/tpls/*.gohtml")
	require.NoError(t, err)

	engine := &GoTemplateEngine{
		T: tpl,
	}
	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.Get("/login", func(ctx *Context) {
		err := ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	h.Start(":8081")
}
