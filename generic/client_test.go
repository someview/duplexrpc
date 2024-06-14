package generic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testmethod"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/server"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
)

type msgMock struct{}

// MarshalToSizedBuffer implements protocol.SizeableMarshaller.
func (m *msgMock) MarshalToSizedBuffer([]byte) (int, error) {
	return 0, nil
}

// Size implements protocol.SizeableMarshaller.
func (m *msgMock) Size() int {
	return 10
}

// Unmarshal implements protocol.SizeableMarshaller.
func (*msgMock) Unmarshal(data []byte) error {
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}
	return nil
}

var _ protocol.SizeableMarshaller = (*msgMock)(nil)

func TestClient(t *testing.T) {
	srv := server.NewServer()
	err := server.RegisterMethodToServer[testmethod.ProcessClientMessage](srv, testmethod.ProcessClientHandler(func(msg *testdata.ClientMessage) error {
		return nil
	}))
	assert.NoError(t, err)
	go func() {
		if err := srv.Run("tcp", ":8080"); err != nil {
			log.Fatalln("err:", err)
		}
	}()

	time.Sleep(time.Second * 1)
	cli := NewMuxClient(defaultOpt)

	if err := cli.AsyncConnect("tcp", "127.0.0.1:8080"); err != nil {
		log.Fatalln("err:", err)
	}
	time.Sleep(time.Second * 1)
	wg := sync.WaitGroup{}
	req := protocol.NewMessage()
	req.SetReq(&testdata.ClientMessage{
		Header: &testdata.Header{TraceId: "123456789"},
	})
	wg.Add(1)
	cli.AsyncSend(req, func(err error) {
		assert.NoError(t, err)
		wg.Done()
	})
	wg.Wait()

}

func TestMuxClient_Close(t *testing.T) {
	srv := server.NewServer()
	err := server.RegisterMethodToServer[testmethod.ProcessClientMessage](srv, testmethod.ProcessClientHandler(func(msg *testdata.ClientMessage) error {

		return nil
	}))
	assert.NoError(t, err)
	go func() {
		if err := srv.Run("tcp", ":8080"); err != nil {
			log.Fatalln("err:", err)
		}
	}()

	time.Sleep(time.Second * 1)

	cli := NewMuxClient(defaultOpt)

	if err := cli.AsyncConnect("tcp", "127.0.0.1:8080"); err != nil {
		log.Fatalln("err:", err)
	}
	time.Sleep(time.Second * 1)

	req := protocol.NewMessage()
	req.SetReq(&testdata.ClientMessage{
		Header: &testdata.Header{TraceId: "123456789"},
	})
	wg := sync.WaitGroup{}
	count := 1000
	wg.Add(count)
	var handleCount int32
	for i := 0; i < count; i++ {
		cli.AsyncSend(req, func(err error) {
			assert.NoError(t, err)
			atomic.AddInt32(&handleCount, 1)
			wg.Done()
		})
	}
	wg.Wait()
	assert.NoError(t, cli.Close())
	assert.Equal(t, int32(count), atomic.LoadInt32(&handleCount))
	cli.AsyncSend(req, func(err error) {
		assert.Error(t, err)
	})

}

func TestMuxClient_Middleware(t *testing.T) {
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
	cli := NewMuxClient(opt)
	cli.AsyncSend(nil, func(err error) {
		assert.Equal(t, err, myErr)
	})

}

func BenchmarkMuxClient_AsyncSend(b *testing.B) {
	srv := server.NewServer()

	err := server.RegisterMethodToServer[testmethod.ProcessClientMessage](srv, testmethod.ProcessClientHandler(func(msg *testdata.ClientMessage) error {
		return nil
	}))
	assert.NoError(b, err)
	go func() {
		if err := srv.Run("tcp", ":8080"); err != nil {
			//log.Fatalln("err:", err)
		}
	}()
	time.Sleep(50 * time.Millisecond)
	// 哎呦,就这么处理吧
	cli := NewMuxClient(Option{writeBufferThreshold: 40960})
	assert.NoError(b, cli.AsyncConnect("tcp", "127.0.0.1:8080"))
	time.Sleep(50 * time.Millisecond)
	b.ResetTimer()
	cb := func(err error) {
		if err != nil {
			assert.NoError(b, err)
		}
	}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := protocol.NewMessage()
			req.SetReq(&testdata.ClientMessage{
				Header: &testdata.Header{TraceId: "123456789"},
			})
			cli.AsyncSend(req, cb)
			req.Recycle()
		}
	})
}

func BenchmarkTimeSysCall(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = time.Now()
		}
	})
}
