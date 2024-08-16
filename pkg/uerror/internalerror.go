package uerror

import (
	"errors"
	"fmt"
	"syscall"

	netpoll "github.com/cloudwego/netpoll"
)

const (
	TransErrCode = 2
	CodecErrCode = 10

	ServiceMethodNotFindCode = 11
)

// InternalError internalError 用于表示所有内部错误（网络、传输、端点、客户端选择器）
type InternalError struct {
	BasicError
	cause error
}

func (e *InternalError) ErrType() ErrorCategory { return e.errorCategory }
func (e *InternalError) Code() int32            { return e.code }
func (e *InternalError) Message() string        { return e.message }
func (e *InternalError) Cause() error           { return e.cause }
func (e *InternalError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("ErrType: %d, code: %d, message: %s, cause: %v", e.errorCategory, e.code, e.message, e.cause)
	}
	return fmt.Sprintf("ErrType: %d, code: %d, message: %s", e.errorCategory, e.code, e.message)
}
func (e *InternalError) WithCause(cause error) error {
	newErr := *e
	newErr.cause = cause
	return &newErr
}

// NewInternalError 创建一个新的通用内部错误
func NewInternalError(tp ErrorCategory, code int32, message string) *InternalError {
	return &InternalError{
		BasicError: BasicError{
			errorCategory: tp,
			code:          code,
			message:       message,
		},
	}
}

// NewNetworkError 创建一个新的网络级别错误
func NewNetworkError(code int32, message string) *InternalError {
	return NewInternalError(ErrNetwork, code, message)
}

// NewTransError 创建一个新的传输级别错误
func NewTransError(code int32, message string) *InternalError {
	return NewInternalError(ErrTrans, code, message)
}

func NewServiceMethodNotFindError(code int32, message string) *InternalError {
	return NewInternalError(ErrServiceMethodNotFind, code, message)
}

// NewEndpointError 创建一个新的端点级别错误
func NewEndpointError(code int32, message string) *InternalError {
	return NewInternalError(ErrEndpoint, code, message)
}

// NewClientSelectorError 创建一个新的客户端选择器错误
func NewClientSelectorError(code int32, message string) *InternalError {
	return NewInternalError(ErrClientSelector, code, message)
}

// IsRemoteClosedErr 连接级别错误
func IsRemoteClosedErr(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, netpoll.ErrConnClosed) || errors.Is(err, syscall.EPIPE)
}

func IsTransError(err error) bool {
	for err != nil {
		switch v := err.(type) {
		case *BasicError:
			return v.errorCategory == ErrTrans
		case *InternalError:
			return v.errorCategory == ErrTrans
		case interface{ Unwrap() error }:
			err = v.Unwrap()
		default:
			return false
		}
	}
	return false
}
