package client

import (
	"context"
	"errors"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
	"sync"
	"sync/atomic"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/util"
)

var (
	// ErrXClientShutdown xclient is shutdown.
	ErrXClientShutdown = errors.New("xClient is shut down")
	// ErrXClientNoServer selector can't found one server.
	ErrXClientNoServer = errors.New("can not found any server")
	// ErrServerUnavailable selected server is unavailable.
	ErrServerUnavailable = errors.New("selected server is unavailable")
	ErrorEndUnavailable  = errors.New("no avaliable endpoint with this url")
)

type xClient struct {
	ctx    context.Context
	cancel context.CancelCauseFunc

	serviceManager service.Manager

	isShutDown atomic.Bool
	opt        option

	targetServiceInfo *service.ServiceInfo

	mu           sync.RWMutex
	cachedClient *util.OrderedMap[string, RPCClient]
	instanceSet  util.Set[string]

	transHandler remote.TransHandler
	callFunc     middleware.GenericCall

	selector.Selector
	resolver.Resolver
}

var _ XClient = (*xClient)(nil)

// Close implements XClient.
func (c *xClient) Close() error {
	if c.isShutDown.Load() {
		return nil
	}
	c.isShutDown.Store(true)
	c.cancel(fmt.Errorf("xClient is shut down"))
	c.Resolver.Stop()
	return nil
}

func (c *xClient) invokeEndpointMethod() middleware.GenericCall {
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

// todo 接口化 将从resolver -> selector -> endpoint，添加node到endpoint的变更接口
// 参考kitex的change接口, 以及grpc的updateTransport 接口定义
func (c *xClient) updateTransport(instances []discovery.Node) {
	num := len(instances)
	if num == 0 {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.cachedClient = util.NewOrderedMap[string, RPCClient]()
		c.instanceSet = util.NewSet[string]()
		c.Selector.Apply(nil)
		return
	}
	c.Selector.Apply(instances)
	nowSet := util.NewSet[string]()
	for _, ins := range instances {
		nowSet.Insert(ins.Address())
	}
	// 求出相同节点
	added := nowSet.Difference(c.instanceSet)
	removed := c.instanceSet.Difference(nowSet)
	commonIns := nowSet.Intersection(c.instanceSet)

	// todo 这里减少锁的粒度, remove的过程可以异步进行
	newCacheClient := util.NewOrderedMap[string, RPCClient]()
	c.mu.Lock()
	commonIns.Range(func(url string) bool {
		cli, ok := c.cachedClient.Get(url)
		if ok {
			newCacheClient.Set(url, cli)
		} else {
			// 记录下来，当前域名解析获取到但是实现上连接不上
		}
		return true
	})
	added.Range(func(addr string) bool {
		newCacheClient.Set(addr, newMuxClient(c.ctx, "tcp", addr, c.transHandler, c.opt))
		return true
	})
	oldCacheClient := c.cachedClient
	c.cachedClient = newCacheClient
	c.mu.Unlock()

	// 删除旧endpoints的过程异步进行,无需加锁
	// todo go1.23版本以后使用使用lockfreeMap, 减小锁的粒度，保证并发安全
	// todo 参考deepflow-server, 以及github lockfreeMap相关实现
	removed.Range(func(url string) bool {
		oldCacheClient.DeleteWithFunc(url, func(url string, cli RPCClient) {
			cli.Close()
		})
		return true
	})

	c.instanceSet = nowSet
}

func NewXClient(serviceName string, targetInfo *service.ServiceInfo, opts ...Option) (XClient, error) {

	// todo check service format
	ctx, cancel := context.WithCancelCause(context.TODO())
	cli := &xClient{
		ctx:               ctx,
		cancel:            cancel,
		serviceManager:    service.Manager{},
		targetServiceInfo: targetInfo,
		cachedClient:      util.NewOrderedMap[string, RPCClient](),
		instanceSet:       util.NewSet[string](),
		opt:               defaultOpt,
		Selector:          selector.NewRoundRobinSelectorBuilder().Build(),
		Resolver:          resolver.NewLocalResolver(serviceName),
	}

	for _, opt := range opts {
		opt(&cli.opt)
	}
	if cli.opt.selector != nil {
		cli.Selector = cli.opt.selector
	}
	if cli.opt.resolver != nil {
		cli.Resolver = cli.opt.resolver
	}

	transHandler, err := remote.NewDefaultMuxTransHandler(
		middleware.Chain(cli.opt.mws...)(cli.invokeEndpointMethod()),
		cli.serviceManager, &cli.opt.remoteOpt,
	)
	if err != nil {
		return nil, err
	}
	cli.transHandler = transHandler
	cli.callFunc = middleware.Chain(cli.opt.mws...)(cli.call)

	go cli.Resolver.Start(cli.updateTransport)
	return cli, nil
}

// NewDuplexClient serviceImpl 必须实现所有的serverCall方法
func NewDuplexClient(serviceName string, targetInfo, serviceInfo *service.ServiceInfo, serviceImpl any, opts ...Option) (XClient, error) {

	ctx, cancel := context.WithCancelCause(context.TODO())
	cli := &xClient{
		ctx:               ctx,
		cancel:            cancel,
		serviceManager:    service.Manager{},
		targetServiceInfo: targetInfo,
		cachedClient:      util.NewOrderedMap[string, RPCClient](),
		instanceSet:       util.NewSet[string](),
		opt:               defaultOpt,
		Selector:          selector.NewRoundRobinSelectorBuilder().Build(),
		Resolver:          resolver.NewLocalResolver(serviceName),
	}

	for _, opt := range opts {
		opt(&cli.opt)
	}
	if cli.opt.selector != nil {
		cli.Selector = cli.opt.selector
	}
	if cli.opt.resolver != nil {
		cli.Resolver = cli.opt.resolver
	}
	svc := service.NewService(serviceInfo, serviceImpl)
	err := cli.serviceManager.AddService(svc)
	if err != nil {
		return nil, err
	}

	transHandler, err := remote.NewDefaultMuxTransHandler(
		middleware.Chain(cli.opt.mws...)(cli.invokeEndpointMethod()),
		cli.serviceManager, &cli.opt.remoteOpt,
	)
	if err != nil {
		return nil, err
	}
	cli.transHandler = transHandler
	cli.callFunc = middleware.Chain(cli.opt.mws...)(cli.call)

	go cli.Resolver.Start(cli.updateTransport)
	return cli, nil
}
