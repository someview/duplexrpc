package server

import (
	"context"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
)

func (s *server) newEventLoop() (netpoll.EventLoop, error) {
	loop, err := netpoll.NewEventLoop(
		s.onRequest,
		netpoll.WithOnPrepare(s.onPrepare),
		netpoll.WithOnConnect(s.onConnect),
		netpoll.WithOnDisconnect(s.onDisconnect),
		netpoll.WithReadBufferThreshold(s.opt.readBufferThreshold),
	)
	return loop, err
}

func (s *server) onPrepare(conn netpoll.Connection) context.Context {
	return s.ctx
}

func (s *server) onConnect(ctx context.Context, conn netpoll.Connection) context.Context {
	wr := bufwriter.NewBufWriter(
		s.opt.remoteOpt.WriterType,
		s.opt.remoteOpt.WriteBufferThreshold,
		s.opt.remoteOpt.WriteMaxMsgNum,
		s.opt.remoteOpt.WriteDelayTime,
		conn)
	_, err := remote.NewEndpoint(ctx, s.transHandler.(remote.ConnectionHandler), conn, wr)
	if err != nil {
		_ = conn.Close()
	}
	return ctx
}

func (s *server) onRequest(ctx context.Context, conn netpoll.Connection) error {
	return nil
}

func (s *server) onDisconnect(ctx context.Context, conn netpoll.Connection) {

}
