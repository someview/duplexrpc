package client

import (
	"context"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/metadata"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

func (c *xClient) selectNode(ctx context.Context) (RPCClient, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	node, err := c.Selector.Select(ctx)
	if err != nil {
		return nil, err
	}
	addr := node.Address()
	end, ok := c.cachedClient.Get(addr)
	if !ok {
		return nil, fmt.Errorf("no avaliable endpoint with this url: %s", addr)
	}
	return end, nil
}

func (c *xClient) selectTargetNode(ctx context.Context, addr string) (RPCClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	end, ok := c.cachedClient.Get(addr)
	if !ok {
		return nil, fmt.Errorf("no avaliable ep with this url: %s", addr)
	}
	return end, nil
}

func (c *xClient) selectFailOverNode(url string) (RPCClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 这个属性也可以设置的instance本身去
	end, ok := c.cachedClient.Next(url)
	if !ok {
		return nil, ErrorEndUnavailable
	}
	return end, nil
}

func (c *xClient) selectAllNode(ctx context.Context) ([]RPCClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	length := c.cachedClient.Len()
	if length <= 0 {
		return nil, fmt.Errorf("no avalibale ep")
	}
	ends := make([]RPCClient, 0, length)
	c.cachedClient.Range(func(k string, v RPCClient) bool {
		ends = append(ends, v)
		return true
	})
	return ends, nil
}

func (c *xClient) CallOneway(ctx context.Context, methodName string, req any) error {
	// TODO 这里需要让Invocation携带method信息
	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, rpcinfo.NewInvocation(c.targetServiceInfo.ServiceName, methodName))
	defer func() {
		ri.Recycle()
		metadata.RecycleContext(ctx)
	}()

	return c.callFunc(ctx, req, nil, ri)
}

func (c *xClient) Call(ctx context.Context, methodName string, req any, resp any) (err error) {
	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, rpcinfo.NewInvocation(c.targetServiceInfo.ServiceName, methodName))
	defer func() {
		ri.Recycle()
		metadata.RecycleContext(ctx)
	}()

	return c.callFunc(ctx, req, resp, ri)
}

func (c *xClient) call(ctx context.Context, req any, resp any, ri rpcinfo.RPCInfo) (err error) {
	var end RPCClient
	ep, ok := metadata.ExtractEndpoint(ctx)
	if !ok {
		end, err = c.selectNode(ctx)
	} else {
		end, err = c.selectTargetNode(ctx, ep.Address())
	}
	if err != nil {
		return err
	}
	switch ri.InteractionMode() {
	case rpcinfo.Oneway:
		return end.CallOneway(ctx, req, ri)
	case rpcinfo.Unary:
		return end.Call(ctx, req, resp, ri)
	default:
		return fmt.Errorf("unsupported interaction mode")
	}

}

func (c *xClient) Broadcast(ctx context.Context, methodName string, req any) error {

	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, rpcinfo.NewInvocation(c.targetServiceInfo.ServiceName, methodName))
	defer ri.Recycle()

	ends, err := c.selectAllNode(ctx)
	if err != nil {
		return err
	}
	for _, end := range ends {
		_ = end.CallOneway(ctx, req, ri)
	}
	return nil
}
