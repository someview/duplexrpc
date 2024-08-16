package bufwriter

import (
	"context"
	"math/rand"
	"testing"
	"time"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func mockEncode(size int, writer netpoll.Writer) {
	writer.WriteBinary(make([]byte, size))
	writer.Flush()
}

// batchContainerSuite is a test suite for batchContainer.
type batchContainerSuite struct {
	suite.Suite
	batchContainer *batchContainer
	TC             *testing.T
	ctrl           *gomock.Controller
	maxMsgNum      int
	maxBufSize     int
}

func (s *batchContainerSuite) SetupSuite() {
	s.maxBufSize = 1024
	s.maxMsgNum = 100
}

func (s *batchContainerSuite) TearDownSuite() {
	s.maxBufSize = 1024
	s.maxMsgNum = 100
}

func (s *batchContainerSuite) SetupTest() {
	s.batchContainer = newBatchContainer(s.maxBufSize, s.maxMsgNum).(*batchContainer)
	s.TC = s.Suite.T()
	s.ctrl = gomock.NewController(s.TC)
}

func (s *batchContainerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *batchContainerSuite) TestAdd_消息条数ok_缓冲区ok() {
	data := netpoll.NewSizedLinkBuffer(1024)

	err := s.batchContainer.Add(context.Background(), data)
	s.Nil(err)
}

func (s *batchContainerSuite) TestAdd_消息条数ok_缓冲区耗尽() {
	data := netpoll.NewSizedLinkBuffer(1024)
	mockEncode(s.maxBufSize+1, data)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.batchContainer.Add(ctx, data)
	s.Nil(err)
	err = s.batchContainer.Add(ctx, data)
	s.NotNil(err)
}

func (s *batchContainerSuite) TestAdd_消息条数err_缓冲区ok() {

	data := netpoll.NewSizedLinkBuffer(1024)
	mockEncode(s.maxBufSize/s.maxMsgNum, data)

	for i := 0; i <= s.maxMsgNum; i++ {
		_ = s.batchContainer.Add(context.Background(), data)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.batchContainer.Add(ctx, data)
	s.Nil(err)
	s.True(s.batchContainer.IsFull(), "只要缓冲区未耗尽就能写入数据,消息条数达到上限缓冲区会被标记为full")
}

func (s *batchContainerSuite) TestIsEmpty() {
	s.True(s.batchContainer.IsEmpty())
	data := netpoll.NewSizedLinkBuffer(1024)
	dataSize := s.maxBufSize
	mockEncode(dataSize, data)

	err := s.batchContainer.Add(context.Background(), data)
	s.Nil(err)
	s.False(s.batchContainer.IsEmpty())
	// 测试flushTo，flush完毕以后，必定是empty, 不管是否成功flush
	conn := NewMockConnection(s.ctrl)
	conn.EXPECT().Writer().Return(netpoll.NewLinkBuffer())

	n, err := s.batchContainer.FlushTo(conn)
	s.Equal(dataSize, n)
	s.Nil(err)
	s.True(s.batchContainer.IsEmpty())
}

func (s *batchContainerSuite) TestFlushCount() {
	var dataS []*netpoll.LinkBuffer
	count := 0
	flushed := 0
	for i := 0; i < 10; i++ {
		lb := netpoll.NewSizedLinkBuffer(1024)
		mockEncode(rand.Intn(4096), lb)
		count += lb.Len()
		dataS = append(dataS, lb)
	}
	conn := NewMockConnection(s.ctrl)
	conn.EXPECT().Writer().Return(netpoll.NewLinkBuffer()).AnyTimes()
	for _, data := range dataS {
	RETRY:
		ok := s.batchContainer.TryAdd(data)
		if !ok {
			n, err := s.batchContainer.FlushTo(conn)
			flushed += n
			s.Nil(err)
			goto RETRY
		}

	}
	n, err := s.batchContainer.FlushTo(conn)
	flushed += n
	s.Nil(err)

	s.Equal(count, flushed)

}

func TestBatchContainerSuite(t *testing.T) {
	suite.Run(t, new(batchContainerSuite))
}
