// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lib/rpc/session.proto

package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type UserInfo struct {
	Name      string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Domain    string `protobuf:"bytes,2,opt,name=Domain" json:"Domain,omitempty"`
	Device    string `protobuf:"bytes,3,opt,name=Device" json:"Device,omitempty"`
	NickName  string `protobuf:"bytes,4,opt,name=NickName" json:"NickName,omitempty"`
	LoginTime uint64 `protobuf:"varint,5,opt,name=LoginTime" json:"LoginTime,omitempty"`
}

func (m *UserInfo) Reset()                    { *m = UserInfo{} }
func (m *UserInfo) String() string            { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()               {}
func (*UserInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *UserInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UserInfo) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *UserInfo) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *UserInfo) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *UserInfo) GetLoginTime() uint64 {
	if m != nil {
		return m.LoginTime
	}
	return 0
}

type ClientInfo struct {
	ClientType    uint32 `protobuf:"varint,1,opt,name=ClientType" json:"ClientType,omitempty"`
	ClientVersion string `protobuf:"bytes,2,opt,name=ClientVersion" json:"ClientVersion,omitempty"`
	ProtoType     string `protobuf:"bytes,3,opt,name=ProtoType" json:"ProtoType,omitempty"`
	ProtoVersion  string `protobuf:"bytes,4,opt,name=ProtoVersion" json:"ProtoVersion,omitempty"`
	OSType        uint32 `protobuf:"varint,5,opt,name=OSType" json:"OSType,omitempty"`
	OSVersion     string `protobuf:"bytes,6,opt,name=OSVersion" json:"OSVersion,omitempty"`
	DeviceType    string `protobuf:"bytes,7,opt,name=DeviceType" json:"DeviceType,omitempty"`
	DeviceNumber  string `protobuf:"bytes,8,opt,name=DeviceNumber" json:"DeviceNumber,omitempty"`
}

func (m *ClientInfo) Reset()                    { *m = ClientInfo{} }
func (m *ClientInfo) String() string            { return proto.CompactTextString(m) }
func (*ClientInfo) ProtoMessage()               {}
func (*ClientInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *ClientInfo) GetClientType() uint32 {
	if m != nil {
		return m.ClientType
	}
	return 0
}

func (m *ClientInfo) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

func (m *ClientInfo) GetProtoType() string {
	if m != nil {
		return m.ProtoType
	}
	return ""
}

func (m *ClientInfo) GetProtoVersion() string {
	if m != nil {
		return m.ProtoVersion
	}
	return ""
}

func (m *ClientInfo) GetOSType() uint32 {
	if m != nil {
		return m.OSType
	}
	return 0
}

func (m *ClientInfo) GetOSVersion() string {
	if m != nil {
		return m.OSVersion
	}
	return ""
}

func (m *ClientInfo) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *ClientInfo) GetDeviceNumber() string {
	if m != nil {
		return m.DeviceNumber
	}
	return ""
}

type Session struct {
	ID     uint64            `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Addr   string            `protobuf:"bytes,2,opt,name=Addr" json:"Addr,omitempty"`
	User   *UserInfo         `protobuf:"bytes,3,opt,name=User" json:"User,omitempty"`
	Client *ClientInfo       `protobuf:"bytes,4,opt,name=Client" json:"Client,omitempty"`
	Extra  map[string]string `protobuf:"bytes,5,rep,name=Extra" json:"Extra,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Session) Reset()                    { *m = Session{} }
func (m *Session) String() string            { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()               {}
func (*Session) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Session) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Session) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *Session) GetUser() *UserInfo {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Session) GetClient() *ClientInfo {
	if m != nil {
		return m.Client
	}
	return nil
}

func (m *Session) GetExtra() map[string]string {
	if m != nil {
		return m.Extra
	}
	return nil
}

type SessionRequest struct {
	ID uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
}

func (m *SessionRequest) Reset()                    { *m = SessionRequest{} }
func (m *SessionRequest) String() string            { return proto.CompactTextString(m) }
func (*SessionRequest) ProtoMessage()               {}
func (*SessionRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *SessionRequest) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type SessionResult struct {
	Ok   bool       `protobuf:"varint,1,opt,name=Ok" json:"Ok,omitempty"`
	Data []*Session `protobuf:"bytes,2,rep,name=Data" json:"Data,omitempty"`
}

func (m *SessionResult) Reset()                    { *m = SessionResult{} }
func (m *SessionResult) String() string            { return proto.CompactTextString(m) }
func (*SessionResult) ProtoMessage()               {}
func (*SessionResult) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *SessionResult) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *SessionResult) GetData() []*Session {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*UserInfo)(nil), "rpc.UserInfo")
	proto.RegisterType((*ClientInfo)(nil), "rpc.ClientInfo")
	proto.RegisterType((*Session)(nil), "rpc.Session")
	proto.RegisterType((*SessionRequest)(nil), "rpc.SessionRequest")
	proto.RegisterType((*SessionResult)(nil), "rpc.SessionResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SessionStore service

type SessionStoreClient interface {
	Store(ctx context.Context, in *Session, opts ...grpc.CallOption) (*SessionResult, error)
	Remove(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*SessionResult, error)
	Get(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*Session, error)
	Find(ctx context.Context, in *UserInfo, opts ...grpc.CallOption) (*SessionResult, error)
}

type sessionStoreClient struct {
	cc *grpc.ClientConn
}

func NewSessionStoreClient(cc *grpc.ClientConn) SessionStoreClient {
	return &sessionStoreClient{cc}
}

func (c *sessionStoreClient) Store(ctx context.Context, in *Session, opts ...grpc.CallOption) (*SessionResult, error) {
	out := new(SessionResult)
	err := grpc.Invoke(ctx, "/rpc.SessionStore/Store", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionStoreClient) Remove(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*SessionResult, error) {
	out := new(SessionResult)
	err := grpc.Invoke(ctx, "/rpc.SessionStore/Remove", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionStoreClient) Get(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := grpc.Invoke(ctx, "/rpc.SessionStore/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionStoreClient) Find(ctx context.Context, in *UserInfo, opts ...grpc.CallOption) (*SessionResult, error) {
	out := new(SessionResult)
	err := grpc.Invoke(ctx, "/rpc.SessionStore/Find", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SessionStore service

type SessionStoreServer interface {
	Store(context.Context, *Session) (*SessionResult, error)
	Remove(context.Context, *SessionRequest) (*SessionResult, error)
	Get(context.Context, *SessionRequest) (*Session, error)
	Find(context.Context, *UserInfo) (*SessionResult, error)
}

func RegisterSessionStoreServer(s *grpc.Server, srv SessionStoreServer) {
	s.RegisterService(&_SessionStore_serviceDesc, srv)
}

func _SessionStore_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionStoreServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SessionStore/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionStoreServer).Store(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionStore_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionStoreServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SessionStore/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionStoreServer).Remove(ctx, req.(*SessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionStore_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionStoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SessionStore/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionStoreServer).Get(ctx, req.(*SessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionStore_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionStoreServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SessionStore/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionStoreServer).Find(ctx, req.(*UserInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _SessionStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.SessionStore",
	HandlerType: (*SessionStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Store",
			Handler:    _SessionStore_Store_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _SessionStore_Remove_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SessionStore_Get_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _SessionStore_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lib/rpc/session.proto",
}

func init() { proto.RegisterFile("lib/rpc/session.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 489 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0xaf, 0x4d, 0x27, 0x49, 0x41, 0xc3, 0xcd, 0x8a, 0x10, 0x32, 0x16, 0x12, 0x11, 0x15,
	0xa9, 0x94, 0xbe, 0x54, 0xbc, 0x55, 0xa4, 0xa0, 0x48, 0x28, 0x41, 0x9b, 0xc2, 0xbb, 0xe3, 0x2c,
	0x68, 0x95, 0xf8, 0xc2, 0xda, 0x89, 0xc8, 0x1f, 0xf0, 0xc4, 0x7f, 0xf1, 0x1b, 0x7c, 0x09, 0xda,
	0xd9, 0x4d, 0x62, 0xa3, 0xf4, 0x6d, 0xcf, 0x99, 0xdb, 0x39, 0x33, 0x36, 0x3c, 0x59, 0x89, 0xf9,
	0xa5, 0x2c, 0x92, 0xcb, 0x92, 0x97, 0xa5, 0xc8, 0xb3, 0x41, 0x21, 0xf3, 0x2a, 0x47, 0x47, 0x16,
	0x49, 0xf4, 0xcb, 0x82, 0xd6, 0x97, 0x92, 0xcb, 0x71, 0xf6, 0x2d, 0x47, 0x04, 0x77, 0x12, 0xa7,
	0x3c, 0xb0, 0x42, 0xab, 0x7f, 0xc6, 0xe8, 0x8d, 0x4f, 0xc1, 0x1f, 0xe5, 0x69, 0x2c, 0xb2, 0xc0,
	0x26, 0xd6, 0x20, 0xe2, 0xf9, 0x46, 0x24, 0x3c, 0x70, 0x0c, 0x4f, 0x08, 0x7b, 0xd0, 0x9a, 0x88,
	0x64, 0x49, 0x7d, 0x5c, 0x8a, 0xec, 0x31, 0x3e, 0x87, 0xb3, 0x4f, 0xf9, 0x77, 0x91, 0xdd, 0x89,
	0x94, 0x07, 0x5e, 0x68, 0xf5, 0x5d, 0x76, 0x20, 0xa2, 0xdf, 0x36, 0xc0, 0xfb, 0x95, 0xe0, 0x59,
	0x45, 0x62, 0x5e, 0xec, 0xd0, 0xdd, 0xb6, 0xd0, 0x92, 0xba, 0xac, 0xc6, 0xe0, 0x2b, 0xe8, 0x6a,
	0xf4, 0x95, 0x4b, 0xe5, 0xca, 0xe8, 0x6b, 0x92, 0x6a, 0xe4, 0x67, 0xe5, 0x96, 0x9a, 0x68, 0xa5,
	0x07, 0x02, 0x23, 0xe8, 0x10, 0xd8, 0xb5, 0xd0, 0x82, 0x1b, 0x9c, 0x32, 0x3a, 0x9d, 0x51, 0xb9,
	0x47, 0x1a, 0x0c, 0x52, 0x9d, 0xa7, 0xb3, 0x5d, 0xa1, 0xaf, 0x3b, 0xef, 0x09, 0xa5, 0x5e, 0x2f,
	0x84, 0x2a, 0x4f, 0x29, 0x5c, 0x63, 0xd4, 0x64, 0x8d, 0x26, 0xeb, 0x74, 0xce, 0x65, 0xd0, 0xd2,
	0x93, 0xeb, 0x5c, 0xf4, 0xd7, 0x82, 0xd3, 0x99, 0x3e, 0x19, 0x9e, 0x83, 0x3d, 0x1e, 0xd1, 0x16,
	0x5c, 0x66, 0x8f, 0x47, 0xea, 0x54, 0x37, 0x8b, 0x85, 0x34, 0xa6, 0xe9, 0x8d, 0x2f, 0xc1, 0x55,
	0xa7, 0x24, 0x9b, 0xed, 0x61, 0x77, 0x20, 0x8b, 0x64, 0xb0, 0xbb, 0x2d, 0xa3, 0x10, 0xbe, 0x06,
	0x5f, 0xef, 0x87, 0xac, 0xb6, 0x87, 0x0f, 0x28, 0xe9, 0xb0, 0x75, 0x66, 0xc2, 0xf8, 0x16, 0xbc,
	0xdb, 0x9f, 0x95, 0x8c, 0x03, 0x2f, 0x74, 0xfa, 0xed, 0xe1, 0x33, 0xca, 0x33, 0x62, 0x06, 0x14,
	0xb9, 0xcd, 0x2a, 0xb9, 0x65, 0x3a, 0xab, 0x77, 0x0d, 0x70, 0x20, 0xf1, 0x21, 0x38, 0x4b, 0xbe,
	0x35, 0x9f, 0x91, 0x7a, 0xe2, 0x63, 0xf0, 0x36, 0xf1, 0x6a, 0xcd, 0x8d, 0x5e, 0x0d, 0xde, 0xd9,
	0xd7, 0x56, 0x14, 0xc2, 0xb9, 0x69, 0xcb, 0xf8, 0x8f, 0x35, 0x2f, 0xab, 0xff, 0xad, 0x46, 0x37,
	0xd0, 0xdd, 0x67, 0x94, 0xeb, 0x15, 0x25, 0x4c, 0x97, 0x94, 0xd0, 0x62, 0xf6, 0x74, 0x89, 0x21,
	0xb8, 0xa3, 0xb8, 0x8a, 0x03, 0x9b, 0xa4, 0x76, 0xea, 0x52, 0x19, 0x45, 0x86, 0x7f, 0x2c, 0xe8,
	0x18, 0x66, 0x56, 0xe5, 0x92, 0xe3, 0x05, 0x78, 0xfa, 0xd1, 0xc8, 0xee, 0x61, 0xa3, 0x96, 0xa6,
	0x45, 0x27, 0x78, 0x05, 0x3e, 0xe3, 0x69, 0xbe, 0xe1, 0xf8, 0xa8, 0x19, 0x27, 0xbd, 0xf7, 0x14,
	0xbd, 0x01, 0xe7, 0x23, 0xaf, 0x8e, 0x57, 0x34, 0x86, 0x46, 0x27, 0x78, 0x01, 0xee, 0x07, 0x91,
	0x2d, 0xb0, 0x79, 0xb2, 0xe3, 0x8d, 0xe7, 0x3e, 0xfd, 0xbd, 0x57, 0xff, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x4e, 0x5c, 0xb1, 0x49, 0xd6, 0x03, 0x00, 0x00,
}
