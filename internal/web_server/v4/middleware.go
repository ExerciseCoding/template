package v4

import "sync"

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
// 用的最多的方式
type Middleware func(next HandlerFunc) HandlerFunc

type MiddlewareV1 interface {
	Invoke(next HandlerFunc) HandlerFunc
}

type Interceptor interface {
	Before(ctx *Context)
	After(ctx *Context)
	Surround(ctx *Context)
}

type ChainV1 struct {
	handlers []HandlerFunc
}

func (c ChainV1) Run(ctx *Context) {
	for _, h := range c.handlers {
		h(ctx)
	}
}

type HandlerFuncV1 func(ctx *Context) (next bool)

type ChainV2 struct {
	hanlers []HandlerFuncV1
}

func (c ChainV2) Run(ctx *Context) {
	for _, h := range c.hanlers {
		next := h(ctx)
		// 中断执行
		if !next {
			return
		}
	}
}

// 环形

type Net struct {
	handlers []HandlerFuncV2
}

func (c Net) Run(ctx *Context) {
	var wg *sync.WaitGroup
	for _, hdl := range c.handlers {
		h := hdl
		if h.concurrent {
			wg.Add(1)
			go func() {
				h.Run(ctx)
				wg.Done()
			}()
		} else {
			h.Run(ctx)
		}
	}
	wg.Wait()
}

type HandlerFuncV2 struct {
	concurrent bool
	handlers   []*HandlerFuncV2
}

func (h HandlerFuncV2) Run(ctx *Context) {
	var wg *sync.WaitGroup
	for _, hdl := range h.handlers {
		h := hdl
		if h.concurrent {
			wg.Add(1)
			go func() {
				h.Run(ctx)
				wg.Done()
			}()
		} else {
			h.Run(ctx)
		}
	}
}
