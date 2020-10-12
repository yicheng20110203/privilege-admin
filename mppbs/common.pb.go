// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package mppbs

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	math "math"
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

// output any
type Resp struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data                 *any.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resp) Reset()         { *m = Resp{} }
func (m *Resp) String() string { return proto.CompactTextString(m) }
func (*Resp) ProtoMessage()    {}
func (*Resp) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *Resp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Resp.Unmarshal(m, b)
}
func (m *Resp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Resp.Marshal(b, m, deterministic)
}
func (m *Resp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resp.Merge(m, src)
}
func (m *Resp) XXX_Size() int {
	return xxx_messageInfo_Resp.Size(m)
}
func (m *Resp) XXX_DiscardUnknown() {
	xxx_messageInfo_Resp.DiscardUnknown(m)
}

var xxx_messageInfo_Resp proto.InternalMessageInfo

func (m *Resp) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Resp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *Resp) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Resp)(nil), "mppbs.Resp")
}

func init() {
	proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206)
}

var fileDescriptor_555bd8c177793206 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0x2d, 0x28, 0x48, 0x2a, 0x96,
	0x92, 0x4c, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x0b, 0x26, 0x95, 0xa6, 0xe9, 0x27, 0xe6,
	0x55, 0x42, 0x54, 0x28, 0x85, 0x71, 0xb1, 0x04, 0xa5, 0x16, 0x17, 0x08, 0x09, 0x71, 0xb1, 0x24,
	0xe7, 0xa7, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x81, 0xd9, 0x42, 0x02, 0x5c, 0xcc,
	0xb9, 0xc5, 0xe9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x20, 0xa6, 0x90, 0x06, 0x17, 0x4b,
	0x4a, 0x62, 0x49, 0xa2, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0xc4, 0x5c, 0x3d,
	0x98, 0xb9, 0x7a, 0x8e, 0x79, 0x95, 0x41, 0x60, 0x15, 0x49, 0x6c, 0x60, 0x31, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x7f, 0x1e, 0x2d, 0x55, 0x90, 0x00, 0x00, 0x00,
}
