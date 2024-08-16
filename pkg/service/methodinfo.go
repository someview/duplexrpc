package service

import (
	"context"
)

// MethodInfo to record meta info of unary method
type MethodInfo interface {
	Handler() MethodHandler

	OneWay() bool
	ServerCall() bool

	ArgsFactory() ArgFactory
	ResultFactory() ArgFactory
}

// MethodHandler is corresponding to the serviceImpl wrapper func that in generated code
type MethodHandler func(handler any, ctx context.Context, args, result interface{}) error

// NewMethodInfo is called in generated code to build method info

func NewMethodInfo(methodHandler MethodHandler, argsFactory, resultFactory ArgFactory, oneWay bool, isServerCall bool) MethodInfo {

	if !oneWay && resultFactory == nil {
		panic("unexpected type for newResult, it should not be nil")
	}

	info := methodInfo{
		handler:     methodHandler,
		argsFactory: argsFactory,

		resultFactory: resultFactory,

		oneWay:     oneWay,
		serverCall: isServerCall,
	}

	return info

}

type methodInfo struct {
	handler     MethodHandler
	argsFactory ArgFactory

	resultFactory ArgFactory

	oneWay     bool
	serverCall bool
}

func (m methodInfo) ArgsFactory() ArgFactory {
	return m.argsFactory
}

func (m methodInfo) ResultFactory() ArgFactory {
	return m.resultFactory
}

// Handler implements the MethodInfo interface.
func (m methodInfo) Handler() MethodHandler {
	return m.handler
}

// OneWay implements the MethodInfo interface.
func (m methodInfo) OneWay() bool {
	return m.oneWay
}

func (m methodInfo) ServerCall() bool {
	return m.serverCall
}
