// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: api/proto/server.proto

package server

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

// SecretKeeperClient is the client API for SecretKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretKeeperClient interface {
	Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	GetAllNames(ctx context.Context, in *GetAllNamesRequest, opts ...grpc.CallOption) (*GetAllNamesResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
}

type secretKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretKeeperClient(cc grpc.ClientConnInterface) SecretKeeperClient {
	return &secretKeeperClient{cc}
}

func (c *secretKeeperClient) Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretKeeperClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretKeeperClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretKeeperClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretKeeperClient) GetAllNames(ctx context.Context, in *GetAllNamesRequest, opts ...grpc.CallOption) (*GetAllNamesResponse, error) {
	out := new(GetAllNamesResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/GetAllNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretKeeperClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/api.SecretKeeper/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretKeeperServer is the server API for SecretKeeper service.
// All implementations must embed UnimplementedSecretKeeperServer
// for forward compatibility
type SecretKeeperServer interface {
	Auth(context.Context, *AuthRequest) (*AuthResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	GetAllNames(context.Context, *GetAllNamesRequest) (*GetAllNamesResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	mustEmbedUnimplementedSecretKeeperServer()
}

// UnimplementedSecretKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedSecretKeeperServer struct {
}

func (UnimplementedSecretKeeperServer) Auth(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (UnimplementedSecretKeeperServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedSecretKeeperServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSecretKeeperServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedSecretKeeperServer) GetAllNames(context.Context, *GetAllNamesRequest) (*GetAllNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNames not implemented")
}
func (UnimplementedSecretKeeperServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedSecretKeeperServer) mustEmbedUnimplementedSecretKeeperServer() {}

// UnsafeSecretKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretKeeperServer will
// result in compilation errors.
type UnsafeSecretKeeperServer interface {
	mustEmbedUnimplementedSecretKeeperServer()
}

func RegisterSecretKeeperServer(s grpc.ServiceRegistrar, srv SecretKeeperServer) {
	s.RegisterService(&SecretKeeper_ServiceDesc, srv)
}

func _SecretKeeper_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).Auth(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretKeeper_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretKeeper_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretKeeper_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretKeeper_GetAllNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).GetAllNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/GetAllNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).GetAllNames(ctx, req.(*GetAllNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretKeeper_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretKeeperServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SecretKeeper/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretKeeperServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretKeeper_ServiceDesc is the grpc.ServiceDesc for SecretKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.SecretKeeper",
	HandlerType: (*SecretKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _SecretKeeper_Auth_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _SecretKeeper_Register_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SecretKeeper_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _SecretKeeper_Delete_Handler,
		},
		{
			MethodName: "GetAllNames",
			Handler:    _SecretKeeper_GetAllNames_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _SecretKeeper_Set_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/server.proto",
}
