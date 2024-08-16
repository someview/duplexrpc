package remote

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
)

// NewDefaultMuxTransHandler to provide default impl of muxTransHandler, it can be reused in netpoll, shm-ipc, framework-sdk extensions
func NewDefaultMuxTransHandler(invokeFunc middleware.GenericCall, serviceManager service.Manager, opt *RemoteOptions) (TransHandler, error) {
	handler := &muxTransHandler{
		invokeFunc:     invokeFunc,
		serviceManager: serviceManager,
		opt:            opt,
	}
	return handler, nil
}

type muxTransHandler struct {
	opt            *RemoteOptions
	invokeFunc     func(ctx context.Context, req, resp any, info rpcinfo.RPCInfo) error
	endpointSet    sync.Map
	serviceManager service.Manager
	rpcTaskWg      sync.WaitGroup
}

// Write implements the ServerTransHandler interface.
func (t *muxTransHandler) Write(ctx context.Context, end Endpoint, sendMsg Message) (err error) {

	msgType := sendMsg.MsgType()
	if msgType == ReqType || msgType == PingType {
		cbManager := end.CallBackManager()
		seqID := sendMsg.RPCInfo().Invocation().SeqID()
		cbManager.Set(seqID, readerChanPool.Get())
		defer func() {
			if err != nil {
				c, ok := cbManager.LoadAndDelete(seqID)
				if ok {
					readerChanPool.Put(c)
				}
			}
		}()
	}

	lb := netpoll.NewSizedLinkBuffer(1 * util.KiB)
	err = sendMsg.Codec().Encode(ctx, lb, sendMsg)
	if err != nil {
		return err
	}
	err = lb.Flush()
	if err != nil {
		return err
	}

	return end.Add(ctx, lb)

}

// Read implements the ServerTransHandler interface.
func (t *muxTransHandler) Read(ctx context.Context, end Endpoint, recvMsg Message) (err error) {
	cbManager := end.CallBackManager()
	seqID := recvMsg.RPCInfo().Invocation().SeqID()
	ch, ok := cbManager.Load(seqID)
	if !ok {
		return fmt.Errorf("RecvMsg Before SendMsg,seqID:%d not found", seqID)
	}
	defer func() {
		cbManager.Delete(seqID)
		if err == nil {
			readerChanPool.Put(ch)
		}
	}()
	select {
	case <-end.Context().Done():
		err = end.Context().Err()
	case <-ctx.Done():
		err = ctx.Err()
	case reader := <-ch:
		err = recvMsg.Codec().Decode(ctx, reader, recvMsg)
	}

	if err != nil {
		return err
	}

	if recvMsg.MsgType() == ErrorResponseType {
		ue := &uerror.BasicError{}
		err = DecodeProtobuf(recvMsg.Payload(), ue)
		if err != nil {
			return err
		}
		switch ue.ErrorCategory() {
		case uerror.ErrBiz:
			return ue
		default:
			return t.handleInternalErrorResponse(ctx, end, uerror.InternalError{BasicError: *ue})
		}

	}

	return DecodeProtobuf(recvMsg.Payload(), recvMsg.Data())

}

// OnRead implements the ServerTransHandler interface.
// The connection should be closed after returning error.
func (t *muxTransHandler) OnRead(ctx context.Context, conn netpoll.Connection) error {
	end := conn.(Endpoint)
	reader := end.Reader()

	for reader.Len() > 0 {
		// 协议级别的错误，直接关连接
		msgType, length, seqID, err := parseHeader(reader)
		if err != nil {
			t.handleProtocolError(ctx, end, err)
			break
		}
		msgReader, err := util.SliceBuf(length, end)
		if err != nil {
			t.handleProtocolError(ctx, end, err)
			break
		}

		switch msgType {
		case ReqType, OnewayType:
			// 处理RPC请求
			t.onRPCRequest(ctx, end, msgReader)
		case ResponseType, ErrorResponseType, PongType:
			// 处理RPC响应
			t.onRPCResponse(ctx, seqID, end, msgReader)
		case CloseType:
			// 处理关闭信号
			err = end.(*endpoint).GracefulShutdown(ctx)
		case PingType:
			// 处理Ping消息，即返回一个pong消息
			pongMsg := NewPongMessage(t.opt.Codec)
			err = t.Write(ctx, end, pongMsg)
			pongMsg.Recycle()
			_ = util.PutBufFromSlice(msgReader)
		}

		if err != nil {
			t.OnError(ctx, end, err)
		}
	}

	return nil
}

// 代表协议解析失败或各种原因收到了不正确，不可能发生的请求/响应
func (t *muxTransHandler) handleProtocolError(ctx context.Context, end Endpoint, err error) {
	_ = end.Close()
	t.OnError(ctx, end, err)
}

func (t *muxTransHandler) handleInternalErrorResponse(ctx context.Context, end Endpoint, iErr uerror.InternalError) error {
	// TODO 处理内部错误再向上抛出
	return nil
}

// 处理RPC请求
func (t *muxTransHandler) onRPCRequest(ctx context.Context, end Endpoint, reader netpoll.Reader) {

	var recvMsg, sendMsg Message
	iv := rpcinfo.NewEmptyInvocation()
	ri := rpcinfo.NewEmptyRPCInfo(iv)

	recvMsg = NewMessage(ri, t.opt.Codec)
	err := t.opt.Codec.Decode(ctx, reader, recvMsg)

	if err != nil {
		slog.Debug("trans: decode failed, remote=%s, error=%s", end.RemoteAddr(), err.Error())
		t.handleProtocolError(ctx, end, err)
		return
	}

	svc, ok := t.serviceManager.GetService(iv.ServiceName())
	if !ok {
		t.handleProtocolError(ctx, end, fmt.Errorf("service not found, serviceName=%s", iv.ServiceName()))
		return
	}

	methodInfo, serviceImpl := svc.GetMethodInfoAndSvcImpl(iv.MethodName())
	iv.SetMethodInfo(methodInfo)
	iv.SetServiceImpl(serviceImpl)

	if !methodInfo.OneWay() {
		sendMsg = NewRespMessage(recvMsg.RPCInfo(), t.opt.Codec)
	}

	t.rpcTaskWg.Add(1)
	task := func() {
		t.procMessage(ctx, end, recvMsg, sendMsg)
		t.rpcTaskDone(ctx, recvMsg, sendMsg, ri)
		t.rpcTaskWg.Done()
	}

	isParallel := t.opt.ParallelDecider(iv.ServiceName(), iv.MethodName())

	if isParallel {
		netpoll.Go(task)
	} else {
		task()
	}
}

// 处理RPC响应
func (t *muxTransHandler) onRPCResponse(ctx context.Context, seqID uint32, end Endpoint, reader netpoll.Reader) {
	var err error
	defer func() {
		if err != nil {
			slog.Debug("trans: onRPCResponse failed, remote=%s, error=%s", end.RemoteAddr(), err.Error())
			t.handleProtocolError(ctx, end, err)
		}
	}()
	cbManager := end.CallBackManager()
	ch, ok := cbManager.Load(seqID)
	if !ok {
		err = fmt.Errorf("OnRead Before SendMsg,seqID:%d not found", seqID)
		return
	}
	select {
	case ch <- reader:
	case <-end.Context().Done():
		err = end.Context().Err()
	case <-ctx.Done():
		err = ctx.Err()
	}
	return

}

func (t *muxTransHandler) rpcTaskDone(ctx context.Context, recvMsg, sendMsg Message, ri rpcinfo.RPCInfo) {
	var req, res any
	req = recvMsg.Data()
	info := ri.Invocation().MethodInfo()
	if !info.OneWay() {
		res = sendMsg.Data()
		defer func() {
			info.ResultFactory().Recycle(res)
			sendMsg.Recycle()
		}()
	}

	t.opt.onRPCDone(ctx, req, res, ri)

	info.ArgsFactory().Recycle(req)
	ri.Recycle()
	recvMsg.Recycle()

}

// this method only called after onMessageCall
func (t *muxTransHandler) writeErrorReplyIfNeeded(ctx context.Context, end Endpoint, recvMsg Message, err error,
) {
	if recvMsg.MsgType() == OnewayType {
		return
	}
	ri := recvMsg.RPCInfo()
	if ri == nil {
		return
	}
	errMsg := NewErrorMessage(ri, err, t.opt.Codec)
	writeErr := t.Write(ctx, end, errMsg)
	if writeErr != nil {
		slog.DebugContext(ctx, "trans: Write error reply failed, remote=%s, error=%s", end.RemoteAddr(), writeErr.Error())
		_ = end.Close()
		return
	}
}
func (t *muxTransHandler) procMessage(ctx context.Context, end Endpoint, recvMsg, sendMsg Message) {

	defer func() {
		if panicErr := recover(); panicErr != nil {
			stack := string(debug.Stack())
			slog.Error("panic error: %v, stack: %v", panicErr, stack)
			_ = end.Close()
		}
	}()
	var err error
	if ctx, err = t.OnMessage(ctx, recvMsg, sendMsg); err != nil {
		t.OnError(ctx, end, err)
		t.writeErrorReplyIfNeeded(ctx, end, recvMsg, err)
		return
	}

	if recvMsg.MsgType() == OnewayType {
		return
	}

	if err = t.Write(ctx, end, sendMsg); err != nil {
		slog.DebugContext(ctx, "trans: write error reply failed, remote=%s, error=%s", end.RemoteAddr(), err.Error())
		_ = end.Close()
	}
}

// OnMessage implements the ServerTransHandler interface.
// msg is the decoded instance, such as Arg and Result.
func (t *muxTransHandler) OnMessage(ctx context.Context, recvMsg, sendMsg Message) (context.Context, error) {
	var args, result any
	ri := recvMsg.RPCInfo()
	iv := ri.Invocation()

	info := iv.MethodInfo()
	args = info.ArgsFactory().New()
	recvMsg.SetData(args)
	if recvMsg.MsgType() != OnewayType {
		result = info.ResultFactory().New()
		sendMsg.SetData(result)
	}
	err := DecodeProtobuf(recvMsg.Payload(), args)
	if err != nil {
		return ctx, err
	}

	err = t.invokeFunc(ctx, args, result, ri)
	return ctx, err
}

// OnActive implements the ServerTransHandler interface.
func (t *muxTransHandler) OnActive(ctx context.Context, conn netpoll.Connection) (context.Context, error) {
	end := conn.(Endpoint)
	_ = conn.SetIdleTimeout(t.opt.MaxConnectionIdleTime)
	_ = conn.SetReadTimeout(t.opt.ReadWriteTimeout)
	_ = conn.SetWriteTimeout(t.opt.ReadWriteTimeout)
	t.endpointSet.Store(end, struct{}{})
	return t.opt.onActive(ctx, end)
}

// OnInactive implements the ServerTransHandler interface.
func (t *muxTransHandler) OnInactive(ctx context.Context, conn netpoll.Connection) {
	end := conn.(Endpoint)
	t.endpointSet.Delete(end)
	t.opt.onInactive(ctx, end)
}

// OnError implements the ServerTransHandler interface.
func (t *muxTransHandler) OnError(ctx context.Context, conn netpoll.Connection, err error) {
	end := conn.(Endpoint)
	t.opt.onError(ctx, end, err)
}

func (t *muxTransHandler) GracefulShutdown(ctx context.Context) error {
	closeMsg := NewCloseMessage(t.opt.Codec)
	done := make(chan struct{})
	netpoll.Go(func() {
		t.endpointSet.Range(func(k any, _ any) (ok bool) {
			ok = true
			ep := k.(Endpoint)
			if !ep.IsActive() {
				return
			}
			err := t.Write(ctx, ep, closeMsg)
			if err != nil {
				// TODO log connection closing error
			}
			return
		})
		t.rpcTaskWg.Wait()

		t.endpointSet.Range(func(k any, _ any) bool {
			ep := k.(Endpoint)
			if ep.IsActive() {
				shut, ok := ep.(GracefulShutdown)
				if ok {
					_ = shut.GracefulShutdown(ctx)
				}
			}
			return true
		})
		time.Sleep(3 * time.Second)
		close(done)
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}

}
