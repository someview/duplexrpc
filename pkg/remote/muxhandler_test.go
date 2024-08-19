package remote

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/stretchr/testify/suite"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
	"go.uber.org/mock/gomock"
)

type muxTransHandlerSuite struct {
	suite.Suite

	ctrl            *gomock.Controller
	muxTransHandler *muxTransHandler
	opt             *RemoteOptions
	invokeFunc      middleware.GenericCall
	serviceManager  service.Manager
}

func (t *muxTransHandlerSuite) SetupSuite() {
	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		return nil
	}
	opt := DefaultRemoteOption
	t.opt = &opt
	t.serviceManager = service.Manager{}
}

func (t *muxTransHandlerSuite) TearDownSuite() {

}

func (t *muxTransHandlerSuite) SetupTest() {
	m, err := NewDefaultMuxTransHandler(t.invokeFunc, t.serviceManager, t.opt)
	t.NoError(err)
	t.muxTransHandler = m.(*muxTransHandler)

	t.ctrl = gomock.NewController(t.T())

}

func (t *muxTransHandlerSuite) TearDownTest() {
	t.ctrl.Finish()
}

func (t *muxTransHandlerSuite) reset() {
	t.muxTransHandler.opt = t.opt
	t.muxTransHandler.invokeFunc = t.invokeFunc
}

func (t *muxTransHandlerSuite) TestEvent() {
	var ons []OnActive
	count := 10
	for i := 0; i < count; i++ {
		idx := i
		ons = append(ons, func(ctx context.Context, end Endpoint) (context.Context, error) {
			ctx = context.WithValue(ctx, idx, idx)
			return ctx, nil
		})
	}
	opt := DefaultRemoteOption

	WithOnActive(func(ctx context.Context, end Endpoint) (context.Context, error) {
		var err error
		for _, on := range ons {
			ctx, err = on(ctx, end)
			if err != nil {
				break
			}
		}
		return ctx, err
	})(&opt)

	h, err := NewDefaultMuxTransHandler(nil, nil, &opt)
	m := h.(*muxTransHandler)
	t.NoError(err)
	ep := NewMockEndpoint(t.ctrl)
	ep.EXPECT().SetIdleTimeout(gomock.Any()).AnyTimes()
	ep.EXPECT().SetReadTimeout(gomock.Any()).AnyTimes()
	ep.EXPECT().SetWriteTimeout(gomock.Any()).AnyTimes()

	ctx, err := m.OnActive(context.TODO(), ep)
	t.NoError(err)
	for i := 0; i < count; i++ {
		t.Equal(i, ctx.Value(i))
	}

}

func (e *muxTransHandlerSuite) TestWrite_Oneway_编码成功() {
	mockMsg := NewMockMessage(e.ctrl)
	mockMsg.EXPECT().MsgType().Return(OnewayType).AnyTimes()
	mockCodec := NewMockCodec(e.ctrl)
	ctx := context.TODO()
	b := make([]byte, rand.Intn(30*util.KiB))
	mockCodec.EXPECT().Encode(ctx, gomock.Any(), mockMsg).Return(nil).Times(1).Do(
		func(_ context.Context, wr netpoll.Writer, _ Message) {
			wr.WriteBinary(b)
			wr.Flush()
		})
	mockMsg.EXPECT().Codec().Return(mockCodec).AnyTimes()

	mockEP := NewMockEndpoint(e.ctrl)
	mockEP.EXPECT().Add(ctx, gomock.Any()).Do(func(_ context.Context, lb *netpoll.LinkBuffer) {
		p, err := lb.ReadBinary(lb.Len())
		e.NoError(err)
		e.Equal(b, p)
	})
	err := e.muxTransHandler.Write(context.TODO(), mockEP, mockMsg)
	e.NoError(err)

}

func (e *muxTransHandlerSuite) TestWrite_Req_编码成功() {
	seqID := uint32(rand.Intn(10000))
	cbManager := newCallBackManager()
	mockEP := NewMockEndpoint(e.ctrl)
	mockEP.EXPECT().CallBackManager().Return(cbManager)

	mockInvocation := rpcinfo.NewMockInvocation(e.ctrl)
	mockRPCInfo := rpcinfo.NewMockRPCInfo(e.ctrl)
	mockRPCInfo.EXPECT().Invocation().Return(mockInvocation).AnyTimes()
	mockInvocation.EXPECT().SeqID().Return(seqID).AnyTimes()

	mockMsg := NewMockMessage(e.ctrl)
	mockMsg.EXPECT().MsgType().Return(ReqType).AnyTimes()
	mockMsg.EXPECT().RPCInfo().Return(mockRPCInfo).AnyTimes()
	mockCodec := NewMockCodec(e.ctrl)
	ctx := context.TODO()
	b := make([]byte, rand.Intn(30*util.KiB))
	mockCodec.EXPECT().Encode(ctx, gomock.Any(), mockMsg).Return(nil).Times(1).Do(
		func(_ context.Context, wr netpoll.Writer, _ Message) {
			wr.WriteBinary(b)
			wr.Flush()
		})
	mockMsg.EXPECT().Codec().Return(mockCodec).AnyTimes()
	mockEP.EXPECT().Add(ctx, gomock.Any()).Do(func(_ context.Context, lb *netpoll.LinkBuffer) {
		p, err := lb.ReadBinary(lb.Len())
		e.NoError(err)
		e.Equal(b, p)
	})
	err := e.muxTransHandler.Write(context.TODO(), mockEP, mockMsg)
	e.NoError(err)

	_, has := cbManager.Load(seqID)
	e.True(has)
}

func (e *muxTransHandlerSuite) TestWrite_Oneway_编码失败() {
	mockEP := NewMockEndpoint(e.ctrl)
	mockMsg := NewMockMessage(e.ctrl)
	mockMsg.EXPECT().MsgType().Return(OnewayType).AnyTimes()
	mockCodec := NewMockCodec(e.ctrl)
	ctx := context.TODO()
	mockCodec.EXPECT().Encode(ctx, gomock.Any(), mockMsg).Return(errors.New("encode error")).Times(1)
	mockMsg.EXPECT().Codec().Return(mockCodec).AnyTimes()

	err := e.muxTransHandler.Write(context.TODO(), mockEP, mockMsg)
	e.Error(err)
}

func (e *muxTransHandlerSuite) TestWrite_Req_编码失败() {
	seqID := uint32(rand.Intn(10000))
	cbManager := newCallBackManager()
	mockEP := NewMockEndpoint(e.ctrl)
	mockEP.EXPECT().CallBackManager().Return(cbManager)

	mockInvocation := rpcinfo.NewMockInvocation(e.ctrl)
	mockRPCInfo := rpcinfo.NewMockRPCInfo(e.ctrl)
	mockRPCInfo.EXPECT().Invocation().Return(mockInvocation).AnyTimes()
	mockInvocation.EXPECT().SeqID().Return(seqID).AnyTimes()

	mockMsg := NewMockMessage(e.ctrl)
	mockMsg.EXPECT().MsgType().Return(ReqType).AnyTimes()
	mockMsg.EXPECT().RPCInfo().Return(mockRPCInfo).AnyTimes()
	mockCodec := NewMockCodec(e.ctrl)
	ctx := context.TODO()

	mockCodec.EXPECT().Encode(ctx, gomock.Any(), mockMsg).Return(errors.New("encode error")).Times(1)
	mockMsg.EXPECT().Codec().Return(mockCodec).AnyTimes()

	err := e.muxTransHandler.Write(context.TODO(), mockEP, mockMsg)
	e.Error(err)

	_, has := cbManager.Load(seqID)
	e.False(has)
}

func (t *muxTransHandlerSuite) TestRead() {
	seqID := uint32(rand.Intn(10000))
	ctx := context.TODO()
	mockData := NewMockSizeableMarshaller(t.ctrl)
	mockData.EXPECT().Unmarshal(gomock.Any()).AnyTimes()

	iv := rpcinfo.NewInvocation("test", "test")
	iv.SetSeqID(seqID)
	info := rpcinfo.NewRPCInfo(rpcinfo.Unary, iv)

	mockCodec := NewMockCodec(t.ctrl)
	mockCodec.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any())

	mockMsg := NewMockMessage(t.ctrl)
	mockMsg.EXPECT().Payload().Return(netpoll.NewLinkBuffer()).Times(1)
	mockMsg.EXPECT().Data().Return(mockData)
	mockMsg.EXPECT().MsgType().Return(ResponseType)
	mockMsg.EXPECT().RPCInfo().Return(info)
	mockMsg.EXPECT().Codec().Return(mockCodec)

	ch := make(chan netpoll.Reader, 1)
	ch <- netpoll.NewLinkBuffer()
	cbManager := newCallBackManager()
	cbManager.Set(seqID, ch)

	mockEndpoint := NewMockEndpoint(t.ctrl)
	mockEndpoint.EXPECT().CallBackManager().Return(cbManager)
	mockEndpoint.EXPECT().Context().Return(ctx)

	err := t.muxTransHandler.Read(ctx, mockEndpoint, mockMsg)
	t.NoError(err)
}

func (t *muxTransHandlerSuite) TestRead_ErrorResp() {
	ctx := context.TODO()

	errCode := int32(100)
	errMsg := "my error"
	errData := uerror.NewBusinessError(errCode, errMsg)
	lb := netpoll.NewLinkBuffer()
	b, err := lb.Malloc(errData.Size())
	t.NoError(err)
	_, err = errData.MarshalToSizedBuffer(b)
	t.NoError(err)
	t.NoError(lb.Flush())

	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, rpcinfo.NewInvocation("", ""))

	mockCodec := NewMockCodec(t.ctrl)
	mockCodec.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any())

	mockMsg := NewMockMessage(t.ctrl)
	mockMsg.EXPECT().Codec().Return(mockCodec)
	mockMsg.EXPECT().Payload().Return(lb).Times(1)
	mockMsg.EXPECT().MsgType().Return(ErrorResponseType)
	mockMsg.EXPECT().RPCInfo().Return(ri)
	ch := make(chan netpoll.Reader, 1)
	ch <- lb

	mockCbManager := NewMockCallBackManager(t.ctrl)
	mockCbManager.EXPECT().Load(gomock.Any()).Return(ch, true)
	mockCbManager.EXPECT().Delete(gomock.Any())
	mockEndpoint := NewMockEndpoint(t.ctrl)
	mockEndpoint.EXPECT().CallBackManager().Return(mockCbManager)
	mockEndpoint.EXPECT().Context().Return(ctx)

	err = t.muxTransHandler.Read(ctx, mockEndpoint, mockMsg)
	t.Error(err)
	ok := uerror.IsBizError(err)
	t.True(ok)
	t.Equal(errData.BasicError.Error(), err.Error())
	t.T().Log(err)
}

func (t *muxTransHandlerSuite) TestOnMessage_Oneway() {
	serviceName := "sdfadf"
	methodName := "asdfasdfash"
	t.initServiceManager(serviceName, methodName)
	svc, ok := t.serviceManager.GetService(serviceName)
	t.True(ok)
	info, impl := svc.GetMethodInfoAndSvcImpl(methodName)

	iv := rpcinfo.NewMockInvocation(t.ctrl)
	iv.EXPECT().ServiceName().Return(serviceName).AnyTimes()
	iv.EXPECT().MethodName().Return(methodName).AnyTimes()
	iv.EXPECT().MethodInfo().Return(info).AnyTimes()
	iv.EXPECT().ServiceImpl().Return(impl).AnyTimes()

	ri := rpcinfo.NewMockRPCInfo(t.ctrl)
	ri.EXPECT().Invocation().Return(iv).AnyTimes()

	mockRecvMsg := NewMockMessage(t.ctrl)
	mockRecvMsg.EXPECT().MsgType().Return(OnewayType).AnyTimes()
	mockRecvMsg.EXPECT().RPCInfo().Return(ri).AnyTimes()
	mockRecvMsg.EXPECT().SetData(gomock.Any()).Times(1)
	mockRecvMsg.EXPECT().Payload().Return(netpoll.NewLinkBuffer()).Times(1)

	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		t.Equal(ri, info)
		return nil
	}
	t.reset()
	_, err := t.muxTransHandler.OnMessage(context.TODO(), mockRecvMsg, nil)
	t.NoError(err)
}

func (t *muxTransHandlerSuite) TestOnMessage_Request() {

	serviceName := "sdfadf"
	methodName := "asdfasdfash"
	t.initServiceManager(serviceName, methodName)
	svc, ok := t.serviceManager.GetService(serviceName)
	t.True(ok)
	info, impl := svc.GetMethodInfoAndSvcImpl(methodName)

	iv := rpcinfo.NewMockInvocation(t.ctrl)
	iv.EXPECT().ServiceName().Return(serviceName).AnyTimes()
	iv.EXPECT().MethodName().Return(methodName).AnyTimes()
	iv.EXPECT().MethodInfo().Return(info).AnyTimes()
	iv.EXPECT().ServiceImpl().Return(impl).AnyTimes()

	ri := rpcinfo.NewMockRPCInfo(t.ctrl)
	ri.EXPECT().Invocation().Return(iv).AnyTimes()

	mockRecvMsg := NewMockMessage(t.ctrl)
	mockRecvMsg.EXPECT().MsgType().Return(ReqType).AnyTimes()
	mockRecvMsg.EXPECT().RPCInfo().Return(ri).AnyTimes()
	mockRecvMsg.EXPECT().SetData(gomock.Any()).Times(1)
	mockRecvMsg.EXPECT().Payload().Return(netpoll.NewLinkBuffer()).Times(1)

	mockSendMsg := NewMockMessage(t.ctrl)
	mockSendMsg.EXPECT().SetData(gomock.Any()).Times(1)

	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		t.Equal(ri, info)
		return nil
	}
	t.reset()
	_, err := t.muxTransHandler.OnMessage(context.TODO(), mockRecvMsg, mockSendMsg)
	t.NoError(err)
}

type newArgsFactory struct {
	ctrl *gomock.Controller
}

func (n newArgsFactory) New() any {
	res := NewMockSizeableMarshaller(n.ctrl)
	res.EXPECT().Size().AnyTimes()
	res.EXPECT().MarshalToSizedBuffer(gomock.Any()).AnyTimes()
	res.EXPECT().Unmarshal(gomock.Any()).AnyTimes()
	return res
}

func (n newArgsFactory) Recycle(a any) {

}

func (t *muxTransHandlerSuite) initServiceManager(serviceName, methodName string) {
	newArgs := newArgsFactory{ctrl: t.ctrl}

	info := service.NewMethodInfo(nil, newArgs, newArgs, false, false)
	i := service.NewServiceInfo(serviceName, nil, map[string]service.MethodInfo{
		methodName: info,
	})
	t.serviceManager[serviceName] = service.NewService(i, nil)
}

func (t *muxTransHandlerSuite) TestOnRead_Oneway() {
	t.opt.ParallelDecider = func(svcName, methodName string) bool {
		return false
	}

	serviceName := "sdfasdf"
	methodName := "asdfas"
	t.initServiceManager(serviceName, methodName)
	py := make([]byte, 1024)

	mockData := NewMockSizeableMarshaller(t.ctrl)
	mockData.EXPECT().Size().Return(len(py)).AnyTimes()
	mockData.EXPECT().MarshalToSizedBuffer(gomock.Any()).Do(func(p []byte) {
		copy(p, py)
	})

	iv := rpcinfo.NewInvocation(serviceName, methodName)

	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, iv)
	msg := NewMessage(ri, t.opt.Codec)
	msg.SetData(mockData)
	lb := netpoll.NewLinkBuffer(1024)

	err := msg.Codec().Encode(context.TODO(), lb, msg)
	t.NoError(err)
	lb.Flush()
	mockEndpoint := NewMockEndpoint(t.ctrl)
	mockEndpoint.EXPECT().Reader().Return(lb)
	mockEndpoint.EXPECT().SliceIntoReader(gomock.Any(), gomock.Any()).DoAndReturn(func(n int, buf *netpoll.LinkBuffer) error {
		t.Equal(lb.Len(), n)
		return lb.SliceInto(n, buf)
	})

	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		inv := info.Invocation()
		t.Equal(serviceName, inv.ServiceName())
		t.Equal(methodName, inv.MethodName())
		t.Equal(iv.SeqID(), inv.SeqID())

		return nil
	}
	t.reset()
	err = t.muxTransHandler.OnRead(context.TODO(), mockEndpoint)
	t.NoError(err)
}

func (t *muxTransHandlerSuite) TestOnRead_Request() {
	t.opt.ParallelDecider = func(svcName, methodName string) bool {
		return false
	}

	serviceName := "sdfasdf"
	methodName := "asdfas"
	t.initServiceManager(serviceName, methodName)
	py := make([]byte, 1024)

	mockData := NewMockSizeableMarshaller(t.ctrl)
	mockData.EXPECT().Size().Return(len(py)).AnyTimes()
	mockData.EXPECT().MarshalToSizedBuffer(gomock.Any()).Do(func(p []byte) {
		copy(p, py)
	})

	iv := rpcinfo.NewInvocation(serviceName, methodName)

	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, iv)
	msg := NewMessage(ri, t.opt.Codec)
	msg.SetData(mockData)

	lb := netpoll.NewLinkBuffer(1024)

	err := msg.Codec().Encode(context.TODO(), lb, msg)
	t.NoError(err)
	lb.Flush()
	mockEndpoint := NewMockEndpoint(t.ctrl)
	mockEndpoint.EXPECT().Reader().Return(lb)
	mockEndpoint.EXPECT().SliceIntoReader(gomock.Any(), gomock.Any()).DoAndReturn(func(n int, buf *netpoll.LinkBuffer) error {
		t.Equal(lb.Len(), n)
		return lb.SliceInto(n, buf)
	})
	mockEndpoint.EXPECT().Add(gomock.Any(), gomock.Any()).Times(1)

	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		inv := info.Invocation()
		t.Equal(serviceName, inv.ServiceName())
		t.Equal(methodName, inv.MethodName())
		t.Equal(iv.SeqID(), inv.SeqID())

		return nil
	}
	t.reset()
	err = t.muxTransHandler.OnRead(context.TODO(), mockEndpoint)
	t.NoError(err)
}

func (t *muxTransHandlerSuite) TestGracefulShutdown() {
	taskCount := atomic.Int32{}
	taskCount.Store(100)

	t.invokeFunc = func(ctx context.Context, req, resp interface{}, info rpcinfo.RPCInfo) (err error) {
		time.Sleep(500 * time.Millisecond)
		taskCount.Add(-1)
		return nil
	}
	t.reset()
	py := make([]byte, 1024)
	mockData := NewMockSizeableMarshaller(t.ctrl)
	mockData.EXPECT().Size().Return(len(py)).AnyTimes()
	mockData.EXPECT().MarshalToSizedBuffer(gomock.Any()).Do(func(p []byte) {
		copy(p, py)
	}).AnyTimes()

	iv := rpcinfo.NewInvocation("", "")

	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, iv)
	t.initServiceManager("", "")

	for i := int32(0); i < taskCount.Load(); i++ {
		msg := NewMessage(ri, t.opt.Codec)
		msg.SetData(mockData)
		lb := netpoll.NewLinkBuffer(1024)
		err := msg.Codec().Encode(context.TODO(), lb, msg)
		t.NoError(err)
		lb.Flush()
		mockEndpoint := NewMockEndpoint(t.ctrl)
		mockEndpoint.EXPECT().Reader().Return(lb).AnyTimes()
		mockEndpoint.EXPECT().SliceIntoReader(gomock.Any(), gomock.Any()).DoAndReturn(func(n int, buf *netpoll.LinkBuffer) error {
			t.Equal(lb.Len(), n)
			return lb.SliceInto(n, buf)
		}).AnyTimes()
		err = t.muxTransHandler.OnRead(context.TODO(), mockEndpoint)
		t.NoError(err)
	}

	err := t.muxTransHandler.GracefulShutdown(context.TODO())
	t.NoError(err)

	t.Equal(int32(0), taskCount.Load())

}

func TestMuxTransHandlerSuite(t *testing.T) {
	suite.Run(t, new(muxTransHandlerSuite))
}

func BenchmarkAssert(b *testing.B) {
	end := &endpoint{}

	var conn netpoll.Connection
	conn = end

	b.Run("断言接口", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = conn.(Endpoint)
			}
		})

	})
	b.Run("断言结构体", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = conn.(*endpoint)
			}
		})
	})
}
