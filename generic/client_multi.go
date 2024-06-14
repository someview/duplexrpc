package generic

import (
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"sync/atomic"
)

type multiClient struct {
	opt Option

	clients []*MuxClient
	cursor  atomic.Int32
	size    int32

	address string
	network string

	isShutDown atomic.Bool
}

func NewMultiClient(size int, opt Option) RPCClient {
	m := &multiClient{
		opt:     opt,
		size:    int32(size),
		clients: make([]*MuxClient, size),
	}
	for i := 0; i < len(m.clients); i++ {
		m.clients[i] = NewMuxClient(m.opt)
	}
	return m
}

func (m *multiClient) AsyncConnect(network, address string) error {
	m.address = address
	m.network = network
	if m.isShutDown.Load() {
		return ErrShutdown
	}
	for _, client := range m.clients {
		err := client.AsyncConnect(network, address)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *multiClient) AsyncSend(req protocol.Message, cb netpoll.CallBack) {
	c := m.clients[m.cursor.Add(1)%m.size]
	c.AsyncSend(req, cb)
}

func (m *multiClient) IsShutdown() bool {
	isShutdown := false
	for _, client := range m.clients {
		isShutdown = client.IsShutdown() || isShutdown
	}
	return isShutdown
}

func (m *multiClient) Address() string {
	return m.address
}

func (m *multiClient) Close() (err error) {
	m.isShutDown.Store(true)
	for _, client := range m.clients {
		err = client.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
