package protocol

import (
	"context"
	"errors"
	"io"
	"net"

	"rpc-oneway/util"
)

type FrameType byte

var (
	bufferPool          = util.NewLimitedPool(512, 4096)
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

const (
	InitialType     FrameType = iota
	InitialDoneType           = 1
	DataType                  = 2
	PingType                  = 6
	PongType                  = 7
	CloseType                 = 8
)

// data类型的帧的结构定义

type Message interface {
	MsgContext
	SetReq(req any)
	Payload() []byte // 编码后的二进制数据
	Len() uint32
	// todo how to make this readable
	Encode() (*[]byte, error)
	Decode(io.Reader) error
	Recycle()
}

// todo 是否需要To方法
type MsgContext interface {
	From() net.Conn
	SetFrom(net.Conn)

	TraceInfo() []byte
	SetTraceInfo([]byte)

	MsgType() byte
	SetMsgType(byte)

	SetLBKey(uint32)
	Metadata() map[string]string
	context.Context // 原始的ctx信息,用于客户端超时控制
}
