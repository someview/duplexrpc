package client

import (
	"context"
	"errors"
	"rpc-oneway/selector"
	"sync"
)

var (
	// ErrXClientShutdown xclient is shutdown.
	ErrXClientShutdown = errors.New("xClient is shut down")
	// ErrXClientNoServer selector can't found one server.
	ErrXClientNoServer = errors.New("can not found any server")
	// ErrServerUnavailable selected server is unavailable.
	ErrServerUnavailable = errors.New("selected server is unavailable")
)

type Closeable interface {
	Close() error
}

// XClient 泛化接口, msgType用于peer之间传输的消息类型, 由应用本身协商即可
type XClient interface {
	Send(ctx context.Context, msgType byte, args any) error
	Recv(ctx context.Context, msgType byte, args any) error
	Closeable
}

type xClient struct {
	mu sync.RWMutex

	servicePath string
	isShutdown  bool

	// used for single connection
	stickyRPCClient RPCClient
	stickyK         string

	selector selector.Selector
	option   Option
}

var _ XClient = (*xClient)(nil)

// Close implements XClient.
func (*xClient) Close() error {
	panic("unimplemented")
}

// Recv implements XClient.
func (*xClient) Recv(ctx context.Context, msgType int32, args any) error {
	panic("unimplemented")
}

// Send implements XClient.
func (c *xClient) Send(ctx context.Context, msgType int32, args any) error {
	if c.isShutdown {
		return ErrXClientShutdown
	}
	// 暂时忽略掉加密
	// if c.auth != "" {
	// 	metadata := ctx.Value(share.ReqMetaDataKey)
	// 	if metadata == nil {
	// 		metadata = map[string]string{}
	// 		ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
	// 	}
	// 	m := metadata.(map[string]string)
	// 	m[share.AuthKey] = c.auth
	// }
	// 存在selectClient的过程
	return nil
}

// todo 这里的
func (c *xClient) selectClient(ctx context.Context, servicePath, serviceMethod string,
	args interface{}) (string, RPCClient, error) {
	// 这儿不应该会被改变
	c.mu.Lock()
	if c.option.Sticky && c.stickyRPCClient != nil {
		c.mu.Unlock()
		return c.stickyK, c.stickyRPCClient, nil
	}
	fn := c.selector.Select

	k := fn(ctx, servicePath, serviceMethod, args)
	c.mu.Unlock()
	// 暂时无可用的服务
	if k == "" {
		return "", nil, ErrXClientNoServer
	}

	client, err := c.getCachedClient(k, servicePath, serviceMethod, args)

	if c.option.Sticky && client != nil {
		c.stickyK = k
		c.stickyRPCClient = client
	}
	return k, client, err
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

	// if this client is broken
	// breaker, ok := c.breakers.Load(k)
	// if ok && !breaker.(Breaker).Ready() {
	// 	return nil, ErrBreakerOpen
	// }

	// c.mu.Lock()
	// client = c.findCachedClient(k, servicePath, serviceMethod)
	// if client != nil {
	// 	if !client.IsClosing() && !client.IsShutdown() {
	// 		c.mu.Unlock()
	// 		return client, nil
	// 	}
	// 	c.deleteCachedClient(client, k, servicePath, serviceMethod)
	// }

	// client = c.findCachedClient(k, servicePath, serviceMethod)
	// c.mu.Unlock()

	return client, nil
}
