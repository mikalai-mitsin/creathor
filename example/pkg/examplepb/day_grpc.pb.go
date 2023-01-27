// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: day.proto

package examplepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DayServiceClient is the client API for DayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DayServiceClient interface {
	Create(ctx context.Context, in *DayCreate, opts ...grpc.CallOption) (*Day, error)
	Get(ctx context.Context, in *DayGet, opts ...grpc.CallOption) (*Day, error)
	Update(ctx context.Context, in *DayUpdate, opts ...grpc.CallOption) (*Day, error)
	Delete(ctx context.Context, in *DayDelete, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *DayFilter, opts ...grpc.CallOption) (*ListDay, error)
}

type dayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDayServiceClient(cc grpc.ClientConnInterface) DayServiceClient {
	return &dayServiceClient{cc}
}

func (c *dayServiceClient) Create(ctx context.Context, in *DayCreate, opts ...grpc.CallOption) (*Day, error) {
	out := new(Day)
	err := c.cc.Invoke(ctx, "/examplepb.DayService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dayServiceClient) Get(ctx context.Context, in *DayGet, opts ...grpc.CallOption) (*Day, error) {
	out := new(Day)
	err := c.cc.Invoke(ctx, "/examplepb.DayService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dayServiceClient) Update(ctx context.Context, in *DayUpdate, opts ...grpc.CallOption) (*Day, error) {
	out := new(Day)
	err := c.cc.Invoke(ctx, "/examplepb.DayService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dayServiceClient) Delete(ctx context.Context, in *DayDelete, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/examplepb.DayService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dayServiceClient) List(ctx context.Context, in *DayFilter, opts ...grpc.CallOption) (*ListDay, error) {
	out := new(ListDay)
	err := c.cc.Invoke(ctx, "/examplepb.DayService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DayServiceServer is the server API for DayService service.
// All implementations should embed UnimplementedDayServiceServer
// for forward compatibility
type DayServiceServer interface {
	Create(context.Context, *DayCreate) (*Day, error)
	Get(context.Context, *DayGet) (*Day, error)
	Update(context.Context, *DayUpdate) (*Day, error)
	Delete(context.Context, *DayDelete) (*emptypb.Empty, error)
	List(context.Context, *DayFilter) (*ListDay, error)
}

// UnimplementedDayServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDayServiceServer struct {
}

func (UnimplementedDayServiceServer) Create(context.Context, *DayCreate) (*Day, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDayServiceServer) Get(context.Context, *DayGet) (*Day, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDayServiceServer) Update(context.Context, *DayUpdate) (*Day, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDayServiceServer) Delete(context.Context, *DayDelete) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDayServiceServer) List(context.Context, *DayFilter) (*ListDay, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeDayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DayServiceServer will
// result in compilation errors.
type UnsafeDayServiceServer interface {
	mustEmbedUnimplementedDayServiceServer()
}

func RegisterDayServiceServer(s grpc.ServiceRegistrar, srv DayServiceServer) {
	s.RegisterService(&DayService_ServiceDesc, srv)
}

func _DayService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DayServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.DayService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DayServiceServer).Create(ctx, req.(*DayCreate))
	}
	return interceptor(ctx, in, info, handler)
}

func _DayService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DayServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.DayService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DayServiceServer).Get(ctx, req.(*DayGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _DayService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DayServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.DayService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DayServiceServer).Update(ctx, req.(*DayUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _DayService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DayServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.DayService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DayServiceServer).Delete(ctx, req.(*DayDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _DayService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DayServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.DayService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DayServiceServer).List(ctx, req.(*DayFilter))
	}
	return interceptor(ctx, in, info, handler)
}

// DayService_ServiceDesc is the grpc.ServiceDesc for DayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examplepb.DayService",
	HandlerType: (*DayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _DayService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DayService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DayService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DayService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DayService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "day.proto",
}