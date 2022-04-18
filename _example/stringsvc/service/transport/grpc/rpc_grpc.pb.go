// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	Buy(ctx context.Context, in *BuyRequest, opts ...grpc.CallOption) (*BuyResponse, error)
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	Join2(ctx context.Context, in *Join2Request, opts ...grpc.CallOption) (*Join2Response, error)
	Join3(ctx context.Context, in *Join3Request, opts ...grpc.CallOption) (*Join3Response, error)
	Lowercase(ctx context.Context, in *LowercaseRequest, opts ...grpc.CallOption) (*LowercaseResponse, error)
	Uppercase(ctx context.Context, in *UppercaseRequest, opts ...grpc.CallOption) (*UppercaseResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Buy(ctx context.Context, in *BuyRequest, opts ...grpc.CallOption) (*BuyResponse, error) {
	out := new(BuyResponse)
	err := c.cc.Invoke(ctx, "/Service/Buy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, "/Service/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Join2(ctx context.Context, in *Join2Request, opts ...grpc.CallOption) (*Join2Response, error) {
	out := new(Join2Response)
	err := c.cc.Invoke(ctx, "/Service/Join2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Join3(ctx context.Context, in *Join3Request, opts ...grpc.CallOption) (*Join3Response, error) {
	out := new(Join3Response)
	err := c.cc.Invoke(ctx, "/Service/Join3", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Lowercase(ctx context.Context, in *LowercaseRequest, opts ...grpc.CallOption) (*LowercaseResponse, error) {
	out := new(LowercaseResponse)
	err := c.cc.Invoke(ctx, "/Service/Lowercase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Uppercase(ctx context.Context, in *UppercaseRequest, opts ...grpc.CallOption) (*UppercaseResponse, error) {
	out := new(UppercaseResponse)
	err := c.cc.Invoke(ctx, "/Service/Uppercase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	Buy(context.Context, *BuyRequest) (*BuyResponse, error)
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	Join2(context.Context, *Join2Request) (*Join2Response, error)
	Join3(context.Context, *Join3Request) (*Join3Response, error)
	Lowercase(context.Context, *LowercaseRequest) (*LowercaseResponse, error)
	Uppercase(context.Context, *UppercaseRequest) (*UppercaseResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) Buy(context.Context, *BuyRequest) (*BuyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Buy not implemented")
}
func (UnimplementedServiceServer) Join(context.Context, *JoinRequest) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedServiceServer) Join2(context.Context, *Join2Request) (*Join2Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join2 not implemented")
}
func (UnimplementedServiceServer) Join3(context.Context, *Join3Request) (*Join3Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join3 not implemented")
}
func (UnimplementedServiceServer) Lowercase(context.Context, *LowercaseRequest) (*LowercaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lowercase not implemented")
}
func (UnimplementedServiceServer) Uppercase(context.Context, *UppercaseRequest) (*UppercaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Uppercase not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_Buy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Buy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Buy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Buy(ctx, req.(*BuyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Join2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Join2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Join2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Join2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Join2(ctx, req.(*Join2Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Join3_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Join3Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Join3(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Join3",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Join3(ctx, req.(*Join3Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Lowercase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LowercaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Lowercase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Lowercase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Lowercase(ctx, req.(*LowercaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Uppercase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UppercaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Uppercase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/Uppercase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Uppercase(ctx, req.(*UppercaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Buy",
			Handler:    _Service_Buy_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _Service_Join_Handler,
		},
		{
			MethodName: "Join2",
			Handler:    _Service_Join2_Handler,
		},
		{
			MethodName: "Join3",
			Handler:    _Service_Join3_Handler,
		},
		{
			MethodName: "Lowercase",
			Handler:    _Service_Lowercase_Handler,
		},
		{
			MethodName: "Uppercase",
			Handler:    _Service_Uppercase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}