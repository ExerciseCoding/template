package v4

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// tplName 模版的名字，按名索引
	// data 渲染页面用的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 渲染页面，数据写入到writer 里面
	// Render(ctx, "aa", map[]{}, responseWriter)
	// Render(ctx context.Context, tplName string, data any, writer io.Writer) error

	// 不需要，让具体实现自己管
	// AddTemplate(tplName string, tpl []byte) error
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}
