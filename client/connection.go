package client

import (
	"context"
	"errors"
	"fmt"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"net"
)

type ConnFactoryFn func(c *MuxClient, network, address string) (net.Conn, error)

var errUnsupportedTransProtocol = errors.New("unsupported transport protocol")

func (c *MuxClient) connect(network, address string) error {
	var conn netpoll.Connection
	var err error
	if network == "tcp" {
		conn, err = newDirectConn(c, c.option.dialer, network, address)
	} else {
		return errUnsupportedTransProtocol
	}

	if err == nil && conn != nil {
		if tc, ok := conn.(*netpoll.TCPConnection); ok && c.option.tcpKeepAlivePeriod > 0 {
			_ = tc.SetKeepAlive(int(c.option.tcpKeepAlivePeriod.Seconds()))
		}

		c.conn = conn
		err = c.conn.SetIdleTimeout(c.option.idleTimeout)
		err = c.conn.SetOnRequest(c.onRequest)
		err = c.conn.AddCloseCallback(c.onConnClose)
		//// start reading and writing since connected
		//if c.option.Heartbeat && c.option.HeartbeatInterval > 0 {
		//	// todo heart message
		//	//go c.heartbeat()
		//}
	}
	return err
}

func (c *MuxClient) AsyncConnect(network, address string) error {

	if c.isShutDown.Load() {
		return ErrShutdown
	}
	c.mu.Lock()
	c.network = network
	c.address = address
	go func() {
		err := c.connect(network, address)
		if err != nil {
			fmt.Println("failed to connect server:", err)
			c.mu.Unlock()
			return
		}
		c.initAsyncWriter()
		c.mu.Unlock()
		if c.isShutDown.Load() {
			_ = c.ShutDown()
		}
	}()
	return nil
}

func newDirectConn(c *MuxClient, dialer netpoll.Dialer, network, address string) (netpoll.Connection, error) {
	var conn netpoll.Connection
	//var tlsConn *tls.Conn
	var err error

	if c == nil {
		err = fmt.Errorf("nil client")
		return nil, err
	}

	if c.option.tLSConfig != nil {
		// TODO 等待netpoll实现SSL/TLS
		return nil, fmt.Errorf("UnSupported TLS")
		//dialer := &net.Dialer{
		//	Timeout: c.option.connectTimeout,
		//}
		//tlsConn, err = tls.DialWithDialer(dialer, network, address, c.option.tLSConfig)
		//// or conn:= tls.Client(netConn, &config)
		//conn = net.Conn(tlsConn)
	} else {
		//conn, err = net.DialTimeout(network, address, c.option.connectTimeout)
		// use netpoll
		conn, err = dialer.DialConnection(network, address, c.option.connectTimeout)
	}

	if err != nil {
		fmt.Println("failed to dial server:", err)
		return nil, err
	}

	return conn, nil
}

func (c *MuxClient) onRequest(ctx context.Context, connection netpoll.Connection) error {

	msg, err := protocol.ParseHeader(connection.Reader(), nil)
	if err != nil {
		c.CloseWithReason(fmt.Errorf("rpcx: client protocol error: %s", err))
		return err
	}
	return c.option.serverMsgHandle(msg)
}

// 当连接主动或被动关闭时，都会调用到这个方法
func (c *MuxClient) onConnClose(connection netpoll.Connection) error {
	c.mu.Lock()
	c.wr = nil
	c.conn = nil
	c.mu.Unlock()
	return nil
}

func (c *MuxClient) CloseWithReason(err error) {
	c.ctx = context.WithValue(c.ctx, "err", err)
	_ = c.conn.Close()
}
