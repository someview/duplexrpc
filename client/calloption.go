package client

import (
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"sync"
)

var callOptionsPool = sync.Pool{
	New: func() any {
		return new(CallOption)
	},
}

type CallOption struct {
	// call data
	info rpcinfo.MethodInfo
	node discovery.Node

	// fail control
	failMode    FailMode
	failedCount int

	// extra options //
	oneway          bool
	disableSelector bool
}

type CallOptFn func(c *CallOption)

func ApplyCallOption(opts ...CallOptFn) *CallOption {
	c := callOptionsPool.Get().(*CallOption)
	for _, optFn := range opts {
		optFn(c)
	}
	return c
}

func (c *CallOption) zero() {
	c.info = nil
	c.node = nil
	c.failMode = 0
	c.failedCount = 0
	c.oneway = false
	c.disableSelector = false
}

func (c *CallOption) Recycle() {
	c.zero()
	callOptionsPool.Put(c)
}

// WithOneway callback will be called after send to server.
//
// it is not needed to wait for the response.
func WithOneway() CallOptFn {
	return func(c *CallOption) {
		c.oneway = true
	}
}

// WithFailMode call with fail mode.
func WithFailMode(mode FailMode) CallOptFn {
	return func(c *CallOption) {
		c.failMode = mode
	}
}

// WithNode call this node without using a Selector.
func WithNode(node discovery.Node) CallOptFn {
	return func(c *CallOption) {
		c.disableSelector = true
		c.node = node
	}
}
