package client

import (
	"context"
	"errors"
	"github.com/cloudwego/netpoll"
	"net"
	"rpc-oneway/protocol"
)

var (
	ErrShutdown         = errors.New("connection is shut down")
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

type muxClient struct {
	// 这是和连接相关的部分
	closing bool // whether the server is going to close this connection
	net.Conn
}

// OnRequest is called when the connection creates.
func (c *muxClient) OnRequest(ctx context.Context, connection netpoll.Connection) (err error) {
	//reader, err := connection.Reader().Slice(length)
	return err
}

func (c *muxClient) handleServerRequest() {

}

// Send 用户层调用接口
func (c *muxClient) Send(ctx context.Context, msgType byte, req any) error {
	//meta := ctx.Value(share.ReqMetaDataKey)
	//if meta != nil { // copy meta in context to meta in requests
	//	call.Metadata = meta.(map[string]string)
	//}
	// todo 采用writeQuota 判断一下
	if c.closing {
		return ErrShutdown
	}
}

func (c *muxClient) sendFrame(ctx context.Context, frame protocol.DataFrame) error {
	c.Conn.Write()
}
