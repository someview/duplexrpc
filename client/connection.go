package client

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"rpc-oneway/protocol"
)

type ConnFactoryFn func(c *MuxClient, network, address string) (net.Conn, error)

var errUnsupportedTransProtocol = errors.New("unsupported transport protocol")

func (c *MuxClient) Connect(network, address string) error {
	var conn net.Conn
	var err error
	if network == "tcp" {
		conn, err = newDirectConn(c, network, address)
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
		c.r = bufio.NewReaderSize(conn, ReaderBufSize)

		// start reading and writing since connected
		go c.readLoop()
		if c.option.Heartbeat && c.option.HeartbeatInterval > 0 {
			// todo heart message
			//go c.heartbeat()
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
		fmt.Println("failed to dial server:", err)
		return nil, err
	}

	return conn, nil
}

func (c *MuxClient) readLoop() {
	var err error

	for err == nil {
		// 需要排除掉非业务类型的消息
		msg := protocol.NewMessage()
		if c.option.IdleTimeout != 0 {
			_ = c.Conn.SetDeadline(time.Now().Add(c.option.IdleTimeout))
		}

		err = msg.Decode(c.r)
		if err != nil {
			break
		}
		if c.option.BidirectionalBlock {
			c.ServerMessageChan <- msg
		} else {
			select {
			case c.ServerMessageChan <- msg:
			default: // put to Pool if we use pool
			}
		}
	}

	c.CloseWithReason(err)
}

func (c *MuxClient) CloseWithReason(err error) {
	if c.closing {
		_ = c.Conn.Close()
		c.closing = true
		fmt.Println("rpcx: client protocol error:", err)
	}
}
