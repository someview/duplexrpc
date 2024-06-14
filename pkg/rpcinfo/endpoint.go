package rpcinfo

import (
	"context"
	"errors"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
)

type rpcEndpoint struct {
	context.Context
	netpoll.Connection
	wr muxextend.AsyncWriter
}

func NewRPCEndpoint(conn netpoll.Connection, wr muxextend.AsyncWriter) RPCEndpoint {
	r := &rpcEndpoint{
		Context:    context.TODO(),
		Connection: conn,
		wr:         wr,
	}
	return r
}

func (r *rpcEndpoint) AsyncWriter() muxextend.AsyncWriter {
	return r.wr
}

func (r *rpcEndpoint) Close() error {
	err1 := r.Connection.Close()
	err2 := r.wr.Close()
	return errors.Join(err1, err2)
}
