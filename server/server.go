package server

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
)

type server struct {
	sync.Mutex
	ctx            context.Context
	cancel         context.CancelCauseFunc
	isRun          bool
	serviceManager service.Manager
	ln             netpoll.Listener
	loop           netpoll.EventLoop
	opt            option
	transHandler   remote.TransHandler
}

func (s *server) NewServerXClient(svcInfo *service.ServiceInfo) ServerXClient {
	return NewServerXClient(svcInfo, s.transHandler, &s.opt)
}

func NewServer(opts ...Option) Server {

	ctx, cancel := context.WithCancelCause(context.TODO())
	s := &server{
		ctx:            ctx,
		cancel:         cancel,
		serviceManager: service.Manager{},
		opt:            defaultOpt,
	}
	for _, opt := range opts {
		opt(&s.opt)
	}
	s.init()
	return s
}

func (s *server) init() {
	transHandler, err := remote.NewDefaultMuxTransHandler(middleware.Chain(s.opt.mws...)(s.invokeEndpointMethod()), s.serviceManager, &s.opt.remoteOpt)
	if err != nil {
		panic(err)
	}
	s.transHandler = transHandler
}

func (s *server) invokeEndpointMethod() middleware.GenericCall {
	return func(ctx context.Context, req any, resp any, ri rpcinfo.RPCInfo) error {
		iv := ri.Invocation()

		methodInfo := iv.MethodInfo()
		serviceImpl := iv.ServiceImpl()
		if methodInfo == nil || serviceImpl == nil {
			return uerror.NewInternalError(uerror.ErrServiceMethodNotFind, 0, fmt.Sprintf("service %s method %s not find", iv.ServiceName(), iv.MethodName()))
		}

		return methodInfo.Handler()(serviceImpl, ctx, req, resp)
	}
}

func (s *server) RegisterService(svcInfo *service.ServiceInfo, serviceImpl interface{}) error {
	s.Lock()
	defer s.Unlock()
	if s.isRun {
		panic("service cannot be registered while server is running")
	}
	if svcInfo == nil {
		panic("svcInfo is nil. please specify non-nil svcInfo")
	}
	if serviceImpl == nil || reflect.ValueOf(serviceImpl).IsNil() {
		panic("serviceImpl is nil. please specify non-nil serviceImpl")
	}
	if err := s.serviceManager.AddService(service.NewService(svcInfo, serviceImpl)); err != nil {
		panic(err)
	}
	return nil
}

func (s *server) GetServiceInfos() map[string]*service.ServiceInfo {
	return s.serviceManager.GetSvcInfoMap()
}

func (s *server) Run() (err error) {
	s.Lock()
	s.isRun = true
	s.Unlock()
	ln, err := netpoll.CreateListener("tcp", s.opt.addr)
	if err != nil {
		return err
	}
	loop, err := s.newEventLoop()
	s.ln, s.loop = ln, loop
	return loop.Serve(ln)
}

func (s *server) Stop() error {
	_ = s.ln.Close()
	if closer, ok := s.transHandler.(remote.GracefulShutdown); ok {
		ctx, cancel := context.WithTimeout(context.TODO(), s.opt.maxExitWaitTime)
		defer func() {
			cancel()
			s.cancel(fmt.Errorf("server shutdown"))
			_ = s.loop.Shutdown(context.TODO())
		}()
		return closer.GracefulShutdown(ctx)
	}
	return nil
}
