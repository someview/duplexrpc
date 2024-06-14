package server

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

type RPCServer interface {
	Run(network, address string) error
	Stop(ctx context.Context) error
	RegisterMethod(mInfo rpcinfo.MethodInfo, handler any) error
}
