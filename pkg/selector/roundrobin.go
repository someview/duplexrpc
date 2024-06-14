package selector

import (
	"context"
	"errors"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"sync/atomic"
)

var ErrNoAvailable = errors.New("no_available_node")

// HashSelectorBuilder is composite selector.
type RoundRobinSelector struct {
	nodes atomic.Value
	index atomic.Int32 // last选择的下标
	count atomic.Int32
}

func (r *RoundRobinSelector) NodeCount() int {
	return int(r.count.Load())
}

// Select is select one node.
func (r *RoundRobinSelector) Select(ctx context.Context) (selected discovery.Node, err error) {
	nodes, ok := r.nodes.Load().([]discovery.Node)
	if !ok {
		return nil, ErrNoAvailable
	}
	if len(nodes) == 0 {
		return nil, ErrNoAvailable
	}
	return nodes[r.index.Add(1)%r.count.Load()], nil
}

// Apply update nodes info.
func (r *RoundRobinSelector) Apply(nodes []discovery.Node) {
	r.nodes.Store(nodes)
	r.index.Store(0) // 重置index
	r.count.Store(int32(len(nodes)))
}

// DefaultBuilder is de
type RoundRobinSelectorBuilder struct {
}

func NewRoundRobinSelectorBuilder() *RoundRobinSelectorBuilder {
	return &RoundRobinSelectorBuilder{}
}

// Build create builder
func (db *RoundRobinSelectorBuilder) Build() Selector {
	res := &RoundRobinSelector{}
	return res
}
