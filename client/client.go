package client

import (
	"bufio"
	"context"
	"errors"
	"rpc-oneway/protocol"
)

var (
	ErrShutdown         = errors.New("connection is shut down")
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 16 * 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 16 * 1024
)

type MuxClient struct {
	// 这是和连接相关的部分
	closing bool // whether the server is going to close this connection

	option Option
	r      *bufio.Reader
	recv   RecvFunc
}

// Send 用户层调用接口
func (c *MuxClient) Send(ctx context.Context, msgType byte, req any) error {
	// todo ctx
	if c.closing {
		return ErrShutdown
	}
	payload, ok := req.(protocol.SizeableMarshaller)
	if !ok {
		return ErrUnsupportedCodec
	}

	msg := protocol.NewMessage()
	buf:= make([]byte, msg)
	payload.MarshalToSizedBuffer()
	msg.Payload =
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
