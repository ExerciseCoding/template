package web

import "github.com/ExerciseCoding/template/internal/web/demo"

type Middleware func(next demo.HandlerFunc) demo.HandlerFunc
