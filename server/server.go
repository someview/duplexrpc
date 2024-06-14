package server

import (
	"context"
	"errors"
	"fmt"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git/muxextend"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/writer"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"

	"github.com/soheilhy/cmux"
)

// ErrServerClosed is returned by the server's Serve, ListenAndServe after a call to Shutdown or Close.
var (
	ErrServerClosed  = errors.New("http: server closed")
	ErrReqReachLimit = errors.New("request reached rate limit")
)

type server struct {
	opt Option

	funcPool sync.Pool

	// eventLoop
	lp atomic.Value

	// 由于msgType是一位byte，所以采用采用数组,索引性能高于map
	// [0]是注册进来的methodInfo, [1]是对应的handler
	router [256][2]any

	inShutdown atomic.Bool

	node *discovery.ServiceInstance
}

func RegisterMethodToServer[T rpcinfo.MethodInfo](s RPCServer, handler any) error {
	return s.RegisterMethod(*new(T), handler)
}

func NewServer(options ...OptionFn) RPCServer {
	s := &server{
		funcPool: sync.Pool{New: func() any {
			fs := make([]func(), 0, 64) // 64 is defined casually, no special meaning
			return &fs
		}},
		opt: defaultOpt,
	}
	for _, option := range options {
		option(&s.opt)
	}

	return s
}

func (s *server) Run(network, address string) error {
	defer func() {
		_ = s.Shutdown(context.TODO())
	}()
	if s.opt.reg != nil {
		node := discovery.NewNode(address, s.opt.metadata)
		ins := node.(*discovery.ServiceInstance)
		err := s.opt.reg.Register(context.TODO(), *ins)
		if err != nil {
			return err
		}
		s.node = ins
	}
	return s.Serve(network, address)
}

func (s *server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

func (s *server) RegisterMethod(mInfo rpcinfo.MethodInfo, handler any) error {

	methodName := mInfo.Type()
	s.router[methodName][0] = mInfo
	s.router[methodName][1] = handler
	return nil
}

func (s *server) sendResp(resp protocol.Message, cb netpoll.CallBack) {
	resp.To().AsyncWriter().AsyncWrite(resp, resp, cb)
}

func (s *server) Serve(network, address string) (err error) {
	ln, err := tcpMakeListener(network)(s, address)
	if err != nil {
		return err
	}
	//use netpoll
	return s.serveListenerWithEventLoop(ln)
}

func (s *server) onPrepare(conn netpoll.Connection) context.Context {
	var wr muxextend.AsyncWriter
	switch s.opt.writer {
	case writer.ShardQueue:
		wr = muxextend.NewShardQueue(s.opt.writeBufferThreshold, conn)
	case writer.BatchWriter:
		wr = muxextend.NewBatchWriter(s.opt.writeBufferThreshold, s.opt.writeMaxMsgNum, s.opt.writeDelayTime, conn)
	}

	return rpcinfo.NewRPCEndpoint(conn, wr)
}

func (s *server) onConnect(ctx context.Context, connection netpoll.Connection) context.Context {

	_ = connection.AddCloseCallback(s.onClose)
	return ctx
}

func (s *server) onRequest(ctx context.Context, connection netpoll.Connection) error {
	endpoint, ok := ctx.(rpcinfo.RPCEndpoint)
	if !ok {
		fmt.Println("ctx is not RPCEndpoint")
		return connection.Close()
	}
	//fs := *s.funcPool.Get().(*[]func())
	r := connection.Reader()
	// 当连接中没有数据时才会退出循环
	for total := r.Len(); total > 0; total = r.Len() {
		// 解析header
		rpcReq, err := protocol.ParseHeader(r, endpoint)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("client has closed this connection: %s\n", connection.RemoteAddr().String())
			} else if errors.Is(err, net.ErrClosed) {
				fmt.Printf("gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git: connection %s is closed\n", connection.RemoteAddr().String())
			} else {
				// debug test
				fmt.Println(err)
			}
			return connection.Close()
		}
		if s.opt.disableHandlerMode {
			s.task(rpcReq)
		} else {
			_ = runTask(func() {
				s.task(rpcReq)
			})
		}

		// 这里表示还有下一个数据包，可以先处理已提交的任务，避免饥饿
		//if total < length && len(fs) > 0 {
		//	s.batchGoTasks(fs)
		//	fs = *s.funcPool.Get().(*[]func())
		//}
		// 提交一个处理任务
		//fs = append(fs, func() {
		//	s.task(rpcReq)
		//})
	}

	//s.batchGoTasks(fs)

	return nil

}

// batchGoTasks centrally creates goroutines to execute tasks.
func (s *server) batchGoTasks(fs []func()) {
	for n := range fs {
		_ = runTask(fs[n])
	}
	fs = fs[:0]
	s.funcPool.Put(&fs)
}

// task contains a complete process about decoding request -> handling -> writing response
func (s *server) task(req protocol.Message) {
	name := req.MethodName()
	defer req.Recycle()
	methodGroup := s.router[name]
	if methodGroup[0] == nil || methodGroup[1] == nil {
		fmt.Printf("method %b no handler\n", name)
		return
	}
	methodInfo := methodGroup[0].(rpcinfo.MethodInfo)
	handler := methodGroup[1]
	args, result := methodInfo.NewArgs(), methodInfo.NewResult()
	err := protocol.DecodeProtobuf(req.Payload(), args)
	if err != nil {
		fmt.Printf("some error happened when decode protobuf: %v\n", err)
		return
	}
	err = methodInfo.Invoke(handler, req, args, result)
	if err != nil {
		fmt.Printf("some error happened when method handler: %v\n", err)
		return
	}
	recycleValue(args)
	if methodInfo.OneWay() {
		return
	}

	resp := protocol.NewMessage()
	resp.SetReq(result)
	resp.SetTo(req.From())
	resp.SetMsgType(req.MsgType())
	resp.SetSeqID(req.SeqID())
	s.sendResp(resp, func(err error) {
		resp.Recycle()
		if err != nil {
			fmt.Printf("some error happened when send response: %v\n", err)
		}
	})
	recycleValue(result)

}

func (s *server) onClose(connection netpoll.Connection) error {
	s.CloseConn(connection)
	return nil
}

func (s *server) serveListenerWithEventLoop(ln net.Listener) error {
	var tempDelay time.Duration
	request := s.onRequest

	eventLoop, err := netpoll.NewEventLoop(request,
		netpoll.WithReadBufferThreshold(int64(s.opt.readBufferThreshold)),
		netpoll.WithOnPrepare(s.onPrepare),
		netpoll.WithOnConnect(s.onConnect),
		netpoll.WithIdleTimeout(s.opt.keepAlivePeriod),
		netpoll.WithReadTimeout(s.opt.readTimeout),
		netpoll.WithWriteTimeout(s.opt.writeTimeout),
	)
	if err != nil {
		return err
	}
	s.lp.Store(eventLoop)
	for {
		tempDelay = 0
		//use netpoll
		e := eventLoop.Serve(ln)
		if e != nil {
			if s.isShutdown() {
				return ErrServerClosed
			}
			// 对于临时性错误，需要重试，例如断电、中断灯。参考grpc-go lis.Accept
			if ne, ok := e.(interface{ Temporary() bool }); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				fmt.Printf("rpcx: Accept error: %v; retrying in %v", ne, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			if errors.Is(e, cmux.ErrListenerClosed) {
				return ErrServerClosed
			}
			return e
		}
		return e
	}

}

func (s *server) isShutdown() bool {
	return s.inShutdown.Load()
}

func (s *server) CloseConn(conn netpoll.Connection) {
	_ = conn.Close()
}

func (s *server) Shutdown(ctx context.Context) error {
	if s.inShutdown.CompareAndSwap(false, true) {
		defer func() {
			s.node = nil
		}()
		if s.opt.reg != nil && s.node != nil {
			err := s.opt.reg.Deregister(ctx, *s.node)
			if err != nil {
				fmt.Println("rpcs: server deregister failed")
			}
		}
		fmt.Println("rpcs: server shut down")
		lp, ok := s.lp.Load().(netpoll.EventLoop)
		if ok {
			return lp.Shutdown(ctx)
		}
	}
	return nil
}
