package generic

import (
	"context"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/callopt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
)

type Closeable interface {
	Close() error
}

// RespCallBack is callback function when receive response.
//
// if you Call with oneway, it will called after send to server, first param is nil.
type RespCallBack func(any, error)

// XClient 泛化接口, msgType用于peer之间传输的消息类型, 由应用本身协商即可
type XClient interface {
	Call(ctx context.Context, req any, cb RespCallBack, opt ...CallOpt)
	Closeable
}

// RPCClient is interface that defines one client to call one server.
type RPCClient interface {
	AsyncConnect(network, address string) error
	AsyncSend(req protocol.Message, cb netpoll.CallBack)
	IsShutdown() bool
	Address() string
	Closeable
}

type GenericClient interface {
	Call(ctx context.Context, method string, req any, opt ...callopt.CallOpt) (res any, err error)
}

type GenericInfo interface {
	GetMethod(method string) rpcinfo.MethodInfo
}
