package protocol

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"math"
	"testing"
)

type testData struct {
	str string
}

func (t *testData) Size() int {
	return len(t.str)
}

func (t *testData) MarshalToSizedBuffer(p []byte) (int, error) {
	l := copy(p, t.str)
	return l, nil
}

func (t *testData) Unmarshal(p []byte) error {
	t.str = string(p)
	return nil
}

func TestProtocolCodeC(t *testing.T) {
	data := &testData{
		str: "testing",
	}
	msgType := byte(128)

	t.Run("无optional", func(t *testing.T) {
		msg := NewMessage()
		defer msg.Recycle()
		msg.SetReq(data)
		msg.SetMsgType(msgType)
		lb := netpoll.NewLinkBuffer()
		assert.NoError(t, msg.EncodeToWriter(lb))
		decodeMsg := NewMessage()
		defer decodeMsg.Recycle()
		_, err := decodeMsg.DecodeFromReader(lb)
		assert.NoError(t, err)
		decodeData := new(testData)
		ml, err := decodeMsg.Payload().Next(decodeMsg.Payload().Len())
		assert.NoError(t, err)
		assert.NoError(t, decodeData.Unmarshal(ml))
		assert.Equal(t, *data, *decodeData)
		assert.Equal(t, msgType, msg.MsgType())
	})

	t.Run("optional: traceInfo", func(t *testing.T) {
		trace := [TraceLength]byte{}
		for i := 0; i < TraceLength; i++ {
			trace[i] = byte(i) + 1
		}
		msg := NewMessage()
		defer msg.Recycle()
		msg.SetReq(data)
		msg.SetTraceInfo(trace[:])
		msg.SetMsgType(msgType)
		lb := netpoll.NewLinkBuffer()
		assert.NoError(t, msg.EncodeToWriter(lb))
		decodeMsg := NewMessage()
		defer decodeMsg.Recycle()
		_, err := decodeMsg.DecodeFromReader(lb)
		assert.NoError(t, err)
		decodeData := new(testData)
		ml, err := decodeMsg.Payload().Next(decodeMsg.Payload().Len())
		assert.NoError(t, err)
		assert.NoError(t, decodeData.Unmarshal(ml))
		assert.Equal(t, *data, *decodeData)
		assert.Equal(t, trace, [TraceLength]byte(decodeMsg.TraceInfo()))
		assert.Equal(t, msgType, msg.MsgType())
	})

	t.Run("optional: seq", func(t *testing.T) {
		seq := uint32(math.MaxUint32)
		msg := NewMessage()
		defer msg.Recycle()
		msg.SetReq(data)
		msg.SetSeqID(seq)
		msg.SetMsgType(msgType)
		lb := netpoll.NewLinkBuffer()
		assert.NoError(t, msg.EncodeToWriter(lb))
		decodeMsg := NewMessage()
		defer decodeMsg.Recycle()
		_, err := decodeMsg.DecodeFromReader(lb)
		assert.NoError(t, err)
		decodeData := new(testData)
		ml, err := decodeMsg.Payload().Next(decodeMsg.Payload().Len())
		assert.NoError(t, err)
		assert.NoError(t, decodeData.Unmarshal(ml))
		assert.Equal(t, *data, *decodeData)
		assert.Equal(t, seq, decodeMsg.SeqID())
		assert.Equal(t, msgType, msg.MsgType())
	})

	t.Run("optional: traceInfo,seq", func(t *testing.T) {
		trace := [TraceLength]byte{}
		for i := 0; i < TraceLength; i++ {
			trace[i] = byte(i) + 1
		}
		seq := uint32(math.MaxUint32)
		msg := NewMessage()
		defer msg.Recycle()
		msg.SetReq(data)
		msg.SetSeqID(seq)
		msg.SetTraceInfo(trace[:])
		msg.SetMsgType(msgType)

		lb := netpoll.NewLinkBuffer()
		assert.NoError(t, msg.EncodeToWriter(lb))
		decodeMsg := NewMessage()
		defer decodeMsg.Recycle()
		_, err := decodeMsg.DecodeFromReader(lb)
		assert.NoError(t, err)
		decodeData := new(testData)
		ml, err := decodeMsg.Payload().Next(decodeMsg.Payload().Len())
		assert.NoError(t, err)
		assert.NoError(t, decodeData.Unmarshal(ml))
		assert.Equal(t, *data, *decodeData)
		assert.Equal(t, trace, [TraceLength]byte(decodeMsg.TraceInfo()))
		assert.Equal(t, seq, decodeMsg.SeqID())
		assert.Equal(t, msgType, msg.MsgType())
	})

}

func BenchmarkCodeC(b *testing.B) {
	data := &testData{
		str: "testing",
	}
	msgType := byte(128)
	encode := func() *netpoll.AsyncLinkBuffer {
		msg := NewMessage()
		msg.SetReq(data)
		msg.SetMsgType(msgType)
		msg.SetSeqID(123456)

		lb := netpoll.NewAsyncLinkBuffer(msg.EncodedLen(), nil)
		msg.EncodeToWriter(lb)
		msg.Recycle()
		return lb
	}
	buf := encode()
	b.Run("Encode", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			encode().Recycle()
		}
	})
	b.Run("Decode", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msg := NewMessage()
			msg.DecodeFromReader(buf)
			msg.Recycle()
		}

	})

	b.Run("GetLen", func(b *testing.B) {
		msg := NewMessage()
		msg.SetReq(data)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msg.Len()
		}
	})

}

func BenchmarkBigLittleEnding(b *testing.B) {

	l := uint32(5120)
	b.Run("大端序", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sli := make([]byte, 4)
			binary.BigEndian.PutUint32(sli, l)
			binary.BigEndian.Uint32(sli)
		}
	})

	b.Run("小端序", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sli := make([]byte, 4)
			binary.LittleEndian.PutUint32(sli, l)
			binary.LittleEndian.Uint32(sli)
		}
	})

}
