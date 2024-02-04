package resolver

import "context"

// Registrar is service registrar.
//go:generate mockgen -source=registry.go -destination=../mocks/registry.go -package=mocks  Registrar
type Registrar interface {
	// Register the registration.
	Register(ctx context.Context, service ServiceInstance) error
	// Deregister the registration.
	Deregister(ctx context.Context, service ServiceInstance) error
}

// Discovery is service discovery.
type Discovery interface {
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context) ([]ServiceInstance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]ServiceInstance, error)
	// Stop close the watcher.
	Stop() error
}

// ServiceInstance is an instance of a service in a discovery system.
type ServiceInstance struct {
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string `json:"metadata"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoint string `json:"endpoints"`
}

type ServiceListener func([]ServiceInstance)

type Resolver interface {
	/// <summary>
	/// Starts listening to resolver for results with the specified callback. Can only be called once.
	/// <para>
	/// The <see cref="ResolverResult"/> passed to the callback has addresses when successful,
	/// otherwise a <see cref="Status"/> details the resolution error.
	/// </para>
	/// </summary>
	/// <param name="listener">The callback used to receive updates on the target.</param>
	Start(ServiceListener)
	/// <summary>
	/// Refresh resolution. Can only be called after <see cref="Start(Action{ResolverResult})"/>.
	/// The default implementation is no-op.
	/// <para>
	/// This is only a hint. Implementation takes it as a signal but may not start resolution.
	/// </para>
	/// </summary>
	Refresh()
	Stop()
}
