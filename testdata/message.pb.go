package testdata

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RoomMemberOp int32

const (
	RoomMemberOp_UnknownOpType RoomMemberOp = 0
	// 增加群成员
	RoomMemberOp_Add RoomMemberOp = 1
	// 删除群成员
	RoomMemberOp_Remove RoomMemberOp = 2
	// 更新禁言群成员
	RoomMemberOp_Gag RoomMemberOp = 3
	// 删除群
	RoomMemberOp_RoomDel RoomMemberOp = 4
)

var RoomMemberOp_name = map[int32]string{
	0: "UnknownOpType",
	1: "Add",
	2: "Remove",
	3: "Gag",
	4: "RoomDel",
}

var RoomMemberOp_value = map[string]int32{
	"UnknownOpType": 0,
	"Add":           1,
	"Remove":        2,
	"Gag":           3,
	"RoomDel":       4,
}

func (x RoomMemberOp) String() string {
	return proto.EnumName(RoomMemberOp_name, int32(x))
}

func (RoomMemberOp) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

type ClientType int32

const (
	ClientType_UnknownType ClientType = 0
	// 订阅或取消订阅
	ClientType_BatchSubOrUnSubType ClientType = 1
	ClientType_BatchPubType        ClientType = 2
	ClientType_UniSubType          ClientType = 3
)

var ClientType_name = map[int32]string{
	0: "UnknownType",
	1: "BatchSubOrUnSubType",
	2: "BatchPubType",
	3: "UniSubType",
}

var ClientType_value = map[string]int32{
	"UnknownType":         0,
	"BatchSubOrUnSubType": 1,
	"BatchPubType":        2,
	"UniSubType":          3,
}

func (x ClientType) String() string {
	return proto.EnumName(ClientType_name, int32(x))
}

func (ClientType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

type ServerType int32

const (
	ServerType_UnknownResponseType ServerType = 0
	ServerType_MessagesType        ServerType = 1
	ServerType_UniSubResponseType  ServerType = 2
	ServerType_PubNotFoundType     ServerType = 5
)

var ServerType_name = map[int32]string{
	0: "UnknownResponseType",
	1: "MessagesType",
	2: "UniSubResponseType",
	5: "PubNotFoundType",
}

var ServerType_value = map[string]int32{
	"UnknownResponseType": 0,
	"MessagesType":        1,
	"UniSubResponseType":  2,
	"PubNotFoundType":     5,
}

func (x ServerType) String() string {
	return proto.EnumName(ServerType_name, int32(x))
}

func (ServerType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

type Header struct {
	// 消息体公共头部
	TraceId              string   `protobuf:"bytes,1,opt,name=traceId,proto3" json:"traceId,omitempty"`
	SpanId               string   `protobuf:"bytes,2,opt,name=spanId,proto3" json:"spanId,omitempty"`
	Timestamp            int64    `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}
func (m *Header) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Header.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return m.Size()
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

func (m *Header) GetSpanId() string {
	if m != nil {
		return m.SpanId
	}
	return ""
}

func (m *Header) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type PayloadMessage struct {
	Type                 RoomMemberOp `protobuf:"varint,2,opt,name=type,proto3,enum=TL.Protobuf.RouteData.RoomMsg.RoomMemberOp" json:"type,omitempty"`
	ManIds               []int64      `protobuf:"fixed64,3,rep,packed,name=manIds,proto3" json:"manIds,omitempty"`
	Scores               []int64      `protobuf:"fixed64,4,rep,packed,name=scores,proto3" json:"scores,omitempty"`
	RoomId               string       `protobuf:"bytes,5,opt,name=roomId,proto3" json:"roomId,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PayloadMessage) Reset()         { *m = PayloadMessage{} }
func (m *PayloadMessage) String() string { return proto.CompactTextString(m) }
func (*PayloadMessage) ProtoMessage()    {}
func (*PayloadMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}
func (m *PayloadMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PayloadMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PayloadMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PayloadMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadMessage.Merge(m, src)
}
func (m *PayloadMessage) XXX_Size() int {
	return m.Size()
}
func (m *PayloadMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadMessage proto.InternalMessageInfo

func (m *PayloadMessage) GetType() RoomMemberOp {
	if m != nil {
		return m.Type
	}
	return RoomMemberOp_UnknownOpType
}

func (m *PayloadMessage) GetManIds() []int64 {
	if m != nil {
		return m.ManIds
	}
	return nil
}

func (m *PayloadMessage) GetScores() []int64 {
	if m != nil {
		return m.Scores
	}
	return nil
}

func (m *PayloadMessage) GetRoomId() string {
	if m != nil {
		return m.RoomId
	}
	return ""
}

type ClientMessage struct {
	Type                 ClientType              `protobuf:"varint,1,opt,name=type,proto3,enum=TL.Protobuf.RouteData.RoomMsg.ClientType" json:"type,omitempty"`
	BatchSubOrUnSub      *BatchSubOrUnSubRequest `protobuf:"bytes,2,opt,name=batchSubOrUnSub,proto3" json:"batchSubOrUnSub,omitempty"`
	BatchPub             *BatchPubRequest        `protobuf:"bytes,3,opt,name=batchPub,proto3" json:"batchPub,omitempty"`
	UniSub               *UniSubRequest          `protobuf:"bytes,4,opt,name=uniSub,proto3" json:"uniSub,omitempty"`
	Header               *Header                 `protobuf:"bytes,10,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ClientMessage) Reset()         { *m = ClientMessage{} }
func (m *ClientMessage) String() string { return proto.CompactTextString(m) }
func (*ClientMessage) ProtoMessage()    {}
func (*ClientMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}
func (m *ClientMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClientMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClientMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClientMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientMessage.Merge(m, src)
}
func (m *ClientMessage) XXX_Size() int {
	return m.Size()
}
func (m *ClientMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ClientMessage proto.InternalMessageInfo

func (m *ClientMessage) GetType() ClientType {
	if m != nil {
		return m.Type
	}
	return ClientType_UnknownType
}

func (m *ClientMessage) GetBatchSubOrUnSub() *BatchSubOrUnSubRequest {
	if m != nil {
		return m.BatchSubOrUnSub
	}
	return nil
}

func (m *ClientMessage) GetBatchPub() *BatchPubRequest {
	if m != nil {
		return m.BatchPub
	}
	return nil
}

func (m *ClientMessage) GetUniSub() *UniSubRequest {
	if m != nil {
		return m.UniSub
	}
	return nil
}

func (m *ClientMessage) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type BatchSubOrUnSubRequest struct {
	IsSub                bool     `protobuf:"varint,1,opt,name=IsSub,proto3" json:"IsSub,omitempty"`
	Channels             []int32  `protobuf:"varint,2,rep,packed,name=channels,proto3" json:"channels,omitempty"`
	TopicList            []string `protobuf:"bytes,3,rep,name=topicList,proto3" json:"topicList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchSubOrUnSubRequest) Reset()         { *m = BatchSubOrUnSubRequest{} }
func (m *BatchSubOrUnSubRequest) String() string { return proto.CompactTextString(m) }
func (*BatchSubOrUnSubRequest) ProtoMessage()    {}
func (*BatchSubOrUnSubRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}
func (m *BatchSubOrUnSubRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BatchSubOrUnSubRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BatchSubOrUnSubRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BatchSubOrUnSubRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchSubOrUnSubRequest.Merge(m, src)
}
func (m *BatchSubOrUnSubRequest) XXX_Size() int {
	return m.Size()
}
func (m *BatchSubOrUnSubRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchSubOrUnSubRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BatchSubOrUnSubRequest proto.InternalMessageInfo

func (m *BatchSubOrUnSubRequest) GetIsSub() bool {
	if m != nil {
		return m.IsSub
	}
	return false
}

func (m *BatchSubOrUnSubRequest) GetChannels() []int32 {
	if m != nil {
		return m.Channels
	}
	return nil
}

func (m *BatchSubOrUnSubRequest) GetTopicList() []string {
	if m != nil {
		return m.TopicList
	}
	return nil
}

type UniSubRequest struct {
	Channels             []int32  `protobuf:"varint,1,rep,packed,name=channels,proto3" json:"channels,omitempty"`
	Topic                string   `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UniSubRequest) Reset()         { *m = UniSubRequest{} }
func (m *UniSubRequest) String() string { return proto.CompactTextString(m) }
func (*UniSubRequest) ProtoMessage()    {}
func (*UniSubRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}
func (m *UniSubRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UniSubRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UniSubRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UniSubRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UniSubRequest.Merge(m, src)
}
func (m *UniSubRequest) XXX_Size() int {
	return m.Size()
}
func (m *UniSubRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UniSubRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UniSubRequest proto.InternalMessageInfo

func (m *UniSubRequest) GetChannels() []int32 {
	if m != nil {
		return m.Channels
	}
	return nil
}

func (m *UniSubRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

type BatchPubRequest struct {
	Channel int32    `protobuf:"varint,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Topic   []string `protobuf:"bytes,2,rep,name=Topic,proto3" json:"Topic,omitempty"`
	Payload []byte   `protobuf:"bytes,3,opt,name=Payload,proto3" json:"Payload,omitempty"`
	//当推送的Topic未找到且设置了此字段,则服务器将此字段返回给客户端
	NotFound             []byte   `protobuf:"bytes,4,opt,name=NotFound,proto3" json:"NotFound,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchPubRequest) Reset()         { *m = BatchPubRequest{} }
func (m *BatchPubRequest) String() string { return proto.CompactTextString(m) }
func (*BatchPubRequest) ProtoMessage()    {}
func (*BatchPubRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}
func (m *BatchPubRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BatchPubRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BatchPubRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BatchPubRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchPubRequest.Merge(m, src)
}
func (m *BatchPubRequest) XXX_Size() int {
	return m.Size()
}
func (m *BatchPubRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchPubRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BatchPubRequest proto.InternalMessageInfo

func (m *BatchPubRequest) GetChannel() int32 {
	if m != nil {
		return m.Channel
	}
	return 0
}

func (m *BatchPubRequest) GetTopic() []string {
	if m != nil {
		return m.Topic
	}
	return nil
}

func (m *BatchPubRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *BatchPubRequest) GetNotFound() []byte {
	if m != nil {
		return m.NotFound
	}
	return nil
}

type ServerMessage struct {
	Type ServerType `protobuf:"varint,1,opt,name=type,proto3,enum=TL.Protobuf.RouteData.RoomMsg.ServerType" json:"type,omitempty"`
	//服务端推送过来的消息
	Messages *Message `protobuf:"bytes,2,opt,name=messages,proto3" json:"messages,omitempty"`
	// 服务端转发的消息类型
	//对应着PubRequest的NotFound字段
	PubNotFound          []byte   `protobuf:"bytes,5,opt,name=pubNotFound,proto3" json:"pubNotFound,omitempty"`
	Header               *Header  `protobuf:"bytes,10,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerMessage) Reset()         { *m = ServerMessage{} }
func (m *ServerMessage) String() string { return proto.CompactTextString(m) }
func (*ServerMessage) ProtoMessage()    {}
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6}
}
func (m *ServerMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ServerMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ServerMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ServerMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerMessage.Merge(m, src)
}
func (m *ServerMessage) XXX_Size() int {
	return m.Size()
}
func (m *ServerMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ServerMessage proto.InternalMessageInfo

func (m *ServerMessage) GetType() ServerType {
	if m != nil {
		return m.Type
	}
	return ServerType_UnknownResponseType
}

func (m *ServerMessage) GetMessages() *Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *ServerMessage) GetPubNotFound() []byte {
	if m != nil {
		return m.PubNotFound
	}
	return nil
}

func (m *ServerMessage) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type Message struct {
	Channel              int32    `protobuf:"varint,1,opt,name=Channel,proto3" json:"Channel,omitempty"`
	Topic                []string `protobuf:"bytes,2,rep,name=Topic,proto3" json:"Topic,omitempty"`
	Payload              []byte   `protobuf:"bytes,3,opt,name=Payload,proto3" json:"Payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetChannel() int32 {
	if m != nil {
		return m.Channel
	}
	return 0
}

func (m *Message) GetTopic() []string {
	if m != nil {
		return m.Topic
	}
	return nil
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterEnum("TL.Protobuf.RouteData.RoomMsg.RoomMemberOp", RoomMemberOp_name, RoomMemberOp_value)
	proto.RegisterEnum("TL.Protobuf.RouteData.RoomMsg.ClientType", ClientType_name, ClientType_value)
	proto.RegisterEnum("TL.Protobuf.RouteData.RoomMsg.ServerType", ServerType_name, ServerType_value)
	proto.RegisterType((*Header)(nil), "TL.Protobuf.RouteData.RoomMsg.Header")
	proto.RegisterType((*PayloadMessage)(nil), "TL.Protobuf.RouteData.RoomMsg.PayloadMessage")
	proto.RegisterType((*ClientMessage)(nil), "TL.Protobuf.RouteData.RoomMsg.ClientMessage")
	proto.RegisterType((*BatchSubOrUnSubRequest)(nil), "TL.Protobuf.RouteData.RoomMsg.BatchSubOrUnSubRequest")
	proto.RegisterType((*UniSubRequest)(nil), "TL.Protobuf.RouteData.RoomMsg.UniSubRequest")
	proto.RegisterType((*BatchPubRequest)(nil), "TL.Protobuf.RouteData.RoomMsg.BatchPubRequest")
	proto.RegisterType((*ServerMessage)(nil), "TL.Protobuf.RouteData.RoomMsg.ServerMessage")
	proto.RegisterType((*Message)(nil), "TL.Protobuf.RouteData.RoomMsg.Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 673 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0xcf, 0x6e, 0xd3, 0x4c,
	0x14, 0xc5, 0xeb, 0x38, 0x4e, 0xd2, 0x9b, 0xa4, 0xf1, 0x37, 0xad, 0xfa, 0x59, 0x08, 0xa2, 0xc8,
	0x12, 0x28, 0x14, 0x94, 0x45, 0x11, 0xcb, 0x0a, 0xf5, 0x8f, 0x80, 0xa0, 0x96, 0x46, 0x93, 0x56,
	0xea, 0x0e, 0xd9, 0xf1, 0xd0, 0x44, 0xc4, 0x1e, 0xe3, 0xb1, 0x5b, 0x75, 0xc7, 0x63, 0xc0, 0x1b,
	0xb1, 0xe4, 0x11, 0x50, 0x79, 0x04, 0x5e, 0x00, 0xcd, 0x9d, 0x71, 0x9c, 0x44, 0x88, 0xa0, 0xb2,
	0x6a, 0xcf, 0xb5, 0xef, 0x6f, 0xce, 0xdc, 0x39, 0x9e, 0x40, 0x33, 0x64, 0x42, 0x78, 0x97, 0xac,
	0x17, 0x27, 0x3c, 0xe5, 0xe4, 0xc1, 0xd9, 0x71, 0x6f, 0x20, 0xff, 0xf3, 0xb3, 0xf7, 0x3d, 0xca,
	0xb3, 0x94, 0x1d, 0x79, 0xa9, 0xd7, 0xa3, 0x9c, 0x87, 0x27, 0xe2, 0xd2, 0xbd, 0x80, 0xca, 0x6b,
	0xe6, 0x05, 0x2c, 0x21, 0x0e, 0x54, 0xd3, 0xc4, 0x1b, 0xb1, 0x7e, 0xe0, 0x18, 0x1d, 0xa3, 0xbb,
	0x4e, 0x73, 0x49, 0xb6, 0xa1, 0x22, 0x62, 0x2f, 0xea, 0x07, 0x4e, 0x09, 0x1f, 0x68, 0x45, 0xee,
	0xc3, 0x7a, 0x3a, 0x09, 0x99, 0x48, 0xbd, 0x30, 0x76, 0xcc, 0x8e, 0xd1, 0x35, 0x69, 0x51, 0x70,
	0xbf, 0x18, 0xb0, 0x31, 0xf0, 0x6e, 0xa6, 0xdc, 0x0b, 0x4e, 0x94, 0x23, 0xf2, 0x02, 0xca, 0xe9,
	0x4d, 0xcc, 0x10, 0xb3, 0xb1, 0xfb, 0xa4, 0xf7, 0x47, 0x6b, 0xea, 0x2f, 0x0b, 0x7d, 0x96, 0x9c,
	0xc6, 0x14, 0x1b, 0xa5, 0x93, 0x50, 0x2e, 0x2d, 0x1c, 0xb3, 0x63, 0x76, 0x6d, 0xaa, 0x15, 0x3a,
	0x1c, 0xf1, 0x84, 0x09, 0xa7, 0xac, 0xea, 0x4a, 0xc9, 0x7a, 0xc2, 0x79, 0xd8, 0x0f, 0x1c, 0x4b,
	0x39, 0x57, 0xca, 0xfd, 0x64, 0x42, 0xf3, 0x70, 0x3a, 0x61, 0x51, 0x9a, 0x5b, 0xdb, 0xd3, 0xd6,
	0x0c, 0xb4, 0xf6, 0x78, 0x85, 0x35, 0xd5, 0x7b, 0x76, 0x13, 0x33, 0x6d, 0xec, 0x1d, 0xb4, 0x7c,
	0x2f, 0x1d, 0x8d, 0x87, 0x99, 0x7f, 0x9a, 0x9c, 0x47, 0xc3, 0xcc, 0xc7, 0x4d, 0xd6, 0x77, 0x9f,
	0xaf, 0x20, 0x1d, 0x2c, 0x76, 0x51, 0xf6, 0x31, 0x63, 0x22, 0xa5, 0xcb, 0x34, 0xf2, 0x06, 0x6a,
	0x58, 0x1a, 0x64, 0x3e, 0x8e, 0xba, 0xbe, 0xdb, 0xfb, 0x1b, 0xf2, 0xa0, 0x40, 0xce, 0xfa, 0xc9,
	0x11, 0x54, 0xb2, 0x68, 0x22, 0x3d, 0x96, 0x91, 0xf4, 0x74, 0x05, 0xe9, 0x1c, 0x5f, 0xce, 0x39,
	0xba, 0x97, 0xec, 0x41, 0x65, 0x8c, 0xc9, 0x71, 0x00, 0x29, 0x0f, 0x57, 0x50, 0x54, 0xcc, 0xa8,
	0x6e, 0x72, 0xc7, 0xb0, 0xfd, 0xfb, 0xbd, 0x93, 0x2d, 0xb0, 0xfa, 0x42, 0xba, 0x93, 0x67, 0x51,
	0xa3, 0x4a, 0x90, 0x7b, 0x50, 0x1b, 0x8d, 0xbd, 0x28, 0x62, 0x53, 0xe1, 0x94, 0x3a, 0x66, 0xd7,
	0xa2, 0x33, 0x8d, 0x41, 0xe4, 0xf1, 0x64, 0x74, 0x3c, 0x11, 0x29, 0x26, 0x63, 0x9d, 0x16, 0x05,
	0x77, 0x1f, 0x9a, 0x0b, 0x3b, 0x58, 0x40, 0x19, 0x4b, 0xa8, 0x2d, 0xb0, 0xb0, 0x53, 0x47, 0x5d,
	0x09, 0xf7, 0x1a, 0x5a, 0x4b, 0xe3, 0x94, 0x9f, 0x8b, 0x6e, 0x42, 0x9f, 0x16, 0xcd, 0xa5, 0x44,
	0x9c, 0x69, 0x84, 0x74, 0xa2, 0x84, 0x7c, 0x5f, 0x7f, 0x0d, 0x78, 0x7e, 0x0d, 0x9a, 0x4b, 0x69,
	0xe7, 0x2d, 0x4f, 0x5f, 0xf2, 0x2c, 0x0a, 0xf0, 0x40, 0x1a, 0x74, 0xa6, 0xdd, 0x9f, 0x06, 0x34,
	0x87, 0x2c, 0xb9, 0x62, 0xc9, 0xdd, 0x82, 0xaa, 0x7a, 0xe7, 0x82, 0x7a, 0x00, 0x35, 0x7d, 0x3f,
	0x08, 0x9d, 0xd0, 0x47, 0x2b, 0x10, 0x7a, 0x61, 0x3a, 0xeb, 0x23, 0x1d, 0xa8, 0xc7, 0x99, 0x3f,
	0xf3, 0x6c, 0xa1, 0xe7, 0xf9, 0xd2, 0xbf, 0x66, 0x63, 0x08, 0xd5, 0x7c, 0xbb, 0x0e, 0x54, 0x0f,
	0x17, 0xc7, 0x7c, 0x78, 0xb7, 0x31, 0xef, 0x1c, 0x43, 0x63, 0xfe, 0x46, 0x21, 0xff, 0xc9, 0x58,
	0x7c, 0x88, 0xf8, 0x75, 0x74, 0x1a, 0xcb, 0x01, 0xd9, 0x6b, 0xa4, 0x0a, 0xe6, 0x7e, 0x10, 0xd8,
	0x06, 0x01, 0xa8, 0x50, 0x16, 0xf2, 0x2b, 0x66, 0x97, 0x64, 0xf1, 0x95, 0x77, 0x69, 0x9b, 0xa4,
	0x0e, 0x55, 0x09, 0x38, 0x62, 0x53, 0xbb, 0xbc, 0x73, 0x01, 0x50, 0x5c, 0x02, 0xa4, 0x05, 0x75,
	0xcd, 0xd2, 0xa4, 0xff, 0x61, 0x73, 0x29, 0xdd, 0xf8, 0xc0, 0x20, 0x36, 0x34, 0xf2, 0x24, 0x61,
	0xa5, 0x44, 0x36, 0x00, 0x54, 0x3c, 0x51, 0x9b, 0x3b, 0x01, 0x40, 0x71, 0x6a, 0x12, 0xa4, 0xc9,
	0x94, 0x89, 0x98, 0x47, 0x82, 0xe9, 0x15, 0x6c, 0x68, 0xe8, 0x19, 0x09, 0x8d, 0xde, 0x06, 0x92,
	0xe7, 0x7c, 0xee, 0xcd, 0x12, 0xd9, 0x84, 0xd6, 0xa0, 0x38, 0x1b, 0x2c, 0x5a, 0x07, 0xf6, 0xd7,
	0xdb, 0xb6, 0xf1, 0xed, 0xb6, 0x6d, 0x7c, 0xbf, 0x6d, 0x1b, 0x9f, 0x7f, 0xb4, 0xd7, 0xfc, 0x0a,
	0xfe, 0x5e, 0x3c, 0xfb, 0x15, 0x00, 0x00, 0xff, 0xff, 0x85, 0xce, 0xc9, 0x78, 0x40, 0x06, 0x00,
	0x00,
}

func (m *Header) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Header) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Header) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Timestamp != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x18
	}
	if len(m.SpanId) > 0 {
		i -= len(m.SpanId)
		copy(dAtA[i:], m.SpanId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.SpanId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.TraceId) > 0 {
		i -= len(m.TraceId)
		copy(dAtA[i:], m.TraceId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.TraceId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PayloadMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PayloadMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PayloadMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.RoomId) > 0 {
		i -= len(m.RoomId)
		copy(dAtA[i:], m.RoomId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.RoomId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Scores) > 0 {
		for iNdEx := len(m.Scores) - 1; iNdEx >= 0; iNdEx-- {
			i -= 8
			encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.Scores[iNdEx]))
		}
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Scores)*8))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ManIds) > 0 {
		for iNdEx := len(m.ManIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= 8
			encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.ManIds[iNdEx]))
		}
		i = encodeVarintMessage(dAtA, i, uint64(len(m.ManIds)*8))
		i--
		dAtA[i] = 0x1a
	}
	if m.Type != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	return len(dAtA) - i, nil
}

func (m *ClientMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClientMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Header != nil {
		{
			size, err := m.Header.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x52
	}
	if m.UniSub != nil {
		{
			size, err := m.UniSub.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.BatchPub != nil {
		{
			size, err := m.BatchPub.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.BatchSubOrUnSub != nil {
		{
			size, err := m.BatchSubOrUnSub.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Type != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BatchSubOrUnSubRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BatchSubOrUnSubRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BatchSubOrUnSubRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.TopicList) > 0 {
		for iNdEx := len(m.TopicList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.TopicList[iNdEx])
			copy(dAtA[i:], m.TopicList[iNdEx])
			i = encodeVarintMessage(dAtA, i, uint64(len(m.TopicList[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Channels) > 0 {
		dAtA6 := make([]byte, len(m.Channels)*10)
		var j5 int
		for _, num1 := range m.Channels {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA6[j5] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j5++
			}
			dAtA6[j5] = uint8(num)
			j5++
		}
		i -= j5
		copy(dAtA[i:], dAtA6[:j5])
		i = encodeVarintMessage(dAtA, i, uint64(j5))
		i--
		dAtA[i] = 0x12
	}
	if m.IsSub {
		i--
		if m.IsSub {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UniSubRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UniSubRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UniSubRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Topic) > 0 {
		i -= len(m.Topic)
		copy(dAtA[i:], m.Topic)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Channels) > 0 {
		dAtA8 := make([]byte, len(m.Channels)*10)
		var j7 int
		for _, num1 := range m.Channels {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA8[j7] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j7++
			}
			dAtA8[j7] = uint8(num)
			j7++
		}
		i -= j7
		copy(dAtA[i:], dAtA8[:j7])
		i = encodeVarintMessage(dAtA, i, uint64(j7))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BatchPubRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BatchPubRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BatchPubRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.NotFound) > 0 {
		i -= len(m.NotFound)
		copy(dAtA[i:], m.NotFound)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.NotFound)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Topic) > 0 {
		for iNdEx := len(m.Topic) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Topic[iNdEx])
			copy(dAtA[i:], m.Topic[iNdEx])
			i = encodeVarintMessage(dAtA, i, uint64(len(m.Topic[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Channel != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Channel))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ServerMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ServerMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ServerMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Header != nil {
		{
			size, err := m.Header.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x52
	}
	if len(m.PubNotFound) > 0 {
		i -= len(m.PubNotFound)
		copy(dAtA[i:], m.PubNotFound)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.PubNotFound)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Messages != nil {
		{
			size, err := m.Messages.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Type != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Topic) > 0 {
		for iNdEx := len(m.Topic) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Topic[iNdEx])
			copy(dAtA[i:], m.Topic[iNdEx])
			i = encodeVarintMessage(dAtA, i, uint64(len(m.Topic[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Channel != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.Channel))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessage(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Header) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TraceId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.SpanId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Timestamp != 0 {
		n += 1 + sovMessage(uint64(m.Timestamp))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PayloadMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovMessage(uint64(m.Type))
	}
	if len(m.ManIds) > 0 {
		n += 1 + sovMessage(uint64(len(m.ManIds)*8)) + len(m.ManIds)*8
	}
	if len(m.Scores) > 0 {
		n += 1 + sovMessage(uint64(len(m.Scores)*8)) + len(m.Scores)*8
	}
	l = len(m.RoomId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ClientMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovMessage(uint64(m.Type))
	}
	if m.BatchSubOrUnSub != nil {
		l = m.BatchSubOrUnSub.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.BatchPub != nil {
		l = m.BatchPub.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.UniSub != nil {
		l = m.UniSub.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Header != nil {
		l = m.Header.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BatchSubOrUnSubRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.IsSub {
		n += 2
	}
	if len(m.Channels) > 0 {
		l = 0
		for _, e := range m.Channels {
			l += sovMessage(uint64(e))
		}
		n += 1 + sovMessage(uint64(l)) + l
	}
	if len(m.TopicList) > 0 {
		for _, s := range m.TopicList {
			l = len(s)
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UniSubRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Channels) > 0 {
		l = 0
		for _, e := range m.Channels {
			l += sovMessage(uint64(e))
		}
		n += 1 + sovMessage(uint64(l)) + l
	}
	l = len(m.Topic)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BatchPubRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Channel != 0 {
		n += 1 + sovMessage(uint64(m.Channel))
	}
	if len(m.Topic) > 0 {
		for _, s := range m.Topic {
			l = len(s)
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.NotFound)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ServerMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovMessage(uint64(m.Type))
	}
	if m.Messages != nil {
		l = m.Messages.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.PubNotFound)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Header != nil {
		l = m.Header.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Channel != 0 {
		n += 1 + sovMessage(uint64(m.Channel))
	}
	if len(m.Topic) > 0 {
		for _, s := range m.Topic {
			l = len(s)
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessage(x uint64) (n int) {
	return sovMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Header) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Header: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Header: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TraceId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpanId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SpanId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PayloadMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PayloadMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PayloadMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= RoomMemberOp(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType == 1 {
				var v int64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
				iNdEx += 8
				m.ManIds = append(m.ManIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthMessage
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthMessage
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 8
				if elementCount != 0 && len(m.ManIds) == 0 {
					m.ManIds = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
					iNdEx += 8
					m.ManIds = append(m.ManIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ManIds", wireType)
			}
		case 4:
			if wireType == 1 {
				var v int64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
				iNdEx += 8
				m.Scores = append(m.Scores, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthMessage
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthMessage
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 8
				if elementCount != 0 && len(m.Scores) == 0 {
					m.Scores = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
					iNdEx += 8
					m.Scores = append(m.Scores, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Scores", wireType)
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoomId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RoomId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ClientMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClientMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ClientType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchSubOrUnSub", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BatchSubOrUnSub == nil {
				m.BatchSubOrUnSub = &BatchSubOrUnSubRequest{}
			}
			if err := m.BatchSubOrUnSub.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchPub", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BatchPub == nil {
				m.BatchPub = &BatchPubRequest{}
			}
			if err := m.BatchPub.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UniSub", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UniSub == nil {
				m.UniSub = &UniSubRequest{}
			}
			if err := m.UniSub.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Header == nil {
				m.Header = &Header{}
			}
			if err := m.Header.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BatchSubOrUnSubRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BatchSubOrUnSubRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BatchSubOrUnSubRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsSub", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsSub = bool(v != 0)
		case 2:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Channels = append(m.Channels, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthMessage
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthMessage
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Channels) == 0 {
					m.Channels = make([]int32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowMessage
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Channels = append(m.Channels, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Channels", wireType)
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TopicList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TopicList = append(m.TopicList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UniSubRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UniSubRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UniSubRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Channels = append(m.Channels, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthMessage
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthMessage
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Channels) == 0 {
					m.Channels = make([]int32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowMessage
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Channels = append(m.Channels, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Channels", wireType)
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Topic = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BatchPubRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BatchPubRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BatchPubRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Channel", wireType)
			}
			m.Channel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Channel |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Topic = append(m.Topic, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NotFound", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NotFound = append(m.NotFound[:0], dAtA[iNdEx:postIndex]...)
			if m.NotFound == nil {
				m.NotFound = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ServerMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ServerMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ServerMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ServerType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Messages == nil {
				m.Messages = &Message{}
			}
			if err := m.Messages.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubNotFound", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubNotFound = append(m.PubNotFound[:0], dAtA[iNdEx:postIndex]...)
			if m.PubNotFound == nil {
				m.PubNotFound = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Header == nil {
				m.Header = &Header{}
			}
			if err := m.Header.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Channel", wireType)
			}
			m.Channel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Channel |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Topic = append(m.Topic, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessage = fmt.Errorf("proto: unexpected end of group")
)
