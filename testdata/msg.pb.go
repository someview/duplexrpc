package testdata

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress uerror if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SubReq struct {
	SubscriptionId       int64    `protobuf:"fixed64,1,opt,name=subscriptionId,proto3" json:"subscriptionId,omitempty"`
	ConnId               int64    `protobuf:"fixed64,2,opt,name=connId,proto3" json:"connId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubReq) Reset()         { *m = SubReq{} }
func (m *SubReq) String() string { return proto.CompactTextString(m) }
func (*SubReq) ProtoMessage()    {}
func (*SubReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}
func (m *SubReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubReq.Merge(m, src)
}
func (m *SubReq) XXX_Size() int {
	return m.Size()
}
func (m *SubReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SubReq.DiscardUnknown(m)
}

var xxx_messageInfo_SubReq proto.InternalMessageInfo

func (m *SubReq) GetSubscriptionId() int64 {
	if m != nil {
		return m.SubscriptionId
	}
	return 0
}

func (m *SubReq) GetConnId() int64 {
	if m != nil {
		return m.ConnId
	}
	return 0
}

type UnsubReq struct {
	SubscriptionId       int64    `protobuf:"fixed64,1,opt,name=subscriptionId,proto3" json:"subscriptionId,omitempty"`
	ConnId               int64    `protobuf:"fixed64,2,opt,name=connId,proto3" json:"connId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnsubReq) Reset()         { *m = UnsubReq{} }
func (m *UnsubReq) String() string { return proto.CompactTextString(m) }
func (*UnsubReq) ProtoMessage()    {}
func (*UnsubReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}
func (m *UnsubReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnsubReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnsubReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnsubReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnsubReq.Merge(m, src)
}
func (m *UnsubReq) XXX_Size() int {
	return m.Size()
}
func (m *UnsubReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UnsubReq.DiscardUnknown(m)
}

var xxx_messageInfo_UnsubReq proto.InternalMessageInfo

func (m *UnsubReq) GetSubscriptionId() int64 {
	if m != nil {
		return m.SubscriptionId
	}
	return 0
}

func (m *UnsubReq) GetConnId() int64 {
	if m != nil {
		return m.ConnId
	}
	return 0
}

// 请求检查某个用户是否订阅的消息
type ExistSubReq struct {
	SubscriptionId       int64    `protobuf:"fixed64,1,opt,name=subscriptionId,proto3" json:"subscriptionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExistSubReq) Reset()         { *m = ExistSubReq{} }
func (m *ExistSubReq) String() string { return proto.CompactTextString(m) }
func (*ExistSubReq) ProtoMessage()    {}
func (*ExistSubReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}
func (m *ExistSubReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExistSubReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExistSubReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExistSubReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistSubReq.Merge(m, src)
}
func (m *ExistSubReq) XXX_Size() int {
	return m.Size()
}
func (m *ExistSubReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistSubReq.DiscardUnknown(m)
}

var xxx_messageInfo_ExistSubReq proto.InternalMessageInfo

func (m *ExistSubReq) GetSubscriptionId() int64 {
	if m != nil {
		return m.SubscriptionId
	}
	return 0
}

// 检查订阅结果的响应消息
type ExistSubRes struct {
	Exists               bool     `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExistSubRes) Reset()         { *m = ExistSubRes{} }
func (m *ExistSubRes) String() string { return proto.CompactTextString(m) }
func (*ExistSubRes) ProtoMessage()    {}
func (*ExistSubRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
}
func (m *ExistSubRes) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExistSubRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExistSubRes.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExistSubRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistSubRes.Merge(m, src)
}
func (m *ExistSubRes) XXX_Size() int {
	return m.Size()
}
func (m *ExistSubRes) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistSubRes.DiscardUnknown(m)
}

var xxx_messageInfo_ExistSubRes proto.InternalMessageInfo

func (m *ExistSubRes) GetExists() bool {
	if m != nil {
		return m.Exists
	}
	return false
}

// 由业务决定有几个Id
type SubInfo struct {
	SubscriptionId       int64    `protobuf:"fixed64,1,opt,name=subscriptionId,proto3" json:"subscriptionId,omitempty"`
	ConnId               int64    `protobuf:"fixed64,2,opt,name=connId,proto3" json:"connId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubInfo) Reset()         { *m = SubInfo{} }
func (m *SubInfo) String() string { return proto.CompactTextString(m) }
func (*SubInfo) ProtoMessage()    {}
func (*SubInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
}
func (m *SubInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubInfo.Merge(m, src)
}
func (m *SubInfo) XXX_Size() int {
	return m.Size()
}
func (m *SubInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SubInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SubInfo proto.InternalMessageInfo

func (m *SubInfo) GetSubscriptionId() int64 {
	if m != nil {
		return m.SubscriptionId
	}
	return 0
}

func (m *SubInfo) GetConnId() int64 {
	if m != nil {
		return m.ConnId
	}
	return 0
}

// 使用proto定义快照,在同步快照的时候，可以直接marshal，而不用生成新的proto
type Snapshot struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
}
func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

// 同步快照
type SyncSnapshotReq struct {
	SubInfos             []*SubInfo `protobuf:"bytes,1,rep,name=subInfos,proto3" json:"subInfos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SyncSnapshotReq) Reset()         { *m = SyncSnapshotReq{} }
func (m *SyncSnapshotReq) String() string { return proto.CompactTextString(m) }
func (*SyncSnapshotReq) ProtoMessage()    {}
func (*SyncSnapshotReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
}
func (m *SyncSnapshotReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncSnapshotReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SyncSnapshotReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SyncSnapshotReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncSnapshotReq.Merge(m, src)
}
func (m *SyncSnapshotReq) XXX_Size() int {
	return m.Size()
}
func (m *SyncSnapshotReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncSnapshotReq.DiscardUnknown(m)
}

var xxx_messageInfo_SyncSnapshotReq proto.InternalMessageInfo

func (m *SyncSnapshotReq) GetSubInfos() []*SubInfo {
	if m != nil {
		return m.SubInfos
	}
	return nil
}

type SubFailedReq struct {
	SubscriptionId       int64    `protobuf:"fixed64,1,opt,name=subscriptionId,proto3" json:"subscriptionId,omitempty"`
	ConnId               int64    `protobuf:"fixed64,2,opt,name=connId,proto3" json:"connId,omitempty"`
	Reason               string   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubFailedReq) Reset()         { *m = SubFailedReq{} }
func (m *SubFailedReq) String() string { return proto.CompactTextString(m) }
func (*SubFailedReq) ProtoMessage()    {}
func (*SubFailedReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
}
func (m *SubFailedReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubFailedReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubFailedReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubFailedReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubFailedReq.Merge(m, src)
}
func (m *SubFailedReq) XXX_Size() int {
	return m.Size()
}
func (m *SubFailedReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SubFailedReq.DiscardUnknown(m)
}

var xxx_messageInfo_SubFailedReq proto.InternalMessageInfo

func (m *SubFailedReq) GetSubscriptionId() int64 {
	if m != nil {
		return m.SubscriptionId
	}
	return 0
}

func (m *SubFailedReq) GetConnId() int64 {
	if m != nil {
		return m.ConnId
	}
	return 0
}

func (m *SubFailedReq) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

// PubReq message for distributing data to subscribers.
// The 'targets' field specifies the subscriber IDs to which the binary 'data' should be sent.
type PubReq struct {
	Targets              []int64  `protobuf:"fixed64,1,rep,packed,name=targets,proto3" json:"targets,omitempty"`
	Content              []byte   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PubReq) Reset()         { *m = PubReq{} }
func (m *PubReq) String() string { return proto.CompactTextString(m) }
func (*PubReq) ProtoMessage()    {}
func (*PubReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{8}
}
func (m *PubReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PubReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PubReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PubReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubReq.Merge(m, src)
}
func (m *PubReq) XXX_Size() int {
	return m.Size()
}
func (m *PubReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PubReq.DiscardUnknown(m)
}

var xxx_messageInfo_PubReq proto.InternalMessageInfo

func (m *PubReq) GetTargets() []int64 {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *PubReq) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*SubReq)(nil), "message.SubReq")
	proto.RegisterType((*UnsubReq)(nil), "message.UnsubReq")
	proto.RegisterType((*ExistSubReq)(nil), "message.ExistSubReq")
	proto.RegisterType((*ExistSubRes)(nil), "message.ExistSubRes")
	proto.RegisterType((*SubInfo)(nil), "message.SubInfo")
	proto.RegisterType((*Snapshot)(nil), "message.Snapshot")
	proto.RegisterType((*SyncSnapshotReq)(nil), "message.SyncSnapshotReq")
	proto.RegisterType((*SubFailedReq)(nil), "message.SubFailedReq")
	proto.RegisterType((*PubReq)(nil), "message.PubReq")
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899) }

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x41, 0x4a, 0xc4, 0x30,
	0x14, 0x86, 0x8d, 0x03, 0x9d, 0xce, 0x9b, 0x41, 0x4b, 0x17, 0xd2, 0x55, 0x19, 0x02, 0x4a, 0x17,
	0xd2, 0x85, 0xe2, 0x4e, 0x10, 0x04, 0xc5, 0xba, 0x92, 0x14, 0x0f, 0x90, 0x74, 0x32, 0xb5, 0xe0,
	0x24, 0xb5, 0x2f, 0x05, 0xbd, 0x89, 0x47, 0x72, 0xe9, 0x11, 0xa4, 0x5e, 0x44, 0xd2, 0xa6, 0x22,
	0xee, 0xc4, 0x59, 0x7e, 0xef, 0xbd, 0xfe, 0xfd, 0xf8, 0x09, 0xcc, 0x36, 0x58, 0xa6, 0x75, 0xa3,
	0x8d, 0x0e, 0xa7, 0x1b, 0x89, 0xc8, 0x4b, 0x49, 0x6f, 0xc0, 0xcb, 0x5b, 0xc1, 0xe4, 0x53, 0x78,
	0x04, 0x7b, 0xd8, 0x0a, 0x2c, 0x9a, 0xaa, 0x36, 0x95, 0x56, 0xd9, 0x2a, 0x22, 0x4b, 0x92, 0x04,
	0xec, 0xd7, 0x34, 0x3c, 0x00, 0xaf, 0xd0, 0xca, 0xee, 0x77, 0xfb, 0xbd, 0x23, 0x7a, 0x0b, 0xfe,
	0xbd, 0xc2, 0xed, 0x64, 0x9d, 0xc1, 0xfc, 0xea, 0xb9, 0x42, 0xf3, 0x37, 0x35, 0x7a, 0xf8, 0xf3,
	0x33, 0xb4, 0xe9, 0xd2, 0x22, 0xf6, 0xe7, 0x3e, 0x73, 0x44, 0x33, 0x98, 0xe6, 0xad, 0xc8, 0xd4,
	0x5a, 0xff, 0x5b, 0x14, 0xc0, 0xcf, 0x15, 0xaf, 0xf1, 0x41, 0x1b, 0x7a, 0x01, 0xfb, 0xf9, 0x8b,
	0x2a, 0x46, 0xb6, 0xe2, 0xc7, 0xe0, 0xe3, 0xf0, 0x27, 0xeb, 0x30, 0x49, 0xe6, 0x27, 0x41, 0xea,
	0x9a, 0x4f, 0x9d, 0x02, 0xfb, 0xbe, 0xa0, 0x6b, 0x58, 0xe4, 0xad, 0xb8, 0xe6, 0xd5, 0xa3, 0x5c,
	0x6d, 0xa1, 0x45, 0x3b, 0x6f, 0x24, 0x47, 0xad, 0xa2, 0xc9, 0x92, 0x24, 0x33, 0xe6, 0x88, 0x9e,
	0x83, 0x77, 0x37, 0x14, 0x1b, 0xc1, 0xd4, 0xf0, 0xa6, 0x94, 0x66, 0xd0, 0x0b, 0xd8, 0x88, 0x76,
	0x53, 0x68, 0x65, 0xa4, 0x32, 0x7d, 0xe8, 0x82, 0x8d, 0x78, 0x19, 0xbc, 0x75, 0x31, 0x79, 0xef,
	0x62, 0xf2, 0xd1, 0xc5, 0xe4, 0xf5, 0x33, 0xde, 0x11, 0x5e, 0xff, 0xa6, 0x4e, 0xbf, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xd6, 0x89, 0x13, 0x86, 0x60, 0x02, 0x00, 0x00,
}

func (m *SubReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ConnId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.ConnId))
		i--
		dAtA[i] = 0x11
	}
	if m.SubscriptionId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.SubscriptionId))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *UnsubReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnsubReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnsubReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ConnId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.ConnId))
		i--
		dAtA[i] = 0x11
	}
	if m.SubscriptionId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.SubscriptionId))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *ExistSubReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExistSubReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExistSubReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.SubscriptionId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.SubscriptionId))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *ExistSubRes) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExistSubRes) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExistSubRes) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Exists {
		i--
		if m.Exists {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SubInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ConnId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.ConnId))
		i--
		dAtA[i] = 0x11
	}
	if m.SubscriptionId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.SubscriptionId))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *SyncSnapshotReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncSnapshotReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncSnapshotReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.SubInfos) > 0 {
		for iNdEx := len(m.SubInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SubInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsg(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SubFailedReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubFailedReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubFailedReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Reason) > 0 {
		i -= len(m.Reason)
		copy(dAtA[i:], m.Reason)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Reason)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ConnId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.ConnId))
		i--
		dAtA[i] = 0x11
	}
	if m.SubscriptionId != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.SubscriptionId))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *PubReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PubReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PubReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Content) > 0 {
		i -= len(m.Content)
		copy(dAtA[i:], m.Content)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Content)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Targets) > 0 {
		for iNdEx := len(m.Targets) - 1; iNdEx >= 0; iNdEx-- {
			i -= 8
			encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.Targets[iNdEx]))
		}
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Targets)*8))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsg(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsg(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SubReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubscriptionId != 0 {
		n += 9
	}
	if m.ConnId != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UnsubReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubscriptionId != 0 {
		n += 9
	}
	if m.ConnId != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ExistSubReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubscriptionId != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ExistSubRes) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Exists {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SubInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubscriptionId != 0 {
		n += 9
	}
	if m.ConnId != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SyncSnapshotReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SubInfos) > 0 {
		for _, e := range m.SubInfos {
			l = e.Size()
			n += 1 + l + sovMsg(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SubFailedReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubscriptionId != 0 {
		n += 9
	}
	if m.ConnId != 0 {
		n += 9
	}
	l = len(m.Reason)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PubReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Targets) > 0 {
		n += 1 + sovMsg(uint64(len(m.Targets)*8)) + len(m.Targets)*8
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMsg(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsg(x uint64) (n int) {
	return sovMsg(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SubReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: SubReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriptionId", wireType)
			}
			m.SubscriptionId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriptionId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnId", wireType)
			}
			m.ConnId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *UnsubReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: UnsubReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnsubReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriptionId", wireType)
			}
			m.SubscriptionId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriptionId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnId", wireType)
			}
			m.ConnId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *ExistSubReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: ExistSubReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExistSubReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriptionId", wireType)
			}
			m.SubscriptionId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriptionId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *ExistSubRes) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: ExistSubRes: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExistSubRes: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exists", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
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
			m.Exists = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *SubInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: SubInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriptionId", wireType)
			}
			m.SubscriptionId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriptionId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnId", wireType)
			}
			m.ConnId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *SyncSnapshotReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: SyncSnapshotReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncSnapshotReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
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
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubInfos = append(m.SubInfos, &SubInfo{})
			if err := m.SubInfos[len(m.SubInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *SubFailedReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: SubFailedReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubFailedReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriptionId", wireType)
			}
			m.SubscriptionId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriptionId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnId", wireType)
			}
			m.ConnId = 0
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnId = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reason", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
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
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Reason = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func (m *PubReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
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
			return fmt.Errorf("proto: PubReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PubReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 1 {
				var v int64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
				iNdEx += 8
				m.Targets = append(m.Targets, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowMsg
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
					return ErrInvalidLengthMsg
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthMsg
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 8
				if elementCount != 0 && len(m.Targets) == 0 {
					m.Targets = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					v = int64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
					iNdEx += 8
					m.Targets = append(m.Targets, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Targets", wireType)
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
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
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = append(m.Content[:0], dAtA[iNdEx:postIndex]...)
			if m.Content == nil {
				m.Content = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
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
func skipMsg(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsg
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
					return 0, ErrIntOverflowMsg
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
					return 0, ErrIntOverflowMsg
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
				return 0, ErrInvalidLengthMsg
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsg
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsg
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsg        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsg          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsg = fmt.Errorf("proto: unexpected end of group")
)
