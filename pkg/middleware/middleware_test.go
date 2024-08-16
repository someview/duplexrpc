package middleware

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"testing"
)

func TestCallMiddleware_Normal(t *testing.T) {
	idx := 0
	var mws1 CallMiddleware
	var mws2 CallMiddleware
	var mws3 CallMiddleware
	mws1 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			idx++
			ctx = context.WithValue(ctx, 1, 1)
			return next(ctx, req, resp, info)
		}

	}
	mws2 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			idx++
			assert.NotNil(t, ctx.Value(1))
			ctx = context.WithValue(ctx, 2, 2)
			return next(ctx, req, resp, info)
		}
	}
	mws3 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			idx++
			assert.NotNil(t, ctx.Value(1))
			assert.NotNil(t, ctx.Value(2))
			return next(ctx, req, resp, info)
		}
	}

	ivFunc := Chain(mws1, mws2, mws3)(func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
		idx++
		return nil
	})

	err := ivFunc(context.TODO(), nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, idx)

}

func TestCallMiddleware_Abort(t *testing.T) {
	mws1Error := errors.New("error1")
	mws2Error := errors.New("error2")
	var mws1 CallMiddleware
	var mws2 CallMiddleware
	var mws3 CallMiddleware
	mws1 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			return errors.Join(next(ctx, req, resp, info), mws1Error)

		}

	}
	mws2 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			// abort
			return mws2Error

		}
	}
	mws3 = func(next GenericCall) GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			assert.Failf(t, "should not be called", "should not be called")
			return nil
		}
	}

	ivFunc := Chain(mws1, mws2, mws3)(func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
		assert.Failf(t, "should not be called", "should not be called")
		return nil
	})

	err := ivFunc(context.TODO(), nil, nil, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, mws2Error))
	assert.True(t, errors.Is(err, mws1Error))
}
