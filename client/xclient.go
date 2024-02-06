package client

import (
	"context"
	"errors"
	"sync"

	"rpc-oneway/pkg/resolver"
	"rpc-oneway/pkg/selector"
	"rpc-oneway/protocol"
	"rpc-oneway/util"

	"golang.org/x/sync/singleflight"
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
	isShutdown bool
	FailMode

	mu           sync.RWMutex
	cachedClient *util.OrderedMap[string, RPCClient]
	instaceSet   util.Set[string]

	slGroup singleflight.Group
	selector.Selector
	resolver.Resolver

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
	if c.isShutdown {
		return ErrXClientShutdown
	}
	end, err := c.selectClient(ctx)
	// 这里可以判断一下错误类型，以便是否继续重试
	if err == nil {
		err = end.Send(ctx, args)
	}
	if err == nil {
		return nil
	}
	if c.FailMode == Failover {
		end, err = c.selectFailoverClient(end.Address())
		if err != nil {
			return err
		} else {
			return end.Send(ctx, args)
		}
	}
	return nil
}

// // todo 这里的
func (c *xClient) selectClient(ctx context.Context) (RPCClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	instance, err := c.Selector.Select(ctx)
	if err != nil {
		return nil, err
	}
	url := instance.Address()
	end, ok := c.cachedClient.Get(url)
	if !ok {
		return nil, ErrorEndUnavailable
	}
	return end, nil
}

func (c *xClient) selectFailoverClient(url string) (RPCClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 这个属性也可以设置的instance本身去
	end, ok := c.cachedClient.Next(url)
	if !ok {
		return nil, ErrorEndUnavailable
	}
	return end, nil
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

func (x *xClient) updateTransport(instances []resolver.Node) {
	num := len(instances)
	if num == 0 {
		x.mu.Lock()
		defer x.mu.Unlock()
		x.cachedClient = util.NewOrderedMap[string, RPCClient]()
		x.instaceSet = util.NewSet[string]()
		x.Selector.Apply(nil)
		return
	}
	nowSet := util.NewSet[string]()
	for _, ins := range instances {
		nowSet.Insert(ins.Address())
	}
	// 求出相同节点
	added := nowSet.Difference(x.instaceSet)
	removed := x.instaceSet.Difference(nowSet)
	x.mu.Lock()
	added.Range(func(url string) bool {
		cli := NewMuxClient()
		x.cachedClient.Set(url, cli)
		return true
	})
	removed.Range(func(url string) bool {
		x.cachedClient.DeleteWithFunc(url, func(url string, cli RPCClient) {
			cli.Close()
		})
		return true
	})
	x.mu.Unlock()

}

// withoutpool
func NewXClient(service string, opt Option) XClient {
	cli := &xClient{}
	go func() {
		opt.Resolver.Start(cli.updateTransport)
	}()
	return cli
}
