package selector

import (
	"context"

	"rpc-oneway/pkg/resolver"
)

type Selector interface {
	Apply([]resolver.ServiceInstance)
	Select(ctx context.Context) (resolver.ServiceInstance, error)
}
