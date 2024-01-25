package server

import (
	"net"
)

type MakeListener func(s *Server, address string) (ln net.Listener, err error)

func tcpMakeListener(network string) MakeListener {
	return func(s *Server, address string) (ln net.Listener, err error) {
		return net.Listen(network, address)
	}
}
