package client

import (
	"context"
	"crypto/tls"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/breaker"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/writer"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
	"time"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
)

// Option contains all options for creating clients.
type Option struct {
	// it is the failure limit for sending requests when using the Call function.
	// whether it takes effect depends on the FailMode
	failureLimit int

	// writeBufferThreshold sets the write quota of the write buffer.
	writeBufferThreshold int

	// tLSConfig for tcp and quic
	tLSConfig *tls.Config

	// connectTimeout sets timeout for dialing
	connectTimeout time.Duration
	// idleTimeout sets max idle time for underlying net.Conns
	idleTimeout time.Duration

	// tcpKeepAlive, if it is zero we don't set keepalive
	tcpKeepAlivePeriod time.Duration

	selector selector.Selector

	resolver resolver.Resolver

	// before invoke
	sendBeforeMWs []middleware.MiddleWare

	// the writer when sending request
	writer writer.Writer

	// writeDelayTime sets the delay time for writer.BatchWriter.
	writeDelayTime time.Duration
	// writeMaxMsgNum sets the max number of messages that can be written in one batch for writer.BatchWriter.
	writeMaxMsgNum int

	muxConn         int
	dialer          netpoll.Dialer
	serverMsgHandle func(protocol.Message) error

	msgH func(resp any, err error)
}

type OptionFn func(opt *Option)

func WithResolver(r resolver.Resolver) OptionFn {
	return func(opt *Option) {
		opt.resolver = r
	}
}

func WithSelector(s selector.Selector) OptionFn {
	return func(opt *Option) {
		opt.selector = s
	}
}

func WithFailureLimit(limit int) OptionFn {
	return func(opt *Option) {
		if limit < 0 {
			return
		}
		opt.failureLimit = limit
	}
}

// WithConnectTimeout sets the timeout for dialing.
func WithConnectTimeout(timeout time.Duration) OptionFn {
	return func(opt *Option) {
		opt.connectTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) OptionFn {
	return func(opt *Option) {
		opt.idleTimeout = timeout
	}
}

// WithSendBeforeMiddleWare adds a middleware to the client.
func WithSendBeforeMiddleWare(m middleware.MiddleWare) OptionFn {
	return func(opt *Option) {
		opt.sendBeforeMWs = append(opt.sendBeforeMWs, m)
	}
}

// WithBreaker adds a breaker middleware to the client.
func WithBreaker(getter func() breaker.Breaker) middleware.MiddleWare {
	return func(next middleware.Next) middleware.Next {
		b := getter()
		return func(ctx context.Context, req protocol.Message, cb netpoll.CallBack) {
			if b.Allow() {
				next(ctx, req, func(err error) {
					if err != nil {
						b.Fail()
					}
				})
			} else {
				cb(ErrBreakerOpen)
			}
		}
	}
}

// WithMsgHandler sets the server-msg handler of the client.
// it will be called when a response is received such seqID=0.
func WithMsgHandler(h func(resp any, err error)) OptionFn {
	return func(opt *Option) {
		opt.msgH = h
	}
}

// WithReadBufferThreshold sets the threshold of the read buffer.
// if size<=0, the threshold will be unlimited.
func WithReadBufferThreshold(size int64) OptionFn {
	return func(opt *Option) {
		opt.dialer = netpoll.NewDialer(netpoll.WithReadBufferThreshold(size))
	}
}

// WithWriteBufferThreshold sets the write quota of the write buffer.
// default is 40960 equal 4K.
func WithWriteBufferThreshold(size int) OptionFn {
	return func(opt *Option) {
		opt.writeBufferThreshold = size
	}
}

// WithTLSConfig set tls config
func WithTLSConfig(conf *tls.Config) OptionFn {
	return func(opt *Option) {
		opt.tLSConfig = conf
	}
}

// WithMuxConn the multiple-connection on one server
func WithMuxConn(size int) OptionFn {
	return func(opt *Option) {
		opt.muxConn = size
	}
}

// WithWriter set the writer when sending request
func WithWriter(wr writer.Writer) OptionFn {
	return func(opt *Option) {
		opt.writer = wr
	}
}

// WithWriteDelayTime sets the delay time for writer.BatchWriter.
func WithWriteDelayTime(delay time.Duration) OptionFn {
	return func(opt *Option) {
		opt.writeDelayTime = delay
	}
}

// WithWriteMaxNum sets the max number of messages that can be written in one batch for writer.BatchWriter.
func WithWriteMaxNum(num int) OptionFn {
	return func(opt *Option) {
		opt.writeMaxMsgNum = num
	}
}

var defaultOpt = Option{
	failureLimit:         1,
	writeBufferThreshold: 40960,
	tLSConfig:            nil,
	connectTimeout:       0,
	idleTimeout:          0,
	tcpKeepAlivePeriod:   0,
	selector:             selector.NewRoundRobinSelectorBuilder().Build(),
	resolver:             resolver.NewLocalResolver("127.0.0.1:8080"),
	sendBeforeMWs:        make([]middleware.MiddleWare, 0),
	muxConn:              1,
	dialer:               netpoll.NewDialer(),
	serverMsgHandle: func(message protocol.Message) error {
		panic("can not recv message from server,serverMsgHandle not defined")
		return nil
	},
	writer:         writer.ShardQueue,
	writeDelayTime: 5 * time.Millisecond,
	writeMaxMsgNum: 100,
}
