// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: plan.proto

package examplepb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PlanCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Repeat      uint64 `protobuf:"varint,2,opt,name=repeat,proto3" json:"repeat,omitempty"`
	EquipmentId string `protobuf:"bytes,3,opt,name=equipment_id,json=equipmentId,proto3" json:"equipment_id,omitempty"`
}

func (x *PlanCreate) Reset() {
	*x = PlanCreate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlanCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlanCreate) ProtoMessage() {}

func (x *PlanCreate) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlanCreate.ProtoReflect.Descriptor instead.
func (*PlanCreate) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{0}
}

func (x *PlanCreate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PlanCreate) GetRepeat() uint64 {
	if x != nil {
		return x.Repeat
	}
	return 0
}

func (x *PlanCreate) GetEquipmentId() string {
	if x != nil {
		return x.EquipmentId
	}
	return ""
}

type PlanGet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PlanGet) Reset() {
	*x = PlanGet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlanGet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlanGet) ProtoMessage() {}

func (x *PlanGet) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlanGet.ProtoReflect.Descriptor instead.
func (*PlanGet) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{1}
}

func (x *PlanGet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PlanUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Repeat      *wrapperspb.UInt64Value `protobuf:"bytes,3,opt,name=repeat,proto3" json:"repeat,omitempty"`
	EquipmentId *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=equipment_id,json=equipmentId,proto3" json:"equipment_id,omitempty"`
}

func (x *PlanUpdate) Reset() {
	*x = PlanUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlanUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlanUpdate) ProtoMessage() {}

func (x *PlanUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlanUpdate.ProtoReflect.Descriptor instead.
func (*PlanUpdate) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{2}
}

func (x *PlanUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PlanUpdate) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *PlanUpdate) GetRepeat() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Repeat
	}
	return nil
}

func (x *PlanUpdate) GetEquipmentId() *wrapperspb.StringValue {
	if x != nil {
		return x.EquipmentId
	}
	return nil
}

type Plan struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Name        string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Repeat      uint64                 `protobuf:"varint,5,opt,name=repeat,proto3" json:"repeat,omitempty"`
	EquipmentId string                 `protobuf:"bytes,6,opt,name=equipment_id,json=equipmentId,proto3" json:"equipment_id,omitempty"`
}

func (x *Plan) Reset() {
	*x = Plan{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Plan) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Plan) ProtoMessage() {}

func (x *Plan) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Plan.ProtoReflect.Descriptor instead.
func (*Plan) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{3}
}

func (x *Plan) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Plan) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Plan) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Plan) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Plan) GetRepeat() uint64 {
	if x != nil {
		return x.Repeat
	}
	return 0
}

func (x *Plan) GetEquipmentId() string {
	if x != nil {
		return x.EquipmentId
	}
	return ""
}

type ListPlan struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Plan `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Count uint64  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ListPlan) Reset() {
	*x = ListPlan{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPlan) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPlan) ProtoMessage() {}

func (x *ListPlan) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPlan.ProtoReflect.Descriptor instead.
func (*ListPlan) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{4}
}

func (x *ListPlan) GetItems() []*Plan {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListPlan) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type PlanDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PlanDelete) Reset() {
	*x = PlanDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlanDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlanDelete) ProtoMessage() {}

func (x *PlanDelete) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlanDelete.ProtoReflect.Descriptor instead.
func (*PlanDelete) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{5}
}

func (x *PlanDelete) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PlanFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNumber *wrapperspb.UInt64Value `protobuf:"bytes,1,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   *wrapperspb.UInt64Value `protobuf:"bytes,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	OrderBy    []string                `protobuf:"bytes,3,rep,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Ids        []string                `protobuf:"bytes,4,rep,name=ids,proto3" json:"ids,omitempty"`
	Search     *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=search,proto3" json:"search,omitempty"`
}

func (x *PlanFilter) Reset() {
	*x = PlanFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plan_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlanFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlanFilter) ProtoMessage() {}

func (x *PlanFilter) ProtoReflect() protoreflect.Message {
	mi := &file_plan_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlanFilter.ProtoReflect.Descriptor instead.
func (*PlanFilter) Descriptor() ([]byte, []int) {
	return file_plan_proto_rawDescGZIP(), []int{6}
}

func (x *PlanFilter) GetPageNumber() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageNumber
	}
	return nil
}

func (x *PlanFilter) GetPageSize() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageSize
	}
	return nil
}

func (x *PlanFilter) GetOrderBy() []string {
	if x != nil {
		return x.OrderBy
	}
	return nil
}

func (x *PlanFilter) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *PlanFilter) GetSearch() *wrapperspb.StringValue {
	if x != nil {
		return x.Search
	}
	return nil
}

var File_plan_proto protoreflect.FileDescriptor

var file_plan_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x0a, 0x50, 0x6c, 0x61, 0x6e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x22, 0x19, 0x0a, 0x07, 0x50, 0x6c, 0x61, 0x6e, 0x47, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xc5, 0x01, 0x0a, 0x0a,
	0x50, 0x6c, 0x61, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x06,
	0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55,
	0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x72, 0x65, 0x70, 0x65,
	0x61, 0x74, 0x12, 0x3f, 0x0a, 0x0c, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x22, 0xdb, 0x01, 0x0a, 0x04, 0x50, 0x6c, 0x61, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x22, 0x47, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x6e, 0x12, 0x25, 0x0a,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x52, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x1c, 0x0a, 0x0a, 0x50, 0x6c,
	0x61, 0x6e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xe9, 0x01, 0x0a, 0x0a, 0x50, 0x6c, 0x61,
	0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55,
	0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74,
	0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x69, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x34,
	0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x32, 0x8a, 0x02, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x15,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x0f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70,
	0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x12, 0x2a, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x12, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x47, 0x65,
	0x74, 0x1a, 0x0f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c,
	0x61, 0x6e, 0x12, 0x30, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x15, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x1a, 0x0f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e,
	0x50, 0x6c, 0x61, 0x6e, 0x12, 0x37, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x15,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x32, 0x0a,
	0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x15, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70,
	0x62, 0x2e, 0x50, 0x6c, 0x61, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x13, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6c, 0x61,
	0x6e, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x30, 0x31, 0x38, 0x62, 0x66, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_plan_proto_rawDescOnce sync.Once
	file_plan_proto_rawDescData = file_plan_proto_rawDesc
)

func file_plan_proto_rawDescGZIP() []byte {
	file_plan_proto_rawDescOnce.Do(func() {
		file_plan_proto_rawDescData = protoimpl.X.CompressGZIP(file_plan_proto_rawDescData)
	})
	return file_plan_proto_rawDescData
}

var file_plan_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_plan_proto_goTypes = []interface{}{
	(*PlanCreate)(nil),             // 0: examplepb.PlanCreate
	(*PlanGet)(nil),                // 1: examplepb.PlanGet
	(*PlanUpdate)(nil),             // 2: examplepb.PlanUpdate
	(*Plan)(nil),                   // 3: examplepb.Plan
	(*ListPlan)(nil),               // 4: examplepb.ListPlan
	(*PlanDelete)(nil),             // 5: examplepb.PlanDelete
	(*PlanFilter)(nil),             // 6: examplepb.PlanFilter
	(*wrapperspb.StringValue)(nil), // 7: google.protobuf.StringValue
	(*wrapperspb.UInt64Value)(nil), // 8: google.protobuf.UInt64Value
	(*timestamppb.Timestamp)(nil),  // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),          // 10: google.protobuf.Empty
}
var file_plan_proto_depIdxs = []int32{
	7,  // 0: examplepb.PlanUpdate.name:type_name -> google.protobuf.StringValue
	8,  // 1: examplepb.PlanUpdate.repeat:type_name -> google.protobuf.UInt64Value
	7,  // 2: examplepb.PlanUpdate.equipment_id:type_name -> google.protobuf.StringValue
	9,  // 3: examplepb.Plan.updated_at:type_name -> google.protobuf.Timestamp
	9,  // 4: examplepb.Plan.created_at:type_name -> google.protobuf.Timestamp
	3,  // 5: examplepb.ListPlan.items:type_name -> examplepb.Plan
	8,  // 6: examplepb.PlanFilter.page_number:type_name -> google.protobuf.UInt64Value
	8,  // 7: examplepb.PlanFilter.page_size:type_name -> google.protobuf.UInt64Value
	7,  // 8: examplepb.PlanFilter.search:type_name -> google.protobuf.StringValue
	0,  // 9: examplepb.PlanService.Create:input_type -> examplepb.PlanCreate
	1,  // 10: examplepb.PlanService.Get:input_type -> examplepb.PlanGet
	2,  // 11: examplepb.PlanService.Update:input_type -> examplepb.PlanUpdate
	5,  // 12: examplepb.PlanService.Delete:input_type -> examplepb.PlanDelete
	6,  // 13: examplepb.PlanService.List:input_type -> examplepb.PlanFilter
	3,  // 14: examplepb.PlanService.Create:output_type -> examplepb.Plan
	3,  // 15: examplepb.PlanService.Get:output_type -> examplepb.Plan
	3,  // 16: examplepb.PlanService.Update:output_type -> examplepb.Plan
	10, // 17: examplepb.PlanService.Delete:output_type -> google.protobuf.Empty
	4,  // 18: examplepb.PlanService.List:output_type -> examplepb.ListPlan
	14, // [14:19] is the sub-list for method output_type
	9,  // [9:14] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_plan_proto_init() }
func file_plan_proto_init() {
	if File_plan_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plan_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlanCreate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlanGet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlanUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Plan); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPlan); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlanDelete); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plan_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlanFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plan_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plan_proto_goTypes,
		DependencyIndexes: file_plan_proto_depIdxs,
		MessageInfos:      file_plan_proto_msgTypes,
	}.Build()
	File_plan_proto = out.File
	file_plan_proto_rawDesc = nil
	file_plan_proto_goTypes = nil
	file_plan_proto_depIdxs = nil
}
