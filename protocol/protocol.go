package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"rpc-oneway/util"
	"runtime"

	"github.com/smallnest/rpcx/log"
)

type FrameType byte

var (
	bufferPool          = util.NewLimitedPool(512, 4096)
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

const (
	InitialType     FrameType = iota
	InitialDoneType           = 1
	DataType                  = 2
	PingType                  = 6
	PongType                  = 7
	CloseType                 = 8
)

// SizeableMarshaller todo 将这里的接口和protobuf解码
type SizeableMarshaller interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

// Message DataFrame [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
// 需要有一个控制信号, 用于在检测时，完成这个处理过程
type Message struct {
	From        net.Conn
	FixedHeader [7]byte
	TraceId     []byte
	SpanId      []byte
	MsgType     byte
	AltKey      string
	Metadata    map[string]string
	Data        any
	DataBuf     []byte
}

func NewMessage() *Message {
	return &Message{}
}

// EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
// [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
const (
	fixedHeaderSize = 7
)

// todo 改进这个过程，让buffer的释放显得没有这么多突兀
func (m *Message) EncodeSlicePointer() (*[]byte, error) {
	// 判断data的编解码类型是否支持
	// todo 决定这里是否可以不适用呢
	payload := m.Data.(SizeableMarshaller)
	if payload == nil {
		return nil, ErrUnsupportedCodec
	}

	dataSize := payload.Size()
	msgSize := fixedHeaderSize + dataSize
	trace := false
	if len(m.TraceId) != 0 { // 说明此时有traceId
		trace = true
		msgSize += 24
	}

	//if m.AltKey != "" {
	//	msgSize += 1 + len(m.AltKey) // 第一个字节表示字符串的长度   // Metadata的处理形式类似
	//}

	bufP := bufferPool.Get(msgSize)
	buf := *bufP
	buf[0] = DataType
	buf[1] = m.MsgType
	binary.BigEndian.PutUint32(buf[2:6], uint32(dataSize))

	startIndex := 7
	if trace {
		buf[6] = 0x1               // set flag
		copy(buf[7:22], m.TraceId) // 16byte
		copy(buf[23:30], m.SpanId) // 8byte
		startIndex = 31
	}
	_, err := payload.MarshalToSizedBuffer(buf[startIndex:])
	if err != nil {
		return nil, err
	}
	return bufP, nil
}

func (m *Message) Reset() {
	m.TraceId = nil
	m.SpanId = nil
}

// PutData puts the byte slice into pool.
func PutData(data *[]byte) {
	bufferPool.Put(data)
}

// Decode decodes a message from reader.
func (m *Message) Decode(r io.Reader) error {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			log.Errorf("panic in message decode: %v, stack: %s", err, errStack[:n])

		}
	}()

	_, err := io.ReadFull(r, m.FixedHeader[:])
	if err != nil {
		return err
	}
	// 解析出msg，traceId, spanId, 并将traceId, spanId设置在ctx中, 用户层从ctx中获取
	m.MsgType = m.FixedHeader[1]
	if traced(m.FixedHeader[6]) {
		_, err := io.ReadFull(r, m.TraceId)
		if err != nil {
			return err
		}
		_, err = io.ReadFull(r, m.SpanId)
		if err != nil {
			return err
		}
	}
	payloadLen := binary.BigEndian.Uint32(m.FixedHeader[2:6])

	if cap(m.DataBuf) >= int(payloadLen) { // reuse data
		m.DataBuf = m.DataBuf[:payloadLen]
	} else {
		m.DataBuf = make([]byte, payloadLen)
	}

	_, err = io.ReadFull(r, m.DataBuf)
	if err != nil {
		return err
	}
	return err
}

func traced(flag byte) bool { // 表示二进制的第一位是1
	return flag&0x80 == 1
}
