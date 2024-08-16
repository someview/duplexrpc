package remote

import (
	"fmt"

	netpoll "github.com/cloudwego/netpoll"
)

// DecodeProtobuf decode protobuf
//
//	example:
//	DecodeProtobuf(Message.Payload(),val)
func DecodeProtobuf(r netpoll.Reader, value any) error {
	proto, ok := value.(SizeableMarshaller)
	if !ok {
		return fmt.Errorf("the value type %T is not a protocol.SizeableMarshaller\n", value)
	}
	py, err := r.Next(r.Len())
	if err != nil {
		return err
	}
	return proto.Unmarshal(py)
}
