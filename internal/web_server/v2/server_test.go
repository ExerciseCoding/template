package v1

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	h := NewHTTPServer()
	h.mdls = []Middleware{
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("one before")
				next(ctx)
				fmt.Println("one after")
			}
		},
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("two before")
				next(ctx)
				fmt.Println("two after")
			}
		},
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("three before")
			}
		},
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("no see")
			}
		},
	}

	h.ServeHTTP(nil, &http.Request{})
}
