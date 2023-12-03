package errhdl

import (
	v2 "github.com/ExerciseCoding/template/internal/web_server/v2"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1> 哈哈哈，走失了 </h1>
	</body>
</html>
`)).AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1> 请求不对 </h1>
	</body>
</html>

`))

	server := v2.NewHTTPServer(v2.ServerWithMiddleware(builder.build()))
	server.Start(":8081")
}
