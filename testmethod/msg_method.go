package testmethod

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
)

type ProcessClientMessage struct {
}

type ProcessClientHandler func(msg *testdata.ClientMessage) error

func (p ProcessClientMessage) Invoke(fn any, ctx context.Context, args any, result any) error {

	// ctx可以用来Handle超时判断
	msg := args.(*testdata.ClientMessage)
	return fn.(ProcessClientHandler)(msg)

}

func (p ProcessClientMessage) NewArgs() any {
	return &testdata.ClientMessage{}
}

func (p ProcessClientMessage) NewResult() any {
	return nil
}

func (p ProcessClientMessage) OneWay() bool {
	return true
}

func (p ProcessClientMessage) Type() rpcinfo.MethodType {
	return 0
}

func (p ProcessClientMessage) Check(handler any) error {
	return nil
}
