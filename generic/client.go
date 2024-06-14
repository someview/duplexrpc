package generic

import (
	"context"
	"errors"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/writer"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"sync"
	"sync/atomic"
)

var (
	ErrShutdown           = errors.New("client is shut down")
	ErrBreakerOpen        = errors.New("breaker is open")
	ErrClientNotConnected = errors.New("client is not connected")
)

type MuxClient struct {
	ctx context.Context

	network string
	address string

	conn   netpoll.Connection
	option Option

	wr muxextend.AsyncWriter

	mu         sync.RWMutex
	isShutDown atomic.Bool

	send middleware.Next
}

func NewMuxClient(option Option) *MuxClient {
	res := &MuxClient{
		ctx:    context.Background(),
		option: option,
	}
	res.initSend()

	return res
}

// ---implements RPCClient---

// Address implements RPCClient.
func (c *MuxClient) Address() string {
	return c.address
}

// IsShutdown implements RPCClient.
func (c *MuxClient) IsShutdown() bool {
	return c.isShutDown.Load()
}

func (c *MuxClient) initAsyncWriter() {
	if c.wr != nil {
		c.wr.Close()
	}
	switch c.option.writer {
	case writer.ShardQueue:
		c.wr = muxextend.NewShardQueue(c.option.writeBufferThreshold, c.conn)
	case writer.BatchWriter:
		c.wr = muxextend.NewBatchWriter(c.option.writeBufferThreshold, c.option.writeMaxMsgNum, c.option.writeDelayTime, c.conn)
	}

}

func (c *MuxClient) initSend() {
	mw := middleware.Chain(c.option.sendBeforeMWs...)
	w := mw(func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
		c.asyncSend(req, cb)
	})
	c.send = w
}

func (c *MuxClient) AsyncSend(req protocol.Message, cb netpoll.CallBack) {
	c.send(req, req, cb)
}

func (c *MuxClient) asyncSend(req protocol.Message, cb netpoll.CallBack) {
	if c.isShutDown.Load() {
		cb(ErrShutdown)
		return
	}

	c.mu.RLock()
	wr := c.wr
	c.mu.RUnlock()
	if wr == nil {
		cb(ErrClientNotConnected)
		return
	}
	wr.AsyncWrite(req, req, cb)
}

func (c *MuxClient) ShutDown() error {
	c.isShutDown.Store(true)
	if c.wr != nil {
		_ = c.wr.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.mu.Lock()
	c.conn = nil
	c.wr = nil
	c.mu.Unlock()
	return nil
}

func (c *MuxClient) Close() error {
	return c.ShutDown()
}
