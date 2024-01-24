package client

import (
	"context"
	"errors"
	"net"
	"rpc-oneway/protocol"
)

var (
	ErrShutdown         = errors.New("connection is shut down")
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

type muxClient struct {
	// 这是和连接相关的部分
	closing bool // whether the server is going to close this connection
	net.Conn
}

// OnRequest is called when the connection creates.

func (c *muxClient) handleServerRequest() {

}

// Send 用户层调用接口
func (c *muxClient) Send(ctx context.Context, msgType byte, req any) error {
	// todo ctx
	if c.closing {
		return ErrShutdown
	}
	payload, ok := req.(protocol.SizeableMarshaller)
	if !ok {
		return ErrUnsupportedCodec
	}

	msg := protocol.NewMessage()
	msg.Payload = payload
	msg.MsgType = msgType
	allData, err := msg.EncodeSlicePointer()
	if err != nil {
		return err
	}

	// todo 抉择是否需要判断
	//if ctx.Err() != nil {
	//	return ctx.Err()
	//}
	if deadline, ok := ctx.Deadline(); ok {
		_ = c.Conn.SetWriteDeadline(deadline)
	}

	_, err = c.Conn.Write(*allData)
	protocol.PutData(allData)
	return err
}



func(c *muxClient)