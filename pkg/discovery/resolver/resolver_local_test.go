package resolver

import (
	"github.com/stretchr/testify/assert"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"sync"
	"testing"
)

func TestLocalResolver_Refresh(t *testing.T) {
	res := NewLocalResolver()
	newS := LocalServices
	wg := sync.WaitGroup{}
	go func() {
		res.Start(func(nodes []discovery.Node) {
			defer wg.Done()
			addr := make([]string, len(nodes))
			for i, node := range nodes {
				addr[i] = node.Address()
			}
			assert.Equal(t, newS, addr)
		})
	}()
	wg.Add(2)
	newS = []string{"127.0.0.1:8081", "127.0.0.1:8082"}
	LocalServices = newS
	res.Refresh()

}
