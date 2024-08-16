package remote

import (
	"context"
	"testing"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/stretchr/testify/suite"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"go.uber.org/mock/gomock"
)

type defaultCodecSuite struct {
	suite.Suite
	codec *defaultCodec

	ctrl *gomock.Controller
}

func (c *defaultCodecSuite) SetupSuite() {
	c.codec = NewDefaultCodeC().(*defaultCodec)
}

func (c *defaultCodecSuite) TearDownSuite() {

}

func (c *defaultCodecSuite) SetupTest() {
	c.ctrl = gomock.NewController(c.T())

}

func (c *defaultCodecSuite) TearDownTest() {
	c.ctrl.Finish()
}

func (c *defaultCodecSuite) TestEncodeDecode() {
	servicePath := "asdfsafasdf"
	methodPath := "asdfpko0ojaf[p"
	seqID := uint32(102)
	dataSize := 1024
	version := Version(1)
	msgType := ReqType
	py := make([]byte, dataSize)
	var trace []byte
	for i := 0; i < 24; i++ {
		trace = append(trace, byte(i))
	}
	var metadata []byte

	mockInvocationSetter := rpcinfo.NewMockInvocationSetter(c.ctrl)
	mockInvocationSetter.EXPECT().SetServiceName(servicePath)
	mockInvocationSetter.EXPECT().SetMethodName(methodPath)
	mockInvocationSetter.EXPECT().SetSeqID(seqID)
	mockInvocationSetter.EXPECT().SetTraceInfo(trace)

	mockInvocation := rpcinfo.NewMockInvocation(c.ctrl)
	mockInvocation.EXPECT().SeqID().Return(seqID).AnyTimes()
	mockInvocation.EXPECT().ServiceName().Return(servicePath).AnyTimes()
	mockInvocation.EXPECT().MethodName().Return(methodPath).AnyTimes()
	mockInvocation.EXPECT().TraceInfo().Return(trace).AnyTimes()

	mockAnyInvocation := struct {
		rpcinfo.Invocation
		rpcinfo.InvocationSetter
	}{
		mockInvocation,
		mockInvocationSetter,
	}

	mockRPCInfo := rpcinfo.NewMockRPCInfo(c.ctrl)
	mockRPCInfo.EXPECT().Invocation().Return(mockAnyInvocation).AnyTimes()

	mockData := NewMockSizeableMarshaller(c.ctrl)
	mockData.EXPECT().Size().Return(dataSize).AnyTimes()
	mockData.EXPECT().MarshalToSizedBuffer(gomock.Any()).Do(func(b []byte) {
		copy(b, py)
	})
	mockData.EXPECT().Unmarshal(py)

	mockMsg := NewMockMessage(c.ctrl)
	mockMsg.EXPECT().RPCInfo().Return(mockRPCInfo).AnyTimes()

	mockMsg.EXPECT().Data().Return(mockData).AnyTimes()
	mockMsg.EXPECT().Metadata().Return(metadata).AnyTimes()
	mockMsg.EXPECT().SetMetadata(metadata).AnyTimes()
	mockMsg.EXPECT().Len().Return(0).AnyTimes()
	mockMsg.EXPECT().SetLen(gomock.Any()).AnyTimes()
	mockMsg.EXPECT().Version().Return(version).AnyTimes()
	mockMsg.EXPECT().SetVersion(version).AnyTimes()
	mockMsg.EXPECT().MsgType().Return(msgType).AnyTimes()
	mockMsg.EXPECT().SetMsgType(msgType).AnyTimes()
	mockMsg.EXPECT().SetPayload(gomock.Any()).Do(func(r netpoll.Reader) {
		b, err := r.ReadBinary(r.Len())
		c.NoError(err)
		err = mockData.Unmarshal(b)
		c.NoError(err)
	})

	lb := netpoll.NewLinkBuffer(1024)
	err := c.codec.Encode(context.TODO(), lb, mockMsg)
	c.NoError(err)
	lb.Flush()
	err = c.codec.Decode(context.TODO(), lb, mockMsg)
	c.NoError(err)
}

func TestDefaultCodecSuite(t *testing.T) {
	suite.Run(t, new(defaultCodecSuite))
}
