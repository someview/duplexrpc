package server

import (
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/registrar"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
	"time"
)

// --- option for RPCServer ---

// option is the option of the server.
type option struct {
	mws                 []middleware.CallMiddleware
	reg                 registrar.Registrar
	remoteOpt           remote.RemoteOptions
	addr                string
	maxExitWaitTime     time.Duration
	readBufferThreshold int64
}

type Option func(opt *option)

// WithRegistrar sets the registrar for the server.
func WithRegistrar(r registrar.Registrar) Option {
	return func(opt *option) {
		opt.reg = r
	}
}

func WithRemoteOption(remoteOpts ...remote.RemoteOption) Option {
	return func(opt *option) {
		for _, remoteOpt := range remoteOpts {
			remoteOpt(&opt.remoteOpt)
		}
	}
}

func WithCallMiddleware(mw middleware.CallMiddleware) Option {
	return func(opt *option) {
		opt.mws = append(opt.mws, mw)
	}
}

// WithReadBufferThreshold sets the threshold of the read buffer.
// if size<=0, the threshold will be unlimited.
func WithReadBufferThreshold(size int64) Option {
	return func(opt *option) {
		opt.readBufferThreshold = size
	}
}

func WithAddress(addr string) Option {
	return func(opt *option) {
		opt.addr = addr
	}
}

var defaultOpt = option{
	reg:                 nil,
	addr:                ":8080",
	maxExitWaitTime:     time.Second * 30,
	readBufferThreshold: 512 * util.MiB,
	remoteOpt:           remote.DefaultRemoteOption,
}
