package registrar

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
)

// Registrar is service registrar.
//
//go:generate mockgen -source=registry.go -destination=../mocks/registry.go -package=mocks  Registrar
type Registrar interface {
	// Register the registration.
	Register(ctx context.Context, service discovery.ServiceInstance) error
	// Deregister the registration.
	Deregister(ctx context.Context, service discovery.ServiceInstance) error
}
