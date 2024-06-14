package testdata

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func BenchmarkProtoMarshal(b *testing.B) {
	c := &ClientMessage{Header: &Header{TraceId: string(make([]byte, 1000))}}
	b.Run("protoMarshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			proto.Marshal(c)
		}
	})
	b.Run("MarshalTo", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Marshal()
		}
	})
}
