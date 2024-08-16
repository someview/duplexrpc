package client

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"io"
)

// XClient XClient 定义了客户端接口，用于注册服务和方法处理程序，以及发起不同类型的RPC调用。
// XClient 所有的方法都不应该被外部使用者直接调用，应通过代码生成器生成的形式调用内部的方法, 以method为单位,组织服务之间的互相调用
type XClient interface {
	// CallOneway 发起一个单向RPC调用。
	// ctx 表示调用上下文，包含超时、取消等信息。
	// methodName 表示调用的方法名称。
	// req 表示请求数据。
	// 返回值表示调用过程中是否发生错误。
	// CallOneway 方法用于发起单向RPC调用。
	// 单向调用意味着客户端只负责发送请求，不等待或期望服务器的响应。
	// 参数:
	//   ctx - 调用上下文，包含调用的超时、取消信号等。
	//   methodName - 调用的方法名称。
	//   req - 请求数据，具体类型取决于调用的方法。
	// 返回值:
	//   error - 如果调用过程中出现错误，则返回错误信息。
	CallOneway(ctx context.Context, methodName string, req any) error

	// Call 发起一个双向RPC调用。
	// ctx 表示调用上下文，包含超时、取消等信息。
	// methodName 表示调用的方法名称。
	// req 表示请求数据。
	// resp 表示响应数据。
	// 返回值表示调用过程中是否发生错误。
	// Call 方法用于发起双向RPC调用。
	// 双向调用意味着客户端发送请求，并等待服务器的响应。
	// 参数:
	//   ctx - 调用上下文，包含调用的超时、取消信号等。
	//   methodName - 调用的方法名称。
	//   req - 请求数据，具体类型取决于调用的方法。
	//   resp - 响应数据，具体类型取决于调用的方法。
	// 返回值:
	//   error - 如果调用过程中出现错误，则返回错误信息。
	Call(ctx context.Context, methodName string, req any, resp any) (err error)

	// Broadcast 发起一个广播RPC调用。
	// ctx 表示调用上下文，包含超时、取消等信息。
	// methodName 表示调用的方法名称。
	// req 表示请求数据。
	// 返回值表示调用过程中是否发生错误。
	// Broadcast 方法用于向所有服务器实例发起广播RPC调用。
	// 广播调用意味着客户端发送请求给所有服务器，不等待或期望每个服务器的单独响应。
	// 参数:
	//   ctx - 调用上下文，包含调用的超时、取消信号等。
	//   methodName - 调用的方法名称。
	//   req - 请求数据，具体类型取决于调用的方法。
	// 返回值:
	//   error - 返回在获取到可用节点之前出现的错误信息
	Broadcast(ctx context.Context, methodName string, req any) error

	// Close 允许客户端作为可关闭资源，例如释放连接
	io.Closer
}

type InstanceManager interface {
	AddInstances(instances []discovery.Node)
	RemoveInstances(instances []discovery.Node)
	ClearInstances()
	GetInstance(node discovery.Node) RPCClient
}

type RPCClient interface {
	IsClosed() bool
	Address() string
	CallOneway(ctx context.Context, args any, ri rpcinfo.RPCInfo) error
	Call(ctx context.Context, args any, resp any, ri rpcinfo.RPCInfo) error
	io.Closer
}
