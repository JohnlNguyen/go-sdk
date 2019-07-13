// Code generated by protoc-gen-go. DO NOT EDIT.
// source: full.proto

package full

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusArgs) Reset()         { *m = StatusArgs{} }
func (m *StatusArgs) String() string { return proto.CompactTextString(m) }
func (*StatusArgs) ProtoMessage()    {}
func (*StatusArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbc2e675beec816d, []int{0}
}

func (m *StatusArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusArgs.Unmarshal(m, b)
}
func (m *StatusArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusArgs.Marshal(b, m, deterministic)
}
func (m *StatusArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusArgs.Merge(m, src)
}
func (m *StatusArgs) XXX_Size() int {
	return xxx_messageInfo_StatusArgs.Size(m)
}
func (m *StatusArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusArgs.DiscardUnknown(m)
}

var xxx_messageInfo_StatusArgs proto.InternalMessageInfo

type StatusResponse struct {
	Version              string   `protobuf:"bytes,10,opt,name=version,proto3" json:"version,omitempty"`
	GitRef               string   `protobuf:"bytes,11,opt,name=gitRef,proto3" json:"gitRef,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbc2e675beec816d, []int{1}
}

func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (m *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(m, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *StatusResponse) GetGitRef() string {
	if m != nil {
		return m.GitRef
	}
	return ""
}

func init() {
	proto.RegisterType((*StatusArgs)(nil), "full.StatusArgs")
	proto.RegisterType((*StatusResponse)(nil), "full.StatusResponse")
}

func init() { proto.RegisterFile("full.proto", fileDescriptor_bbc2e675beec816d) }

var fileDescriptor_bbc2e675beec816d = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2b, 0xcd, 0xc9,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x78, 0xb8, 0xb8, 0x82, 0x4b,
	0x12, 0x4b, 0x4a, 0x8b, 0x1d, 0x8b, 0xd2, 0x8b, 0x95, 0x9c, 0xb8, 0xf8, 0x20, 0xbc, 0xa0, 0xd4,
	0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x09, 0x2e, 0xf6, 0xb2, 0xd4, 0xa2, 0xe2, 0xcc, 0xfc,
	0x3c, 0x09, 0x2e, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0x48, 0x8c, 0x8b, 0x2d, 0x3d, 0xb3,
	0x24, 0x28, 0x35, 0x4d, 0x82, 0x1b, 0x2c, 0x01, 0xe5, 0x19, 0xd9, 0x70, 0xb1, 0x41, 0xcc, 0x10,
	0x32, 0x82, 0xb3, 0x04, 0xf4, 0xc0, 0x16, 0x23, 0x6c, 0x92, 0x12, 0x41, 0x16, 0x81, 0xd9, 0xa6,
	0xc4, 0x90, 0xc4, 0x06, 0x76, 0x9c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x63, 0xba, 0x71, 0x87,
	0xaa, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StatusClient is the client API for Status service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StatusClient interface {
	Status(ctx context.Context, in *StatusArgs, opts ...grpc.CallOption) (*StatusResponse, error)
}

type statusClient struct {
	cc *grpc.ClientConn
}

func NewStatusClient(cc *grpc.ClientConn) StatusClient {
	return &statusClient{cc}
}

func (c *statusClient) Status(ctx context.Context, in *StatusArgs, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/full.Status/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatusServer is the server API for Status service.
type StatusServer interface {
	Status(context.Context, *StatusArgs) (*StatusResponse, error)
}

func RegisterStatusServer(s *grpc.Server, srv StatusServer) {
	s.RegisterService(&_Status_serviceDesc, srv)
}

func _Status_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/full.Status/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServer).Status(ctx, req.(*StatusArgs))
	}
	return interceptor(ctx, in, info, handler)
}

var _Status_serviceDesc = grpc.ServiceDesc{
	ServiceName: "full.Status",
	HandlerType: (*StatusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _Status_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "full.proto",
}
