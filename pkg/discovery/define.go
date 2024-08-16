package discovery

import (
	"context"
)

// Discovery is service discovery.
type Discovery interface {
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context) ([]Node, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]Node, error)
	// Stop close the watcher.
	Stop() error
}

type Node interface {
	Address() string
	MD() map[string]string
	Weight() *int
}

// ServiceInstance is an instance of a service in a discovery system.
type ServiceInstance struct {
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string `json:"metadata"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoint string `json:"endpoint"`

	Wei *int
}

func NewNode(url string, readonlyMD map[string]string) Node {
	return &ServiceInstance{
		Endpoint: url,
		Metadata: readonlyMD,
	}
}

// Address implements Node.
func (s *ServiceInstance) Address() string {
	return s.Endpoint
}

// MD implements Node.
func (s *ServiceInstance) MD() map[string]string {
	return s.Metadata
}

// Weight implements Node.
func (s *ServiceInstance) Weight() *int {
	return s.Wei
}

var _ Node = (*ServiceInstance)(nil)

type ServiceListener func([]Node)
