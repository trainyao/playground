// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hostname.proto

package trainyao_hostname

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type HostnameRequest struct {
	Test                 string   `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostnameRequest) Reset()         { *m = HostnameRequest{} }
func (m *HostnameRequest) String() string { return proto.CompactTextString(m) }
func (*HostnameRequest) ProtoMessage()    {}
func (*HostnameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_412c89e20cb3d4bd, []int{0}
}

func (m *HostnameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostnameRequest.Unmarshal(m, b)
}
func (m *HostnameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostnameRequest.Marshal(b, m, deterministic)
}
func (m *HostnameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostnameRequest.Merge(m, src)
}
func (m *HostnameRequest) XXX_Size() int {
	return xxx_messageInfo_HostnameRequest.Size(m)
}
func (m *HostnameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HostnameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HostnameRequest proto.InternalMessageInfo

func (m *HostnameRequest) GetTest() string {
	if m != nil {
		return m.Test
	}
	return ""
}

type HostnameResponse struct {
	Hostname             string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostnameResponse) Reset()         { *m = HostnameResponse{} }
func (m *HostnameResponse) String() string { return proto.CompactTextString(m) }
func (*HostnameResponse) ProtoMessage()    {}
func (*HostnameResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_412c89e20cb3d4bd, []int{1}
}

func (m *HostnameResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostnameResponse.Unmarshal(m, b)
}
func (m *HostnameResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostnameResponse.Marshal(b, m, deterministic)
}
func (m *HostnameResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostnameResponse.Merge(m, src)
}
func (m *HostnameResponse) XXX_Size() int {
	return xxx_messageInfo_HostnameResponse.Size(m)
}
func (m *HostnameResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HostnameResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HostnameResponse proto.InternalMessageInfo

func (m *HostnameResponse) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func init() {
	proto.RegisterType((*HostnameRequest)(nil), "trainyao.hostname.HostnameRequest")
	proto.RegisterType((*HostnameResponse)(nil), "trainyao.hostname.HostnameResponse")
}

func init() { proto.RegisterFile("hostname.proto", fileDescriptor_412c89e20cb3d4bd) }

var fileDescriptor_412c89e20cb3d4bd = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xc8, 0x2f, 0x2e,
	0xc9, 0x4b, 0xcc, 0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2c, 0x29, 0x4a, 0xcc,
	0xcc, 0xab, 0x4c, 0xcc, 0xd7, 0x83, 0x49, 0x28, 0xa9, 0x72, 0xf1, 0x7b, 0x40, 0xd9, 0x41, 0xa9,
	0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c, 0x2c, 0x25, 0xa9, 0xc5, 0x25, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92, 0x1e, 0x97, 0x00, 0x42, 0x59, 0x71, 0x41, 0x7e, 0x5e,
	0x71, 0xaa, 0x90, 0x14, 0x17, 0x07, 0xcc, 0x18, 0xa8, 0x5a, 0x38, 0xdf, 0x28, 0x85, 0x8b, 0x03,
	0xa6, 0x5e, 0x28, 0x82, 0x8b, 0xdb, 0x3d, 0xb5, 0x04, 0xce, 0x55, 0xd2, 0xc3, 0x70, 0x85, 0x1e,
	0x9a, 0x13, 0xa4, 0x94, 0xf1, 0xaa, 0x81, 0xd8, 0xaf, 0xc4, 0x90, 0xc4, 0x06, 0xf6, 0x96, 0x31,
	0x20, 0x00, 0x00, 0xff, 0xff, 0x3b, 0x65, 0x39, 0x4d, 0xe8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HostnameClient is the client API for Hostname service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HostnameClient interface {
	GetHostname(ctx context.Context, in *HostnameRequest, opts ...grpc.CallOption) (*HostnameResponse, error)
}

type hostnameClient struct {
	cc *grpc.ClientConn
}

func NewHostnameClient(cc *grpc.ClientConn) HostnameClient {
	return &hostnameClient{cc}
}

func (c *hostnameClient) GetHostname(ctx context.Context, in *HostnameRequest, opts ...grpc.CallOption) (*HostnameResponse, error) {
	out := new(HostnameResponse)
	err := c.cc.Invoke(ctx, "/trainyao.hostname.Hostname/GetHostname", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HostnameServer is the server API for Hostname service.
type HostnameServer interface {
	GetHostname(context.Context, *HostnameRequest) (*HostnameResponse, error)
}

// UnimplementedHostnameServer can be embedded to have forward compatible implementations.
type UnimplementedHostnameServer struct {
}

func (*UnimplementedHostnameServer) GetHostname(ctx context.Context, req *HostnameRequest) (*HostnameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHostname not implemented")
}

func RegisterHostnameServer(s *grpc.Server, srv HostnameServer) {
	s.RegisterService(&_Hostname_serviceDesc, srv)
}

func _Hostname_GetHostname_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostnameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HostnameServer).GetHostname(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trainyao.hostname.Hostname/GetHostname",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HostnameServer).GetHostname(ctx, req.(*HostnameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Hostname_serviceDesc = grpc.ServiceDesc{
	ServiceName: "trainyao.hostname.Hostname",
	HandlerType: (*HostnameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHostname",
			Handler:    _Hostname_GetHostname_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hostname.proto",
}
