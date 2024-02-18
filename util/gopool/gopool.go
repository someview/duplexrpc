package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	// Name returns the corresponding pool name.
	Name() string
	// SetCap sets the goroutine capacity of the pool.
	SetCap(cap int32)
	// Go executes f.
	Go(f func())
	// CtxGo executes f and accepts the context.
	CtxGo(ctx context.Context, f func())
	// SetPanicHandler sets the panic handler.
	SetPanicHandler(f func(context.Context, interface{}))
}

var taskPool sync.Pool

const (
	defaultScalaThreshold = 1
)

func init() {
	taskPool.New = newTask
}

type Config struct {
	// threshold for scale.
	// new goroutine is created if len(task chan) > ScaleThreshold.
	// defaults to defaultScalaThreshold.
	ScaleThreshold int32
}

type pool struct {
	name   string
	cap    int32
	config *Config

	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	// Record the number of running workers
	workerCount int32
	// This method will be called when the worker panic
	panicHandler func(context.Context, interface{})
}

// CtxGo implements Pool.
func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)
	// The following two conditions are met:
	// 1. the number of tasks is greater than the threshold.
	// 2. The current number of workers is less than the upper limit p.cap.
	// or there are currently no workers.
	if (atomic.LoadInt32(&p.taskCount) >= p.config.ScaleThreshold && p.WorkerCount() < atomic.LoadInt32(&p.cap)) || p.WorkerCount() == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

// Go implements Pool.
func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

// Name implements Pool.
func (p *pool) Name() string {
	return p.name
}

// SetCap implements Pool.
func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

// SetPanicHandler implements Pool.
func (*pool) SetPanicHandler(f func(context.Context, interface{})) {
	panic("unimplemented")
}

func NewPool(name string, cap int32, config *Config) Pool {
	p := &pool{
		name:   name,
		cap:    cap,
		config: config,
	}
	return p
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
