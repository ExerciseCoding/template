package context

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// context包使用注意事项
// - 一般只用做方法参数，而是作为第一个参数
// - 所有公共方法，除非是util,helper之类的方法，否则都加上context参数
// - 不要用作结构体字段，除非结构体本身也是表达一个上下文的概念

// 面试要点
// - context.Context使用场景：上下文传递和超时控制
// - context.Context原理：
//   - 父亲如何控制儿子：通过儿子主动加入到父亲的children里面，父亲只需要遍历就可以
//   - valueCtx 和timeCtx原理
func TestParentValueCtx(t *testing.T) {
	ctx := context.Background()
	childCtx := context.WithValue(ctx, "map", map[string]string{})
	ccChild := context.WithValue(childCtx, "key1", "value1")
	//m := childCtx.Value("map").(map[string]string) 类型断言，不使用.(map[string]string) m["key1"] = "val1" 赋值会有问题
	m := ccChild.Value("map").(map[string]string)
	m["key1"] = "val1"
	val := childCtx.Value("key1")
	fmt.Println(val)
	val = childCtx.Value("map")
	fmt.Println(val)

	//type valueCtx struct {
	//	Context
	//	key, val any
	//}

	// context包提供的三个控制方法 withCancel、withDeadline、withTimeout
	// 没有过期时间，但又需要再必要的时候取消，使用WithCancel
	// 在固定时间点过期，使用WithDeadline
	// 在一段时间后过期，使用withTimeout
	// 而后监听Done() 返回的chaneel,不管是主动调用cancel()还是超时，都能从这个channel里面取出数据，后面可以用Err()方法来判断究竟是哪种情况

}

func TestContext_timeout(t *testing.T) {
	bg := context.Background()
	timeoutCtx, cancel1 := context.WithTimeout(bg, time.Second)
	subCtx, cancel2 := context.WithTimeout(timeoutCtx, 3*time.Second)
	go func() {
		<-subCtx.Done()
		fmt.Println("timeout")
	}()
	time.Sleep(2 * time.Second)
	cancel2()
	cancel1()
}

func TestBussinessTimeout(t *testing.T) {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	end := make(chan struct{}, 1)
	go func() {
		MyBussiness()
		end <- struct{}{}
	}()
	ch := timeoutCtx.Done()
	select {
	case <-ch:
		fmt.Println("timeout")
	case <-end:
		fmt.Println("bussiness end")
	}
}

func MyBussiness() {
	time.Sleep(2 * time.Second)
	fmt.Println("my bussiness")
}

// 另外一种超时控制: time.AfterFunc（一般这种用法认为是定时任务，不是超时控制）
// 弊端: 1. 如果不主动取消，AfterFunc是必然会执行的 2.如果主动取消，在业务正常结束到主动取消之间，有一个短的时间差
func TestTimeoutTimeAfter(t *testing.T) {
	bsChan := make(chan struct{})
	go func() {
		MyBussiness()
		bsChan <- struct{}{}
	}()

	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("timeout")
	})
	<-bsChan
	timer.Stop()

}

// 例子：errgroup.WithContext 利用context来传递信号
// 1. WithContext会返回一个context.Context实例
// 2.如果errgroup.Group的Wait返回，或者任何一个Group执行的函数返回error，context.Context实例都会被取消
// 3.所有用户可以通过监听context.Context来判断errgroup.Group的执行情况
// 典型的将context.Context来判断errgroup.Group的执行情况

// 经典案例kratos 利用errgroup特性来优雅启动服务实例，并且监听服务实例启动情况: https://github.com/go-kratos/kratos/blob/main/app.go

func TestErrgroup(t *testing.T) {
	//eg := errgroup.Group{}
	eg, ctx := errgroup.WithContext(context.Background())
	var result int64 = 0
	for i := 0; i < 10; i++ {
		delta := i
		eg.Go(func() error {
			atomic.AddInt64(&result, int64(delta))
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
	ctx.Err() // 通过返回的err判断是超时结束还是业务逻辑错误
	fmt.Println(result)
}

func TestSourceCode(t *testing.T) {
	//type cancelCtx struct {
	//	Context
	//
	//	mu       sync.Mutex            // protects following fields
	//	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	//	children map[canceler]struct{} // set to nil by the first cancel call
	//	err      error                 // set to non-nil by the first cancel call
	//	cause    error                 // set to non-nil by the first cancel call
	//}

	// cancelCtx典型的装饰器模式：在已有Context的基础上，加上取消的功能
	// 核心实现:
	// - Done方法是通过类似double-check的机制写的。结合原子操作和锁(思考：能不能换成读写锁)
	// - 利用children来维护了所有的衍生节点

	//func (c *cancelCtx) Done() <-chan struct{} {
	//	d := c.done.Load()
	//	if d != nil {
	//	return d.(chan struct{})
	//}
	//	c.mu.Lock()
	//	defer c.mu.Unlock()
	//	d = c.done.Load()
	//	if d == nil {
	//	d = make(chan struct{})
	//	c.done.Store(d)
	//}
	//	return d.(chan struct{})
	//}

	//ctx := context.WithCancel(context.Background())

	//if p, ok := parentCancelCtx(parent); ok {
	//	p.mu.Lock()
	//	if p.err != nil {  找到最近的是cancelCtx类型的祖先，然后将child加进去祖先的children里面
	//		// parent has already been canceled
	//		child.cancel(false, p.err, p.cause)
	//	} else {
	//		if p.children == nil {
	//			p.children = make(map[canceler]struct{})
	//		}
	//		p.children[child] = struct{}{}
	//	}
	//	p.mu.Unlock()
	//} else {  找不到就只需要监听parent的信号，或者自己的信号，这些信号源自cancel或者超时
	//	goroutines.Add(1)
	//	go func() {
	//		select {
	//		case <-parent.Done():
	//			child.cancel(false, parent.Err(), Cause(parent))
	//		case <-child.Done():
	//		}
	//	}()
}

// timeCtx是装饰器模式：在已有cancelCtx的基础上增加超时功能
// - WithTimeout和WithDeadline 本质一样
// - WithDeadline 里面，在创建timeCtx的时候利用time.AfterFunc来实现超时

//type timerCtx struct {
//	*cancelCtx
//	timer *time.Timer // Under cancelCtx.mu.
//
//	deadline time.Time
//}

//func TestTimerCtx(t *testing.T) {
//	ctx,cancel := context.WithCancel(context.Background())
//}

// Go装饰器
type Cache interface {
	Get(key string) (string, error)
}

// 已有的，不是线程安全
// 要改造为线程安全
type memoryMap struct {
	m map[string]string
}

func (m *memoryMap) Get(key string) (string, error) {
	return m.m[key], nil
}

// 要改造为线程安全
// 无侵入式的改造
type SafeCache struct {
	Cache
	lock sync.RWMutex
}

func (s *SafeCache) Get(key string) (string, error) {
	s.lock.RLock()
	defer s.lock.Unlock()
	return s.Cache.Get(key)
}

var s = &SafeCache{
	Cache: &memoryMap{},
}

// 装饰器核心：在一个已有功能的基础上再增加一些功能

// 适配器: 适用于新老兼容，比如代码版本v2需要适配v1,防止修改代码后使用v1版本的代码调用出错
type OtherCache interface {
	GetValue(ctx context.Context, key string) (any, error)
}

type CacheAdapter struct {
	Cache
}

func (c *CacheAdapter) GetValue(ctx context.Context, key string) (any, error) {
	return c.Cache.Get(key)
}
