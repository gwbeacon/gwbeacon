// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lib/rpc/cluster.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	lib/rpc/cluster.proto
	lib/rpc/session.proto

It has these top-level messages:
	ServiceInfo
	ServerInfo
	ServerInfoIDMap
	RegisterReturn
	UserInfo
	ClientInfo
	Session
	SessionResult
	SessionStatRequest
	SessionStat
*/
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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ServiceInfo struct {
	Name    string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Version int32  `protobuf:"varint,2,opt,name=Version" json:"Version,omitempty"`
}

func (m *ServiceInfo) Reset()                    { *m = ServiceInfo{} }
func (m *ServiceInfo) String() string            { return proto.CompactTextString(m) }
func (*ServiceInfo) ProtoMessage()               {}
func (*ServiceInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServiceInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ServiceInfo) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type ServerInfo struct {
	Type     string         `protobuf:"bytes,1,opt,name=Type" json:"Type,omitempty"`
	ID       uint32         `protobuf:"varint,2,opt,name=ID" json:"ID,omitempty"`
	IP       string         `protobuf:"bytes,3,opt,name=IP" json:"IP,omitempty"`
	Port     int32          `protobuf:"varint,4,opt,name=Port" json:"Port,omitempty"`
	TimeBase int32          `protobuf:"varint,5,opt,name=TimeBase" json:"TimeBase,omitempty"`
	Services []*ServiceInfo `protobuf:"bytes,6,rep,name=Services" json:"Services,omitempty"`
}

func (m *ServerInfo) Reset()                    { *m = ServerInfo{} }
func (m *ServerInfo) String() string            { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()               {}
func (*ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServerInfo) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ServerInfo) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ServerInfo) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *ServerInfo) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *ServerInfo) GetTimeBase() int32 {
	if m != nil {
		return m.TimeBase
	}
	return 0
}

func (m *ServerInfo) GetServices() []*ServiceInfo {
	if m != nil {
		return m.Services
	}
	return nil
}

type ServerInfoIDMap struct {
	Servers map[uint32]*ServerInfo `protobuf:"bytes,1,rep,name=Servers" json:"Servers,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ServerInfoIDMap) Reset()                    { *m = ServerInfoIDMap{} }
func (m *ServerInfoIDMap) String() string            { return proto.CompactTextString(m) }
func (*ServerInfoIDMap) ProtoMessage()               {}
func (*ServerInfoIDMap) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ServerInfoIDMap) GetServers() map[uint32]*ServerInfo {
	if m != nil {
		return m.Servers
	}
	return nil
}

type RegisterReturn struct {
	ID          uint32                      `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	NewServers  map[string]*ServerInfoIDMap `protobuf:"bytes,2,rep,name=NewServers" json:"NewServers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	DownServers map[string]*ServerInfoIDMap `protobuf:"bytes,3,rep,name=DownServers" json:"DownServers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *RegisterReturn) Reset()                    { *m = RegisterReturn{} }
func (m *RegisterReturn) String() string            { return proto.CompactTextString(m) }
func (*RegisterReturn) ProtoMessage()               {}
func (*RegisterReturn) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RegisterReturn) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *RegisterReturn) GetNewServers() map[string]*ServerInfoIDMap {
	if m != nil {
		return m.NewServers
	}
	return nil
}

func (m *RegisterReturn) GetDownServers() map[string]*ServerInfoIDMap {
	if m != nil {
		return m.DownServers
	}
	return nil
}

func init() {
	proto.RegisterType((*ServiceInfo)(nil), "rpc.ServiceInfo")
	proto.RegisterType((*ServerInfo)(nil), "rpc.ServerInfo")
	proto.RegisterType((*ServerInfoIDMap)(nil), "rpc.ServerInfoIDMap")
	proto.RegisterType((*RegisterReturn)(nil), "rpc.RegisterReturn")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Cluster service

type ClusterClient interface {
	Register(ctx context.Context, opts ...grpc.CallOption) (Cluster_RegisterClient, error)
	Sync(ctx context.Context, in *ServerInfo, opts ...grpc.CallOption) (*RegisterReturn, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) Register(ctx context.Context, opts ...grpc.CallOption) (Cluster_RegisterClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Cluster_serviceDesc.Streams[0], c.cc, "/rpc.Cluster/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &clusterRegisterClient{stream}
	return x, nil
}

type Cluster_RegisterClient interface {
	Send(*ServerInfo) error
	Recv() (*RegisterReturn, error)
	grpc.ClientStream
}

type clusterRegisterClient struct {
	grpc.ClientStream
}

func (x *clusterRegisterClient) Send(m *ServerInfo) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clusterRegisterClient) Recv() (*RegisterReturn, error) {
	m := new(RegisterReturn)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *clusterClient) Sync(ctx context.Context, in *ServerInfo, opts ...grpc.CallOption) (*RegisterReturn, error) {
	out := new(RegisterReturn)
	err := grpc.Invoke(ctx, "/rpc.Cluster/Sync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Cluster service

type ClusterServer interface {
	Register(Cluster_RegisterServer) error
	Sync(context.Context, *ServerInfo) (*RegisterReturn, error)
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClusterServer).Register(&clusterRegisterServer{stream})
}

type Cluster_RegisterServer interface {
	Send(*RegisterReturn) error
	Recv() (*ServerInfo, error)
	grpc.ServerStream
}

type clusterRegisterServer struct {
	grpc.ServerStream
}

func (x *clusterRegisterServer) Send(m *RegisterReturn) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clusterRegisterServer) Recv() (*ServerInfo, error) {
	m := new(ServerInfo)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Cluster_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Cluster/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Sync(ctx, req.(*ServerInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _Cluster_Sync_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _Cluster_Register_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "lib/rpc/cluster.proto",
}

func init() { proto.RegisterFile("lib/rpc/cluster.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 408 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xdd, 0x6a, 0xa3, 0x40,
	0x14, 0xc7, 0x77, 0x34, 0x9f, 0xc7, 0xcd, 0x07, 0xb3, 0xbb, 0x30, 0x78, 0x95, 0x75, 0x77, 0x41,
	0x96, 0xc5, 0x2c, 0x29, 0x94, 0xd2, 0xdc, 0x35, 0xb6, 0x20, 0xa5, 0x21, 0x4c, 0x42, 0xef, 0x8d,
	0x4c, 0x8b, 0x34, 0x51, 0x3b, 0x9a, 0x04, 0x9f, 0xa5, 0xd0, 0xe7, 0xea, 0xe3, 0x14, 0xc7, 0x68,
	0x26, 0x92, 0x8b, 0x42, 0xef, 0xce, 0x99, 0x39, 0xff, 0xdf, 0x9c, 0xff, 0x39, 0x0a, 0x3f, 0x56,
	0xfe, 0x72, 0xc8, 0x23, 0x6f, 0xe8, 0xad, 0x36, 0x71, 0xc2, 0xb8, 0x15, 0xf1, 0x30, 0x09, 0xb1,
	0xca, 0x23, 0xcf, 0x18, 0x83, 0x36, 0x67, 0x7c, 0xeb, 0x7b, 0xcc, 0x09, 0x1e, 0x42, 0x8c, 0xa1,
	0x36, 0x75, 0xd7, 0x8c, 0xa0, 0x01, 0x32, 0xdb, 0x54, 0xc4, 0x98, 0x40, 0xf3, 0x9e, 0xf1, 0xd8,
	0x0f, 0x03, 0xa2, 0x0c, 0x90, 0x59, 0xa7, 0x45, 0x6a, 0xbc, 0x22, 0x80, 0x4c, 0xcd, 0x78, 0x21,
	0x5e, 0xa4, 0x51, 0x29, 0xce, 0x62, 0xdc, 0x05, 0xc5, 0xb1, 0x85, 0xae, 0x43, 0x15, 0xc7, 0x16,
	0xf9, 0x8c, 0xa8, 0xa2, 0x42, 0x71, 0x66, 0x99, 0x66, 0x16, 0xf2, 0x84, 0xd4, 0x04, 0x59, 0xc4,
	0x58, 0x87, 0xd6, 0xc2, 0x5f, 0xb3, 0x2b, 0x37, 0x66, 0xa4, 0x2e, 0xce, 0xcb, 0x1c, 0xff, 0x83,
	0xd6, 0xbe, 0xdf, 0x98, 0x34, 0x06, 0xaa, 0xa9, 0x8d, 0xfa, 0x16, 0x8f, 0x3c, 0x4b, 0x32, 0x41,
	0xcb, 0x0a, 0xe3, 0x05, 0x41, 0xef, 0xd0, 0xa0, 0x63, 0xdf, 0xb9, 0x11, 0x1e, 0x43, 0x33, 0x3f,
	0x8a, 0x09, 0x12, 0x80, 0x9f, 0x25, 0x40, 0x2a, 0xdb, 0xe7, 0xf1, 0x75, 0x90, 0xf0, 0x94, 0x16,
	0x0a, 0xfd, 0x16, 0xbe, 0xca, 0x17, 0xb8, 0x0f, 0xea, 0x13, 0x4b, 0x85, 0xe3, 0x0e, 0xcd, 0x42,
	0xfc, 0x07, 0xea, 0x5b, 0x77, 0xb5, 0x61, 0xc2, 0xb3, 0x36, 0xea, 0x55, 0xe0, 0x34, 0xbf, 0xbd,
	0x54, 0x2e, 0x90, 0xf1, 0xa6, 0x40, 0x97, 0xb2, 0x47, 0x3f, 0xdb, 0x09, 0x65, 0xc9, 0x86, 0x07,
	0xfb, 0x71, 0xa1, 0x72, 0x5c, 0x13, 0x80, 0x29, 0xdb, 0x15, 0xfd, 0x2a, 0xa2, 0xdf, 0x5f, 0x02,
	0x79, 0x2c, 0xb4, 0x0e, 0x55, 0x79, 0xc7, 0x92, 0x0c, 0xdf, 0x80, 0x66, 0x87, 0xbb, 0xa0, 0xa0,
	0xa8, 0x82, 0xf2, 0xfb, 0x14, 0x45, 0x2a, 0xcb, 0x31, 0xb2, 0x50, 0x9f, 0x43, 0xaf, 0xf2, 0x8c,
	0xec, 0xbf, 0x9d, 0xfb, 0xff, 0x7b, 0xec, 0xff, 0xfb, 0xa9, 0xe1, 0x4a, 0x43, 0xd0, 0x17, 0xd0,
	0xaf, 0xbe, 0xfa, 0x79, 0xea, 0xe8, 0x19, 0x9a, 0x93, 0xfc, 0x63, 0xc7, 0xe7, 0xd0, 0x2a, 0x5c,
	0xe2, 0xea, 0x36, 0xf4, 0x6f, 0x27, 0xa6, 0x60, 0x7c, 0x31, 0xd1, 0x7f, 0x84, 0x2d, 0xa8, 0xcd,
	0xd3, 0xc0, 0xfb, 0xa8, 0x66, 0xd9, 0x10, 0x7f, 0xd5, 0xd9, 0x7b, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x63, 0xeb, 0x56, 0xb5, 0x6e, 0x03, 0x00, 0x00,
}
