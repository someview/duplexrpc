package remote

import (
	"context"
	"errors"
	"fmt"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/bufwriter"
)

type endpoint struct {
	netpoll.Connection
	bufwriter.BatchWriter

	ctx     context.Context
	cancel  context.CancelCauseFunc
	handler ConnectionHandler
	addr    string

	cbManager CallBackManager
}

func NewEndpoint(ctx context.Context, h ConnectionHandler, conn netpoll.Connection, wr bufwriter.BufWriter) (Endpoint, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	ep := &endpoint{
		ctx:         ctx,
		cancel:      cancel,
		Connection:  conn,
		handler:     h,
		BatchWriter: wr,
		addr:        conn.RemoteAddr().String(),
		cbManager:   newCallBackManager(),
	}

	ctx, err := h.OnActive(ctx, ep)
	if err != nil {
		_ = ep.Close()
		return nil, err
	}
	ep.ctx = ctx

	err = conn.AddCloseCallback(func(connection netpoll.Connection) error {
		_ = ep.Close()
		h.OnInactive(ep.ctx, ep)
		return nil
	})
	if err != nil {
		_ = ep.Close()
		return nil, err
	}

	err = conn.SetOnRequest(ep.OnRequest)
	if err != nil {
		_ = ep.Close()
		return nil, err
	}
	return ep, nil
}

func (e *endpoint) CallBackManager() CallBackManager {
	return e.cbManager
}
func (e *endpoint) Address() string {
	return e.addr
}

func (e *endpoint) Context() context.Context {
	return e.ctx
}

func (e *endpoint) OnRequest(_ context.Context, _ netpoll.Connection) error {
	return e.handler.OnRead(e.ctx, e)
}

// GracefulShutdown 优雅关闭
func (e *endpoint) GracefulShutdown(ctx context.Context) error {
	//只关闭内部writer，即等待所有数据发送完毕
	return e.BatchWriter.Close()
}

func (e *endpoint) Close() error {
	e.cancel(fmt.Errorf("endpoint closed"))

	err1 := e.BatchWriter.Close()
	err2 := e.Connection.Close()
	return errors.Join(err1, err2)
}
