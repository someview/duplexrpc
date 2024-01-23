# rpc-oneway
a rpc framework, sendmsg oneway, support both client send and server send 


# 帧格式
- 握手帧
```
[Frame Type (1 byte)] [Handshake Payload Length (4 bytes)] [Handshake Payload (variable length)]
```
- 握手响应帧
```
[Frame Type (1 byte)] [Handshake Payload Length (4 bytes)] [Handshake Payload (variable length)]
```

- 关闭帧
```
[Frame Type (1 byte)] [Close Code (2 bytes)] [Close Reason Length (2 bytes)] [Close Reason (variable length)]
```

- ping帧
```
[Frame Type (1 byte)] [Ping Data Length (4 bytes)] [Ping Data (variable length)]
```


- pong帧
```
[Frame Type (1 byte)] [Ping Data Length (4 bytes)] [Ping Data (variable length)]
```

- close帧


- data帧
```
// msgType全局共享, 1个字节足够表达了256中消息
[Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
```

- windowupdate
```
[Frame Type (1 byte)] [Window Size Increase (4 bytes)]
```


## 代码设计与实现
// 参考grpc-go、kitex、netpoll、rpcx