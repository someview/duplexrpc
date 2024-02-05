package client

import (
	"context"

	"rpc-oneway/protocol"
)

type Closeable interface {
	Close() error
}

// XClient 泛化接口, msgType用于peer之间传输的消息类型, 由应用本身协商即可
type XClient interface {
	Send(ctx context.Context, args any) error
	ServerMessageHandler() func(protocol.Message) // you1 should release protocol.Message after no longer used it
	Closeable
}

// RPCClient is interface that defines one client to call one server.
type RPCClient interface {
	Connect(network, address string) error
	IsShutdown() bool
	Address() string
	XClient
}
