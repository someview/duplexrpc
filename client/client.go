package client

import (
	"bufio"
	"context"
	"errors"
	"net"

	"rpc-oneway/protocol"

	"rpc-oneway/pkg/breaker"
)

var (
	ErrShutdown    = errors.New("connection is shut down")
	ErrBreakerOpen = errors.New("breaker is open")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBufSize = 16 * 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 16 * 1024
)

type MuxClient struct {
	// 这是和连接相关的部分
	closing bool // whether the server is going to close this connection

	breaker breaker.Breaker
	net.Conn
	option Option
	r      *bufio.Reader

	ServerMessageChan chan<- protocol.Message
}

// Address implements RPCClient.
func (*MuxClient) Address() string {
	panic("unimplemented")
}

// IsShutdown implements RPCClient.
func (*MuxClient) IsShutdown() bool {
	panic("unimplemented")
}

// ServerMessageHandler implements RPCClient.
func (*MuxClient) ServerMessageHandler() func(protocol.Message) {
	panic("unimplemented")
}

func NewMuxClient() *MuxClient {
	return &MuxClient{}
}

// Send 用户层调用接口
func (c *MuxClient) Send(cc context.Context, req any) error {
	if c.closing {
		return ErrShutdown
	}
	// todo 设置成中间件的形式
	if c.breaker != nil && !c.breaker.Allow() {
		return ErrBreakerOpen
	}
	msg := protocol.NewMessage()
	ctx := cc.(protocol.MsgContext)
	msg.SetMsgType(ctx.MsgType())
	msg.SetReq(req)

	allData, err := msg.Encode()
	if err != nil {
		return err
	}

	if deadline, ok := ctx.Deadline(); ok {
		_ = c.Conn.SetWriteDeadline(deadline)
	}
	_, err = c.Conn.Write(*allData)
	// todo 这个地方非常突兀，需要使用linked buffer类似的思路来替代掉
	protocol.PutData(allData)
	return err
}
