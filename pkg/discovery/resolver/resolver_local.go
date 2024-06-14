package resolver

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"time"
)

var LocalServices = []string{"127.0.0.1:8080"}

type localResolver struct {
	services []string
	ctx      context.Context
	can      context.CancelFunc
	re       chan struct{}
}

// NewLocalResolver 相当于直连
//
//	services为本地srv地址
func NewLocalResolver(services ...string) Resolver {
	ctx := context.TODO()
	ctx, can := context.WithCancel(ctx)
	r := &localResolver{
		ctx: ctx,
		can: can,
		re:  make(chan struct{}, 1),
	}
	if len(services) <= 0 {
		r.services = []string{"127.0.0.1:8080", "127.0.0.1:8081"}
	} else {
		r.services = services
	}
	return r
}

func (l *localResolver) Start(listener discovery.ServiceListener) {
	ticker := time.NewTicker(5 * time.Second)
redo:
	var nodes []discovery.Node
	for _, s := range l.services {
		node := discovery.NewNode(s, make(map[string]string))

		nodes = append(nodes, node)

	}
	listener(nodes)
	select {
	case <-ticker.C:
		goto redo
	case <-l.re:
		goto redo
	case <-l.ctx.Done():
		return
	}
}

func (l *localResolver) Refresh() {
	l.services = LocalServices
	l.re <- struct{}{}
}

func (l *localResolver) Stop() {
	l.can()
}
