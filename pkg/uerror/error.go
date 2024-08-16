package uerror

type ErrorCategory int32

const (
	// ErrNetwork 连接级别错误, tcp字节流到buffer, 或者buffer到tcp字节流错误、连接错误
	ErrNetwork ErrorCategory = 1

	// ErrTrans ErrCodec trans codec  qps limiter ..., 处理从buffer到protocol.message的转换，以及tcp连接上错的各种事件
	ErrTrans ErrorCategory = 4

	// ErrEndpoint 在单个endpoint上处理消息时，发生错误, 限流错误，无实际可用单端等等
	ErrEndpoint ErrorCategory = 7

	// ErrClientSelector 客户端选择器错误
	ErrClientSelector ErrorCategory = 11

	ErrServiceMethodNotFind ErrorCategory = 15

	ErrBiz ErrorCategory = 20
)

type BaseError interface {
	error
	ErrType() ErrorCategory
	Code() int32
	Message() string
}

type FrameworkError interface {
	BaseError
	Cause() error
}

type BizError interface {
	BaseError
}
