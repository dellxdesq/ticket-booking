// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/notifications.proto

package notifications

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NotificationService_SendEmail_FullMethodName = "/notifications.NotificationService/SendEmail"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, NotificationService_SendEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	SendEmail(context.Context, *EmailRequest) (*EmailResponse, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) SendEmail(context.Context, *EmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SendEmail(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notifications.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _NotificationService_SendEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/notifications.proto",
}
