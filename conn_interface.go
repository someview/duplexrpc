package rpconeway

import (
	"context"
	"io"
)

type InvokeOneWay func(ctx context.Context, req any) error

type MessageHandler func(ctx context.Context, req any) error

type ConnInterface interface {
	Invokes() []InvokeOneWay
	MessageHandlers() []MessageHandler
}

type GoPrivateService interface {
}

type GoPrivateService2 interface {
}

type req struct{}

type res struct{}

type UnaryConnInterface interface {
	ConnInterface
}

type Stream interface {
	// Context the stream context.Context
	Context() context.Context
	// RecvMsg recvive message from peer
	// will block until an error or a message received
	// not concurrent-safety
	RecvMsg(m interface{}) error
	// SendMsg send message to peer
	// will block until an error or enough buffer to send
	// not concurrent-safety
	SendMsg(m interface{}) error
	// not concurrent-safety with SendMsg
	io.Closer
}
