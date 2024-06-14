package rpcinfo

import (
	"context"
	"fmt"
	"reflect"
)

// InvokeFunc 定义如何将args传入注册的handler中,且返回值如何传入result中
type InvokeFunc[Arg, Res, H any] func(fn H, ctx context.Context, args Arg, result Res) error

// InvokeOnewayFunc 定义如何将args传入注册的handler中
type InvokeOnewayFunc[Arg, H any] func(fn H, ctx context.Context, args Arg) error

// Getter 是请求/响应的参数生成器
type Getter[T any] func() T

type onewayMethodInfo[Arg, H any] struct {
	methodInfo[Arg, Arg, H]
}

// NewOnewayMethodInfo 创建一个oneway的MethodInfo
func NewOnewayMethodInfo[Arg, H any](typ MethodType, newArgs Getter[Arg], handler InvokeOnewayFunc[Arg, H]) MethodInfo {
	info := &onewayMethodInfo[Arg, H]{
		methodInfo: *newMethodInfo[Arg, Arg, H](typ, newArgs, nil, handler, true),
	}
	return info
}

func (o *onewayMethodInfo[Arg, H]) NewResult() any {
	return nil
}

func (o *onewayMethodInfo[Arg, H]) Invoke(handler any, ctx context.Context, args any, result any) error {
	return o.methodHandler.(InvokeOnewayFunc[Arg, H])(handler.(H), ctx, args.(Arg))
}

// NewMethodInfo 创建一个request-response的MethodInfo
func NewMethodInfo[Arg, Res, H any](typ MethodType, newArgs Getter[Arg], newResult Getter[Res], handler InvokeFunc[Arg, Res, H]) MethodInfo {
	info := newMethodInfo[Arg, Res, H](typ, newArgs, newResult, handler, false)
	return info
}

type methodInfo[Arg, Res, H any] struct {
	typ           MethodType
	newArgs       Getter[Arg]
	newResult     Getter[Res]
	methodHandler any
	oneway        bool
}

func (o *methodInfo[Arg, Res, H]) NewArg() any {
	return nil
}

func (o *methodInfo[Arg, Res, H]) Name() string {
	return ""
}

func newMethodInfo[Arg, Res, H any](typ MethodType, newArgs Getter[Arg], newResult Getter[Res], methodHandler any, oneway bool) *methodInfo[Arg, Res, H] {
	info := &methodInfo[Arg, Res, H]{
		typ:           typ,
		newArgs:       newArgs,
		newResult:     newResult,
		methodHandler: methodHandler,
		oneway:        oneway,
	}
	return info
}

func (o *methodInfo[Arg, Res, H]) Invoke(handler any, ctx context.Context, args any, result any) error {
	return o.methodHandler.(InvokeFunc[Arg, Res, H])(handler.(H), ctx, args.(Arg), result.(Res))
}

func (o *methodInfo[Arg, Res, H]) NewArgs() any {
	return o.newArgs()
}

func (o *methodInfo[Arg, Res, H]) NewResult() any {
	return o.newResult()
}

func (o *methodInfo[Arg, Res, H]) OneWay() bool {
	return o.oneway
}

func (o *methodInfo[Arg, Res, H]) Type() MethodType {
	return o.typ
}

func (o *methodInfo[Arg, Res, H]) Check(handler any) error {
	_, ok := handler.(H)
	if !ok {
		return fmt.Errorf("handler type is not match,want %v,but got %v", reflect.TypeFor[H](), reflect.TypeOf(handler))
	}
	return nil
}
