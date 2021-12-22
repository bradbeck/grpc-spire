// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package adder

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

// AddServiceClient is the client API for AddService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AddServiceClient interface {
	Compute(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
}

type addServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddServiceClient(cc grpc.ClientConnInterface) AddServiceClient {
	return &addServiceClient{cc}
}

func (c *addServiceClient) Compute(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	out := new(AddResponse)
	err := c.cc.Invoke(ctx, "/AddService/Compute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddServiceServer is the server API for AddService service.
// All implementations must embed UnimplementedAddServiceServer
// for forward compatibility
type AddServiceServer interface {
	Compute(context.Context, *AddRequest) (*AddResponse, error)
	mustEmbedUnimplementedAddServiceServer()
}

// UnimplementedAddServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAddServiceServer struct {
}

func (UnimplementedAddServiceServer) Compute(context.Context, *AddRequest) (*AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Compute not implemented")
}
func (UnimplementedAddServiceServer) mustEmbedUnimplementedAddServiceServer() {}

// UnsafeAddServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AddServiceServer will
// result in compilation errors.
type UnsafeAddServiceServer interface {
	mustEmbedUnimplementedAddServiceServer()
}

func RegisterAddServiceServer(s grpc.ServiceRegistrar, srv AddServiceServer) {
	s.RegisterService(&AddService_ServiceDesc, srv)
}

func _AddService_Compute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).Compute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AddService/Compute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).Compute(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AddService_ServiceDesc is the grpc.ServiceDesc for AddService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AddService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AddService",
	HandlerType: (*AddServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Compute",
			Handler:    _AddService_Compute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/adder/v1/adder.proto",
}
