package gopool

import (
	"context"
	"sync"
)

type task struct {
	ctx  context.Context
	f    func()
	next *task
}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

func newTask() interface{} {
	return &task{}
}

type TaskList struct {
	sync.Mutex
	taskHead *task
	taskTail *task
}
