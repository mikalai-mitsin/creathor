// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: equipment.proto

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

// EquipmentServiceClient is the client API for EquipmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EquipmentServiceClient interface {
	Create(ctx context.Context, in *EquipmentCreate, opts ...grpc.CallOption) (*Equipment, error)
	Get(ctx context.Context, in *EquipmentGet, opts ...grpc.CallOption) (*Equipment, error)
	Update(ctx context.Context, in *EquipmentUpdate, opts ...grpc.CallOption) (*Equipment, error)
	Delete(ctx context.Context, in *EquipmentDelete, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *EquipmentFilter, opts ...grpc.CallOption) (*ListEquipment, error)
}

type equipmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEquipmentServiceClient(cc grpc.ClientConnInterface) EquipmentServiceClient {
	return &equipmentServiceClient{cc}
}

func (c *equipmentServiceClient) Create(ctx context.Context, in *EquipmentCreate, opts ...grpc.CallOption) (*Equipment, error) {
	out := new(Equipment)
	err := c.cc.Invoke(ctx, "/examplepb.EquipmentService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *equipmentServiceClient) Get(ctx context.Context, in *EquipmentGet, opts ...grpc.CallOption) (*Equipment, error) {
	out := new(Equipment)
	err := c.cc.Invoke(ctx, "/examplepb.EquipmentService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *equipmentServiceClient) Update(ctx context.Context, in *EquipmentUpdate, opts ...grpc.CallOption) (*Equipment, error) {
	out := new(Equipment)
	err := c.cc.Invoke(ctx, "/examplepb.EquipmentService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *equipmentServiceClient) Delete(ctx context.Context, in *EquipmentDelete, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/examplepb.EquipmentService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *equipmentServiceClient) List(ctx context.Context, in *EquipmentFilter, opts ...grpc.CallOption) (*ListEquipment, error) {
	out := new(ListEquipment)
	err := c.cc.Invoke(ctx, "/examplepb.EquipmentService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EquipmentServiceServer is the server API for EquipmentService service.
// All implementations should embed UnimplementedEquipmentServiceServer
// for forward compatibility
type EquipmentServiceServer interface {
	Create(context.Context, *EquipmentCreate) (*Equipment, error)
	Get(context.Context, *EquipmentGet) (*Equipment, error)
	Update(context.Context, *EquipmentUpdate) (*Equipment, error)
	Delete(context.Context, *EquipmentDelete) (*emptypb.Empty, error)
	List(context.Context, *EquipmentFilter) (*ListEquipment, error)
}

// UnimplementedEquipmentServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEquipmentServiceServer struct {
}

func (UnimplementedEquipmentServiceServer) Create(context.Context, *EquipmentCreate) (*Equipment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedEquipmentServiceServer) Get(context.Context, *EquipmentGet) (*Equipment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedEquipmentServiceServer) Update(context.Context, *EquipmentUpdate) (*Equipment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEquipmentServiceServer) Delete(context.Context, *EquipmentDelete) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEquipmentServiceServer) List(context.Context, *EquipmentFilter) (*ListEquipment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeEquipmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EquipmentServiceServer will
// result in compilation errors.
type UnsafeEquipmentServiceServer interface {
	mustEmbedUnimplementedEquipmentServiceServer()
}

func RegisterEquipmentServiceServer(s grpc.ServiceRegistrar, srv EquipmentServiceServer) {
	s.RegisterService(&EquipmentService_ServiceDesc, srv)
}

func _EquipmentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EquipmentCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EquipmentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.EquipmentService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EquipmentServiceServer).Create(ctx, req.(*EquipmentCreate))
	}
	return interceptor(ctx, in, info, handler)
}

func _EquipmentService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EquipmentGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EquipmentServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.EquipmentService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EquipmentServiceServer).Get(ctx, req.(*EquipmentGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _EquipmentService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EquipmentUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EquipmentServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.EquipmentService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EquipmentServiceServer).Update(ctx, req.(*EquipmentUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _EquipmentService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EquipmentDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EquipmentServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.EquipmentService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EquipmentServiceServer).Delete(ctx, req.(*EquipmentDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _EquipmentService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EquipmentFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EquipmentServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examplepb.EquipmentService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EquipmentServiceServer).List(ctx, req.(*EquipmentFilter))
	}
	return interceptor(ctx, in, info, handler)
}

// EquipmentService_ServiceDesc is the grpc.ServiceDesc for EquipmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EquipmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examplepb.EquipmentService",
	HandlerType: (*EquipmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _EquipmentService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _EquipmentService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EquipmentService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EquipmentService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _EquipmentService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "equipment.proto",
}