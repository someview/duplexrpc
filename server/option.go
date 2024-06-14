package server

import (
	"context"
	"github.com/panjf2000/ants/v2"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/registrar"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/writer"
	"time"
)

// --- setting for muxextend ---

func init() {
	SetAntsRunner(-1)
	//SetRunner(func(task func()) error {
	//	go task()
	//	return nil
	//})
}

// Runner is a function that runs the task.
var runTask Runner

type Runner func(task func()) error

// SetRunner sets the runner for the server to run handle-task.
// By default, the runner is a goroutine.
// runner can resolve the problem of goroutine leak and stack expansion.
func SetRunner(r Runner) {
	runTask = r
	netpoll.SetRunner(func(ctx context.Context, f func()) {
		_ = runTask(func() {
			f()
		})
	})
}

// SetAntsRunner uses the ants pool to SetRunner.
// it equals to SetRunner(antsPool.Submit).
func SetAntsRunner(size int) {
	pool, _ := ants.NewPool(size)
	SetRunner(pool.Submit)
}

// SetNumLoops is used to set the number of pollers, generally do not need to actively set.
// By default, the number of pollers is equal to runtime.GOMAXPROCS(0)/20+1.
// If the number of cores in your service process is less than 20c, theoretically only one poller is needed.
// Otherwise you may need to adjust the number of pollers to achieve the best results.
// Experience recommends assigning a poller every 20c.
//
// You can only use SetNumLoops before any connection is created. An example usage:
//
//	func init() {
//	    netpoll.SetNumLoops(...)
//	}
func SetNumLoops(numLoops int) error {
	return netpoll.SetNumLoops(numLoops)
}

// --- Option for RPCServer ---

// Option is the option of the server.
type Option struct {
	// read and write timeout of the server.
	readTimeout  time.Duration
	writeTimeout time.Duration

	// tcp keep alive.
	keepAlivePeriod time.Duration

	// throttle
	readBufferThreshold  int
	writeBufferThreshold int

	// the writer when sending response
	writer writer.Writer

	// writeDelayTime sets the delay time for writer.BatchWriter.
	writeDelayTime time.Duration
	// writeMaxMsgNum sets the max number of messages that can be written in one batch for writer.BatchWriter.
	writeMaxMsgNum int

	// registrar
	reg registrar.Registrar
	// server metadata
	metadata    map[string]string
	serviceName string

	disableHandlerMode bool
}

type OptionFn func(opt *Option)

// WithServiceName sets the service name of the server.
func WithServiceName(service string) OptionFn {
	return func(opt *Option) {
		opt.serviceName = service
	}
}

// WithRegistrar sets the registrar for the server.
func WithRegistrar(r registrar.Registrar) OptionFn {
	return func(opt *Option) {
		opt.reg = r
	}
}

// WithMetaData sets the metadata of the server.
func WithMetaData(key string, value string) OptionFn {
	return func(opt *Option) {
		opt.metadata[key] = value
	}
}

// WithReadTimeout sets the read timeout of the server.
func WithReadTimeout(timeout time.Duration) OptionFn {
	return func(opt *Option) {
		opt.readTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout of the server.
func WithWriteTimeout(timeout time.Duration) OptionFn {
	return func(opt *Option) {
		opt.writeTimeout = timeout
	}
}

// WithKeepAlivePeriod sets the keep alive period of the server.
func WithKeepAlivePeriod(timeout time.Duration) OptionFn {
	return func(opt *Option) {
		opt.keepAlivePeriod = timeout
	}
}

// WithReadBufferThreshold sets the threshold of the read buffer.
// if size<=0, the threshold will be unlimited.
func WithReadBufferThreshold(size int) OptionFn {
	return func(opt *Option) {
		opt.readBufferThreshold = size
	}
}

// WithWriteBufferThreshold sets the write quota of the write buffer.
// default is 40960.
func WithWriteBufferThreshold(size int) OptionFn {
	return func(opt *Option) {
		opt.writeBufferThreshold = size
	}
}

// WithWriter set the writer when sending response
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

// WithDisableHandlerMode the rpc request will handle on one goroutine.
//
// such like:
//
//	func(req){
//	 handle(req) // now
//	 go handle(req) // before
//	}
//
// it mean your rpc handle cannot block for too long, otherwise it will impact subsequent requests.
func WithDisableHandlerMode() OptionFn {
	return func(opt *Option) {
		opt.disableHandlerMode = true
	}
}

var defaultOpt = Option{
	readTimeout:          0,
	writeTimeout:         0,
	keepAlivePeriod:      3 * time.Minute,
	readBufferThreshold:  -1,
	writeBufferThreshold: 40960,
	serviceName:          "",
	writer:               writer.ShardQueue,
	reg:                  nil,
	metadata:             make(map[string]string),
	writeDelayTime:       5 * time.Millisecond,
	writeMaxMsgNum:       100,
}
