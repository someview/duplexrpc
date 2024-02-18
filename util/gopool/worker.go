package gopool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var workerPool sync.Pool

type worker struct {
	pool *pool
}

func init() {
	workerPool.New = newWorker
}

func newWorker() interface{} {
	return &worker{}
}

func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}
			if t == nil {
				// if there's no task to do, exit
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}

			w.pool.taskLock.Unlock()
			func() {
				defer func() {
					if r := recover(); r != nil {
						msg := fmt.Sprintf("GOPOOL: panic in pool: %s: %v: %s", w.pool.name, r, debug.Stack())
						//logger.CtxErrorf(t.ctx, msg)
						fmt.Println("msg:", msg)
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						}
					}
				}()
				t.f()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) close() {
	w.pool.decWorkerCount()
}
func (w *worker) Recycle() {
	w.zero()
	workerPool.Put(w)
}

func (w *worker) zero() {
	w.pool = nil
}
