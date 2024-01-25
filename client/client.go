package client

import (
	"bufio"
	"context"
	"errors"
	"net"
	"rpc-oneway/protocol"
)

var (
	ErrShutdown = errors.New("connection is shut down")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBufsize = 16 * 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 16 * 1024
)

type MuxClient struct {
	// 这是和连接相关的部分
	closing bool // whether the server is going to close this connection

	net.Conn
	option Option
	r      *bufio.Reader

	ServerMessageChan chan<- *protocol.Message
}

// Send 用户层调用接口
func (c *MuxClient) Send(ctx context.Context, msgType byte, req any) error {
	// todo ctx
	if c.closing {
		return ErrShutdown
	}

	msg := protocol.NewMessage()
	msg.MsgType = msgType
	msg.Data = req

	allData, err := msg.EncodeSlicePointer()
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
