package middleware

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

type CallMiddleware func(next GenericCall) GenericCall

type GenericCall func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error)

func Chain(mws ...CallMiddleware) CallMiddleware {
	return func(next GenericCall) GenericCall {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}
