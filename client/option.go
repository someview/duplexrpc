package client

import (
	"crypto/tls"
	"time"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
)

// option contains all options for creating clients.
type option struct {
	tLSConfig *tls.Config

	mws []middleware.CallMiddleware

	remoteOpt remote.RemoteOptions

	// 连接级别相关参数
	connectTimeout time.Duration

	selector selector.Selector
	resolver resolver.Resolver

	dialer netpoll.Dialer
}

type Option func(opt *option)

func WithResolver(r resolver.Resolver) Option {
	return func(opt *option) {
		opt.resolver = r
	}
}

// WithCallMiddleware Call中间件，server-call同样适用.
func WithCallMiddleware(mw middleware.CallMiddleware) Option {
	return func(opt *option) {
		opt.mws = append(opt.mws, mw)
	}
}

func WithSelector(s selector.Selector) Option {
	return func(opt *option) {
		opt.selector = s
	}
}

// WithConnectTimeout sets the timeout for dialing.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(opt *option) {
		opt.connectTimeout = timeout
	}
}

// WithReadBufferThreshold sets the threshold of the read buffer.
// if size<=0, the threshold will be unlimited.
func WithReadBufferThreshold(size int64) Option {
	return func(opt *option) {
		opt.dialer = netpoll.NewDialer(netpoll.WithReadBufferThreshold(size))
	}
}

func WithRemoteOption(remoteOpts ...remote.RemoteOption) Option {
	return func(opt *option) {
		for _, remoteOpt := range remoteOpts {
			remoteOpt(&opt.remoteOpt)
		}
	}
}

var defaultOpt = option{
	connectTimeout: 0,
	selector:       selector.NewRoundRobinSelectorBuilder().Build(),
	resolver:       nil,
	dialer:         netpoll.NewDialer(netpoll.WithReadBufferThreshold(512 * util.MiB)),
	remoteOpt:      remote.DefaultRemoteOption,
}
