package client

import (
	"bufio"
	"errors"
	"net"

	"rpc-oneway/protocol"
)

var (
	ErrShutdown = errors.New("connection is shut down")
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

	net.Conn
	option Option
	r      *bufio.Reader

	ServerMessageChan chan<- protocol.Message
}

// Send 用户层调用接口
func (c *MuxClient) Send(ctx protocol.MsgContext, req any) error {
	if c.closing {
		return ErrShutdown
	}

	msg := protocol.NewMessage()
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
