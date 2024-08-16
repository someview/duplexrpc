package remote

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
	"time"
)

// EventTypes
type (
	// OnActive endpoint建立时
	OnActive func(ctx context.Context, end Endpoint) (context.Context, error)
	// OnInactive endpoint断开时
	OnInactive func(ctx context.Context, end Endpoint)
	// OnError 发生错误时
	OnError func(ctx context.Context, end Endpoint, err error)
	// OnRPCDone 任务执行完成后
	OnRPCDone func(ctx context.Context, req, res any, ri rpcinfo.RPCInfo)
)

type eventHooks struct {
	onActive   OnActive
	onInactive OnInactive
	onError    OnError
	onRPCDone  OnRPCDone
}

// RemoteOptions contains option that is used to init the remote server.
type RemoteOptions struct {
	eventHooks

	WriterType           bufwriter.BufWriterType
	WriteDelayTime       time.Duration
	WriteMaxMsgNum       int
	WriteBufferThreshold int
	Codec                Codec
	ParallelDecider      ParallelDecider

	MaxConnectionIdleTime time.Duration
	ReadWriteTimeout      time.Duration
}

type RemoteOption func(opt *RemoteOptions)

// WithOnActive 在endpoint建立时执行的回调
func WithOnActive(fn OnActive) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.onActive = fn
	}
}

// WithOnInactive 在endpoint断开时执行的回调
func WithOnInactive(fn OnInactive) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.onInactive = fn
	}
}

// WithOnError 发生错误时执行的回调
func WithOnError(fn OnError) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.onError = fn
	}
}

// WithOnRPCDone 当一个RPC请求执行完毕后执行的回调
func WithOnRPCDone(fn OnRPCDone) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.onRPCDone = fn
	}
}

// WithWriterType 设置bufWriter的类型
func WithWriterType(typ bufwriter.BufWriterType) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.WriterType = typ
	}
}

// WithWriteDelayTime 设置bufWriter的写延时时间，仅当WriterType为DelayQueueType时有效
func WithWriteDelayTime(t time.Duration) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.WriteDelayTime = t
	}
}

// WithWriteMaxMsgNum 设置bufWriter一次最多写多少个消息，仅当WriterType为DelayQueueType时有效
func WithWriteMaxMsgNum(num int) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.WriteMaxMsgNum = num
	}
}

// WithWriteBufferThreshold 设置bufWriter写缓冲区的阈值
func WithWriteBufferThreshold(threshold int) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.WriteBufferThreshold = threshold
	}
}

// WithCodec 设置RPC协议的编解码器，
func WithCodec(c Codec) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.Codec = c
	}
}

func WithParallelDecider(decider ParallelDecider) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.ParallelDecider = decider
	}
}

func WithMaxConnectionIdleTime(t time.Duration) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.MaxConnectionIdleTime = t
	}
}

func WithReadWriteTimeout(t time.Duration) RemoteOption {
	return func(opt *RemoteOptions) {
		opt.ReadWriteTimeout = t
	}
}

var DefaultRemoteOption = RemoteOptions{
	eventHooks: eventHooks{
		onActive: func(ctx context.Context, end Endpoint) (context.Context, error) {
			return ctx, nil
		},
		onInactive: func(ctx context.Context, end Endpoint) {

		},
		onError: func(ctx context.Context, end Endpoint, err error) {

		},
		onRPCDone: func(ctx context.Context, req, res any, ri rpcinfo.RPCInfo) {

		},
	},

	WriterType:           bufwriter.DelayQueueType,
	WriteDelayTime:       1 * time.Millisecond,
	WriteMaxMsgNum:       100,
	WriteBufferThreshold: 64 * util.MiB,
	Codec:                NewDefaultCodeC(),
	ParallelDecider: func(svcName, methodName string) bool {
		return true
	},
	MaxConnectionIdleTime: 5 * time.Minute,
	ReadWriteTimeout:      5 * time.Minute,
}
