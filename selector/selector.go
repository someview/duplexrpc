package selector

import (
	"context"
)

// todo 这里的接口还需要进一步设计
type SelectFunc func(ctx context.Context, servicePath, serviceMethod string, args interface{}) string

// Selector defines selector that selects one service from candidates.
type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // SelectFunc
	UpdateServer(servers map[string]string)
}
