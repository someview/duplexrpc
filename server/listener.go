package server

import (
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"net"
)

type MakeListener func(s *server, address string) (ln net.Listener, err error)

func tcpMakeListener(network string) MakeListener {
	return func(s *server, address string) (ln net.Listener, err error) {
		return netpoll.CreateListener(network, address)
	}
}
