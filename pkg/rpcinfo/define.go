package rpcinfo

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
	"net"
)

// FlagBit 两位Byte
type FlagBit []byte

var (
	TraceBit FlagBit = []byte{0b1, 0b0}
)

const (
	TraceLength = 24
)

func HasFlagBit(b []byte, f FlagBit) bool {
	for i := 0; i < 2; i++ {
		if b[i]&f[i] != 0 {
			return true
		}
	}
	return false
}

func SetFlagBit(b []byte, f FlagBit) []byte {
	for i := 0; i < 2; i++ {
		b[i] |= f[i]
	}
	return b
}

//go:generate mockgen -source=define.go -destination ./define_mock.go -package rpcinfo

type FlagInfo interface {
	TraceInfo() []byte
}
type FlagInfoSetter interface {
	SetTraceInfo(traceInfo []byte)
}

// Invocation contains specific information about the call.
type Invocation interface {
	FlagInfo
	MethodInfo() service.MethodInfo
	ServiceImpl() any
	PackageName() string
	ServiceName() string
	MethodName() string
	SeqID() uint32
	BizError() uerror.BizError
}

// InvocationSetter is used to set information about an RPC.
type InvocationSetter interface {
	FlagInfoSetter
	SetMethodInfo(service.MethodInfo)
	SetServiceImpl(any)
	SetPackageName(name string)
	SetServiceName(name string)
	SetMethodName(name string)
	SetSeqID(seqID uint32)
	SetBizError(err uerror.BizError)
	SetExtra(key string, value interface{})
	Reset()
}

// EndpointInfo contains info for endpoint.
type EndpointInfo interface {
	ServiceName() string
	Method() string
	Address() net.Addr
	Tag(key string) (value string, exist bool)
	DefaultTag(key, def string) string
}

// InvokeInfo 包含了一次调用的信息
type InvokeInfo interface {
	// ServicePath 服务名
	ServicePath() string
	// MethodPath 调用的方法路径
	MethodPath() string
}

// RPCInfo is the core abstraction of information about an RPC in Kitex.
type RPCInfo interface {
	Invocation() Invocation
	EndpointContext() context.Context
	InteractionMode() InteractionMode

	Recycle()
}

// InteractionMode RPC请求的交互模式
type InteractionMode int32

const (

	// Oneway 单向交互，对端不会给予回复
	Oneway InteractionMode = iota + 1
	// Unary 双向交互，正常情况下对端会给予回复
	Unary
)
