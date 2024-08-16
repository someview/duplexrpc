package bufwriter

import (
	"context"
	"io"

	netpoll "github.com/cloudwego/netpoll"
)

const (
	active  = 0
	closing = 1
	closed  = 2
)

type BufWriterType int

const (
	// ShardQueueType 使用ShardQueue
	ShardQueueType BufWriterType = 1
	// DelayQueueType 使用DelayQueue
	DelayQueueType BufWriterType = 2
)

type BufWriter interface {
	// Add 将数据写入缓冲区
	Add(ctx context.Context, lb *netpoll.LinkBuffer) error
	io.Closer
}

type BufFlusher interface {
	FlushTo(netpoll.Connection) (n int, err error)
}

// FlowControl 用于流量控制
//
//go:generate mockgen -source=define.go -destination ./interface_mock.go -package bufwriter
type FlowControl interface {
	// GetWithCtx 阻塞获取，直到有可用配额和或ctx.Done()
	GetWithCtx(ctx context.Context, sz int) error
	// TryGet 尝试获取配额，非阻塞
	TryGet(sz int) bool
	// Release 归还配额
	Release(sz int)
	// Available 剩余配额
	Available() int
}

// BatchContainer 用于批量写入数据，除了FlushToWriter，其他都是线程安全的
type BatchContainer interface {
	IsFull() bool
	IsEmpty() bool
	Len() int
	TryAdd(lb *netpoll.LinkBuffer) (ok bool)

	BatchWriter
	BufFlusher
}

// 智能写接口,自动将数据flush到connection中
type BatchWriter interface {
	io.Closer
	BufWriter
}

type Connection interface {
	netpoll.Connection
}
