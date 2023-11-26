package pool

import (
	"fmt"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			fmt.Println("hhh,new")
			return []byte{}
		},
	}
	for i := 0; i < 100; i++ {
		val := pool.Get()
		pool.Put(val)
	}
}

func TestPoolUser(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			fmt.Println("hhh,new")
			return &User{}
		},
	}
	u1 := pool.Get().(*User)
	u1.Id = 12
	u1.Name = "Tom"
	u1.Reset()
	pool.Put(u1)
	u2 := pool.Get().(*User)
	fmt.Println(u2)
	u3 := pool.Get().(*User)
	fmt.Println(u3)
}

type User struct {
	Id   int
	Name string
}

func (u *User) Reset() {
	u.Id = 0
	u.Name = ""
}

func BenchmarkPool_Get(b *testing.B) {
	p := New[string](func() string {
		return ""
	})

	sp := &sync.Pool{
		New: func() any {
			return ""
		},
	}
	b.Run("pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p.Get()
		}
	})

	b.Run("sync.Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sp.Get()
		}
	})
}

func TestFor(t *testing.T) {
	batchSize := 10
	sourceLen := 22
	count := 0
	for start := 0; start < sourceLen; start += batchSize {
		end := start + batchSize
		if end > sourceLen {
			end = sourceLen
			count = sourceLen - start
		} else {
			count = 10
		}
		fmt.Println(start, end, count)
	}
}
