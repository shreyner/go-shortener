// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: proto/shortener.proto

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	CreateShort(ctx context.Context, in *CreateShortRequest, opts ...grpc.CallOption) (*CreateShortResponse, error)
	CreateBatchShort(ctx context.Context, in *CreateBatchShortRequest, opts ...grpc.CallOption) (*CreateBatchShortResponse, error)
	ListUserURLs(ctx context.Context, in *ListUserURLsRequest, opts ...grpc.CallOption) (*ListUserURLsResponse, error)
	DeleteByIDs(ctx context.Context, in *DeleteByIDsRequest, opts ...grpc.CallOption) (*DeleteByIDsResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) CreateShort(ctx context.Context, in *CreateShortRequest, opts ...grpc.CallOption) (*CreateShortResponse, error) {
	out := new(CreateShortResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/CreateShort", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) CreateBatchShort(ctx context.Context, in *CreateBatchShortRequest, opts ...grpc.CallOption) (*CreateBatchShortResponse, error) {
	out := new(CreateBatchShortResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/CreateBatchShort", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) ListUserURLs(ctx context.Context, in *ListUserURLsRequest, opts ...grpc.CallOption) (*ListUserURLsResponse, error) {
	out := new(ListUserURLsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/ListUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteByIDs(ctx context.Context, in *DeleteByIDsRequest, opts ...grpc.CallOption) (*DeleteByIDsResponse, error) {
	out := new(DeleteByIDsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/DeleteByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	CreateShort(context.Context, *CreateShortRequest) (*CreateShortResponse, error)
	CreateBatchShort(context.Context, *CreateBatchShortRequest) (*CreateBatchShortResponse, error)
	ListUserURLs(context.Context, *ListUserURLsRequest) (*ListUserURLsResponse, error)
	DeleteByIDs(context.Context, *DeleteByIDsRequest) (*DeleteByIDsResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) CreateShort(context.Context, *CreateShortRequest) (*CreateShortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShort not implemented")
}
func (UnimplementedShortenerServer) CreateBatchShort(context.Context, *CreateBatchShortRequest) (*CreateBatchShortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBatchShort not implemented")
}
func (UnimplementedShortenerServer) ListUserURLs(context.Context, *ListUserURLsRequest) (*ListUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserURLs not implemented")
}
func (UnimplementedShortenerServer) DeleteByIDs(context.Context, *DeleteByIDsRequest) (*DeleteByIDsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteByIDs not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_CreateShort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateShort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/CreateShort",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateShort(ctx, req.(*CreateShortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_CreateBatchShort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBatchShortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateBatchShort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/CreateBatchShort",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateBatchShort(ctx, req.(*CreateBatchShortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_ListUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).ListUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/ListUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).ListUserURLs(ctx, req.(*ListUserURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteByIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/DeleteByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteByIDs(ctx, req.(*DeleteByIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShort",
			Handler:    _Shortener_CreateShort_Handler,
		},
		{
			MethodName: "CreateBatchShort",
			Handler:    _Shortener_CreateBatchShort_Handler,
		},
		{
			MethodName: "ListUserURLs",
			Handler:    _Shortener_ListUserURLs_Handler,
		},
		{
			MethodName: "DeleteByIDs",
			Handler:    _Shortener_DeleteByIDs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}
