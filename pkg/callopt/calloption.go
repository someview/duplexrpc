package callopt

import (
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"sync"
)

var callOptionsPool = sync.Pool{
	New: func() any {
		return new(CallOption)
	},
}

type CallOption struct {
	node            discovery.Node
	failedCount     int
	disableSelector bool
}

type CallOpt func(c *CallOption)

func ApplyCallOption(opts ...CallOpt) *CallOption {
	c := callOptionsPool.Get().(*CallOption)
	for _, optFn := range opts {
		optFn(c)
	}
	return c
}

func (c *CallOption) zero() {
}

func (c *CallOption) Recycle() {
	c.zero()
	callOptionsPool.Put(c)
}
