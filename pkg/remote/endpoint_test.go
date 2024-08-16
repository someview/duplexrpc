package remote

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
	"go.uber.org/mock/gomock"
	"testing"
)

type mockAddr struct {
}

func (m *mockAddr) Network() string {
	return "tcp"
}

func (m *mockAddr) String() string {
	return "localhost"
}

type endpointSuite struct {
	suite.Suite

	endpoint *endpoint

	ctrl                  *gomock.Controller
	rootCtx               context.Context
	mockConn              *bufwriter.MockConnection
	mockTransHandler      *MockTransHandler
	mockConnectionHandler *MockConnectionHandler
	mockBufWriter         *bufwriter.MockBufWriter
}

func (e *endpointSuite) SetupSuite() {
	e.rootCtx = context.TODO()
}

func (e *endpointSuite) TearDownSuite() {

}

func (e *endpointSuite) initEndpoint() {

	e.mockConnectionHandler.EXPECT().OnActive(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
		func(ctx context.Context, end Endpoint) (context.Context, error) {
			ctx = context.WithValue(ctx, "test", "test")
			return ctx, nil
		},
	)

	e.mockConn.EXPECT().AddCloseCallback(gomock.Any()).Times(1)
	e.mockConn.EXPECT().SetOnRequest(gomock.Any()).Times(1)
	e.mockConn.EXPECT().RemoteAddr().Return(&mockAddr{}).Times(1)

	ep, err := NewEndpoint(e.rootCtx, e.mockConnectionHandler, e.mockConn, e.mockBufWriter)
	if err != nil {
		e.T().Fatal(err)
	}
	e.endpoint = ep.(*endpoint)
}

func (e *endpointSuite) SetupTest() {
	e.ctrl = gomock.NewController(e.T())
	e.mockConn = bufwriter.NewMockConnection(e.ctrl)
	e.mockTransHandler = NewMockTransHandler(e.ctrl)
	e.mockBufWriter = bufwriter.NewMockBufWriter(e.ctrl)
	e.mockConnectionHandler = NewMockConnectionHandler(e.ctrl)
}

func (e *endpointSuite) TearDownTest() {
	e.ctrl.Finish()
}

func (e *endpointSuite) TestContext() {
	e.initEndpoint()
	e.NotNil(e.endpoint.Context())
	e.NotNil(e.endpoint.Context().Value("test"))
}

func (e *endpointSuite) TestOnRequest() {
	e.initEndpoint()
	e.mockConnectionHandler.EXPECT().OnRead(e.endpoint.ctx, e.endpoint).Return(nil).Times(1)
	err := e.endpoint.OnRequest(context.TODO(), e.mockConn)
	e.NoError(err)
}

func (e *endpointSuite) TestClose() {
	e.initEndpoint()

	e.mockConn.EXPECT().Close().Times(1).Return(nil)
	e.mockBufWriter.EXPECT().Close().Times(1).Return(nil)

	e.NoError(e.endpoint.Close())
	e.Error(e.endpoint.ctx.Err())

}

func (e *endpointSuite) TestGracefulShutdown() {
	e.initEndpoint()

	e.mockBufWriter.EXPECT().Close().Times(1)

	err := e.endpoint.GracefulShutdown(context.TODO())
	e.NoError(err)
}

func TestEndpointSuite(t *testing.T) {
	suite.Run(t, new(endpointSuite))
}
