package metadata

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
)

type endpointKey struct{}

func WithEndpoint(parent context.Context, endpoint remote.Endpoint) context.Context {
	return WithValue(parent, endpointKey{}, endpoint)
}

func ExtractEndpoint(ctx context.Context) (end remote.Endpoint, ok bool) {
	end, ok = ctx.Value(endpointKey{}).(remote.Endpoint)
	return
}
