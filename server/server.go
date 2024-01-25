package server

import (
	"bufio"
	"context"
	"errors"
	"github.com/smallnest/rpcx/log"
	"github.com/soheilhy/cmux"
	"io"
	"net"
	"rpc-oneway/protocol"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ErrServerClosed is returned by the Server's Serve, ListenAndServe after a call to Shutdown or Close.
var (
	ErrServerClosed  = errors.New("http: Server closed")
	ErrReqReachLimit = errors.New("request reached rate limit")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBufsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024

	// // WriteChanSize is used for response.
	// WriteChanSize = 1024 * 1024
)

type Handler func(ctx *ClientRequestContext, msg any) error

type Server struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	pool         WorkerPool
	ln           net.Listener

	mu sync.RWMutex

	activeConn map[net.Conn]struct{}
	options    map[string]interface{}
	router     map[byte]Handler

	doneChan   chan struct{}
	inShutdown int32
}

func NewServer(options ...OptionFn) *Server {
	s := &Server{}
	// 设置保活时间
	if s.options["TCPKeepAlivePeriod"] == nil {
		s.options["TCPKeepAlivePeriod"] = 3 * time.Minute
	}
	return s
}

func (s *Server) AddHandler(msgType byte, handler func(requestContext *ClientRequestContext, msg any) error) {
	s.router[msgType] = handler
}

// SendMessage 反向发送的时候需要知道向哪一个conn发送消息
func (s *Server) SendMessage(conn net.Conn, msgType byte, msg any) error {
	req := protocol.NewMessage()
	req.MsgType = msgType
	req.Data = msg
	allData, err := req.EncodeSlicePointer()
	if err != nil {
		return err
	}
	_, err = conn.Write(*allData)
	protocol.PutData(allData)
	return err
}

func (s *Server) Serve(network, address string) (err error) {
	ln, err := tcpMakeListener(network)(s, address)
	if err != nil {
		return err
	}
	return s.serveListener(ln)
}

// serveListener accepts incoming connections on the Listener ln,
// creating a new service goroutine for each.
// The service goroutines read requests and then call services to reply to them.
func (s *Server) serveListener(ln net.Listener) error {
	var tempDelay time.Duration

	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if s.isShutdown() {
				<-s.doneChan
				return ErrServerClosed
			}

			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Errorf("rpcx: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			if errors.Is(e, cmux.ErrListenerClosed) {
				return ErrServerClosed
			}
			return e
		}
		tempDelay = 0

		if tc, ok := conn.(*net.TCPConn); ok {
			period := s.options["TCPKeepAlivePeriod"]
			if period != nil {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(period.(time.Duration))
				tc.SetLinger(10)
			}
		}

		//conn, ok := s.Plugins.DoPostConnAccept(conn)
		//if !ok {
		//	conn.Close()
		//	continue
		//}

		s.mu.Lock()
		s.activeConn[conn] = struct{}{}
		s.mu.Unlock()

		//if share.Trace {
		//	log.Debugf("server accepted an conn: %v", conn.RemoteAddr().String())
		//}

		go s.serveConn(conn)
	}
}

func (s *Server) isShutdown() bool {
	return atomic.LoadInt32(&s.inShutdown) == 1
}

// 接收到1个连接后开启1个协程
func (s *Server) serveConn(conn net.Conn) {
	// 这儿应该没有必要判断是否shutDown
	if s.isShutdown() {
		s.CloseConn(conn)
	}
	// make sure all inflight requests are handled and all drained
	defer func() {
		if s.isShutdown() {
			<-s.doneChan
		}
		s.CloseConn(conn)
	}()
	r := bufio.NewReaderSize(conn, ReaderBufsize)
	for {
		if s.readTimeout != 0 {
			_ = conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		}

		ctx := &ClientRequestContext{}
		req, err := s.readRequest(ctx, r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
			} else if errors.Is(err, net.ErrClosed) {
				log.Infof("rpc-oneway: connection %s is closed", conn.RemoteAddr().String())
			}
			return
		}

		if s.pool != nil {
			s.pool.Submit(func() {
				s.processOneRequest(ctx, req)
			})
		} else {
			s.processOneRequest(ctx, req)
		}
	}
}

func (s *Server) CloseConn(conn net.Conn) {
	s.mu.Lock()
	delete(s.activeConn, conn)
	s.mu.Unlock()
	_ = conn.Close()
}

func (s *Server) readRequest(ctx context.Context, r io.Reader) (*protocol.Message, error) {
	req := protocol.NewMessage()
	err := req.Decode(r)
	if err != nil {
		return nil, err
	}
	if err == io.EOF {
		return req, err
	}
	return req, err
}

func (s *Server) processOneRequest(ctx *ClientRequestContext, req *protocol.Message) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, true)]

			log.Errorf("failed to handle the request: %v， stacks: %s", r, buf)
		}
	}()

	//// 心跳请求，直接处理返回
	//if req.IsHeartbeat() {
	//	s.Plugins.DoHeartbeatRequest(ctx, req)
	//	req.SetMessageType(protocol.Response)
	//	data := req.EncodeSlicePointer()
	//
	//	if s.writeTimeout != 0 {
	//		conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	//	}
	//	conn.Write(*data)
	//
	//	protocol.PutData(data)
	//
	//	return
	//}

	cancelFunc := parseServerTimeout(ctx, req)
	if cancelFunc != nil {
		defer cancelFunc()
	}

	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}

	// use handlers first
	if handler, ok := s.router[req.MsgType]; ok {
		err := handler(ctx, req.Data)
		if err != nil {
			log.Errorf("[handler internal error]: servicepath: %s, servicemethod, err: %v", req.MsgType, err)
		}
		return
	}
	return
}

// 聊天系统中暂时不需要同步的方法调用，因此这个方法暂时不考虑
func parseServerTimeout(ctx *ClientRequestContext, req *protocol.Message) context.CancelFunc {
	//if req == nil || req.Metadata == nil {
	//	return nil
	//}
	//
	//st := req.Metadata[share.ServerTimeout]
	//if st == "" {
	//	return nil
	//}
	//
	//timeout, err := strconv.ParseInt(st, 10, 64)
	//if err != nil {
	//	return nil
	//}
	//
	//newCtx, cancel := context.WithTimeout(ctx.Context, time.Duration(timeout)*time.Millisecond)
	//ctx.Context = newCtx
	return nil
}