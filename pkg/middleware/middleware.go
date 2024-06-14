package middleware

import (
	"context"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
)

// MiddleWare is a function that wraps a Next to control the flow of the request.
// for example
// this is a MiddleWare function:
//
//		func middleware(next Next) Next {
//		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
//	 1. next(ctx, req, cb) // it will pass the request, you must be put it in the last line of the MiddleWare function(before return)
//	 2. cb(fmt.Error("myMiddlewareError")) // it will abort the request, and return error to request callback
//	 3. next(ctx, req, func(err error){cb(err)}) // it will pass the request, and you can do something after the request is done
//	    }
//
// }
type MiddleWare func(next Next) Next

// Next is a function that is called to pass the request.
type Next func(ctx context.Context, req protocol.Message, cb netpoll.CallBack)

// Chain chains multiple MiddleWare functions into one.
func Chain(mws ...MiddleWare) MiddleWare {
	if len(mws) == 0 {
		return DummyMiddleWare
	}
	return func(next Next) Next {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}

// DummyMiddleWare is a dummy MiddleWare function that does nothing.
func DummyMiddleWare(next Next) Next {
	return next
}
