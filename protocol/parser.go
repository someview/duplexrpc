package protocol

import (
	"fmt"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

var defaultParser = NewParser()

type parser struct{}

func NewParser() Parser {
	p := &parser{}

	return p
}

// ParseHeader parse header.
func ParseHeader(r netpoll.Reader, from rpcinfo.RPCEndpoint) (Message, error) {
	return defaultParser.ParseHeader(r, from)
}

func (p *parser) ParseHeader(r netpoll.Reader, from rpcinfo.RPCEndpoint) (Message, error) {
	msg := NewMessage()
	msg.SetFrom(from)
	_, err := msg.DecodeFromReader(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

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
