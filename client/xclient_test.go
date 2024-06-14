package client

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/server"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testmethod"
	"sync"
	"testing"
	"time"
)

type msgEcho struct {
	testdata.ClientMessage
	real *testdata.ClientMessage
}

func (m msgEcho) Invoke(handler any, ctx context.Context, args, result any) error {

	res := handler.(func(msg *testdata.ClientMessage) *testdata.ClientMessage)(args.(*testdata.ClientMessage))
	result.(*msgEcho).real = res
	return nil

}

func (m msgEcho) NewArgs() any {
	return new(testdata.ClientMessage)
}

func (m msgEcho) NewResult() any {
	return new(msgEcho)
}

func (m msgEcho) OneWay() bool {
	return false
}

func (m msgEcho) Type() rpcinfo.MethodType {
	return 0
}

func (m msgEcho) Check(handler any) error {
	return nil
}

func TestXClient_Send(t *testing.T) {
	wg := sync.WaitGroup{}
	svr1 := server.NewServer()
	svr2 := server.NewServer()
	err := server.RegisterMethodToServer[msgEcho](svr1, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		t.Logf("server %s, recive msg: %s", "8080", msg.String())
		return msg
	})
	assert.NoError(t, err)

	err = server.RegisterMethodToServer[msgEcho](svr2, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		t.Logf("server %s, recive msg: %s", "8081", msg.String())
		return msg
	})
	assert.NoError(t, err)

	go func() {
		_ = svr1.Run("tcp", ":8080")
	}()
	go func() {
		_ = svr2.Run("tcp", ":8081")
	}()

	time.Sleep(50 * time.Millisecond)

	res := resolver.NewLocalResolver("127.0.0.1:8080", "127.0.0.1:8080", "127.0.0.1:8081")
	ops := []OptionFn{WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
		WithMsgHandler(func(resp any, err error) {
			assert.NoError(t, err)
			defer wg.Done()
			t.Logf("echo form server")
		}),
	}
	mInfo := new(msgEcho)
	xCli := NewClient("", mInfo, ops...)

	time.Sleep(500 * time.Millisecond)

	sendData := &testdata.ClientMessage{Header: &testdata.Header{TraceId: "123456789", Timestamp: time.Now().Unix()}}
	ctx := context.Background()

	wg.Add(6)

	for i := 0; i < 6; i++ {
		xCli.Call(ctx, sendData, func(a any, err error) {
			assert.NoError(t, err)
		})
	}

	wg.Wait()
}

func TestXClient_SendFailSelect(t *testing.T) {
	wg := sync.WaitGroup{}
	svr1 := server.NewServer()
	svr2 := server.NewServer()
	err := server.RegisterMethodToServer[msgEcho](svr1, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		panic("server 8080 should not receive msg")
	})
	assert.NoError(t, err)
	err = server.RegisterMethodToServer[msgEcho](svr2, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		defer wg.Done()
		t.Logf("server %s, recive msg: %s", "8081", msg.String())
		return msg
	})
	assert.NoError(t, err)

	go func() {
		_ = svr1.Run("tcp", ":8080")
	}()
	go func() {
		_ = svr2.Run("tcp", ":8081")
	}()

	time.Sleep(50 * time.Millisecond)
	res := resolver.NewLocalResolver("127.0.0.1:8080", "127.0.0.1:8081")
	ops := []OptionFn{
		WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
	}
	mInfo := new(msgEcho)
	xCli := NewClient("", mInfo, ops...)

	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, svr1.Stop(context.Background()))

	sendData := &testdata.ClientMessage{Header: &testdata.Header{TraceId: "123456789", Timestamp: time.Now().Unix()}}
	ctx := context.Background()

	wg.Add(6)

	for i := 0; i < 6; i++ {
		xCli.Call(ctx, sendData, func(a any, err error) {
			assert.NoError(t, err)
		})
	}
	wg.Wait()
}

func TestXClient_SendFail(t *testing.T) {
	wg := sync.WaitGroup{}
	svr1 := server.NewServer()
	svr2 := server.NewServer()
	err := server.RegisterMethodToServer[msgEcho](svr1, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		panic("server 8080 should not receive msg")
	})
	assert.NoError(t, err)
	err = server.RegisterMethodToServer[msgEcho](svr2, func(msg *testdata.ClientMessage) *testdata.ClientMessage {
		defer wg.Done()
		t.Logf("server %s, recive msg: %s", "8081", msg.String())
		return msg
	})
	assert.NoError(t, err)

	go func() {
		_ = svr1.Run("tcp", ":8080")
	}()
	go func() {
		_ = svr2.Run("tcp", ":8081")
	}()

	time.Sleep(50 * time.Millisecond)
	res := resolver.NewLocalResolver("127.0.0.1:8080", "127.0.0.1:8081")
	ops := []OptionFn{
		WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
	}
	mInfo := new(msgEcho)
	xCli := NewClient("", mInfo, ops...)

	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, svr1.Stop(context.Background()))

	sendData := &testdata.ClientMessage{Header: &testdata.Header{TraceId: "123456789", Timestamp: time.Now().Unix()}}
	ctx := context.TODO()

	wg.Add(6)

	for i := 0; i < 6; i++ {
		xCli.Call(ctx, sendData, func(a any, err error) {
			assert.NoError(t, err)
		})
	}
	wg.Wait()
	assert.NoError(t, svr2.Stop(context.Background()))
	wg.Add(3)
	for i := 0; i < 3; i++ {
		xCli.Call(ctx, sendData, func(a any, err error) {
			assert.Error(t, err)
			wg.Done()
		})
	}
	wg.Wait()
}

func TestXClient_CallOptNode(t *testing.T) {
	svr1 := server.NewServer()
	go func() {
		_ = svr1.Run("tcp", ":8989")
	}()
	defer svr1.Stop(context.TODO())
	time.Sleep(100 * time.Millisecond)
	res := resolver.NewLocalResolver("127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8989")
	ops := []OptionFn{
		WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
		WithSendBeforeMiddleWare(func(next middleware.Next) middleware.Next {
			return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
				cb(nil)
			}
		}),
	}
	mInfo := new(msgEcho)
	xCli := NewClient("", mInfo, ops...)
	time.Sleep(100 * time.Millisecond)
	xCli.Call(context.TODO(), mInfo.NewArgs(), func(a any, err error) {
		assert.NoError(t, err)
	}, WithNode(discovery.NewNode("127.0.0.1:8989", nil)))

}

func TestXClient_FailMode(t *testing.T) {
	svr1 := server.NewServer()
	go func() {
		_ = svr1.Run("tcp", ":8989")
	}()
	defer svr1.Stop(context.TODO())
	time.Sleep(100 * time.Millisecond)
	res := resolver.NewLocalResolver("127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8989")
	failCount := 0
	ops := []OptionFn{
		WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
		WithFailureLimit(2),
		WithSendBeforeMiddleWare(func(next middleware.Next) middleware.Next {
			return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
				failCount++
				if failCount > 4 {
					assert.Fail(t, "fail")
				}
				next(ctx, req, cb)
			}
		}),
	}
	mInfo := new(msgEcho)
	xCli := NewClient("", mInfo, ops...)
	time.Sleep(100 * time.Millisecond)
	wg := make(chan struct{}, 1)
	xCli.Call(context.TODO(), mInfo.NewArgs(), func(a any, err error) {
		assert.NoError(t, err)
		wg <- struct{}{}
	}, WithOneway(), WithFailMode(Failover))
	<-wg
	failCount = 0
	xCli.Call(context.TODO(), mInfo.NewArgs(), func(a any, err error) {
		assert.Error(t, err)
		wg <- struct{}{}
	}, WithOneway(), WithFailMode(Failfast))
	<-wg
	for i := 0; i < 2; i++ {
		failCount = 0
		xCli.Call(context.TODO(), mInfo.NewArgs(), func(a any, err error) {
			if i == 0 {
				// here, the client selector was selected the true node
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			wg <- struct{}{}
		}, WithOneway(), WithFailMode(Failtry))
		<-wg
	}

	failCount = 0
	xCli.Call(context.TODO(), mInfo.NewArgs(), func(a any, err error) {
		assert.Error(t, err)
		fmt.Println(err)
		wg <- struct{}{}
	}, WithOneway(), WithFailMode(Failbackup))
	<-wg

}

func BenchmarkXClient_AsyncSend(b *testing.B) {

	svr1 := server.NewServer()
	err := server.RegisterMethodToServer[testmethod.ProcessClientMessage](svr1, testmethod.ProcessClientHandler(func(msg *testdata.ClientMessage) error {
		return nil
	}))
	assert.NoError(b, err)
	go func() {
		_ = svr1.Run("tcp", ":8080")
	}()

	time.Sleep(50 * time.Millisecond)
	rInfo := new(testmethod.ProcessClientMessage)
	res := resolver.NewLocalResolver("127.0.0.1:8080")
	ops := []OptionFn{
		WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
		WithResolver(res),
	}
	time.Sleep(50 * time.Millisecond)
	xCli := NewClient("", rInfo, ops...)

	time.Sleep(50 * time.Millisecond)

	sendData := &testdata.ClientMessage{Header: &testdata.Header{TraceId: "123456789", Timestamp: time.Now().Unix()}}

	b.ResetTimer()

	cb := RespCallBack(func(resp any, err error) {
		if err != nil {
			panic(err)
		}
	})
	for i := 0; i < b.N; i++ {
		xCli.Call(context.Background(), sendData, cb)
	}
}
