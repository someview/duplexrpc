package selector

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
)

type Selector interface {
	Apply([]discovery.Node)
	Select(ctx context.Context) (discovery.Node, error)
}
