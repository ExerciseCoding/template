package web

import "github.com/template/internal/web/demo"

//import "template/internal/web/demo"

type Middleware func(next demo.HandlerFunc) demo.HandlerFunc
