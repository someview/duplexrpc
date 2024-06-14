package generic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"testing"
)

func TestNewMultiClient(t *testing.T) {
	_ = NewMultiClient(2, defaultOpt)
}

func TestMultiClient_IsShutdown(t *testing.T) {
	cli := NewMultiClient(2, defaultOpt)
	assert.NoError(t, cli.AsyncConnect("tcp", "127.0.0.1:8080"))
	assert.Equal(t, false, cli.IsShutdown())
	assert.NoError(t, cli.Close())
	assert.Equal(t, true, cli.IsShutdown())
	assert.Error(t, cli.AsyncConnect("tcp", ""))
	assert.Equal(t, "", cli.Address())

}

func TestMultiClient_Middleware(t *testing.T) {
	myErr := fmt.Errorf("error")
	opt := defaultOpt
	optFns := []OptionFn{
		WithSendBeforeMiddleWare(func(next middleware.Next) middleware.Next {
			return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
				cb(myErr)
			}
		}),
	}
	for _, fn := range optFns {
		fn(&opt)
	}
	cli := NewMultiClient(10, opt)
	cli.AsyncSend(nil, func(err error) {
		assert.Equal(t, err, myErr)
	})

}
