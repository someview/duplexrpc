package protocol

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git/dt/pool"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

var msgContextPool = pool.NewSyncPool[*msgContext](func() any {
	return &msgContext{}
})

// dataFrame传输的消息的上下文
type msgContext struct {
	context.Context
	from      rpcinfo.RPCEndpoint
	to        rpcinfo.RPCEndpoint
	traceInfo []byte // 24byte
	metadata  map[string]string
	seqID     uint32

	lbKey uint32
	// 注意结构体字段的排列对内存对齐的影响，msgType这样的字段应该放置在最后
	fixedHeader [7]byte
}

func (m *msgContext) MethodName() rpcinfo.MethodType {
	return rpcinfo.MethodType(m.MsgType())
}

func (m *msgContext) SeqID() uint32 {
	return m.seqID
}

func (m *msgContext) SetSeqID(id uint32) {
	m.seqID = id
}

func (m *msgContext) From() rpcinfo.RPCEndpoint {
	return m.from
}

func (m *msgContext) To() rpcinfo.RPCEndpoint {
	return m.to
}

// SetLBKey implements MsgContext.
func (m *msgContext) SetLBKey(key uint32) {
	m.lbKey = key
}

func (m *msgContext) SetFrom(from rpcinfo.RPCEndpoint) {
	m.from = from
}

func (m *msgContext) SetTo(to rpcinfo.RPCEndpoint) {
	m.to = to
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

func (m *msgContext) Recycle() {
	m.zero()
	msgContextPool.Put(m)
}

func (m *msgContext) zero() {
	m.Context = nil
	m.from = nil
	m.to = nil
	m.traceInfo = nil
	m.metadata = nil
	//m.lbKey = 0
	//m.fixedHeader = [7]byte{}
	m.seqID = 0
}

var _ MsgContext = (*msgContext)(nil)

func NewMsgContext(ctx context.Context) *msgContext {
	mc := msgContextPool.Get()
	mc.Context = ctx
	return mc
}

func NewContextWithMsgMD(ctx context.Context) MsgContext {
	return NewMsgContext(ctx)
}
