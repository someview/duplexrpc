package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
)

// SizeableMarshaller todo 将这里的接口和protobuf解码
type SizeableMarshaller interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

type BufGetter func(msgSize int) *[]byte

// message DataFrame [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
type message struct {
	data    any    // 原始数据
	payload []byte // 编码后的二进制数据
	len     uint32 // 当前msg的长度
	msgContext
}

// Recycle free msg to pool
func (m *message) Recycle() {

}

// SetReq implements Message.
func (m *message) SetReq(req any) {
	m.data = req
}

// Encode this method should be called after all msg info has been setted
// bufP must be larger than m.Len()

// Len implements Message.
func (m *message) Len() uint32 {
	return m.len
}

// Payload implements Message.
func (m *message) Payload() []byte {
	return m.payload
}

// // EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
// // [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
const (
	fixedHeaderSize = 7
)

// // todo 改进这个过程，让buffer的释放显得没有这么多突兀
func (m *message) Encode() (*[]byte, error) {
	// 判断data的编解码类型是否支持
	// todo 添加对binary类型的支持
	payload := m.data.(SizeableMarshaller)
	if payload == nil {
		return nil, ErrUnsupportedCodec
	}

	dataSize := payload.Size()
	msgSize := fixedHeaderSize + dataSize
	trace := false
	if len(m.traceInfo) != 0 { // 说明此时有traceId
		trace = true
		msgSize += 24
	}

	//if m.AltKey != "" {
	//	msgSize += 1 + len(m.AltKey) // 第一个字节表示字符串的长度   // Metadata的处理形式类似
	//}

	bufP := bufferPool.Get(msgSize)
	buf := *bufP
	buf[0] = DataType
	buf[1] = m.fixedHeader[1]
	binary.BigEndian.PutUint32(buf[2:6], uint32(dataSize))

	startIndex := 7
	if trace {
		buf[6] = 0x1                 // set flag
		copy(buf[7:30], m.traceInfo) // 24byte
		startIndex = 31
	}
	_, err := payload.MarshalToSizedBuffer(buf[startIndex:])
	if err != nil {
		return nil, err
	}
	return bufP, nil
}

func (m *message) Reset() {
	m.traceInfo = nil
}

// // PutData puts the byte slice into pool.

// // Decode decodes a message from reader.
func (m *message) Decode(r io.Reader) error {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			fmt.Printf("panic in message decode: %v, stack: %s", err, errStack[:n])

		}
	}()

	_, err := io.ReadFull(r, m.fixedHeader[:])
	if err != nil {
		return err
	}
	// 解析出msg，traceId, spanId, 并将traceId, spanId设置在ctx中, 用户层从ctx中获取
	if traced(m.fixedHeader[6]) {
		_, err := io.ReadFull(r, m.traceInfo)
		if err != nil {
			return err
		}
	}
	payloadLen := binary.BigEndian.Uint32(m.fixedHeader[2:6])

	if cap(m.payload) >= int(payloadLen) { // reuse data
		m.payload = m.payload[:payloadLen]
	} else {
		m.payload = make([]byte, payloadLen)
	}

	_, err = io.ReadFull(r, m.payload)
	if err != nil {
		return err
	}
	return err
}

func traced(flag byte) bool { // 表示二进制的第一位是1
	return flag&0x80 == 1
}

var _ Message = (*message)(nil)

func NewMessage() *message {
	return &message{}
}

// fixme  make thie readable
func PutData(data *[]byte) {
	bufferPool.Put(data)
}
