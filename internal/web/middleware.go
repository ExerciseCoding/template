package web

import "template/internal/web/demo"

type Middleware func(next demo.HandlerFunc) demo.HandlerFunc
