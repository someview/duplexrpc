package codec

import "github.com/cloudwego/netpoll"

type MsgGetter func() any

type Codec interface {
	Encode(writer netpoll.Writer, msg any) error
	Decode(reader netpoll.Reader, getter MsgGetter) (msg any, err error)
}
