package registrar

import (
	"github.com/stretchr/testify/assert"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"testing"
)

func TestLocalRegistrar_Register(t *testing.T) {
	r := NewLocalRegistrar("http://127.0.0.1:5004")
	node := discovery.ServiceInstance{Endpoint: "127.0.0.1:8080", Metadata: make(map[string]string)}
	assert.NoError(t, r.Register(nil, node))
}

func TestLocalRegistrar_Deregister(t *testing.T) {
	r := NewLocalRegistrar("http://127.0.0.1:5004")
	node := discovery.ServiceInstance{Endpoint: "127.0.0.1:8080", Metadata: make(map[string]string)}
	assert.NoError(t, r.Deregister(nil, node))
}
