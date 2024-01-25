package client

import (
	"context"
	"rpc-oneway/protocol"
)

// RPCClient is interface that defines one client to call one server.
type RPCClient interface {
	Connect(network, address string) error
	Send(ctx context.Context, msgType int32, msg any) error
	Recv() chan *protocol.Message
	Close() error
	RemoteAddr() string
	IsClosing() bool
	IsShutdown() bool
}
