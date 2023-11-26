package sync

import (
	"fmt"
	"sync"
	"testing"
)

type safeResource struct {
	resource map[string]string
	lock  sync.RWMutex
}


func (s *safeResource) Add(key string, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.resource[key] = value
}

type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock sync.RWMutex
}

// 使用RWMutex实现double-check
// 加读锁先检查一遍
// 释放读锁
// 加写锁
// 再检查一遍


// 已经有key, 返回对应的值，然后loaded = true
// 没有,则放进去，返回loaded false
// goroutine 1 => ("key1",1)
// goroutine 2 => ("key1",2)
func(s *SafeMap[K, V]) LoadOrStore(key K, newValue V)( V,  bool) {
	s.lock.RLock()
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	// 为了解决问题1
	if ok {
		return oldVal, true
	}
	// 问题1：
	// goroutine1 先进来，那么这里就会变成key1 => 1
	// goroutine2 进来，那么这里就会变成key1 => 2
	s.values[key] = newValue
	return newValue, false
}


func TestLoadOrStore(t *testing.T) {
	sm := SafeMap[string, string]{
		values: make(map[string]string, 4),
	}
	sm.LoadOrStore("a", "b")
	fmt.Print("hello")
}


// 例子：实现一个线程安全的ArrayList
// 思路：切片本身不是线程安全的，所以最简单的做法就是利用读写锁封装一下。典型的装饰器模式的应用
// 如果考虑扩展性，那么需要预先定义一个List接口，后续可以有ArrayList， LinkedList,锁实现的线程安全List,以及无锁实现的线程安全List
// 任何非线程安全的类型，接口都可以利用读写锁+装饰器模式无侵入式地改造为线程安全的类型、接口

// Go的读写锁是写优先
//type SafeList[T any] struct {
//	List[T]
//	lock sync.RWMutex
//}
//
//func (s *SafeList[T]) Get(index int) (T, error) {
//	s.lock.RLock()
//	defer s.lock.Unlock()
//	return s.List.Get(index)
//}
//
//func (s *SafeList[T]) Append(t T) error {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	return s.List.Append(t)
//}


type OnceClose struct {
	close sync.Once
}

func (o *OnceClose) Close() error {
	o.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}

func TestOnceClose(t *testing.T) {
	o := &OnceClose{}
	o.Close()
	o.Close()
}