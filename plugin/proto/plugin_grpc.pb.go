// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginClient interface {
	HealthProbe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthStatus, error)
	Describe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginDescription, error)
	GetActions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ActionList, error)
	ExecuteAction(ctx context.Context, in *ExecuteActionRequest, opts ...grpc.CallOption) (*ExecuteActionResponse, error)
	TestCredentials(ctx context.Context, in *TestCredentialsRequest, opts ...grpc.CallOption) (*TestCredentialsResponse, error)
	GetAssets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Assets, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) HealthProbe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthStatus, error) {
	out := new(HealthStatus)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/HealthProbe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Describe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginDescription, error) {
	out := new(PluginDescription)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/Describe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetActions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ActionList, error) {
	out := new(ActionList)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/GetActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) ExecuteAction(ctx context.Context, in *ExecuteActionRequest, opts ...grpc.CallOption) (*ExecuteActionResponse, error) {
	out := new(ExecuteActionResponse)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/ExecuteAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) TestCredentials(ctx context.Context, in *TestCredentialsRequest, opts ...grpc.CallOption) (*TestCredentialsResponse, error) {
	out := new(TestCredentialsResponse)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/TestCredentials", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetAssets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Assets, error) {
	out := new(Assets)
	err := c.cc.Invoke(ctx, "/integration_pack.Plugin/GetAssets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
// All implementations must embed UnimplementedPluginServer
// for forward compatibility
type PluginServer interface {
	HealthProbe(context.Context, *Empty) (*HealthStatus, error)
	Describe(context.Context, *Empty) (*PluginDescription, error)
	GetActions(context.Context, *Empty) (*ActionList, error)
	ExecuteAction(context.Context, *ExecuteActionRequest) (*ExecuteActionResponse, error)
	TestCredentials(context.Context, *TestCredentialsRequest) (*TestCredentialsResponse, error)
	GetAssets(context.Context, *Empty) (*Assets, error)
	mustEmbedUnimplementedPluginServer()
}

// UnimplementedPluginServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (UnimplementedPluginServer) HealthProbe(context.Context, *Empty) (*HealthStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthProbe not implemented")
}
func (UnimplementedPluginServer) Describe(context.Context, *Empty) (*PluginDescription, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Describe not implemented")
}
func (UnimplementedPluginServer) GetActions(context.Context, *Empty) (*ActionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActions not implemented")
}
func (UnimplementedPluginServer) ExecuteAction(context.Context, *ExecuteActionRequest) (*ExecuteActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteAction not implemented")
}
func (UnimplementedPluginServer) TestCredentials(context.Context, *TestCredentialsRequest) (*TestCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestCredentials not implemented")
}
func (UnimplementedPluginServer) GetAssets(context.Context, *Empty) (*Assets, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAssets not implemented")
}
func (UnimplementedPluginServer) mustEmbedUnimplementedPluginServer() {}

// UnsafePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServer will
// result in compilation errors.
type UnsafePluginServer interface {
	mustEmbedUnimplementedPluginServer()
}

func RegisterPluginServer(s grpc.ServiceRegistrar, srv PluginServer) {
	s.RegisterService(&Plugin_ServiceDesc, srv)
}

func _Plugin_HealthProbe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).HealthProbe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/HealthProbe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).HealthProbe(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Describe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Describe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/Describe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Describe(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/GetActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetActions(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_ExecuteAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).ExecuteAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/ExecuteAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).ExecuteAction(ctx, req.(*ExecuteActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_TestCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).TestCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/TestCredentials",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).TestCredentials(ctx, req.(*TestCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetAssets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integration_pack.Plugin/GetAssets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetAssets(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Plugin_ServiceDesc is the grpc.ServiceDesc for Plugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Plugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "integration_pack.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthProbe",
			Handler:    _Plugin_HealthProbe_Handler,
		},
		{
			MethodName: "Describe",
			Handler:    _Plugin_Describe_Handler,
		},
		{
			MethodName: "GetActions",
			Handler:    _Plugin_GetActions_Handler,
		},
		{
			MethodName: "ExecuteAction",
			Handler:    _Plugin_ExecuteAction_Handler,
		},
		{
			MethodName: "TestCredentials",
			Handler:    _Plugin_TestCredentials_Handler,
		},
		{
			MethodName: "GetAssets",
			Handler:    _Plugin_GetAssets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin/proto/plugin.proto",
}
