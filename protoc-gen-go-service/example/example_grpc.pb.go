// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package example

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExampleServiceClient is the client API for ExampleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleServiceClient interface {
	Unary(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	NoReturn(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	ServerStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (ExampleService_ServerStreamClient, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (ExampleService_ClientStreamClient, error)
	BiDirectionalStream(ctx context.Context, opts ...grpc.CallOption) (ExampleService_BiDirectionalStreamClient, error)
}

type exampleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleServiceClient(cc grpc.ClientConnInterface) ExampleServiceClient {
	return &exampleServiceClient{cc}
}

func (c *exampleServiceClient) Unary(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/v1.ExampleService/Unary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleServiceClient) NoReturn(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/v1.ExampleService/NoReturn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleServiceClient) ServerStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (ExampleService_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExampleService_ServiceDesc.Streams[0], "/v1.ExampleService/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleServiceServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExampleService_ServerStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type exampleServiceServerStreamClient struct {
	grpc.ClientStream
}

func (x *exampleServiceServerStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exampleServiceClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (ExampleService_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExampleService_ServiceDesc.Streams[1], "/v1.ExampleService/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleServiceClientStreamClient{stream}
	return x, nil
}

type ExampleService_ClientStreamClient interface {
	Send(*Request) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type exampleServiceClientStreamClient struct {
	grpc.ClientStream
}

func (x *exampleServiceClientStreamClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exampleServiceClientStreamClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exampleServiceClient) BiDirectionalStream(ctx context.Context, opts ...grpc.CallOption) (ExampleService_BiDirectionalStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExampleService_ServiceDesc.Streams[2], "/v1.ExampleService/BiDirectionalStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleServiceBiDirectionalStreamClient{stream}
	return x, nil
}

type ExampleService_BiDirectionalStreamClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type exampleServiceBiDirectionalStreamClient struct {
	grpc.ClientStream
}

func (x *exampleServiceBiDirectionalStreamClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exampleServiceBiDirectionalStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExampleServiceServer is the server API for ExampleService service.
// All implementations must embed UnimplementedExampleServiceServer
// for forward compatibility
type ExampleServiceServer interface {
	Unary(context.Context, *Request) (*Response, error)
	NoReturn(context.Context, *empty.Empty) (*empty.Empty, error)
	ServerStream(*Request, ExampleService_ServerStreamServer) error
	ClientStream(ExampleService_ClientStreamServer) error
	BiDirectionalStream(ExampleService_BiDirectionalStreamServer) error
	mustEmbedUnimplementedExampleServiceServer()
}

// UnimplementedExampleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExampleServiceServer struct {
}

func (UnimplementedExampleServiceServer) Unary(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unary not implemented")
}
func (UnimplementedExampleServiceServer) NoReturn(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoReturn not implemented")
}
func (UnimplementedExampleServiceServer) ServerStream(*Request, ExampleService_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (UnimplementedExampleServiceServer) ClientStream(ExampleService_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedExampleServiceServer) BiDirectionalStream(ExampleService_BiDirectionalStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method BiDirectionalStream not implemented")
}
func (UnimplementedExampleServiceServer) mustEmbedUnimplementedExampleServiceServer() {}

// UnsafeExampleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleServiceServer will
// result in compilation errors.
type UnsafeExampleServiceServer interface {
	mustEmbedUnimplementedExampleServiceServer()
}

func RegisterExampleServiceServer(s grpc.ServiceRegistrar, srv ExampleServiceServer) {
	s.RegisterService(&ExampleService_ServiceDesc, srv)
}

func _ExampleService_Unary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).Unary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExampleService/Unary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).Unary(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExampleService_NoReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).NoReturn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExampleService/NoReturn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).NoReturn(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExampleService_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExampleServiceServer).ServerStream(m, &exampleServiceServerStreamServer{stream})
}

type ExampleService_ServerStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type exampleServiceServerStreamServer struct {
	grpc.ServerStream
}

func (x *exampleServiceServerStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _ExampleService_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExampleServiceServer).ClientStream(&exampleServiceClientStreamServer{stream})
}

type ExampleService_ClientStreamServer interface {
	SendAndClose(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type exampleServiceClientStreamServer struct {
	grpc.ServerStream
}

func (x *exampleServiceClientStreamServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exampleServiceClientStreamServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExampleService_BiDirectionalStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExampleServiceServer).BiDirectionalStream(&exampleServiceBiDirectionalStreamServer{stream})
}

type ExampleService_BiDirectionalStreamServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type exampleServiceBiDirectionalStreamServer struct {
	grpc.ServerStream
}

func (x *exampleServiceBiDirectionalStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exampleServiceBiDirectionalStreamServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExampleService_ServiceDesc is the grpc.ServiceDesc for ExampleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExampleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ExampleService",
	HandlerType: (*ExampleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unary",
			Handler:    _ExampleService_Unary_Handler,
		},
		{
			MethodName: "NoReturn",
			Handler:    _ExampleService_NoReturn_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _ExampleService_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _ExampleService_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BiDirectionalStream",
			Handler:       _ExampleService_BiDirectionalStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "example/example.proto",
}
