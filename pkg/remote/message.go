package remote

import (
	netpoll "github.com/cloudwego/netpoll"

	"github.com/someview/dt/pool"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
)

// ServiceSearcher is used to search the service info by service name and method name,
// strict equals to true means the service name must match the registered service name.
type ServiceSearcher interface {
	SearchService(svcName, methodName string, strict bool) *service.ServiceInfo
}

var messagePool = pool.NewSyncPool[*message](func() any {
	msg := new(message)
	return msg
})

type message struct {
	ri rpcinfo.RPCInfo

	data    any            // 原始数据
	payload netpoll.Reader // 编码后的二进制数据
	len     int            // 当前msg的长度

	version Version
	typ     MessageType

	metadata []byte

	codec Codec
}

func (m *message) Codec() Codec {
	return m.codec
}

func (m *message) SetPayload(reader netpoll.Reader) {
	m.payload = reader
}

func (m *message) SetLen(l int) {
	m.len = l
}

func (m *message) RPCInfo() rpcinfo.RPCInfo {
	return m.ri
}

func (m *message) SetData(data any) {
	m.data = data
}

func (m *message) Data() any {
	return m.data
}

func (m *message) SetVersion(v Version) {
	m.version = v
}

func (m *message) Version() Version {
	return m.version
}

func (m *message) SetMetadata(metadata []byte) {
	m.metadata = metadata
}

func (m *message) Metadata() []byte {
	return m.metadata
}

func (m *message) SetMsgType(typ MessageType) {
	m.typ = typ
}

func (m *message) MsgType() MessageType {
	return m.typ
}

func (m *message) zero() {
	if m.payload != nil {
		_ = util.PutBufFromSlice(m.payload)
		m.payload = nil
	}

	m.ri = nil

	m.data = nil
	m.len = 0
	m.version = 0
	m.typ = 0

	m.metadata = nil
	m.codec = nil

}

// Recycle free msg to pool
func (m *message) Recycle() {
	m.zero()
	messagePool.Put(m)
}

// Len implements Message.
func (m *message) Len() int {
	return m.len
}

// Payload implements Message.
func (m *message) Payload() netpoll.Reader {
	p := m.payload
	return p
}

var _ Message = (*message)(nil)

func NewMessage(ri rpcinfo.RPCInfo, codec Codec) Message {
	msg := messagePool.Get()
	msg.ri = ri
	msg.codec = codec

	switch ri.InteractionMode() {
	case rpcinfo.Oneway:
		msg.SetMsgType(OnewayType)
	case rpcinfo.Unary:
		msg.SetMsgType(ReqType)
	}

	return msg
}

func NewRespMessage(ri rpcinfo.RPCInfo, codec Codec) Message {
	msg := messagePool.Get()
	msg.ri = ri
	msg.codec = codec
	msg.SetMsgType(ResponseType)
	return msg
}

func NewErrorMessage(ri rpcinfo.RPCInfo, err error, codec Codec) Message {
	msg := messagePool.Get()
	msg.ri = ri
	msg.codec = codec
	msg.SetMsgType(ErrorResponseType)
	msg.data = err
	return msg
}

func NewCloseMessage(codec Codec) Message {
	msg := messagePool.Get()
	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, rpcinfo.NewInvocation("none", "none"))
	msg.ri = ri
	msg.SetMsgType(CloseType)
	msg.codec = codec
	return msg
}

type heartbeatMessage struct{}

func (h heartbeatMessage) Size() int {
	return 0
}

func (h heartbeatMessage) MarshalToSizedBuffer(bytes []byte) (int, error) {
	return 0, nil
}

func (h heartbeatMessage) Unmarshal(bytes []byte) error {
	return nil
}

func NewPingMessage(codec Codec) Message {
	msg := messagePool.Get()
	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, rpcinfo.NewInvocation("heartbeat", "heartbeat"))
	msg.ri = ri
	msg.SetMsgType(PingType)
	msg.codec = codec
	msg.data = heartbeatMessage{}
	return msg
}

func NewPongMessage(codec Codec) Message {
	msg := messagePool.Get()
	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, rpcinfo.NewInvocation("heartbeat", "heartbeat"))
	msg.ri = ri
	msg.SetMsgType(PongType)
	msg.codec = codec
	msg.data = heartbeatMessage{}
	return msg

}
