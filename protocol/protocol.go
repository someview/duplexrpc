package protocol

import (
	"encoding/binary"
	"github.com/smallnest/rpcx/util"
	"io"
)

type FrameType byte

const (
	InitialType     FrameType = iota
	InitialDoneType           = 1
	DataType                  = 2
	PingType                  = 6
	PongType                  = 7
	CloseType                 = 8
)

type Message interface {
	EncodeSlicePointer() *[]byte
	Reset()
	io.WriterTo
}

// DataFrame [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
// 需要有一个控制信号, 用于在检测时，完成这个处理过程
type DataFrame struct {
	TraceId  []byte
	SpanId   []byte
	AltKey   string
	MetaData map[string]string
	Payload  []byte
	msgType  byte
}

type Stream struct {
	Size int64
	Done chan struct{}
}

func NewDataFrame() *DataFrame {
	return &DataFrame{}
}

// EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
func (m Message) EncodeSlicePointer() *[]byte {

	m.Payload

	encodeMetadata(m.Metadata, bb)
	meta := bb.Bytes()

	spL := len(m.ServicePath)
	smL := len(m.ServiceMethod)

	var err error
	payload := m.Payload
	if m.CompressType() != None {
		compressor := Compressors[m.CompressType()]
		if compressor == nil {
			m.SetCompressType(None)
		} else {
			payload, err = compressor.Zip(m.Payload)
			if err != nil {
				m.SetCompressType(None)
				payload = m.Payload
			}
		}
	}

	totalL := (4 + spL) + (4 + smL) + (4 + len(meta)) + (4 + len(payload))

	// header + dataLen + spLen + sp + smLen + sm + metaL + meta + payloadLen + payload
	metaStart := 12 + 4 + (4 + spL) + (4 + smL)

	payLoadStart := metaStart + (4 + len(meta))
	l := 12 + 4 + totalL

	data := bufferPool.Get(l)
	copy(*data, m.Header[:])

	// totalLen
	binary.BigEndian.PutUint32((*data)[12:16], uint32(totalL))

	binary.BigEndian.PutUint32((*data)[16:20], uint32(spL))
	copy((*data)[20:20+spL], util.StringToSliceByte(m.ServicePath))

	binary.BigEndian.PutUint32((*data)[20+spL:24+spL], uint32(smL))
	copy((*data)[24+spL:metaStart], util.StringToSliceByte(m.ServiceMethod))

	binary.BigEndian.PutUint32((*data)[metaStart:metaStart+4], uint32(len(meta)))
	copy((*data)[metaStart+4:], meta)

	binary.BigEndian.PutUint32((*data)[payLoadStart:payLoadStart+4], uint32(len(payload)))
	copy((*data)[payLoadStart+4:], payload)

	return data
}
