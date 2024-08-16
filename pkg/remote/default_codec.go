package remote

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	netpoll "github.com/cloudwego/netpoll"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
)

var (
	ErrUnsupportedPayloadCodec = errors.New("unsupported payloadCodec")
)

type defaultCodec struct {
}

func NewDefaultCodeC() Codec {
	return &defaultCodec{}
}

func (c *defaultCodec) Encode(ctx context.Context, writer netpoll.Writer, m Message) error {
	payload, ok := m.Data().(SizeableMarshaller)
	if !ok {
		return ErrUnsupportedPayloadCodec
	}
	iv := m.RPCInfo().Invocation()
	// 整个包大小，留在最后写入
	packageLengthBuf, err := writer.Malloc(4)
	if err != nil {
		return err
	}

	// 写入Version
	err = writer.WriteByte(byte(m.Version()))
	if err != nil {
		return err
	}

	// 写入MsgType
	err = writer.WriteByte(byte(m.MsgType()))
	if err != nil {
		return err
	}

	// 写入SeqId
	seqIDBuf, err := writer.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(seqIDBuf, iv.SeqID())

	// 写入servicePath
	serviceLen := len(iv.ServiceName())
	servicePathBuf, err := writer.Malloc(serviceLen + 1)
	if err != nil {
		return err
	}
	servicePathBuf[0] = byte(serviceLen)
	copy(servicePathBuf[1:], unsafe.Slice(unsafe.StringData(iv.ServiceName()), serviceLen))

	// 写入MethodPath
	methodLen := len(iv.MethodName())
	methodPathBuf, err := writer.Malloc(methodLen + 1)
	if err != nil {
		return err
	}
	methodPathBuf[0] = byte(methodLen)
	copy(methodPathBuf[1:], unsafe.Slice(unsafe.StringData(iv.MethodName()), methodLen))

	// 写入Flags
	flags, err := writer.Malloc(2)
	if err != nil {
		return err
	}
	flags[0] = 0
	flags[1] = 0

	// 写入TraceInfo
	if len(iv.TraceInfo()) > 0 {
		flags = rpcinfo.SetFlagBit(flags, rpcinfo.TraceBit)
		_, err = writer.WriteBinary(iv.TraceInfo())
		if err != nil {
			return err
		}

	}

	// ...待扩展的Flags

	// 写入metadata
	metadataLen := len(m.Metadata())
	metadataLenBuf, err := writer.Malloc(2)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(metadataLenBuf, uint16(metadataLen))
	_, err = writer.WriteBinary(m.Metadata())
	if err != nil {
		return err
	}

	// 写入payload
	payloadLen := payload.Size()
	payloadBuf, err := writer.Malloc(payloadLen)
	if err != nil {
		return err
	}
	_, err = payload.MarshalToSizedBuffer(payloadBuf)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(packageLengthBuf, uint32(writer.MallocLen()))
	return nil
}

func (c *defaultCodec) Decode(ctx context.Context, reader netpoll.Reader, m Message) error {

	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			fmt.Printf("panic in message decode: %v, stack: %s\n", err, errStack[:n])

		}
	}()

	iv, ok := m.RPCInfo().Invocation().(rpcinfo.InvocationSetter)
	if !ok {
		return fmt.Errorf("invocation is not InvocationSetter")
	}

	packageLengthBuf, err := reader.Next(4)
	if err != nil {
		return err
	}
	packageLength := binary.BigEndian.Uint32(packageLengthBuf)

	// 解析Version
	versionBuf, err := reader.Next(1)
	if err != nil {
		return err
	}
	m.SetVersion(Version(versionBuf[0]))

	// 解析MsgType
	msgTypeBuf, err := reader.Next(1)
	if err != nil {
		return err
	}
	m.SetMsgType(MessageType(msgTypeBuf[0]))

	// 解析SeqID
	seqIDBuf, err := reader.Next(4)
	if err != nil {
		return err
	}
	seqID := binary.BigEndian.Uint32(seqIDBuf)
	iv.SetSeqID(seqID)

	// 解析ServicePath
	serviceLenBuf, err := reader.ReadByte()
	if err != nil {
		return err
	}
	serviceLen := int(serviceLenBuf)

	servicePathBuf, err := reader.Next(serviceLen)
	if err != nil {
		return err
	}
	iv.SetServiceName(unsafe.String(unsafe.SliceData(servicePathBuf), serviceLen))

	// 解析MethodPath
	methodLenBuf, err := reader.ReadByte()
	if err != nil {
		return err
	}
	methodLen := int(methodLenBuf)

	methodPathBuf, err := reader.Next(methodLen)
	if err != nil {
		return err
	}
	iv.SetMethodName(unsafe.String(unsafe.SliceData(methodPathBuf), methodLen))

	// 解析Flags
	flags, err := reader.Next(2)
	if rpcinfo.HasFlagBit(flags, rpcinfo.TraceBit) {
		traceInfo, err := reader.Next(rpcinfo.TraceLength)
		if err != nil {
			return err
		}
		iv.SetTraceInfo(traceInfo)
	}
	// 待扩展的Flags...

	// 解析Metadata
	metadataLenBuf, err := reader.Next(2)
	if err != nil {
		return err
	}
	metadataLen := int(binary.BigEndian.Uint16(metadataLenBuf))

	metadata, err := reader.Next(metadataLen)
	if err != nil {
		return err
	}
	m.SetMetadata(metadata)

	// 剩下的就是Payload了

	m.SetPayload(reader)

	m.SetLen(int(packageLength))

	return nil
}

func (c *defaultCodec) Name() string {
	return "default"
}
