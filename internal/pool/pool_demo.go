package pool

import (
	"fmt"
	"sync"
)

type MyCache struct {
	pool sync.Pool
}

func NewMyCache() *MyCache {
	return &MyCache{
		pool: sync.Pool{
			New: func() any {
				fmt.Println("hhh, new")
				return []byte{}
			},
		},
	}
}


type Pool[T any] struct {
	p *sync.Pool
}

// New 创建一个Pool实例
// factory 必须返回T类型的值，并且不能返回nil
func New[T any](factory func() T) *Pool[T] {
	return &Pool[T]{
		p: &sync.Pool{
			New: func() any {
				return factory()
			},
		},
	}
}


// Get 取出一个元素
func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

// Put 放回去一个元素
func (p *Pool[T]) Put(t T) {
	p.p.Put(t)
}


