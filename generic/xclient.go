package generic

import (
	"context"
	"errors"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git/dt/cmap"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"sync"
	"sync/atomic"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/selector"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
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

var idGenerator = protocol.NewSeqIdGenerator()

func genSeqID() uint32 {
	return idGenerator.Next()
}

type xClient struct {
	FailMode

	opt Option

	isShutdown   atomic.Bool
	mu           sync.RWMutex
	cachedClient *util.OrderedMap[string, RPCClient]
	instanceSet  util.Set[string]

	selector.Selector
	resolver.Resolver

	msgChan chan protocol.Message

	cbMap cmap.ConcurrentMap[uint32, RespCallBack]
	info  rpcinfo.MethodInfo

	GenericInfo
	msgH func(resp any, err error)
}

var _ XClient = (*xClient)(nil)

// Close implements XClient.
func (c *xClient) Close() error {
	close(c.msgChan)
	c.isShutdown.Store(true)
	return nil
}

func (c *xClient) Call(ctx context.Context, args any, cb RespCallBack, optFns ...CallOpt) {
	if c.isShutdown.Load() {
		cb(nil, ErrXClientShutdown)
		return
	}

	callOpt := ApplyCallOption(optFns...)
	if c.info.OneWay() {
		callOpt.oneway = true
	}
	end, err := c.selectClient(ctx, callOpt)
	if err != nil {
		cb(nil, err)
		callOpt.Recycle()
		return
	}

	req := protocol.NewMessageWithCtx(ctx)
	req.SetReq(args)
	req.SetMsgType(byte(c.info.Type()))
	seqID := genSeqID()
	req.SetSeqID(seqID)

	var sendCallBack netpoll.CallBack
	sendCallBack = func(err error) {

		if callOpt.oneway && err == nil {
			// onewayCall
			cb(nil, nil)
		} else if !callOpt.oneway && err == nil {
			// set callback, to wait server response
			c.cbMap.Set(seqID, cb)
		} else {
			// has error
			if callOpt.failedCount > c.opt.failureLimit {
				cb(nil, fmt.Errorf("failure limit :%w", err))
			} else {
				// record failedCount
				callOpt.failedCount++
				switch callOpt.failMode {
				// you must return after any retry!!!! if you don't, req.Recycle will cause panic
				case Failover:
					end, err = c.selectFailoverClient(end.Address())
					if err != nil {
						cb(nil, err)
					} else {
						end.AsyncSend(req, sendCallBack)
						return
					}
				case Failfast:
					cb(nil, err)
				case Failtry:
					end.AsyncSend(req, sendCallBack)
					return
				case Failbackup:
					cb(nil, fmt.Errorf("failBackup is not implemented"))
				}
			}

		}
		callOpt.Recycle()
		req.Recycle()
	}

	end.AsyncSend(req, sendCallBack)

}

func (c *xClient) recv() {
	// TODO use handler mode
	for {
		msg := <-c.msgChan
		if msg == nil {
			return
		}
		seqID := msg.SeqID()
		if seqID == 0 {
			// seqID 0 just for server message
			if c.msgH != nil {
				if c.info.Type() != msg.MethodName() {
					c.msgH(nil, errors.New("method not match"))
				} else {
					result := c.info
					resp := result.NewResult()
					err := protocol.DecodeProtobuf(msg.Payload(), resp)
					c.msgH(resp, err)
				}
			}
			msg.Recycle()
			continue
		}
		cb, ok := c.cbMap.Pop(seqID)
		if !ok {
			// throw this message if no callback
			msg.Recycle()
			continue
		}
		if c.info.Type() != msg.MethodName() {
			cb(nil, errors.New("method not match"))
			msg.Recycle()
			continue
		}
		result := c.info
		resp := result.NewResult()
		err := protocol.DecodeProtobuf(msg.Payload(), resp)
		cb(resp, err)
		msg.Recycle()
	}

}

func (c *xClient) selectClient(ctx context.Context, callOpt *CallOption) (RPCClient, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	var err error
	if !callOpt.disableSelector {
		callOpt.node, err = c.Selector.Select(ctx)
	}
	if err != nil {
		return nil, err
	}
	addr := callOpt.node.Address()
	end, ok := c.cachedClient.Get(addr)
	if !ok {
		return nil, fmt.Errorf("no avaliable endpoint with this url: %s", addr)
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
		var cli RPCClient
		if c.opt.muxConn > 1 {
			cli = NewMultiClient(c.opt.muxConn, c.opt)
		} else {
			cli = NewMuxClient(c.opt)
		}

		_ = cli.AsyncConnect("tcp", addr)
		newCacheClient.Set(addr, cli)
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

func NewClient(service string, info rpcinfo.MethodInfo, opts ...OptionFn) XClient {

	cli := &xClient{
		cachedClient: util.NewOrderedMap[string, RPCClient](),
		instanceSet:  util.NewSet[string](),
		msgChan:      make(chan protocol.Message, 2048),
		cbMap: cmap.New[uint32, RespCallBack](cmap.WithHashFunc[uint32, RespCallBack](func(k uint32) uint32 {
			return k
		})),
		info: info,
		opt:  defaultOpt,
	}
	cli.opt.serverMsgHandle = func(message protocol.Message) error {
		// TODO while full ,block it? or use linked list
		cli.msgChan <- message
		return nil
	}
	for _, opt := range opts {
		opt(&cli.opt)
	}
	cli.Selector = cli.opt.selector
	cli.Resolver = cli.opt.resolver
	cli.msgH = cli.opt.msgH
	go cli.Resolver.Start(cli.updateTransport)
	go cli.recv()
	return cli
}
