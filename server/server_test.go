package server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/client"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var info = &mahProcessMethod{}

func TestServer_Pressure(t *testing.T) {
	SetAntsRunner(-1)

	start := true
	var sTime time.Time
	// 发送的消息条数
	num := 10000 * 1000
	clientNum := 1
	wg := sync.WaitGroup{}
	wg.Add(num * clientNum)
	svr1 := NewServer(WithReadBufferThreshold(1024*1024), WithDisableHandlerMode())
	err := svr1.RegisterMethod(info, func() {
		if start {
			start = false
			sTime = time.Now()
		}

		wg.Done()
	})
	assert.NoError(t, err)

	go func() {
		err = svr1.Run("tcp", ":8080")
	}()
	time.Sleep(50 * time.Millisecond)
	t.Run("XClient", func(t *testing.T) {
		RunXClient(t, num, clientNum)
	})

	wg.Wait()
	fmt.Println("全部接收完毕！")
	fmt.Println("接收用时：", time.Since(sTime).String())
	svr1.Stop(context.Background())
}

func RunXClient(t *testing.T, num int, clientNum int) {
	var sTime time.Time
	var isStart bool = true
	wgS := sync.WaitGroup{}
	//sendData := &testdata.ClientMessage{Header: &testdata.Header{TraceId: string(make([]byte, 100)), Timestamp: time.Now().Unix()}}
	sendData := &mah{data: make([]byte, 10)}
	for i := 0; i < clientNum; i++ {
		wgS.Add(1)
		go func() {
			defer wgS.Done()
			res := resolver.NewLocalResolver("127.0.0.1:8080")
			xCli := client.NewClient("", info, client.WithResolver(res), client.WithSelector(selector.NewRoundRobinSelectorBuilder().Build()),
				client.WithSendBeforeMiddleWare(func(next middleware.Next) middleware.Next {
					counter := &atomic.Int32{}
					return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
						if count := counter.Add(1); count%10000 == 0 {
							//fmt.Println("已发送：", count)
						}
						next(ctx, req, cb)
					}
				}),
			)
			time.Sleep(100 * time.Millisecond)

			wg := sync.WaitGroup{}
			wg.Add(num)
			if isStart {
				sTime = time.Now()
				isStart = false
			}
			testCount := int32(0)
			callFns := []client.CallOptFn{
				client.WithOneway(),
			}
			for j := 0; j < num; j++ {
				xCli.Call(context.TODO(), sendData, client.RespCallBack(func(a any, err error) {
					assert.NoError(t, err)
					atomic.AddInt32(&testCount, 1)
					wg.Done()
				}), callFns...)

			}
			wg.Wait()
		}()
	}

	wgS.Wait()
	fmt.Println("全部发送完毕！")
	fmt.Println("发送用时：", time.Since(sTime).String())

}

func BenchmarkMapAndSlice(b *testing.B) {
	m := make(map[byte]struct{}, 255)
	s := make([]struct{}, 1, 255)
	for i := 0; i < 255; i++ {
		m[byte(i)] = struct{}{}
		s = append(s, struct{}{})
	}

	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = m[byte(i%255)]
		}
	})
	b.Run("Slice", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = s[i%255]
		}
	})

}

type mah struct {
	data []byte
}

func (m *mah) Size() int {
	return len(m.data)
}

func (m *mah) MarshalToSizedBuffer(bytes []byte) (int, error) {
	if m.data == nil {
		return 0, nil
	}
	copy(bytes, m.data)
	return m.Size(), nil
}

func (m *mah) Unmarshal(p []byte) error {
	return nil
}

type mahProcessMethod struct {
	oneWay atomic.Bool
}

func (m *mahProcessMethod) Invoke(fn any, ctx context.Context, args any, result any) error {

	fn.(func())()
	result.(*mah).data = args.(*mah).data
	return nil

}

func (m *mahProcessMethod) NewArgs() any {
	return new(mah)
}

func (m *mahProcessMethod) NewResult() any {
	return new(mah)
}

func (m *mahProcessMethod) OneWay() bool {
	return true
}

func (m *mahProcessMethod) Type() rpcinfo.MethodType {
	return 0
}

func (m *mahProcessMethod) Check(handler any) error {
	return nil
}
