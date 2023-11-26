package task

import (
	"context"
	"time"
)

// Task代表一个任务

type Task interface {
	// Run 执行任务
	// 如果ctx设置了超时时间，那么实现者需要自己决定是否进行超时控制
	Run(ctx context.Context) error
}

// TaskPool 任务池
type TaskPool interface {
	// Submit 执行一个任务
	// 如果任务池提供了阻塞的功能，那么如果在ctx过期都没有提交成功，那么应该返回错误
	// 调用Start之后能否继续提交任务，则取决于具体的实现
	// 调用Shutdown 或者 ShutdownNow之后提交任务都会返回错误
	Submit(ctx context.Context, task Task) error

	// Start 开始调度任务执行。在调用Start之前，所有的任务都不会被调度执行。
	// Start 之后，能否继续调用Sumbit提交任务，取决于具体的实现
	Start() error

	// Shutdown 关闭任务池。如果此时尚未调用Start方法，那么将会立刻返回。
	// 任务池将会停止接收新的任务，但是会继续执行剩下来的任务
	// 在所有任务执行完毕之后，用户可以从返回的chan中得到通知
	// 任务池在发出通知之后会关闭chan struct{}
	Shutdown() (<-chan struct{}, error)

	// ShutdownNow 立刻关闭线程池
	// 任务池能否中断当前正在执行的任务，取决于 TaskPool 的具体实现，以及 Task 的具体实现
	// 该方法会返回所有剩下的任务，剩下的任务是否包含正在执行的任务，也取决于具体的实现
	ShutdownNow() ([]Task, error)

	// States 暴露 TaskPool 生命周期内的运行状态
	// ctx 是让用户来控制什么时候退出采样。那么最基本的两个退出机制：一个是 ctx 被 cancel 了或者超时了，一个是TaskPool 被关闭了
	// error 仅仅表示创建 chan state 是否成功
	// interval 表示获取TaskPool运行期间内部状态的周期/时间间隔
	States(ctx context.Context, interval time.Duration) (<-chan State, error)
}

type State struct {
	PoolState       int32
	GoCnt           int32
	WaitingTasksCnt int
	QueueSize       int
	RunningTasksCnt int32
	Timestamp       int64
}
