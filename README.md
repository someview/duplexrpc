# TL.GoRpc

为聊天系统oneway的rpc框架，通用属性和业务信息相分离，性能至少和grpc stream持平

```
{
  MsgLen:       4byte  // 包括数据包剩余部分的字节大小，包含 MsgLen 自身长度
  Version:      1byte  // 协议版本，为future扩展预留
  MsgType:      1byte  // onewayReq, normarlreq, normalres, normalerror, heart
  SeqId:        4byte
  ServiceLen:   1byte                        // generic service 信息
  ServicePath:  变长 (长度由 ServiceLen 决定)  // serviceName
  MethodLen:    1byte
  MethodPath:   变长 (长度由 MethodLen 决定)     
  FlagsLen:     2byte  // Flags 和 FlagOptions 的总长度  bit 0: traceInfo 24byte
  Flags:        变长   // Flags 和对应的 FlagOptions     
  MetadataLen:  2byte
  Metadata:     变长 (长度由 MetadataLen 决定)
  Payload:      变长 (长度由 MsgLen 减去之前所有字段的长度决定)
}
```

// msgType全局共享, 1个字节足够表达了256中消息
[Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 bytes)][flag(1 bytes)][optional][payload]
```
trace: [flag = 0x1] [optional 24byte]
seq: [flag = 0x2] [optional 4byte]
```


## 代码设计与实现
### 客户端call流程

```
1. call(ctx,req,res)error
2. initCallContext
3. 全局中间件
4. 重试策略
5. selectMuxClient
6. 实例级别熔断器、实例级别中间件, 根据level指定先后顺序
7. MuxClient.Call(ctx,req,res) error
8. IoErrorhandler(ctx,err)
```

### 客户端详解
XClient -> MuxClient -> Endpoint -> TransHandler -> EventLoop

#### XClient
- 接口定义
```
Call(ctx context.Context, req, res interface{}) error
```
- 功能说明
 包括双向call，既能用作客户端，也能用作服务端
- 组件说明
  全局中间件 -> 重试 -> selector -> muxClient

#### MuxClient
```
Call(ctx context.Context, req, res interface{}) error
```
endpoint + endpoint连接状态管理, 处理req,res到protocol.Message的转换

#### Endpoint OnRequest
```
Send(ctx,protocol.Message)error
Recv(ctx)(protocol.Message,error)
OnMessage(ctx CallContext, args, result protocol.Message) error
```

#### TransHandler
对接tcp框架和endpoint接口,提供transport到protocol.Message的转换
```
Write(ctx context.Context, conn Connection, msg protocol.Message)(ctx, error)
Read(ctx context.Context, conn Connection, msg protol.Message) (ctx,error)
OnConnect(ctx context.Context, conn Connection) (ctx, error) 
OnClose(ctx context.Context, conn Connection) 
OnError(ctx context.Context, conn Connection, err error)
OnMessage(ctx context.Context, args, result Message) (context.Context, error)
```






