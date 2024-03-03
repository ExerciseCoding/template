package v4

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestUpload(t *testing.T) {
	tpl, err := template.ParseGlob("../../testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}

	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.Get("/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	fu := FileUploader{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("/Users/halisliu/go/src/template/internal/testdata", "upload", header.Filename)
		},
	}
	h.Post("/upload", fu.Handle())
	h.Start(":8081")
}

func TestDownload(t *testing.T) {
	h := NewHTTPServer()
	fu := FileDownloader{
		Dir: filepath.Join("/Users/halisliu/go/src/template/internal/testdata", "download"),
	}

	h.Get("/download", fu.Handle())
	h.Start(":8081")
}


func TestStaticResourceHandler_Handle(t *testing.T) {
	h := NewHTTPServer()
	s, err := NewStaticResourceHandler(filepath.Join("/Users/halisliu/go/src/template/internal/testdata","static"))
	require.NoError(t, err)
	h.Get("/static/:file", s.Handle)
	h.Start(":8089")
}