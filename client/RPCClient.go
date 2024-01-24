package client

import (
	"context"
)

// RPCClient is interface that defines one client to call one server.
type RPCClient interface {
	Connect(network, address string) error
	Send(ctx context.Context, msgType int32, msg any) error
	Recv(ctx context.Context, msgType int32, msg any) error
	Close() error
	RemoteAddr() string
	IsClosing() bool
	IsShutdown() bool
}

type SendFunc func(ctx context.Context, msgType int32, msg any) error
type RecvFunc func(ctx context.Context, msgType int32, msg any) error
