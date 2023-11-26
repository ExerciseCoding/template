package pool

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MyPool struct {
	p sync.Pool
	maxCnt int32
	cnt int32
}

func (p *MyPool) Get() any {
	return p.p.Get()
}

func (p *MyPool) Put(val any) {
	// 大对象不放回去
	if unsafe.Sizeof(val) > 1024 {
		return
	}
	// 超过数量不放回去
	// 这里明面上控制了数量，其实有可能被GC已经清理掉，所以这里的数量是不准确的
	cnt := atomic.AddInt32(&p.cnt, 1)
	if cnt >= p.maxCnt {
		atomic.AddInt32(&p.cnt, -1)
		return
	}
	p.p.Put(val)
}

//type Pool struct {
//	calls       [steps]uint64
//	calibrating uint64
//
//	defaultSize uint64
//	maxSize     uint64
//
//	pool sync.Pool
//}

// 开源实例- bytebufferpool 实现要点
// Github: https://github.com/valyala/bytebufferpool
// 也是依托于sync.Pool进行二次封装
// defaultSize 是每次创建的buffer的默认大小，超过maxSize的buffer就不会放回去
// 统计不同大小的buffer的使用次数，例如0-64, bytes的buffer被使用了多少次。这个我们称为分组统计使用次数
// 引入了所谓的校准机制，其实就是动态计算 defaultSize 和maxSize



// bytebufferpool就根据使用次数来决定：
// - 新创建的多大
// - 超过多大的就没必要放回去

//func (p *Pool) Put(b *ByteBuffer) {
//	idx := index(len(b.B))
//
//	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold {
//		p.calibrate()
//	}
//
//	maxSize := int(atomic.LoadUint64(&p.maxSize))
//	if maxSize == 0 || cap(b.B) <= maxSize {
//		b.Reset()
//		p.pool.Put(b)
//	}
//}

