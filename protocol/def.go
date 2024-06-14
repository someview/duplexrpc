package protocol

import (
	"context"
	"errors"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

type FrameType byte

var (
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

type FlagBit byte

const (
	TraceBit FlagBit = 0b1
	SeqBit   FlagBit = 0b10
)

const (
	// FixHeaderLength EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
	//  [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
	FixHeaderLength = 7
	TraceLength     = 24
	SeqLength       = 4
)

func HasFlagBit(b byte, f FlagBit) bool {
	switch f {
	case TraceBit:
		return b&byte(TraceBit) != 0
	case SeqBit:
		return b&byte(SeqBit) != 0
	default:
		return false
	}
}

func SetFlagBit(b byte, f FlagBit) byte {
	switch f {
	case TraceBit:
		return b | byte(TraceBit)
	case SeqBit:
		return b | byte(SeqBit)
	default:
		return b
	}
}

// data类型的帧的结构定义

type Message interface {
	MsgContext
	muxextend.SizedEncodable
	muxextend.SizedDecodable
	SetReq(req any)
	GetReq() any
	Payload() netpoll.Reader // 编码后的二进制数据
	Len() uint32
	Recycle()
}

type MsgContext interface {
	context.Context // 原始的ctx信息,用于客户端超时控制

	rpcinfo.InvokeInfo

	TraceInfo() []byte
	SetTraceInfo([]byte)

	MsgType() byte
	SetMsgType(byte)

	SetLBKey(uint32)
	Metadata() map[string]string
}

// Parser 用于解析RPC-Header
type Parser interface {
	// ParseHeader 解析header
	ParseHeader(r netpoll.Reader, from rpcinfo.RPCEndpoint) (Message, error)
}

// SeqIdGenerator 用于生成SeqId
type SeqIdGenerator interface {
	// Next 并发安全
	Next() uint32
}
