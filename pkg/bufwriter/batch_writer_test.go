package bufwriter

import (
	"context"
	"testing"
	"time"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type batchWriterSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	batchWriter    *batchWriter
	maxBufSize     int
	maxMsgNum      int
	flushDelay     time.Duration
	MockConnection *MockConnection
}

func (s *batchWriterSuite) SetupSuite() {
	s.maxBufSize = 1024
	s.maxMsgNum = 10
	s.flushDelay = time.Millisecond
}

func (s *batchWriterSuite) TearDownSuite() {

}

func (s *batchWriterSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.Suite.T())
	s.MockConnection = NewMockConnection(s.ctrl)
	s.batchWriter = NewBatchWriter(s.maxBufSize, s.maxMsgNum,
		s.flushDelay, s.MockConnection, func(n int, err error) {

		}).(*batchWriter)
}

func (s *batchWriterSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *batchWriterSuite) TestAdd_正常写入() {

	dataSize := s.maxBufSize
	lb := netpoll.NewSizedLinkBuffer(dataSize)
	mockEncode(dataSize, lb)

	s.MockConnection.EXPECT().Writer().Return(netpoll.NewLinkBuffer(dataSize)).AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.batchWriter.Add(ctx, lb)
	time.Sleep(time.Millisecond)
	s.Equal(nil, err)
}

func TestBatchWriterSuite(t *testing.T) {
	suite.Run(t, new(batchWriterSuite))
}
