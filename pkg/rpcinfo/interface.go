package rpcinfo

import (
	"context"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
)

// MethodInfo to record meta info of unary method
//	MethodInfo是为了将Method和RPCServer-Client解耦分离，可以通过一个统一的包来定义Method
//	使用Method来约束Server-Client，而不是通过Server-Client去规定Method

type FuncCall func(ctx context.Context, req any) (res any, err error)

type MethodInfo interface {
	// Invoke 定义如何将args传入注册的handler中,且返回值如何传入result中
	Invoke(fn any, ctx context.Context, args, result any) error
	// NewArg NewArgs 生成一个新的args
	NewArg() any
	// NewResult 生成一个新的result
	NewResult() any
	// OneWay 是否是单向调用
	OneWay() bool
	// Type method对应的消息类型,和Name形成一一映射,在泛化调用的时候
	Type() MethodType
	Name() string
}

type MethodType byte

// InvokeInfo 包含了一次调用的信息
type InvokeInfo interface {
	// MethodName 调用的方法名
	MethodName() MethodType
	From() RPCEndpoint
	To() RPCEndpoint
	// SeqID 序列号
	SeqID() uint32
}

// RPCEndpoint 调用双方的信息
type RPCEndpoint interface {
	// Context 这个context接口是配合netpoll的API使用的
	context.Context
	netpoll.Connection
	// AsyncWriter 返回一个异步writer
	AsyncWriter() muxextend.AsyncWriter
}
