package middleware

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"sync/atomic"
	"testing"
)

func TestMiddleWare(t *testing.T) {
	init := Next(func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
	})

	c := &atomic.Int32{}
	mw1 := MiddleWare(func(next Next) Next {
		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
			c.Add(1)
			next(ctx, req, cb)
		}
	})
	mw2 := MiddleWare(func(next Next) Next {
		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
			c.Add(3)
			next(ctx, req, cb)
		}
	})
	mw := Chain(mw1, mw2)
	w := mw(init)
	w(nil, nil, func(err error) {
		assert.NoError(t, err)
	})
	assert.Equal(t, c.Load(), int32(4))

	mw3 := MiddleWare(func(next Next) Next {
		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
			cb(fmt.Errorf("error"))
		}
	})

	mw = Chain(mw, mw3)
	w = mw(init)
	w(nil, nil, func(err error) {
		assert.Error(t, err)
	})
	assert.Equal(t, c.Load(), int32(8))
}

func BenchmarkMiddleware(b *testing.B) {
	mw1 := MiddleWare(func(next Next) Next {
		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
			next(ctx, req, cb)
		}
	})
	mws := make([]MiddleWare, 0, 10)
	for i, _ := range mws {
		mws[i] = mw1
	}
	mw := Chain(mws...)
	m := mw(func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
		cb(nil)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m(nil, nil, func(err error) {

		})
	}
}
