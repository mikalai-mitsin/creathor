// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: examplepb/v1/session.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

type SessionCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Weight      uint64                 `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Versions    []uint64               `protobuf:"varint,4,rep,packed,name=versions,proto3" json:"versions,omitempty"`
	Release     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=tested,proto3" json:"tested,omitempty"`
}

func (x *SessionCreate) Reset() {
	*x = SessionCreate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionCreate) ProtoMessage() {}

func (x *SessionCreate) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionCreate.ProtoReflect.Descriptor instead.
func (*SessionCreate) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{0}
}

func (x *SessionCreate) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *SessionCreate) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SessionCreate) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *SessionCreate) GetVersions() []uint64 {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *SessionCreate) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *SessionCreate) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

type SessionGet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SessionGet) Reset() {
	*x = SessionGet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionGet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionGet) ProtoMessage() {}

func (x *SessionGet) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionGet.ProtoReflect.Descriptor instead.
func (*SessionGet) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{1}
}

func (x *SessionGet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SessionUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title       *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Weight      *wrapperspb.UInt64Value `protobuf:"bytes,4,opt,name=weight,proto3" json:"weight,omitempty"`
	Versions    *structpb.ListValue     `protobuf:"bytes,5,opt,name=versions,proto3" json:"versions,omitempty"`
	Release     *timestamppb.Timestamp  `protobuf:"bytes,6,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp  `protobuf:"bytes,7,opt,name=tested,proto3" json:"tested,omitempty"`
}

func (x *SessionUpdate) Reset() {
	*x = SessionUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionUpdate) ProtoMessage() {}

func (x *SessionUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionUpdate.ProtoReflect.Descriptor instead.
func (*SessionUpdate) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{2}
}

func (x *SessionUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SessionUpdate) GetTitle() *wrapperspb.StringValue {
	if x != nil {
		return x.Title
	}
	return nil
}

func (x *SessionUpdate) GetDescription() *wrapperspb.StringValue {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *SessionUpdate) GetWeight() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Weight
	}
	return nil
}

func (x *SessionUpdate) GetVersions() *structpb.ListValue {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *SessionUpdate) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *SessionUpdate) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

type Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Title       string                 `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Description string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Weight      uint64                 `protobuf:"varint,6,opt,name=weight,proto3" json:"weight,omitempty"`
	Versions    []uint64               `protobuf:"varint,7,rep,packed,name=versions,proto3" json:"versions,omitempty"`
	Release     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=tested,proto3" json:"tested,omitempty"`
}

func (x *Session) Reset() {
	*x = Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{3}
}

func (x *Session) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Session) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Session) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Session) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Session) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Session) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Session) GetVersions() []uint64 {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *Session) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *Session) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

type ListSession struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Session `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Count uint64     `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ListSession) Reset() {
	*x = ListSession{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSession) ProtoMessage() {}

func (x *ListSession) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSession.ProtoReflect.Descriptor instead.
func (*ListSession) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{4}
}

func (x *ListSession) GetItems() []*Session {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListSession) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type SessionDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SessionDelete) Reset() {
	*x = SessionDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionDelete) ProtoMessage() {}

func (x *SessionDelete) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionDelete.ProtoReflect.Descriptor instead.
func (*SessionDelete) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{5}
}

func (x *SessionDelete) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SessionFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNumber *wrapperspb.UInt64Value `protobuf:"bytes,1,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   *wrapperspb.UInt64Value `protobuf:"bytes,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	OrderBy    []string                `protobuf:"bytes,3,rep,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Ids        []string                `protobuf:"bytes,4,rep,name=ids,proto3" json:"ids,omitempty"`
	Search     *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=search,proto3" json:"search,omitempty"`
}

func (x *SessionFilter) Reset() {
	*x = SessionFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_session_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionFilter) ProtoMessage() {}

func (x *SessionFilter) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_session_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionFilter.ProtoReflect.Descriptor instead.
func (*SessionFilter) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_session_proto_rawDescGZIP(), []int{6}
}

func (x *SessionFilter) GetPageNumber() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageNumber
	}
	return nil
}

func (x *SessionFilter) GetPageSize() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageSize
	}
	return nil
}

func (x *SessionFilter) GetOrderBy() []string {
	if x != nil {
		return x.OrderBy
	}
	return nil
}

func (x *SessionFilter) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *SessionFilter) GetSearch() *wrapperspb.StringValue {
	if x != nil {
		return x.Search
	}
	return nil
}

var File_examplepb_v1_session_proto protoreflect.FileDescriptor

var file_examplepb_v1_session_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x04, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x1c, 0x0a, 0x0a,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x47, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xeb, 0x02, 0x0a, 0x0d, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x3e, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x34, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06,
	0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x34,
	0x0a, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x72, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0xe5, 0x02, 0x0a, 0x07, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x04, 0x52, 0x08, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x06,
	0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x22, 0x50, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x2b, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x1f, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0xec, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e,
	0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x34, 0x0a, 0x06,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x32, 0xda, 0x03, 0x0a, 0x0e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x1b, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x15, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x55, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x47, 0x65,
	0x74, 0x1a, 0x15, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17,
	0x12, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x5e, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x1b, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x15,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a,
	0x32, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x5c, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x1b, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x2a, 0x15,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x58, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1b, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x42,
	0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x30, 0x31,
	0x38, 0x62, 0x66, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examplepb_v1_session_proto_rawDescOnce sync.Once
	file_examplepb_v1_session_proto_rawDescData = file_examplepb_v1_session_proto_rawDesc
)

func file_examplepb_v1_session_proto_rawDescGZIP() []byte {
	file_examplepb_v1_session_proto_rawDescOnce.Do(func() {
		file_examplepb_v1_session_proto_rawDescData = protoimpl.X.CompressGZIP(file_examplepb_v1_session_proto_rawDescData)
	})
	return file_examplepb_v1_session_proto_rawDescData
}

var file_examplepb_v1_session_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_examplepb_v1_session_proto_goTypes = []interface{}{
	(*SessionCreate)(nil),          // 0: examplepb.v1.SessionCreate
	(*SessionGet)(nil),             // 1: examplepb.v1.SessionGet
	(*SessionUpdate)(nil),          // 2: examplepb.v1.SessionUpdate
	(*Session)(nil),                // 3: examplepb.v1.Session
	(*ListSession)(nil),            // 4: examplepb.v1.ListSession
	(*SessionDelete)(nil),          // 5: examplepb.v1.SessionDelete
	(*SessionFilter)(nil),          // 6: examplepb.v1.SessionFilter
	(*timestamppb.Timestamp)(nil),  // 7: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil), // 8: google.protobuf.StringValue
	(*wrapperspb.UInt64Value)(nil), // 9: google.protobuf.UInt64Value
	(*structpb.ListValue)(nil),     // 10: google.protobuf.ListValue
	(*emptypb.Empty)(nil),          // 11: google.protobuf.Empty
}
var file_examplepb_v1_session_proto_depIdxs = []int32{
	7,  // 0: examplepb.v1.SessionCreate.release:type_name -> google.protobuf.Timestamp
	7,  // 1: examplepb.v1.SessionCreate.tested:type_name -> google.protobuf.Timestamp
	8,  // 2: examplepb.v1.SessionUpdate.title:type_name -> google.protobuf.StringValue
	8,  // 3: examplepb.v1.SessionUpdate.description:type_name -> google.protobuf.StringValue
	9,  // 4: examplepb.v1.SessionUpdate.weight:type_name -> google.protobuf.UInt64Value
	10, // 5: examplepb.v1.SessionUpdate.versions:type_name -> google.protobuf.ListValue
	7,  // 6: examplepb.v1.SessionUpdate.release:type_name -> google.protobuf.Timestamp
	7,  // 7: examplepb.v1.SessionUpdate.tested:type_name -> google.protobuf.Timestamp
	7,  // 8: examplepb.v1.Session.updated_at:type_name -> google.protobuf.Timestamp
	7,  // 9: examplepb.v1.Session.created_at:type_name -> google.protobuf.Timestamp
	7,  // 10: examplepb.v1.Session.release:type_name -> google.protobuf.Timestamp
	7,  // 11: examplepb.v1.Session.tested:type_name -> google.protobuf.Timestamp
	3,  // 12: examplepb.v1.ListSession.items:type_name -> examplepb.v1.Session
	9,  // 13: examplepb.v1.SessionFilter.page_number:type_name -> google.protobuf.UInt64Value
	9,  // 14: examplepb.v1.SessionFilter.page_size:type_name -> google.protobuf.UInt64Value
	8,  // 15: examplepb.v1.SessionFilter.search:type_name -> google.protobuf.StringValue
	0,  // 16: examplepb.v1.SessionService.Create:input_type -> examplepb.v1.SessionCreate
	1,  // 17: examplepb.v1.SessionService.Get:input_type -> examplepb.v1.SessionGet
	2,  // 18: examplepb.v1.SessionService.Update:input_type -> examplepb.v1.SessionUpdate
	5,  // 19: examplepb.v1.SessionService.Delete:input_type -> examplepb.v1.SessionDelete
	6,  // 20: examplepb.v1.SessionService.List:input_type -> examplepb.v1.SessionFilter
	3,  // 21: examplepb.v1.SessionService.Create:output_type -> examplepb.v1.Session
	3,  // 22: examplepb.v1.SessionService.Get:output_type -> examplepb.v1.Session
	3,  // 23: examplepb.v1.SessionService.Update:output_type -> examplepb.v1.Session
	11, // 24: examplepb.v1.SessionService.Delete:output_type -> google.protobuf.Empty
	4,  // 25: examplepb.v1.SessionService.List:output_type -> examplepb.v1.ListSession
	21, // [21:26] is the sub-list for method output_type
	16, // [16:21] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_examplepb_v1_session_proto_init() }
func file_examplepb_v1_session_proto_init() {
	if File_examplepb_v1_session_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examplepb_v1_session_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionCreate); i {
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
		file_examplepb_v1_session_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionGet); i {
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
		file_examplepb_v1_session_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionUpdate); i {
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
		file_examplepb_v1_session_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Session); i {
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
		file_examplepb_v1_session_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListSession); i {
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
		file_examplepb_v1_session_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionDelete); i {
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
		file_examplepb_v1_session_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionFilter); i {
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
			RawDescriptor: file_examplepb_v1_session_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examplepb_v1_session_proto_goTypes,
		DependencyIndexes: file_examplepb_v1_session_proto_depIdxs,
		MessageInfos:      file_examplepb_v1_session_proto_msgTypes,
	}.Build()
	File_examplepb_v1_session_proto = out.File
	file_examplepb_v1_session_proto_rawDesc = nil
	file_examplepb_v1_session_proto_goTypes = nil
	file_examplepb_v1_session_proto_depIdxs = nil
}
