// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.11.4
// source: userGroups.proto

package userGroups

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//群基本信息
type Group struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid        string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`               //群主id
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`               //群名称
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"` //群介绍
	Avatar      string `protobuf:"bytes,4,opt,name=avatar,proto3" json:"avatar,omitempty"`           //群头像
	Gid         string `protobuf:"bytes,5,opt,name=gid,proto3" json:"gid,omitempty"`                 //群id
}

func (x *Group) Reset() {
	*x = Group{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Group) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Group) ProtoMessage() {}

func (x *Group) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Group.ProtoReflect.Descriptor instead.
func (*Group) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{0}
}

func (x *Group) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Group) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Group) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Group) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *Group) GetGid() string {
	if x != nil {
		return x.Gid
	}
	return ""
}

//创建群
type CreateGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *CreateGroupRequest) Reset() {
	*x = CreateGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGroupRequest) ProtoMessage() {}

func (x *CreateGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGroupRequest.ProtoReflect.Descriptor instead.
func (*CreateGroupRequest) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{1}
}

func (x *CreateGroupRequest) GetGroup() *Group {
	if x != nil {
		return x.Group
	}
	return nil
}

type CreateGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gid string `protobuf:"bytes,1,opt,name=gid,proto3" json:"gid,omitempty"`
}

func (x *CreateGroupResponse) Reset() {
	*x = CreateGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGroupResponse) ProtoMessage() {}

func (x *CreateGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGroupResponse.ProtoReflect.Descriptor instead.
func (*CreateGroupResponse) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{2}
}

func (x *CreateGroupResponse) GetGid() string {
	if x != nil {
		return x.Gid
	}
	return ""
}

//删除群
type DelGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Gid  string `protobuf:"bytes,2,opt,name=gid,proto3" json:"gid,omitempty"`
}

func (x *DelGroupRequest) Reset() {
	*x = DelGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelGroupRequest) ProtoMessage() {}

func (x *DelGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelGroupRequest.ProtoReflect.Descriptor instead.
func (*DelGroupRequest) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{3}
}

func (x *DelGroupRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *DelGroupRequest) GetGid() string {
	if x != nil {
		return x.Gid
	}
	return ""
}

type DelGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *DelGroupResponse) Reset() {
	*x = DelGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelGroupResponse) ProtoMessage() {}

func (x *DelGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelGroupResponse.ProtoReflect.Descriptor instead.
func (*DelGroupResponse) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{4}
}

func (x *DelGroupResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

//获取群列表
type GroupsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *GroupsRequest) Reset() {
	*x = GroupsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupsRequest) ProtoMessage() {}

func (x *GroupsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupsRequest.ProtoReflect.Descriptor instead.
func (*GroupsRequest) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{5}
}

func (x *GroupsRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type GroupsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Groups []*Group `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty"`
}

func (x *GroupsResponse) Reset() {
	*x = GroupsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupsResponse) ProtoMessage() {}

func (x *GroupsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupsResponse.ProtoReflect.Descriptor instead.
func (*GroupsResponse) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{6}
}

func (x *GroupsResponse) GetGroups() []*Group {
	if x != nil {
		return x.Groups
	}
	return nil
}

//根据gid查找群
type FindByGidRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gid string `protobuf:"bytes,1,opt,name=gid,proto3" json:"gid,omitempty"`
}

func (x *FindByGidRequest) Reset() {
	*x = FindByGidRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindByGidRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindByGidRequest) ProtoMessage() {}

func (x *FindByGidRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindByGidRequest.ProtoReflect.Descriptor instead.
func (*FindByGidRequest) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{7}
}

func (x *FindByGidRequest) GetGid() string {
	if x != nil {
		return x.Gid
	}
	return ""
}

type FindByGidResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *FindByGidResponse) Reset() {
	*x = FindByGidResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindByGidResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindByGidResponse) ProtoMessage() {}

func (x *FindByGidResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindByGidResponse.ProtoReflect.Descriptor instead.
func (*FindByGidResponse) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{8}
}

func (x *FindByGidResponse) GetGroup() *Group {
	if x != nil {
		return x.Group
	}
	return nil
}

//修改群信息
type EditGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *EditGroupRequest) Reset() {
	*x = EditGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditGroupRequest) ProtoMessage() {}

func (x *EditGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditGroupRequest.ProtoReflect.Descriptor instead.
func (*EditGroupRequest) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{9}
}

func (x *EditGroupRequest) GetGroup() *Group {
	if x != nil {
		return x.Group
	}
	return nil
}

type EditGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *EditGroupResponse) Reset() {
	*x = EditGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userGroups_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditGroupResponse) ProtoMessage() {}

func (x *EditGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userGroups_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditGroupResponse.ProtoReflect.Descriptor instead.
func (*EditGroupResponse) Descriptor() ([]byte, []int) {
	return file_userGroups_proto_rawDescGZIP(), []int{10}
}

func (x *EditGroupResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_userGroups_proto protoreflect.FileDescriptor

var file_userGroups_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x05, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x10, 0x0a,
	0x03, 0x67, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22,
	0x32, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x22, 0x27, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x37, 0x0a, 0x0f,
	0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x24, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x23, 0x0a, 0x0d, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x22, 0x30, 0x0a, 0x0e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1e, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x06, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x06, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x22, 0x24, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x47, 0x69, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x31, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64,
	0x42, 0x79, 0x47, 0x69, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a,
	0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x30, 0x0a, 0x10, 0x45,
	0x64, 0x69, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06,
	0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x25, 0x0a,
	0x11, 0x45, 0x64, 0x69, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x32, 0x91, 0x02, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x13, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x10, 0x2e, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12,
	0x0e, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0f, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x32, 0x0a, 0x09, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x47, 0x69, 0x64, 0x12, 0x11, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x47, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x47, 0x69, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x45, 0x64, 0x69, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x12, 0x11, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x3b, 0x75,
	0x73, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_userGroups_proto_rawDescOnce sync.Once
	file_userGroups_proto_rawDescData = file_userGroups_proto_rawDesc
)

func file_userGroups_proto_rawDescGZIP() []byte {
	file_userGroups_proto_rawDescOnce.Do(func() {
		file_userGroups_proto_rawDescData = protoimpl.X.CompressGZIP(file_userGroups_proto_rawDescData)
	})
	return file_userGroups_proto_rawDescData
}

var file_userGroups_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_userGroups_proto_goTypes = []interface{}{
	(*Group)(nil),               // 0: Group
	(*CreateGroupRequest)(nil),  // 1: CreateGroupRequest
	(*CreateGroupResponse)(nil), // 2: CreateGroupResponse
	(*DelGroupRequest)(nil),     // 3: DelGroupRequest
	(*DelGroupResponse)(nil),    // 4: DelGroupResponse
	(*GroupsRequest)(nil),       // 5: GroupsRequest
	(*GroupsResponse)(nil),      // 6: GroupsResponse
	(*FindByGidRequest)(nil),    // 7: FindByGidRequest
	(*FindByGidResponse)(nil),   // 8: FindByGidResponse
	(*EditGroupRequest)(nil),    // 9: EditGroupRequest
	(*EditGroupResponse)(nil),   // 10: EditGroupResponse
}
var file_userGroups_proto_depIdxs = []int32{
	0,  // 0: CreateGroupRequest.group:type_name -> Group
	0,  // 1: GroupsResponse.groups:type_name -> Group
	0,  // 2: FindByGidResponse.group:type_name -> Group
	0,  // 3: EditGroupRequest.group:type_name -> Group
	1,  // 4: UserGroupsService.CreateGroup:input_type -> CreateGroupRequest
	3,  // 5: UserGroupsService.DelGroup:input_type -> DelGroupRequest
	5,  // 6: UserGroupsService.Groups:input_type -> GroupsRequest
	7,  // 7: UserGroupsService.FindByGid:input_type -> FindByGidRequest
	9,  // 8: UserGroupsService.EditGroup:input_type -> EditGroupRequest
	2,  // 9: UserGroupsService.CreateGroup:output_type -> CreateGroupResponse
	4,  // 10: UserGroupsService.DelGroup:output_type -> DelGroupResponse
	6,  // 11: UserGroupsService.Groups:output_type -> GroupsResponse
	8,  // 12: UserGroupsService.FindByGid:output_type -> FindByGidResponse
	10, // 13: UserGroupsService.EditGroup:output_type -> EditGroupResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_userGroups_proto_init() }
func file_userGroups_proto_init() {
	if File_userGroups_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_userGroups_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Group); i {
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
		file_userGroups_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateGroupRequest); i {
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
		file_userGroups_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateGroupResponse); i {
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
		file_userGroups_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelGroupRequest); i {
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
		file_userGroups_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelGroupResponse); i {
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
		file_userGroups_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupsRequest); i {
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
		file_userGroups_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupsResponse); i {
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
		file_userGroups_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindByGidRequest); i {
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
		file_userGroups_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindByGidResponse); i {
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
		file_userGroups_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditGroupRequest); i {
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
		file_userGroups_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditGroupResponse); i {
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
			RawDescriptor: file_userGroups_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userGroups_proto_goTypes,
		DependencyIndexes: file_userGroups_proto_depIdxs,
		MessageInfos:      file_userGroups_proto_msgTypes,
	}.Build()
	File_userGroups_proto = out.File
	file_userGroups_proto_rawDesc = nil
	file_userGroups_proto_goTypes = nil
	file_userGroups_proto_depIdxs = nil
}
