package _chan


// 控制goroutine的数量
// 应用场景：任务分离

type TashPool struct {
	ch chan struct{}
}


func (t *TashPool) Do(f func()) {
	token := <- t.ch
	// 异步执行
	go func() {
		f()
		t.ch <- token
	}()
	// 同步执行
	//f()
	//t.ch <- token

}


func NewTaskPool(limit int) *TashPool {
	t := &TashPool{
		ch: make(chan struct{}, limit),
	}
	// 提前转呗好了令牌
	for i := 0; i < limit; i++ {
		t.ch <- struct{}{}
	}
	return t
}


type TaskPoolWithCache struct {
	cache chan func()
}


func NewTaskPoolWithCache(limit int, cacheSize int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		cache: make(chan func(), cacheSize),
		//ch: make(chan struct{}, limit),
	}
	// 直接把goroutine 开好
	for i := 0; i < limit; i++ {
		go func() {
			for {
				// 在goroutine 里面不断尝试cache 里面拿到任务
				select {
				case task, ok := <-t.cache:
					if !ok {
						return
					}
					task()
				}
			}
		}()
	}
	return t
}


func (t *TaskPoolWithCache) Do(f func()) {
	t.cache <- f
}


// 显式控制生命周期

//func (t *TaskPoolWithCache) Start() {
//	for i := 0; i < limit; i++ {
//		go func() {
//			for {
//				// 在goroutine 里面不断尝试cache 里面拿到任务
//				select {
//				case task, ok := <-t.cache:
//					if !ok {
//						return
//					}
//					task()
//				}
//			}
//		}()
//	}
//	return t
//}

