package remote

import (
	"context"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

//go:generate mockgen -source=define.go -destination ./define_mock.go -package remote

type SizeableMarshaller interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

type Codec interface {
	Encode(ctx context.Context, writer netpoll.Writer, msg Message) error
	Decode(ctx context.Context, reader netpoll.Reader, msg Message) error
	Name() string
}

// TransHandler 双向Call的rpc框架, 对于单个连接上的处理方法，必须是相同的
type TransHandler interface {
	Write(ctx context.Context, end Endpoint, send Message) (err error)
	Read(ctx context.Context, end Endpoint, recv Message) (err error)
}

type ConnectionHandler interface {
	OnError(ctx context.Context, conn netpoll.Connection, err error)
	OnActive(ctx context.Context, conn netpoll.Connection) (context.Context, error)
	OnInactive(ctx context.Context, conn netpoll.Connection)
	OnRead(ctx context.Context, conn netpoll.Connection) error
}

// GracefulShutdown supports closing connections in a graceful manner.
type GracefulShutdown interface {
	GracefulShutdown(ctx context.Context) error
}

type Endpoint interface {
	Address() string
	Context() context.Context
	CallBackManager() CallBackManager
	netpoll.Connection
	bufwriter.BatchWriter
}

type CallBackManager interface {
	Set(seqID uint32, reader chan netpoll.Reader)
	Delete(seqID uint32)
	Load(seqID uint32) (reader chan netpoll.Reader, ok bool)
	LoadAndDelete(seqID uint32) (reader chan netpoll.Reader, ok bool)
}

type Message interface {
	RPCInfo() rpcinfo.RPCInfo

	SetMetadata(metadata []byte)
	Metadata() []byte

	// 原始请求
	Data() any
	SetData(data any)
	// message转换后的字节流
	SetPayload(reader netpoll.Reader)
	Payload() netpoll.Reader // 编码后的二进制数据

	Version() Version
	SetVersion(v Version)
	SetMsgType(typ MessageType)
	MsgType() MessageType

	Len() int
	SetLen(l int)

	Codec() Codec

	Recycle()
}

type (
	// MessageType 消息类型
	MessageType     byte
	Version         byte
	ParallelDecider func(svcName, methodName string) bool
)

const (
	OnewayType MessageType = iota + 1
	ReqType
	ResponseType
	ErrorResponseType
	PingType
	PongType
	CloseType
)
