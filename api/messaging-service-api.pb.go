// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messaging-service-api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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

type MessagingRequest struct {
	UserId               []string                   `protobuf:"bytes,1,rep,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Template             string                     `protobuf:"bytes,2,opt,name=template,proto3" json:"template,omitempty"`
	Data                 string                     `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Override             *MessagingRequest_Override `protobuf:"bytes,4,opt,name=override,proto3" json:"override,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *MessagingRequest) Reset()         { *m = MessagingRequest{} }
func (m *MessagingRequest) String() string { return proto.CompactTextString(m) }
func (*MessagingRequest) ProtoMessage()    {}
func (*MessagingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d683ee4dabdf70fd, []int{0}
}

func (m *MessagingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessagingRequest.Unmarshal(m, b)
}
func (m *MessagingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessagingRequest.Marshal(b, m, deterministic)
}
func (m *MessagingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessagingRequest.Merge(m, src)
}
func (m *MessagingRequest) XXX_Size() int {
	return xxx_messageInfo_MessagingRequest.Size(m)
}
func (m *MessagingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessagingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessagingRequest proto.InternalMessageInfo

func (m *MessagingRequest) GetUserId() []string {
	if m != nil {
		return m.UserId
	}
	return nil
}

func (m *MessagingRequest) GetTemplate() string {
	if m != nil {
		return m.Template
	}
	return ""
}

func (m *MessagingRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *MessagingRequest) GetOverride() *MessagingRequest_Override {
	if m != nil {
		return m.Override
	}
	return nil
}

type MessagingRequest_Override struct {
	OverridePref         bool     `protobuf:"varint,1,opt,name=override_pref,json=overridePref,proto3" json:"override_pref,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessagingRequest_Override) Reset()         { *m = MessagingRequest_Override{} }
func (m *MessagingRequest_Override) String() string { return proto.CompactTextString(m) }
func (*MessagingRequest_Override) ProtoMessage()    {}
func (*MessagingRequest_Override) Descriptor() ([]byte, []int) {
	return fileDescriptor_d683ee4dabdf70fd, []int{0, 0}
}

func (m *MessagingRequest_Override) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessagingRequest_Override.Unmarshal(m, b)
}
func (m *MessagingRequest_Override) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessagingRequest_Override.Marshal(b, m, deterministic)
}
func (m *MessagingRequest_Override) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessagingRequest_Override.Merge(m, src)
}
func (m *MessagingRequest_Override) XXX_Size() int {
	return xxx_messageInfo_MessagingRequest_Override.Size(m)
}
func (m *MessagingRequest_Override) XXX_DiscardUnknown() {
	xxx_messageInfo_MessagingRequest_Override.DiscardUnknown(m)
}

var xxx_messageInfo_MessagingRequest_Override proto.InternalMessageInfo

func (m *MessagingRequest_Override) GetOverridePref() bool {
	if m != nil {
		return m.OverridePref
	}
	return false
}

func (m *MessagingRequest_Override) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func init() {
	proto.RegisterType((*MessagingRequest)(nil), "influenzanet.messaging_service_api.MessagingRequest")
	proto.RegisterType((*MessagingRequest_Override)(nil), "influenzanet.messaging_service_api.MessagingRequest.Override")
}

func init() { proto.RegisterFile("messaging-service-api.proto", fileDescriptor_d683ee4dabdf70fd) }

var fileDescriptor_d683ee4dabdf70fd = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x50, 0x4f, 0x4b, 0xc3, 0x30,
	0x1c, 0xa5, 0x6e, 0xd6, 0x2e, 0x2a, 0x4a, 0x94, 0x59, 0xba, 0x4b, 0x99, 0x97, 0x5e, 0x96, 0xc1,
	0x14, 0x6f, 0x1e, 0x14, 0x44, 0x76, 0x10, 0xa5, 0x3b, 0x29, 0x48, 0xc9, 0xec, 0xaf, 0x35, 0xd0,
	0x36, 0x31, 0xf9, 0x75, 0x30, 0x3f, 0x93, 0xdf, 0x51, 0x59, 0xff, 0x81, 0x22, 0x08, 0xde, 0xf2,
	0x5e, 0x5e, 0x5e, 0xde, 0x7b, 0x64, 0x94, 0x83, 0x31, 0x3c, 0x15, 0x45, 0x3a, 0x31, 0xa0, 0x57,
	0xe2, 0x05, 0x26, 0x5c, 0x09, 0xa6, 0xb4, 0x44, 0x49, 0xc7, 0xa2, 0x48, 0xb2, 0x12, 0x8a, 0x77,
	0x5e, 0x00, 0xb2, 0x4e, 0x19, 0x35, 0xca, 0x88, 0x2b, 0xe1, 0xd1, 0x34, 0x93, 0x4b, 0x9e, 0x4d,
	0x70, 0xad, 0xc0, 0xd4, 0xef, 0xbc, 0x51, 0x2a, 0x65, 0x9a, 0xc1, 0xb4, 0x42, 0xcb, 0x32, 0x99,
	0x42, 0xae, 0x70, 0x5d, 0x5f, 0x8e, 0x3f, 0x2d, 0x72, 0x78, 0xd7, 0x5a, 0x85, 0xf0, 0x56, 0x82,
	0x41, 0x7a, 0x42, 0x76, 0x4a, 0x03, 0x3a, 0x12, 0xb1, 0x6b, 0xf9, 0xbd, 0x60, 0x10, 0xda, 0x1b,
	0x38, 0x8f, 0xa9, 0x47, 0x1c, 0x84, 0x5c, 0x65, 0x1c, 0xc1, 0xdd, 0xf2, 0xad, 0x60, 0x10, 0x76,
	0x98, 0x52, 0xd2, 0x8f, 0x39, 0x72, 0xb7, 0x57, 0xf1, 0xd5, 0x99, 0x3e, 0x12, 0x47, 0xae, 0x40,
	0x6b, 0x11, 0x83, 0xdb, 0xf7, 0xad, 0x60, 0x77, 0x76, 0xc9, 0xfe, 0x6e, 0xc1, 0x7e, 0x06, 0x62,
	0xf7, 0x8d, 0x49, 0xd8, 0xd9, 0x79, 0xb7, 0xc4, 0x69, 0x59, 0x7a, 0x4a, 0xf6, 0x5b, 0x3e, 0x52,
	0x1a, 0x12, 0xd7, 0xf2, 0xad, 0xc0, 0x09, 0xf7, 0x5a, 0xf2, 0x41, 0x43, 0x42, 0x87, 0xc4, 0xce,
	0x01, 0x5f, 0x65, 0xdc, 0x24, 0x6f, 0xd0, 0xec, 0xc3, 0x22, 0x47, 0xdd, 0x87, 0x8b, 0x3a, 0xc5,
	0x95, 0x12, 0xf4, 0x82, 0xd8, 0x0b, 0xe4, 0x58, 0x1a, 0x3a, 0x64, 0xf5, 0x82, 0xac, 0x5d, 0x90,
	0xdd, 0x6c, 0x16, 0xf4, 0x8e, 0xbf, 0x77, 0x69, 0xd4, 0xcf, 0xe4, 0x60, 0x5e, 0x08, 0x14, 0x1c,
	0xa1, 0xb6, 0x05, 0x7a, 0xfe, 0x9f, 0xd2, 0xbf, 0xdb, 0x5f, 0x6f, 0x3f, 0xf5, 0xb8, 0x12, 0x4b,
	0xbb, 0xca, 0x72, 0xf6, 0x15, 0x00, 0x00, 0xff, 0xff, 0x77, 0x1b, 0xe0, 0xe7, 0x32, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MessagingServiceApiClient is the client API for MessagingServiceApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessagingServiceApiClient interface {
	Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Status, error)
	InitiateMessage(ctx context.Context, in *MessagingRequest, opts ...grpc.CallOption) (*Status, error)
}

type messagingServiceApiClient struct {
	cc *grpc.ClientConn
}

func NewMessagingServiceApiClient(cc *grpc.ClientConn) MessagingServiceApiClient {
	return &messagingServiceApiClient{cc}
}

func (c *messagingServiceApiClient) Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/influenzanet.messaging_service_api.MessagingServiceApi/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagingServiceApiClient) InitiateMessage(ctx context.Context, in *MessagingRequest, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/influenzanet.messaging_service_api.MessagingServiceApi/InitiateMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessagingServiceApiServer is the server API for MessagingServiceApi service.
type MessagingServiceApiServer interface {
	Status(context.Context, *empty.Empty) (*Status, error)
	InitiateMessage(context.Context, *MessagingRequest) (*Status, error)
}

// UnimplementedMessagingServiceApiServer can be embedded to have forward compatible implementations.
type UnimplementedMessagingServiceApiServer struct {
}

func (*UnimplementedMessagingServiceApiServer) Status(ctx context.Context, req *empty.Empty) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (*UnimplementedMessagingServiceApiServer) InitiateMessage(ctx context.Context, req *MessagingRequest) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitiateMessage not implemented")
}

func RegisterMessagingServiceApiServer(s *grpc.Server, srv MessagingServiceApiServer) {
	s.RegisterService(&_MessagingServiceApi_serviceDesc, srv)
}

func _MessagingServiceApi_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagingServiceApiServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/influenzanet.messaging_service_api.MessagingServiceApi/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagingServiceApiServer).Status(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessagingServiceApi_InitiateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessagingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagingServiceApiServer).InitiateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/influenzanet.messaging_service_api.MessagingServiceApi/InitiateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagingServiceApiServer).InitiateMessage(ctx, req.(*MessagingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MessagingServiceApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "influenzanet.messaging_service_api.MessagingServiceApi",
	HandlerType: (*MessagingServiceApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _MessagingServiceApi_Status_Handler,
		},
		{
			MethodName: "InitiateMessage",
			Handler:    _MessagingServiceApi_InitiateMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messaging-service-api.proto",
}