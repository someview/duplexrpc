package client

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"

	netpoll "github.com/cloudwego/netpoll"
)

var (
	ErrShutdown           = errors.New("client is shut down")
	ErrConnecting         = errors.New("client is connecting")
	ErrBreakerOpen        = errors.New("breaker is open")
	ErrClientNotConnected = errors.New("client is not connected")
)

type MuxClient struct {
	ctx context.Context

	network string
	address string

	transHandler remote.TransHandler

	// 单个rpcEndpint节点
	rpcEndpoint atomic.Value
	connecting  atomic.Int32
	// 表示当前muxClient是否已经关闭
	closed atomic.Bool

	opt option
}

func newMuxClient(ctx context.Context, network, address string, transHandler remote.TransHandler, opt option) *MuxClient {
	res := &MuxClient{
		ctx:          ctx,
		transHandler: transHandler,
		opt:          opt,
		network:      network,
		address:      address,
	}
	res.connectAsync()
	return res
}

// ---implements RPCClient---
// Address implements RPCClient.
func (c *MuxClient) Address() string {
	return c.address
}

func (c *MuxClient) IsClosed() bool {
	return c.closed.Load()
}

func (c *MuxClient) Close() (err error) {
	if c.IsClosed() {
		return nil
	}
	c.closed.Store(true)
	return nil
}

var errUnsupportedTransProtocol = errors.New("unsupported transport protocol")

func (c *MuxClient) connect(network, address string) (remote.Endpoint, error) {
	var conn netpoll.Connection
	var err error
	var ep remote.Endpoint
	if network == "tcp" {
		conn, err = newDirectConn(c, c.opt.dialer, network, address)
	} else {
		return nil, errUnsupportedTransProtocol
	}

	if err == nil && conn != nil {
		wr := bufwriter.NewBufWriter(
			c.opt.remoteOpt.WriterType,
			c.opt.remoteOpt.WriteBufferThreshold,
			c.opt.remoteOpt.WriteMaxMsgNum,
			c.opt.remoteOpt.WriteDelayTime,
			conn)

		ep, err = remote.NewEndpoint(c.ctx, c.transHandler.(remote.ConnectionHandler), conn, wr)
		if err != nil {
			return nil, err
		}

	}
	return ep, err
}

func (c *MuxClient) connectAsync() {
	if c.connecting.CompareAndSwap(0, 1) {
		go func() {
			defer func() {
				c.connecting.Store(0)
				if c.IsClosed() {
					c.rpcEndpoint.Load().(remote.Endpoint).Close()
				}
			}()
			ep, err := c.connect(c.network, c.address)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.rpcEndpoint.Store(ep)
		}()
	}
}

func newDirectConn(c *MuxClient, dialer netpoll.Dialer, network, address string) (netpoll.Connection, error) {
	var conn netpoll.Connection
	//var tlsConn *tls.Conn
	var err error

	if c == nil {
		err = fmt.Errorf("nil client")
		return nil, err
	}

	if c.opt.tLSConfig != nil {
		// TODO 等待netpoll实现SSL/TLS
		return nil, fmt.Errorf("UnSupported TLS")
		//dialer := &net.Dialer{
		//	Timeout: c.opt.connectTimeout,
		//}
		//tlsConn, err = tls.DialWithDialer(dialer, network, address, c.opt.tLSConfig)
		//// or conn:= tls.Client(netConn, &config)
		//conn = net.Conn(tlsConn)
	} else {
		//conn, err = net.DialTimeout(network, address, c.opt.connectTimeout)
		// use netpoll
		conn, err = dialer.DialConnection(network, address, c.opt.connectTimeout)
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

var ErrClientShutdown = errors.New("client is shut down")

func (c *MuxClient) CallOneway(ctx context.Context, args any, ri rpcinfo.RPCInfo) (err error) {
	if c.IsClosed() {
		return ErrClientShutdown
	}

	return c.callCommon(ctx, args, nil, ri)
}

func (c *MuxClient) Call(ctx context.Context, args any, result any, ri rpcinfo.RPCInfo) (err error) {
	if c.IsClosed() {
		return ErrClientShutdown
	}
	return c.callCommon(ctx, args, result, ri)
}

func (c *MuxClient) callCommon(ctx context.Context, args any, result any, ri rpcinfo.RPCInfo) (err error) {

	end, ok := c.rpcEndpoint.Load().(remote.Endpoint)
	if !ok || !end.IsActive() {
		c.connectAsync()
		return fmt.Errorf("endpoint is not avalibale")
	}

	sendMsg := remote.NewMessage(ri, c.opt.remoteOpt.Codec)
	defer sendMsg.Recycle()
	sendMsg.SetData(args)

	err = c.transHandler.Write(ctx, end, sendMsg)
	if err != nil {
		return err
	}
	if ri.InteractionMode() == rpcinfo.Oneway {
		return nil
	}

	recvMsg := remote.NewMessage(ri, c.opt.remoteOpt.Codec)
	defer recvMsg.Recycle()
	recvMsg.SetData(result)
	return c.transHandler.Read(ctx, end, recvMsg)
}
