// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: admin_movies_service_v1.proto

package protos

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_admin_movies_service_v1_proto protoreflect.FileDescriptor

var file_admin_movies_service_v1_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x14, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x26, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x76, 0x31, 0x5f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76,
	0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xe4, 0x16, 0x0a, 0x0f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x31, 0x12, 0x61, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x25, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x22, 0x11, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x12,
	0x82, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14,
	0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x65, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x12, 0x26, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12,
	0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x12, 0x78, 0x0a, 0x0b, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x28, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22,
	0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x62, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x12, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x2a, 0x09,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x0d, 0x49, 0x73,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x2a, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x49, 0x73, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49,
	0x73, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x63, 0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x76,
	0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x6c,
	0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x28, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x85, 0x01, 0x0a,
	0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x50, 0x69, 0x63, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x12, 0x30, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x50, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x24,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x22, 0x19, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2d, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x73, 0x3a, 0x01, 0x2a, 0x12, 0x62, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x67, 0x65, 0x52, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x20, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22,
	0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65,
	0x2d, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x73, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x2c, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x67, 0x65, 0x2d, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x70, 0x0a,
	0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x2c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x67,
	0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x2a, 0x0f,
	0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x2d, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x6c, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x12, 0x25, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x22,
	0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e,
	0x72, 0x65, 0x2f, 0x7b, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x74, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x2b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x42,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x12, 0x12, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x12, 0x78, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x65, 0x6e,
	0x72, 0x65, 0x12, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x47, 0x65, 0x6e, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x22,
	0x09, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x55, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x73,
	0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65,
	0x6e, 0x72, 0x65, 0x73, 0x12, 0x6c, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x65,
	0x6e, 0x72, 0x65, 0x12, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x22, 0x10, 0x2f,
	0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x3a,
	0x01, 0x2a, 0x12, 0x6d, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x47, 0x65, 0x6e, 0x72,
	0x65, 0x12, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x47,
	0x65, 0x6e, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x2a, 0x14, 0x2f, 0x76, 0x31,
	0x2f, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x2f, 0x7b, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x5f, 0x69, 0x64,
	0x7d, 0x12, 0x7e, 0x0a, 0x0e, 0x49, 0x73, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x73, 0x12, 0x2b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x73, 0x47, 0x65, 0x6e,
	0x72, 0x65, 0x73, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11,
	0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x2f, 0x65, 0x78, 0x69, 0x73, 0x74,
	0x73, 0x12, 0x76, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x27, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12,
	0x18, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x2f, 0x7b, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x7c, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x42,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x1a, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79,
	0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x80, 0x01, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x2a, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x63, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x3a, 0x01, 0x2a, 0x12, 0x5e, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x1f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x72, 0x0a, 0x0d, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x2a, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x22, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x72, 0x79, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x75,
	0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x2a, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x2a, 0x18, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x2f, 0x7b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x79, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x87, 0x01, 0x0a, 0x11, 0x49, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x2e, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x49, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63,
	0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x42,
	0xb8, 0x02, 0x5a, 0x1e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x92, 0x41, 0x94, 0x02, 0x12, 0x5c, 0x0a, 0x14, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x20,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x3f,
	0x0a, 0x07, 0x46, 0x61, 0x6c, 0x6f, 0x6b, 0x75, 0x74, 0x12, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x73,
	0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x46, 0x61,
	0x6c, 0x6f, 0x6b, 0x75, 0x74, 0x1a, 0x18, 0x74, 0x69, 0x6d, 0x75, 0x72, 0x2e, 0x73, 0x69, 0x6e,
	0x65, 0x6c, 0x6e, 0x69, 0x6b, 0x40, 0x79, 0x61, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x72, 0x75, 0x32,
	0x03, 0x31, 0x2e, 0x30, 0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x52, 0x50, 0x0a, 0x03, 0x34,
	0x30, 0x34, 0x12, 0x49, 0x0a, 0x2a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77,
	0x68, 0x65, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x20, 0x64, 0x6f, 0x65, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x2e,
	0x12, 0x1b, 0x0a, 0x19, 0x1a, 0x17, 0x23, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x3b, 0x0a,
	0x03, 0x35, 0x30, 0x30, 0x12, 0x34, 0x0a, 0x15, 0x53, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e,
	0x67, 0x20, 0x77, 0x65, 0x6e, 0x74, 0x20, 0x77, 0x72, 0x6f, 0x6e, 0x67, 0x2e, 0x12, 0x1b, 0x0a,
	0x19, 0x1a, 0x17, 0x23, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x72, 0x70, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var file_admin_movies_service_v1_proto_goTypes = []interface{}{
	(*GetMovieRequest)(nil),            // 0: admin_movies_service.GetMovieRequest
	(*GetMovieDurationRequest)(nil),    // 1: admin_movies_service.GetMovieDurationRequest
	(*GetMoviesRequest)(nil),           // 2: admin_movies_service.GetMoviesRequest
	(*CreateMovieRequest)(nil),         // 3: admin_movies_service.CreateMovieRequest
	(*DeleteMovieRequest)(nil),         // 4: admin_movies_service.DeleteMovieRequest
	(*IsMovieExistsRequest)(nil),       // 5: admin_movies_service.IsMovieExistsRequest
	(*UpdateMovieRequest)(nil),         // 6: admin_movies_service.UpdateMovieRequest
	(*UpdateMoviePicturesRequest)(nil), // 7: admin_movies_service.UpdateMoviePicturesRequest
	(*emptypb.Empty)(nil),              // 8: google.protobuf.Empty
	(*CreateAgeRatingRequest)(nil),     // 9: admin_movies_service.CreateAgeRatingRequest
	(*DeleteAgeRatingRequest)(nil),     // 10: admin_movies_service.DeleteAgeRatingRequest
	(*GetGenreRequest)(nil),            // 11: admin_movies_service.GetGenreRequest
	(*GetGenreByNameRequest)(nil),      // 12: admin_movies_service.GetGenreByNameRequest
	(*CreateGenreRequest)(nil),         // 13: admin_movies_service.CreateGenreRequest
	(*UpdateGenreRequest)(nil),         // 14: admin_movies_service.UpdateGenreRequest
	(*DeleteGenreRequest)(nil),         // 15: admin_movies_service.DeleteGenreRequest
	(*IsGenresExistsRequest)(nil),      // 16: admin_movies_service.IsGenresExistsRequest
	(*GetCountryRequest)(nil),          // 17: admin_movies_service.GetCountryRequest
	(*GetCountryByNameRequest)(nil),    // 18: admin_movies_service.GetCountryByNameRequest
	(*CreateCountryRequest)(nil),       // 19: admin_movies_service.CreateCountryRequest
	(*UpdateCountryRequest)(nil),       // 20: admin_movies_service.UpdateCountryRequest
	(*DeleteCountryRequest)(nil),       // 21: admin_movies_service.DeleteCountryRequest
	(*IsCountriesExistsRequest)(nil),   // 22: admin_movies_service.IsCountriesExistsRequest
	(*Movie)(nil),                      // 23: admin_movies_service.Movie
	(*MovieDuration)(nil),              // 24: admin_movies_service.MovieDuration
	(*Movies)(nil),                     // 25: admin_movies_service.Movies
	(*CreateMovieResponce)(nil),        // 26: admin_movies_service.CreateMovieResponce
	(*IsMovieExistsResponce)(nil),      // 27: admin_movies_service.IsMovieExistsResponce
	(*AgeRatings)(nil),                 // 28: admin_movies_service.AgeRatings
	(*Genre)(nil),                      // 29: admin_movies_service.Genre
	(*CreateGenreResponce)(nil),        // 30: admin_movies_service.CreateGenreResponce
	(*Genres)(nil),                     // 31: admin_movies_service.Genres
	(*ExistsResponce)(nil),             // 32: admin_movies_service.ExistsResponce
	(*Country)(nil),                    // 33: admin_movies_service.Country
	(*CreateCountryResponce)(nil),      // 34: admin_movies_service.CreateCountryResponce
	(*Countries)(nil),                  // 35: admin_movies_service.Countries
}
var file_admin_movies_service_v1_proto_depIdxs = []int32{
	0,  // 0: admin_movies_service.moviesServiceV1.GetMovie:input_type -> admin_movies_service.GetMovieRequest
	1,  // 1: admin_movies_service.moviesServiceV1.GetMovieDuration:input_type -> admin_movies_service.GetMovieDurationRequest
	2,  // 2: admin_movies_service.moviesServiceV1.GetMovies:input_type -> admin_movies_service.GetMoviesRequest
	3,  // 3: admin_movies_service.moviesServiceV1.CreateMovie:input_type -> admin_movies_service.CreateMovieRequest
	4,  // 4: admin_movies_service.moviesServiceV1.DeleteMovie:input_type -> admin_movies_service.DeleteMovieRequest
	5,  // 5: admin_movies_service.moviesServiceV1.IsMovieExists:input_type -> admin_movies_service.IsMovieExistsRequest
	6,  // 6: admin_movies_service.moviesServiceV1.UpdateMovie:input_type -> admin_movies_service.UpdateMovieRequest
	7,  // 7: admin_movies_service.moviesServiceV1.UpdateMoviePictures:input_type -> admin_movies_service.UpdateMoviePicturesRequest
	8,  // 8: admin_movies_service.moviesServiceV1.GetAgeRatings:input_type -> google.protobuf.Empty
	9,  // 9: admin_movies_service.moviesServiceV1.CreateAgeRating:input_type -> admin_movies_service.CreateAgeRatingRequest
	10, // 10: admin_movies_service.moviesServiceV1.DeleteAgeRating:input_type -> admin_movies_service.DeleteAgeRatingRequest
	11, // 11: admin_movies_service.moviesServiceV1.GetGenre:input_type -> admin_movies_service.GetGenreRequest
	12, // 12: admin_movies_service.moviesServiceV1.GetGenreByName:input_type -> admin_movies_service.GetGenreByNameRequest
	13, // 13: admin_movies_service.moviesServiceV1.CreateGenre:input_type -> admin_movies_service.CreateGenreRequest
	8,  // 14: admin_movies_service.moviesServiceV1.GetGenres:input_type -> google.protobuf.Empty
	14, // 15: admin_movies_service.moviesServiceV1.UpdateGenre:input_type -> admin_movies_service.UpdateGenreRequest
	15, // 16: admin_movies_service.moviesServiceV1.DeleteGenre:input_type -> admin_movies_service.DeleteGenreRequest
	16, // 17: admin_movies_service.moviesServiceV1.IsGenresExists:input_type -> admin_movies_service.IsGenresExistsRequest
	17, // 18: admin_movies_service.moviesServiceV1.GetCountry:input_type -> admin_movies_service.GetCountryRequest
	18, // 19: admin_movies_service.moviesServiceV1.GetCountryByName:input_type -> admin_movies_service.GetCountryByNameRequest
	19, // 20: admin_movies_service.moviesServiceV1.CreateCountry:input_type -> admin_movies_service.CreateCountryRequest
	8,  // 21: admin_movies_service.moviesServiceV1.GetCountries:input_type -> google.protobuf.Empty
	20, // 22: admin_movies_service.moviesServiceV1.UpdateCountry:input_type -> admin_movies_service.UpdateCountryRequest
	21, // 23: admin_movies_service.moviesServiceV1.DeleteCountry:input_type -> admin_movies_service.DeleteCountryRequest
	22, // 24: admin_movies_service.moviesServiceV1.IsCountriesExists:input_type -> admin_movies_service.IsCountriesExistsRequest
	23, // 25: admin_movies_service.moviesServiceV1.GetMovie:output_type -> admin_movies_service.Movie
	24, // 26: admin_movies_service.moviesServiceV1.GetMovieDuration:output_type -> admin_movies_service.MovieDuration
	25, // 27: admin_movies_service.moviesServiceV1.GetMovies:output_type -> admin_movies_service.Movies
	26, // 28: admin_movies_service.moviesServiceV1.CreateMovie:output_type -> admin_movies_service.CreateMovieResponce
	8,  // 29: admin_movies_service.moviesServiceV1.DeleteMovie:output_type -> google.protobuf.Empty
	27, // 30: admin_movies_service.moviesServiceV1.IsMovieExists:output_type -> admin_movies_service.IsMovieExistsResponce
	8,  // 31: admin_movies_service.moviesServiceV1.UpdateMovie:output_type -> google.protobuf.Empty
	8,  // 32: admin_movies_service.moviesServiceV1.UpdateMoviePictures:output_type -> google.protobuf.Empty
	28, // 33: admin_movies_service.moviesServiceV1.GetAgeRatings:output_type -> admin_movies_service.AgeRatings
	8,  // 34: admin_movies_service.moviesServiceV1.CreateAgeRating:output_type -> google.protobuf.Empty
	8,  // 35: admin_movies_service.moviesServiceV1.DeleteAgeRating:output_type -> google.protobuf.Empty
	29, // 36: admin_movies_service.moviesServiceV1.GetGenre:output_type -> admin_movies_service.Genre
	29, // 37: admin_movies_service.moviesServiceV1.GetGenreByName:output_type -> admin_movies_service.Genre
	30, // 38: admin_movies_service.moviesServiceV1.CreateGenre:output_type -> admin_movies_service.CreateGenreResponce
	31, // 39: admin_movies_service.moviesServiceV1.GetGenres:output_type -> admin_movies_service.Genres
	8,  // 40: admin_movies_service.moviesServiceV1.UpdateGenre:output_type -> google.protobuf.Empty
	8,  // 41: admin_movies_service.moviesServiceV1.DeleteGenre:output_type -> google.protobuf.Empty
	32, // 42: admin_movies_service.moviesServiceV1.IsGenresExists:output_type -> admin_movies_service.ExistsResponce
	33, // 43: admin_movies_service.moviesServiceV1.GetCountry:output_type -> admin_movies_service.Country
	33, // 44: admin_movies_service.moviesServiceV1.GetCountryByName:output_type -> admin_movies_service.Country
	34, // 45: admin_movies_service.moviesServiceV1.CreateCountry:output_type -> admin_movies_service.CreateCountryResponce
	35, // 46: admin_movies_service.moviesServiceV1.GetCountries:output_type -> admin_movies_service.Countries
	8,  // 47: admin_movies_service.moviesServiceV1.UpdateCountry:output_type -> google.protobuf.Empty
	8,  // 48: admin_movies_service.moviesServiceV1.DeleteCountry:output_type -> google.protobuf.Empty
	32, // 49: admin_movies_service.moviesServiceV1.IsCountriesExists:output_type -> admin_movies_service.ExistsResponce
	25, // [25:50] is the sub-list for method output_type
	0,  // [0:25] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_admin_movies_service_v1_proto_init() }
func file_admin_movies_service_v1_proto_init() {
	if File_admin_movies_service_v1_proto != nil {
		return
	}
	file_admin_movies_service_v1_messages_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_movies_service_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_movies_service_v1_proto_goTypes,
		DependencyIndexes: file_admin_movies_service_v1_proto_depIdxs,
	}.Build()
	File_admin_movies_service_v1_proto = out.File
	file_admin_movies_service_v1_proto_rawDesc = nil
	file_admin_movies_service_v1_proto_goTypes = nil
	file_admin_movies_service_v1_proto_depIdxs = nil
}
