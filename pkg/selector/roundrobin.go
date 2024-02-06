package selector

import (
	"context"
	"errors"
	"rpc-oneway/pkg/resolver"
	"sync/atomic"
)

var ErrNoAvailable = errors.New("no_available_node")

// HashSelectorBuilder is composite selector.
type RoundRobinSelector struct {
	nodes atomic.Value
	index int // last选择的下标
	count int
}

func (r *RoundRobinSelector) NodeCount() int {
	return r.count
}

// Select is select one node.
func (r *RoundRobinSelector) Select(ctx context.Context) (selected resolver.Node, err error) {
	nodes, ok := r.nodes.Load().([]resolver.Node)
	if !ok {
		return nil, ErrNoAvailable
	}
	if len(nodes) == 0 {
		return nil, ErrNoAvailable
	}
	r.index++
	if r.index >= len(nodes) {
		r.index = 0
		return nodes[0], nil
	}
	return nodes[r.index], nil
}

// Apply update nodes info.
func (r *RoundRobinSelector) Apply(nodes []resolver.Node) {
	r.nodes.Store(nodes)
	r.index = 0 // 重置index
	r.count = len(nodes)
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
