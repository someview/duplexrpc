package selector

import (
	"context"

	"rpc-oneway/pkg/resolver"
)

type Selector interface {
	Apply([]resolver.Node)
	Select(ctx context.Context) (resolver.Node, error)
}
