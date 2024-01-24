package protocol

import (
	"encoding/binary"
	"github.com/smallnest/rpcx/log"
	"io"
	"rpc-oneway/util"
	"runtime"
)

type FrameType byte

var bufferPool = util.NewLimitedPool(512, 4096)

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
	FixedHeader [7]byte
	TraceId     []byte
	SpanId      []byte
	MsgType     byte
	AltKey      string
	Metadata    map[string]string
	Payload     []byte
}

func NewMessage() *Message {
	return &Message{}
}

// EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
// [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
const (
	fixedHeaderSize = 7
)

func (m *Message) EncodeSlicePointer() (*[]byte, error) {
	msgSize := fixedHeaderSize

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

	startIndex := 8
	if trace {
		buf[7] = 0x1               // set flag
		copy(buf[8:24], m.TraceId) // 16byte
		copy(buf[24:32], m.SpanId) // 8byte
		startIndex = 33
	}
	copy(buf[startIndex:], m.Payload)
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
	m.MsgType = m.FixedHeader[2]
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
	lenData := make([]byte, payloadLen)
	_, err = io.ReadFull(r, lenData)
	if err != nil {
		return err
	}

	return err
}

func traced(flag byte) bool { // 表示二进制的第一位是1
	return flag&0x80 == 1
}
