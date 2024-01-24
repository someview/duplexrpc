package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/log"
	"net"
	"rpc-oneway/protocol"
	"time"
)

type ConnFactoryFn func(c *MuxClient, network, address string) (net.Conn, error)

var errUnsupportedTransProtocol = errors.New("unsupported transport protocol")

func (c *MuxClient) Connect(network, address string) error {
	var conn net.Conn
	var err error
	if network == "tcp" {

	} else {
		return errUnsupportedTransProtocol
	}

	if err == nil && conn != nil {
		if tc, ok := conn.(*net.TCPConn); ok && c.option.TCPKeepAlivePeriod > 0 {
			_ = tc.SetKeepAlive(true)
			_ = tc.SetKeepAlivePeriod(c.option.TCPKeepAlivePeriod)
		}

		if c.option.IdleTimeout != 0 {
			_ = conn.SetDeadline(time.Now().Add(c.option.IdleTimeout))
		}

		c.Conn = conn
		c.r = bufio.NewReaderSize(conn, ReaderBuffsize)

		// start reading and writing since connected
		go c.readLoop()
		if c.option.Heartbeat && c.option.HeartbeatInterval > 0 {
			go c.heartbeat()
		}
	}

	return err
}

func newDirectConn(c *MuxClient, network, address string) (net.Conn, error) {
	var conn net.Conn
	var tlsConn *tls.Conn
	var err error

	if c == nil {
		err = fmt.Errorf("nil client")
		return nil, err
	}

	if c.option.TLSConfig != nil {
		dialer := &net.Dialer{
			Timeout: c.option.ConnectTimeout,
		}
		tlsConn, err = tls.DialWithDialer(dialer, network, address, c.option.TLSConfig)
		// or conn:= tls.Client(netConn, &config)
		conn = net.Conn(tlsConn)
	} else {
		conn, err = net.DialTimeout(network, address, c.option.ConnectTimeout)
	}

	if err != nil {
		log.Warnf("failed to dial server: %v", err)
		return nil, err
	}

	return conn, nil
}

func (c *MuxClient) readLoop() {
	var err error

	for err == nil {
		msg := protocol.NewMessage()
		if c.option.IdleTimeout != 0 {
			_ = c.Conn.SetDeadline(time.Now().Add(c.option.IdleTimeout))
		}

		err = msg.Decode(c.r)
		if err != nil {
			break
		}
		c.recv
	}

	// todo 添加error类型的判断
}

// RecvMsg 当前暂时可以不用传递
func (c *MuxClient) RecvMsg(ctx context.Context, msgType int32, m any) error {
	msg := protocol.NewMessage()
	if c.option.IdleTimeout != 0 {
		_ = c.Conn.SetDeadline(time.Now().Add(c.option.IdleTimeout))
	}

	err := msg.Decode(c.r)
	if err != nil {
		return err
	}
	res, ok := m.(protocol.SizeableMarshaller)
	if !ok {
		return ErrUnsupportedCodec
	}
	_, err = res.Unmarshal(msg.Payload)
	return err

}
