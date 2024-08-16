package metadata

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/internal"
	"sync"
	"time"
)

const (
	valuesInOneContext = 10
)

var valueCtxPool = sync.Pool{New: func() any {
	return &valueContext{}
}}

type valueContext struct {
	parent context.Context
	idx    int
	keys   [valuesInOneContext]any
	values [valuesInOneContext]any
}

func newValueContext(parent context.Context) *valueContext {
	valueCtx := valueCtxPool.Get().(*valueContext)
	valueCtx.parent = parent
	return valueCtx
}

func WithValue(parent context.Context, key, val any) context.Context {
	valueCtx, ok := parent.(*valueContext)
	if !ok || valueCtx.idx >= valuesInOneContext {
		valueCtx = newValueContext(parent)
	}
	valueCtx.keys[valueCtx.idx] = key
	valueCtx.values[valueCtx.idx] = val
	valueCtx.idx++
	return valueCtx
}

func (v *valueContext) Deadline() (deadline time.Time, ok bool) {
	return v.parent.Deadline()
}

func (v *valueContext) Done() <-chan struct{} {
	return v.parent.Done()
}

func (v *valueContext) Err() error {
	return v.parent.Err()
}

func (v *valueContext) Value(key any) any {
	for i := 0; i < len(v.keys); i++ {
		if v.keys[i] == key {
			return v.values[i]
		}
	}
	return v.parent.Value(key)
}

func (v *valueContext) Recycle() {
	v.idx = 0
	v.keys = [valuesInOneContext]any{}
	v.values = [valuesInOneContext]any{}
	RecycleContext(v.parent)
	v.parent = nil

	valueCtxPool.Put(v)
}

func RecycleContext(ctx context.Context) {
	if v, ok := ctx.(internal.Reusable); ok {
		v.Recycle()
	}
}
