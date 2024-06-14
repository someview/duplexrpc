package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git/dt/pool"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
	"runtime"
)

var messagePool = pool.NewSyncPool[*message](func() any {
	msg := new(message)
	return msg
})

// SizeableMarshaller todo 将这里的接口和protobuf解码
type SizeableMarshaller interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

// message DataFrame [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
type message struct {
	data       any            // 原始数据
	payload    netpoll.Reader // 编码后的二进制数据
	len        uint32         // 当前msg的长度
	payloadLen uint32         // payload长度

	*msgContext
	muxextend.SizedEncodable
}

func (m *message) zero() {
	if m.payload != nil {
		_ = util.PutBufFromSlice(m.payload)
		m.payload = nil
	}
	m.data = nil
	m.len = 0
	m.msgContext = nil
	m.payloadLen = 0
}

// Recycle free msg to pool
func (m *message) Recycle() {
	m.msgContext.Recycle()
	m.zero()
	messagePool.Put(m)
}

// SetReq implements Message.
func (m *message) SetReq(req any) {
	m.data = req
}

func (m *message) GetReq() any {
	return m.data
}

// Len implements Message.
func (m *message) Len() uint32 {
	return uint32(m.EncodedLen())
}

// Payload implements Message.
func (m *message) Payload() netpoll.Reader {
	p := m.payload
	return p
}

// --- 实现mux.SizedEncodable接口 ---

func (m *message) EncodedLen() int {
	// 在SetReq时就计算len会不会更好？
	if m.len <= 0 {
		// 判断data的编解码类型是否支持
		// todo 添加对binary类型的支持

		dataSize := m.payloadSize()
		msgSize := FixHeaderLength + dataSize
		if len(m.traceInfo) != 0 {
			msgSize += TraceLength
		}
		if m.seqID != 0 {
			msgSize += SeqLength
		}
		m.len = msgSize
	}

	return int(m.len)
}

func (m *message) payloadSize() uint32 {
	if m.payloadLen <= 0 {
		if m.payload != nil {
			m.payloadLen = uint32(m.payload.Len())
		} else {
			// fixme: 断言失败时会发生panic
			m.payloadLen = uint32(m.data.(SizeableMarshaller).Size())
		}
	}
	return m.payloadLen
}

func (m *message) EncodeToWriter(writer netpoll.Writer) error {
	payload, ok := m.data.(SizeableMarshaller)
	if !ok {
		return ErrUnsupportedCodec
	}

	buf, err := writer.Malloc(m.EncodedLen())
	if err != nil {
		return err
	}

	buf[0] = DataType                                     // FrameType
	buf[1] = m.fixedHeader[1]                             // MsgType
	binary.BigEndian.PutUint32(buf[2:6], m.payloadSize()) // payload Len

	// 这段内存会被LinkBuffer重用，一定要置0 ！！！！！！！！
	// 非常难找到的Bug
	buf[6] = 0
	startIndex := FixHeaderLength
	if len(m.traceInfo) != 0 {
		buf[6] = SetFlagBit(buf[6], TraceBit)                     // set flag
		copy(buf[startIndex:startIndex+TraceLength], m.traceInfo) // 24byte
		startIndex += TraceLength
	}
	if m.seqID != 0 {
		buf[6] = SetFlagBit(buf[6], SeqBit)
		binary.BigEndian.PutUint32(buf[startIndex:startIndex+SeqLength], m.seqID) // 4byte
		startIndex += SeqLength

	}

	_, err = payload.MarshalToSizedBuffer(buf[startIndex:])
	if err != nil {
		return err
	}
	return writer.Flush()

}

func (m *message) DecodeFromReader(reader netpoll.Reader) (any, error) {
	sli := reader.(netpoll.Sliceable)
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			fmt.Printf("panic in message decode: %v, stack: %s", err, errStack[:n])

		}
	}()

	var headerLength int
	fh, err := reader.Peek(FixHeaderLength)
	if err != nil {
		return nil, err
	}
	m.fixedHeader = [7]byte(fh)
	headerLength += 7

	// 解析optional
	var optionLength int
	var (
		hasTrace = false
		hasSeq   = false
	)

	if HasFlagBit(m.fixedHeader[6], TraceBit) {
		hasTrace = true
		optionLength += TraceLength
	}
	if HasFlagBit(m.fixedHeader[6], SeqBit) {
		hasSeq = true
		optionLength += SeqLength
	}
	headerBuf, err := reader.Peek(headerLength + optionLength)
	if err != nil {
		return m, err
	}
	optionBuf := headerBuf[headerLength:]
	optionIdx := 0
	if hasTrace {
		m.traceInfo = optionBuf[optionIdx : optionIdx+TraceLength]
		optionIdx += TraceLength
	}
	if hasSeq {
		m.seqID = binary.BigEndian.Uint32(optionBuf[optionIdx : optionIdx+SeqLength])
		optionIdx += SeqLength
	}
	headerLength += optionLength
	payloadLen := binary.BigEndian.Uint32(m.fixedHeader[2:6])

	// 这里slice出来的整段内存是包括header的
	r, err := util.SliceBuf(headerLength+int(payloadLen), sli)
	if err != nil {
		return nil, err
	}
	// 跳过header的内容
	err = r.Skip(headerLength)
	if err != nil {
		return nil, err
	}
	m.payload = r
	m.payloadLen = payloadLen
	m.len = uint32(headerLength + int(payloadLen))
	return m, nil
}

func (m *message) Reset() {
	m.traceInfo = nil
}

var _ Message = (*message)(nil)

func NewMessage() *message {
	msg := messagePool.Get()
	msg.msgContext = NewMsgContext(context.TODO())
	return msg
}

func NewMessageWithCtx(ctx context.Context) Message {
	msgCtx, ok := ctx.(*msgContext)
	if ok {
		return newMessageWithMsgCtx(msgCtx)
	}
	msg := messagePool.Get()
	msg.msgContext = NewMsgContext(ctx)
	return msg
}

func newMessageWithMsgCtx(ctx *msgContext) Message {
	msg := messagePool.Get()
	msg.msgContext = ctx
	return msg
}
