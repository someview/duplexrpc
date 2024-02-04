package protocol

import (
	"context"
	"net"
)

// dataFrame传输的消息的上下文
type msgContext struct {
	context.Context
	from      net.Conn
	traceInfo []byte // 24byte
	metadata  map[string]string

	lbKey uint32
	// 注意结构体字段的排列对内存对齐的影响，msgType这样的字段应该放置在最后
	fixedHeader [7]byte
}

// SetLBKey implements MsgContext.
func (m *msgContext) SetLBKey(key uint32) {
	m.lbKey = key
}

func (m *msgContext) SetFrom(conn net.Conn) {
	m.from = conn
}

func (m *msgContext) From() net.Conn {
	return m.from
}

func (m *msgContext) TraceInfo() []byte {
	return m.traceInfo
}

func (m *msgContext) SetTraceInfo(traceInfo []byte) {
	m.traceInfo = traceInfo
}

func (m *msgContext) MsgType() byte {
	return m.fixedHeader[1]
}

func (m *msgContext) SetMsgType(b byte) {
	m.fixedHeader[1] = b
}

func (m *msgContext) Metadata() map[string]string {
	if m.metadata == nil {
		m.metadata = make(map[string]string)
	}
	return m.metadata
}

var _ MsgContext = (*msgContext)(nil)

func GetMsgContext(ctx context.Context) MsgContext {
	return &msgContext{
		Context: ctx,
	}
}

func PutMsgContext(ctx context.Context) {
}

func NewContextWithMsgMD(ctx context.Context) MsgContext {
	return GetMsgContext(ctx)
}
