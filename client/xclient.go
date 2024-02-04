package client

import (
	"context"
	"errors"
	"sync"

	"rpc-oneway/pkg/resolver"
	"rpc-oneway/protocol"

	"golang.org/x/sync/singleflight"
)

var (
	// ErrXClientShutdown xclient is shutdown.
	ErrXClientShutdown = errors.New("xClient is shut down")
	// ErrXClientNoServer selector can't found one server.
	ErrXClientNoServer = errors.New("can not found any server")
	// ErrServerUnavailable selected server is unavailable.
	ErrServerUnavailable = errors.New("selected server is unavailable")
)

type xClient struct {
	isShutdown bool

	mu           sync.RWMutex
	cachedClient map[string]RPCClient
	breakers     sync.Map
	slGroup      singleflight.Group

	// used for single connection
	stickyRPCClient RPCClient
	stickyK         string
	Option
}

var _ XClient = (*xClient)(nil)

// Close implements XClient.
func (x *xClient) Close() error {
	return nil
}

// Recv implements XClient.
func (x *xClient) ServerMessageHandler() func(protocol.Message) {
	return x.serverMsgHandler
}

// Send implements XClient.
func (c *xClient) Send(ctx context.Context, args any) error {
	// 	if c.isShutdown {
	// 		return ErrXClientShutdown
	// 	}
	// 	// 暂时忽略掉加密
	// 	// if c.auth != "" {
	// 	// 	metadata := ctx.Value(share.ReqMetaDataKey)
	// 	// 	if metadata == nil {
	// 	// 		metadata = map[string]string{}
	// 	// 		ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
	// 	// 	}
	// 	// 	m := metadata.(map[string]string)
	// 	// 	m[share.AuthKey] = c.auth
	// 	// }
	// 	// 存在selectClient的过程
	// 	return nil
	// }

	// // todo 这里的
	// func (c *xClient) selectClient(ctx context.Context) (RPCClient, error) {
	// 	c.mu.RLock()
	// 	defer c.mu.RUnlock()
	// 	instance, err := c.Selector.Select(ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	k := instance.Endpoint

	// 	breaker, ok := c.breakers.Load(k)
	// 	if ok && !breaker.(Breaker).Ready() {
	// 		return nil, ErrBreakerOpen
	// 	}

	// 	c.cachedClient
	// 	// fn := c.selector.Select

	// 	// k := fn(ctx, servicePath, serviceMethod, args)
	// 	k := ""
	// 	c.mu.Unlock()
	// 	// 暂时无可用的服务
	// 	if k == "" {
	// 		return "", nil, ErrXClientNoServer
	// 	}

	// 	client, err := c.getCachedClient(k, servicePath, serviceMethod, args)

	// 	if c.option.Sticky && client != nil {
	// 		c.stickyK = k
	// 		c.stickyRPCClient = client
	// 	}
	return nil
}

func (c *xClient) getCachedClient(k string, servicePath, serviceMethod string, args interface{}) (RPCClient, error) {
	// TODO: improve the lock
	var client RPCClient
	// var needCallPlugin bool
	defer func() {
		// if needCallPlugin {
		// 	c.Plugins.DoClientConnected(client.GetConn())
		// }
	}()

	if c.isShutdown {
		return nil, errors.New("this xclient is closed")
	}

	return client, nil
}

func (x *xClient) updateTransport(serviceInstances []resolver.ServiceInstance) {

}

// withoutpool
func NewXClient(service string, opt Option) XClient {
	cli := &xClient{}
	go func() {
		opt.Resolver.Start(cli.updateTransport)
	}()
	return cli
}
