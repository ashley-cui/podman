// Copyright 2016 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: trillian_admin_api.proto

package trillian

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ListTrees request.
// No filters or pagination options are provided.
type ListTreesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If true, deleted trees are included in the response.
	ShowDeleted bool `protobuf:"varint,1,opt,name=show_deleted,json=showDeleted,proto3" json:"show_deleted,omitempty"`
}

func (x *ListTreesRequest) Reset() {
	*x = ListTreesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTreesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTreesRequest) ProtoMessage() {}

func (x *ListTreesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTreesRequest.ProtoReflect.Descriptor instead.
func (*ListTreesRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{0}
}

func (x *ListTreesRequest) GetShowDeleted() bool {
	if x != nil {
		return x.ShowDeleted
	}
	return false
}

// ListTrees response.
// No pagination is provided, all trees the requester has access to are
// returned.
type ListTreesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Trees matching the list request filters.
	Tree []*Tree `protobuf:"bytes,1,rep,name=tree,proto3" json:"tree,omitempty"`
}

func (x *ListTreesResponse) Reset() {
	*x = ListTreesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTreesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTreesResponse) ProtoMessage() {}

func (x *ListTreesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTreesResponse.ProtoReflect.Descriptor instead.
func (*ListTreesResponse) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{1}
}

func (x *ListTreesResponse) GetTree() []*Tree {
	if x != nil {
		return x.Tree
	}
	return nil
}

// GetTree request.
type GetTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the tree to retrieve.
	TreeId int64 `protobuf:"varint,1,opt,name=tree_id,json=treeId,proto3" json:"tree_id,omitempty"`
}

func (x *GetTreeRequest) Reset() {
	*x = GetTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTreeRequest) ProtoMessage() {}

func (x *GetTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTreeRequest.ProtoReflect.Descriptor instead.
func (*GetTreeRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{2}
}

func (x *GetTreeRequest) GetTreeId() int64 {
	if x != nil {
		return x.TreeId
	}
	return 0
}

// CreateTree request.
type CreateTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Tree to be created. See Tree and CreateTree for more details.
	Tree *Tree `protobuf:"bytes,1,opt,name=tree,proto3" json:"tree,omitempty"`
}

func (x *CreateTreeRequest) Reset() {
	*x = CreateTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTreeRequest) ProtoMessage() {}

func (x *CreateTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTreeRequest.ProtoReflect.Descriptor instead.
func (*CreateTreeRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{3}
}

func (x *CreateTreeRequest) GetTree() *Tree {
	if x != nil {
		return x.Tree
	}
	return nil
}

// UpdateTree request.
type UpdateTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Tree to be updated.
	Tree *Tree `protobuf:"bytes,1,opt,name=tree,proto3" json:"tree,omitempty"`
	// Fields modified by the update request.
	// For example: "tree_state", "display_name", "description".
	UpdateMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateTreeRequest) Reset() {
	*x = UpdateTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTreeRequest) ProtoMessage() {}

func (x *UpdateTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTreeRequest.ProtoReflect.Descriptor instead.
func (*UpdateTreeRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateTreeRequest) GetTree() *Tree {
	if x != nil {
		return x.Tree
	}
	return nil
}

func (x *UpdateTreeRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

// DeleteTree request.
type DeleteTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the tree to delete.
	TreeId int64 `protobuf:"varint,1,opt,name=tree_id,json=treeId,proto3" json:"tree_id,omitempty"`
}

func (x *DeleteTreeRequest) Reset() {
	*x = DeleteTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTreeRequest) ProtoMessage() {}

func (x *DeleteTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTreeRequest.ProtoReflect.Descriptor instead.
func (*DeleteTreeRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteTreeRequest) GetTreeId() int64 {
	if x != nil {
		return x.TreeId
	}
	return 0
}

// UndeleteTree request.
type UndeleteTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the tree to undelete.
	TreeId int64 `protobuf:"varint,1,opt,name=tree_id,json=treeId,proto3" json:"tree_id,omitempty"`
}

func (x *UndeleteTreeRequest) Reset() {
	*x = UndeleteTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trillian_admin_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UndeleteTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UndeleteTreeRequest) ProtoMessage() {}

func (x *UndeleteTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trillian_admin_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UndeleteTreeRequest.ProtoReflect.Descriptor instead.
func (*UndeleteTreeRequest) Descriptor() ([]byte, []int) {
	return file_trillian_admin_api_proto_rawDescGZIP(), []int{6}
}

func (x *UndeleteTreeRequest) GetTreeId() int64 {
	if x != nil {
		return x.TreeId
	}
	return 0
}

var File_trillian_admin_api_proto protoreflect.FileDescriptor

var file_trillian_admin_api_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x74, 0x72, 0x69, 0x6c,
	0x6c, 0x69, 0x61, 0x6e, 0x1a, 0x0e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x72,
	0x65, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x68,
	0x6f, 0x77, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x73, 0x68, 0x6f, 0x77, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x37, 0x0a,
	0x11, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x72, 0x65, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x72, 0x65, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x54, 0x72, 0x65, 0x65,
	0x52, 0x04, 0x74, 0x72, 0x65, 0x65, 0x22, 0x29, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x65,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x72, 0x65, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x72, 0x65, 0x65, 0x49,
	0x64, 0x22, 0x47, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x72, 0x65, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e,
	0x54, 0x72, 0x65, 0x65, 0x52, 0x04, 0x74, 0x72, 0x65, 0x65, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x52, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x22, 0x74, 0x0a, 0x11, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x22, 0x0a, 0x04, 0x74, 0x72, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x52, 0x04, 0x74,
	0x72, 0x65, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61,
	0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x4d, 0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b,
	0x22, 0x2c, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x72, 0x65, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x72, 0x65, 0x65, 0x49, 0x64, 0x22, 0x2e,
	0x0a, 0x13, 0x55, 0x6e, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x72, 0x65, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x72, 0x65, 0x65, 0x49, 0x64, 0x32, 0x86,
	0x03, 0x0a, 0x0d, 0x54, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x12, 0x46, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x72, 0x65, 0x65, 0x73, 0x12, 0x1a, 0x2e,
	0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x72, 0x65,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x74, 0x72, 0x69, 0x6c,
	0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x72, 0x65, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x54,
	0x72, 0x65, 0x65, 0x12, 0x18, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x47,
	0x65, 0x74, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e,
	0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x22, 0x00, 0x12,
	0x3b, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x1b, 0x2e,
	0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x72, 0x69,
	0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0a,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x1b, 0x2e, 0x74, 0x72, 0x69,
	0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69,
	0x61, 0x6e, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0a, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x1b, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69,
	0x61, 0x6e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e,
	0x54, 0x72, 0x65, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x0c, 0x55, 0x6e, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x1d, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61,
	0x6e, 0x2e, 0x55, 0x6e, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e,
	0x2e, 0x54, 0x72, 0x65, 0x65, 0x22, 0x00, 0x42, 0x50, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x42, 0x15, 0x54, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1a, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_trillian_admin_api_proto_rawDescOnce sync.Once
	file_trillian_admin_api_proto_rawDescData = file_trillian_admin_api_proto_rawDesc
)

func file_trillian_admin_api_proto_rawDescGZIP() []byte {
	file_trillian_admin_api_proto_rawDescOnce.Do(func() {
		file_trillian_admin_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_trillian_admin_api_proto_rawDescData)
	})
	return file_trillian_admin_api_proto_rawDescData
}

var file_trillian_admin_api_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_trillian_admin_api_proto_goTypes = []interface{}{
	(*ListTreesRequest)(nil),      // 0: trillian.ListTreesRequest
	(*ListTreesResponse)(nil),     // 1: trillian.ListTreesResponse
	(*GetTreeRequest)(nil),        // 2: trillian.GetTreeRequest
	(*CreateTreeRequest)(nil),     // 3: trillian.CreateTreeRequest
	(*UpdateTreeRequest)(nil),     // 4: trillian.UpdateTreeRequest
	(*DeleteTreeRequest)(nil),     // 5: trillian.DeleteTreeRequest
	(*UndeleteTreeRequest)(nil),   // 6: trillian.UndeleteTreeRequest
	(*Tree)(nil),                  // 7: trillian.Tree
	(*fieldmaskpb.FieldMask)(nil), // 8: google.protobuf.FieldMask
}
var file_trillian_admin_api_proto_depIdxs = []int32{
	7,  // 0: trillian.ListTreesResponse.tree:type_name -> trillian.Tree
	7,  // 1: trillian.CreateTreeRequest.tree:type_name -> trillian.Tree
	7,  // 2: trillian.UpdateTreeRequest.tree:type_name -> trillian.Tree
	8,  // 3: trillian.UpdateTreeRequest.update_mask:type_name -> google.protobuf.FieldMask
	0,  // 4: trillian.TrillianAdmin.ListTrees:input_type -> trillian.ListTreesRequest
	2,  // 5: trillian.TrillianAdmin.GetTree:input_type -> trillian.GetTreeRequest
	3,  // 6: trillian.TrillianAdmin.CreateTree:input_type -> trillian.CreateTreeRequest
	4,  // 7: trillian.TrillianAdmin.UpdateTree:input_type -> trillian.UpdateTreeRequest
	5,  // 8: trillian.TrillianAdmin.DeleteTree:input_type -> trillian.DeleteTreeRequest
	6,  // 9: trillian.TrillianAdmin.UndeleteTree:input_type -> trillian.UndeleteTreeRequest
	1,  // 10: trillian.TrillianAdmin.ListTrees:output_type -> trillian.ListTreesResponse
	7,  // 11: trillian.TrillianAdmin.GetTree:output_type -> trillian.Tree
	7,  // 12: trillian.TrillianAdmin.CreateTree:output_type -> trillian.Tree
	7,  // 13: trillian.TrillianAdmin.UpdateTree:output_type -> trillian.Tree
	7,  // 14: trillian.TrillianAdmin.DeleteTree:output_type -> trillian.Tree
	7,  // 15: trillian.TrillianAdmin.UndeleteTree:output_type -> trillian.Tree
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_trillian_admin_api_proto_init() }
func file_trillian_admin_api_proto_init() {
	if File_trillian_admin_api_proto != nil {
		return
	}
	file_trillian_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_trillian_admin_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTreesRequest); i {
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
		file_trillian_admin_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTreesResponse); i {
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
		file_trillian_admin_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTreeRequest); i {
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
		file_trillian_admin_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTreeRequest); i {
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
		file_trillian_admin_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTreeRequest); i {
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
		file_trillian_admin_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTreeRequest); i {
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
		file_trillian_admin_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UndeleteTreeRequest); i {
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
			RawDescriptor: file_trillian_admin_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trillian_admin_api_proto_goTypes,
		DependencyIndexes: file_trillian_admin_api_proto_depIdxs,
		MessageInfos:      file_trillian_admin_api_proto_msgTypes,
	}.Build()
	File_trillian_admin_api_proto = out.File
	file_trillian_admin_api_proto_rawDesc = nil
	file_trillian_admin_api_proto_goTypes = nil
	file_trillian_admin_api_proto_depIdxs = nil
}