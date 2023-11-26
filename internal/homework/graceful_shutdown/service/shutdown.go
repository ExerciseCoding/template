package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Option 典型的Option 设计模式
type Option func(*App)

// ShutdownCallback 采用context.Context来控制超时，不是使用time.After是因为
// - 超时的本质上是使用这个回调的人控制
// - 我们还希望用户知道，他的回调必须要在一定时间内处理完毕，而且必须显示处理超时错误
type ShutdownCallback func(ctx context.Context)

// WithShutdownCallbacks 需要实现的方法
func WithShutdownCallbacks(cbs ...ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}

type App struct {
	servers []*Server

	// 优雅退出整个超时时间，默认30秒
	shutdownTimeout time.Duration

	// 优雅退出的时候等待处理已有的请求时间，默认10秒钟
	waitTime time.Duration

	//  自定义回调超时时间，默认三分钟
	cbTimeout time.Duration

	cbs []ShutdownCallback
}

// NewApp 创建App实例，注意设置默认值，同时使用这些选项

func NewApp(servers []*Server, opts ...Option) *App {
	res := &App{
		waitTime:        10 * time.Second,
		cbTimeout:       3 * time.Second,
		shutdownTimeout: 30 * time.Second,
		servers:         servers,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

var signals = []os.Signal{
	os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
	syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
	syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM,
}

// StartAndServe 主要实现的方法
func (app *App) StartAndServe() {
	for _, s := range app.servers {
		srv := s
		go func() {
			if err := srv.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器已关闭")
				} else {
					log.Printf("服务器异常退出")
				}
			}
		}()
	}
	// 从这里开始优雅退出监听系统信号，强制退出以及超时强制退出
	// 优雅退出的具体步骤 shutdown里面实现
	// 所有你需要在这里恰当的位置，调用shutdown
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, signals...)
	<-ch
	print("hello")
	go func() {
		select {
		case <-ch:
			log.Printf("强制退出")
			os.Exit(1)
		case <-time.After(app.shutdownTimeout):
			log.Printf("超时强制退出")
			os.Exit(1)
		}
	}()
	app.shutdown()
}

func (app *App) shutdown() {
	log.Println("开始关闭应用，停止接收新请求")
	for _, s := range app.servers {
		// 思考：这里为什么可以不用并发控制，即不用锁也不用原子操作
		s.rejectReq()
	}
	// 需要在这里让所有的server拒绝信请求
	log.Println("等待正在执行请求完结")
	//这里可以改造为实时统计正在处理的请求数量，为0，则下一步
	time.Sleep(app.waitTime)

	// 在这里等待一段时间
	log.Println("开始关闭服务器")
	var wg sync.WaitGroup
	wg.Add(len(app.servers))
	// 并发关闭服务器，同时要注意协调所有的server都关闭之后才能步入下一个阶段
	for _, srv := range app.servers {
		srvcp := srv
		go func() {
			if err := srvcp.stop(context.Background()); err != nil {
				log.Printf("关闭服务失败:%s \n", srvcp.name)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("开始释放资源")
	// 并发执行回调，要注意协调所有的回调都执行完才会进入下一个阶段
	wg.Add(len(app.cbs))
	for _, cb := range app.cbs {
		c := cb
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.cbTimeout)
			c(ctx)
			cancel()
			wg.Done()
		}()
	}
	wg.Wait()
	//释放资源
	log.Println("开始释放资源")
	app.close()
}

func (app *App) close() {
	// 在这里释放掉一些可能得资源
	time.Sleep(time.Second)
	log.Println("应用关闭")
}

type Server struct {
	srv  *http.Server
	name string
	mux  *serverMux
}

// serverMux既可以看做是装饰器模式，也可以看做是委托模式
type serverMux struct {
	reject bool
	*http.ServeMux
}

func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只是在考虑到CPU高速缓存的时候，会存在短时间的不一致性
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

func NewServer(name string, addr string) *Server {
	mux := &serverMux{ServeMux: http.NewServeMux()}
	return &Server{
		name: name,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) rejectReq() {
	s.mux.reject = true
}

func (s *Server) stop(ctx context.Context) error {
	log.Printf("服务%s关闭中", s.name)
	return s.srv.Shutdown(ctx)
}
