package client

import (
	"context"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/client"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/metadata"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/server"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
	"go.uber.org/mock/gomock"
)

type mockSubService struct {
	existSubRes testdata.ExistSubRes
	subReq      testdata.SubReq
}

func (c *mockSubService) SetExistSubscriptionRes(res *testdata.ExistSubRes) {
	if res != nil {
		c.existSubRes = *res
	}
}

func (c *mockSubService) GetSubReq() *testdata.SubReq {
	return &c.subReq
}

func (c *mockSubService) SubBroadcast(ctx context.Context, req *testdata.SubReq) error {
	c.subReq.ConnId = req.ConnId
	c.subReq.SubscriptionId = req.SubscriptionId
	return nil
}

func (c *mockSubService) UnSubBroadcast(ctx context.Context, req *testdata.UnsubReq) error {
	return nil
}

func (c *mockSubService) ExistSubscription(ctx context.Context, req *testdata.ExistSubReq, res *testdata.ExistSubRes) error {
	res.Exists = c.existSubRes.Exists
	return nil
}

func (c *mockSubService) SyncSnapshot(ctx context.Context, req *testdata.SyncSnapshotReq) error {
	return nil
}

var _ testdata.SubService = (*mockSubService)(nil)

type mockPubListener struct {
	testdata.PubReq
	ServerExistSubReq *testdata.ExistSubReq
	ServerExistSubRes *testdata.ExistSubRes
}

func (m *mockPubListener) ExistSub(ctx context.Context, req *testdata.ExistSubReq, res *testdata.ExistSubRes) error {
	m.ServerExistSubReq = req
	res.Exists = m.ServerExistSubRes.Exists
	return nil
}

func (m *mockPubListener) GetServerPubReq() *testdata.PubReq {
	return &m.PubReq
}
func (m *mockPubListener) GetServerExistSubReq() *testdata.ExistSubReq {
	return m.ServerExistSubReq
}

func (m *mockPubListener) SetServerExistSubRes(res *testdata.ExistSubRes) {
	if res != nil {
		m.ServerExistSubRes = res
	}
}

func (m *mockPubListener) Pub(ctx context.Context, req *testdata.PubReq) error {
	m.PubReq = *req
	return nil
}

func (m *mockPubListener) SubFailed(ctx context.Context, req *testdata.SubFailedReq) error {
	return nil
}

var _ testdata.PubEmitter = (*mockPubListener)(nil)

type APISuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	serverClient   *server.ServerXClient
	xClient        *server.ServerXClient
	server         server.Server
	serverEndpoint remote.Endpoint

	mockConsumerService *mockSubService
	mockPubListener     *mockPubListener
	consumerServer      testdata.ConsumerServiceServer
	consumerClient      testdata.ConsumerServiceClient
}

func (s *APISuite) SetupSuite() {

}

func (s *APISuite) SetupTest() {
	s.mockPubListener = new(mockPubListener)
	s.mockConsumerService = new(mockSubService)
	s.ctrl = gomock.NewController(s.T())
	s.server = server.NewServer(server.WithAddress(":8080"), server.WithRemoteOption(remote.WithOnActive(func(ctx context.Context, end remote.Endpoint) (context.Context, error) {
		s.serverEndpoint = end
		return ctx, nil
	})), server.WithCallMiddleware(func(next middleware.GenericCall) middleware.GenericCall {
		return func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) (err error) {
			fmt.Printf("call middleware REQ: %v\n", req)
			defer func() {
				fmt.Printf("call middleware RESP: %v\n", resp)
			}()
			return next(ctx, req, resp, info)
		}
	}))
	s.consumerServer = testdata.NewConsumerServer(s.server, s.mockConsumerService)
	go func() {
		if err := s.server.Run(); err != nil {
			s.FailNow("server run error:", err)
		}
	}()
	time.Sleep(time.Second)
	s.consumerClient = testdata.NewConsumerServiceClient("consumerService:5002", s.mockPubListener, client.WithResolver(testdata.NewDirectResolver(
		discovery.NewNode("127.0.0.1:8080", nil))))
	time.Sleep(time.Second)

}

func (s *APISuite) TearDownTest() {
	s.ctrl.Finish()
	s.Nil(s.server.Stop(), "server关闭时错误为nil")
}

func (s *APISuite) TestCallOneway() {
	ctx, cancel := context.WithTimeoutCause(context.TODO(), time.Second, fmt.Errorf("发送超时"))
	defer cancel()
	req := &testdata.SubReq{SubscriptionId: 1, ConnId: 1}
	s.Nil(s.consumerClient.SubBroadcast(ctx, req))
	// oneway method 需要等待一会儿对端才能够接收到消息
	time.Sleep(time.Millisecond * 10)
	s.T().Log("connId:", s.mockConsumerService.GetSubReq().ConnId, req.ConnId)
	s.T().Log("subId:", s.mockConsumerService.GetSubReq().SubscriptionId, req.SubscriptionId)
	s.T().Log("ctx err:", ctx.Err())
}

func (s *APISuite) TestCall() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	retVal := &testdata.ExistSubRes{
		Exists: true,
	}
	s.mockConsumerService.SetExistSubscriptionRes(retVal)
	res := new(testdata.ExistSubRes)
	s.Nil(s.consumerClient.ExistSubscription(ctx, &testdata.ExistSubReq{}, res))
	s.Equal(res.Exists, retVal.Exists)
}

func (s *APISuite) TestServerOneway() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	// 1. 客户端建立连接，发送请求
	s.Nil(s.consumerClient.ExistSubscription(ctx, &testdata.ExistSubReq{}, new(testdata.ExistSubRes)))

	// 2. server端获取到连接的信息
	ctx = metadata.WithEndpoint(ctx, s.serverEndpoint)
	str := "hello"
	req := &testdata.PubReq{Content: []byte(str)}
	s.Nil(s.consumerServer.Pub(ctx, req), "serverClient成功发送消息")
	time.Sleep(time.Millisecond * 10)
	s.T().Log("客户端接收到的payload:", string(s.mockPubListener.GetServerPubReq().Content),
		"服务端发送的payload:", str)
}

func (s *APISuite) TestServerCall() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	// 1. 客户端建立连接，发送请求
	s.Nil(s.consumerClient.ExistSubscription(ctx, &testdata.ExistSubReq{}, new(testdata.ExistSubRes)))
	time.Sleep(time.Millisecond * 100)
	// 2. server端获取到连接的信息
	ctx2, cancel2 := context.WithTimeout(context.TODO(), time.Second)
	defer cancel2()
	endpointCtx := metadata.WithEndpoint(ctx2, s.serverEndpoint)
	req := &testdata.ExistSubReq{SubscriptionId: 1}
	res := &testdata.ExistSubRes{}

	retVal := &testdata.ExistSubRes{
		Exists: true,
	}
	s.mockPubListener.SetServerExistSubRes(retVal)
	err := s.consumerServer.ExistSub(endpointCtx, req, res)
	s.Nil(err, "serverClient成功发送消息")
	time.Sleep(time.Millisecond * 100)
	s.Equal(s.mockPubListener.GetServerExistSubReq().SubscriptionId, req.SubscriptionId)
	s.Equal(retVal.Exists, res.Exists)
}

func TestXClientSuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
