// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: data.proto

package proto

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Login     string `protobuf:"bytes,3,opt,name=Login,proto3" json:"Login,omitempty"`
	Pass      string `protobuf:"bytes,4,opt,name=Pass,proto3" json:"Pass,omitempty"`
	Text      string `protobuf:"bytes,5,opt,name=Text,proto3" json:"Text,omitempty"`
	CardNum   string `protobuf:"bytes,6,opt,name=CardNum,proto3" json:"CardNum,omitempty"`
	Meta      string `protobuf:"bytes,7,opt,name=Meta,proto3" json:"Meta,omitempty"`
	FileName  string `protobuf:"bytes,8,opt,name=FileName,proto3" json:"FileName,omitempty"`
	FileChunk []byte `protobuf:"bytes,9,opt,name=FileChunk,proto3" json:"FileChunk,omitempty"`
	Version   uint64 `protobuf:"varint,10,opt,name=Version,proto3" json:"Version,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Data) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Data) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *Data) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

func (x *Data) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Data) GetCardNum() string {
	if x != nil {
		return x.CardNum
	}
	return ""
}

func (x *Data) GetMeta() string {
	if x != nil {
		return x.Meta
	}
	return ""
}

func (x *Data) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Data) GetFileChunk() []byte {
	if x != nil {
		return x.FileChunk
	}
	return nil
}

func (x *Data) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type DataList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *DataList) Reset() {
	*x = DataList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataList) ProtoMessage() {}

func (x *DataList) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataList.ProtoReflect.Descriptor instead.
func (*DataList) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{1}
}

func (x *DataList) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DataList) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type SaveDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *Data `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *SaveDataRequest) Reset() {
	*x = SaveDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveDataRequest) ProtoMessage() {}

func (x *SaveDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveDataRequest.ProtoReflect.Descriptor instead.
func (*SaveDataRequest) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{2}
}

func (x *SaveDataRequest) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *GetDataRequest) Reset() {
	*x = GetDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataRequest) ProtoMessage() {}

func (x *GetDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataRequest.ProtoReflect.Descriptor instead.
func (*GetDataRequest) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{3}
}

func (x *GetDataRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *DeleteDataRequest) Reset() {
	*x = DeleteDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDataRequest) ProtoMessage() {}

func (x *DeleteDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDataRequest.ProtoReflect.Descriptor instead.
func (*DeleteDataRequest) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteDataRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataId      uint64 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId,omitempty"`
	FileId      uint64 `protobuf:"varint,2,opt,name=FileId,proto3" json:"FileId,omitempty"`
	DataVersion uint64 `protobuf:"varint,3,opt,name=DataVersion,proto3" json:"DataVersion,omitempty"`
	FileName    string `protobuf:"bytes,4,opt,name=FileName,proto3" json:"FileName,omitempty"`
	FileChunk   []byte `protobuf:"bytes,5,opt,name=FileChunk,proto3" json:"FileChunk,omitempty"`
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{5}
}

func (x *UploadFileRequest) GetDataId() uint64 {
	if x != nil {
		return x.DataId
	}
	return 0
}

func (x *UploadFileRequest) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *UploadFileRequest) GetDataVersion() uint64 {
	if x != nil {
		return x.DataVersion
	}
	return 0
}

func (x *UploadFileRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadFileRequest) GetFileChunk() []byte {
	if x != nil {
		return x.FileChunk
	}
	return nil
}

type GetDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  *Data  `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetDataResponse) Reset() {
	*x = GetDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataResponse) ProtoMessage() {}

func (x *GetDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataResponse.ProtoReflect.Descriptor instead.
func (*GetDataResponse) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{6}
}

func (x *GetDataResponse) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GetDataResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type SaveDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataId      uint64 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId,omitempty"`
	DataVersion uint64 `protobuf:"varint,2,opt,name=DataVersion,proto3" json:"DataVersion,omitempty"`
	Error       string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SaveDataResponse) Reset() {
	*x = SaveDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveDataResponse) ProtoMessage() {}

func (x *SaveDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveDataResponse.ProtoReflect.Descriptor instead.
func (*SaveDataResponse) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{7}
}

func (x *SaveDataResponse) GetDataId() uint64 {
	if x != nil {
		return x.DataId
	}
	return 0
}

func (x *SaveDataResponse) GetDataVersion() uint64 {
	if x != nil {
		return x.DataVersion
	}
	return 0
}

func (x *SaveDataResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type DataListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataList []*DataList `protobuf:"bytes,1,rep,name=DataList,proto3" json:"DataList,omitempty"`
}

func (x *DataListResponse) Reset() {
	*x = DataListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataListResponse) ProtoMessage() {}

func (x *DataListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataListResponse.ProtoReflect.Descriptor instead.
func (*DataListResponse) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{8}
}

func (x *DataListResponse) GetDataList() []*DataList {
	if x != nil {
		return x.DataList
	}
	return nil
}

type DeleteDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *DeleteDataResponse) Reset() {
	*x = DeleteDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDataResponse) ProtoMessage() {}

func (x *DeleteDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDataResponse.ProtoReflect.Descriptor instead.
func (*DeleteDataResponse) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteDataResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type FileUploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId      uint64 `protobuf:"varint,1,opt,name=FileId,proto3" json:"FileId,omitempty"`
	DataVersion uint64 `protobuf:"varint,2,opt,name=DataVersion,proto3" json:"DataVersion,omitempty"`
	Size        uint32 `protobuf:"varint,3,opt,name=Size,proto3" json:"Size,omitempty"`
}

func (x *FileUploadResponse) Reset() {
	*x = FileUploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadResponse) ProtoMessage() {}

func (x *FileUploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadResponse.ProtoReflect.Descriptor instead.
func (*FileUploadResponse) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{10}
}

func (x *FileUploadResponse) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *FileUploadResponse) GetDataVersion() uint64 {
	if x != nil {
		return x.DataVersion
	}
	return 0
}

func (x *FileUploadResponse) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

var File_data_proto protoreflect.FileDescriptor

var file_data_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x6f,
	0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xf6, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xba, 0x48, 0x07, 0x72, 0x05,
	0x10, 0x01, 0x18, 0xff, 0x01, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x50, 0x61, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x61, 0x72,
	0x64, 0x4e, 0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x61, 0x72, 0x64,
	0x4e, 0x75, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e,
	0x6b, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x2e, 0x0a, 0x08, 0x44,
	0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x37, 0x0a, 0x0f, 0x53,
	0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24,
	0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67,
	0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x22, 0x29, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x42, 0x07, 0xba, 0x48, 0x04, 0x22, 0x02, 0x20, 0x00, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x23, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x49, 0x64, 0x22, 0xc6, 0x01, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x44, 0x61,
	0x74, 0x61, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xba, 0x48, 0x04, 0x32,
	0x02, 0x20, 0x00, 0x52, 0x06, 0x44, 0x61, 0x74, 0x61, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x46,
	0x69, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x46, 0x69, 0x6c,
	0x65, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xba, 0x48, 0x04, 0x32, 0x02, 0x20,
	0x00, 0x52, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26,
	0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0a, 0xba, 0x48, 0x07, 0x72, 0x05, 0x10, 0x01, 0x18, 0xff, 0x01, 0x52, 0x08, 0x46, 0x69,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xba, 0x48, 0x04, 0x7a, 0x02,
	0x10, 0x01, 0x52, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x4d, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x62, 0x0a, 0x10,
	0x53, 0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x44, 0x61, 0x74, 0x61, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x44, 0x61, 0x74, 0x61, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x44,
	0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x44, 0x0a, 0x10, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65,
	0x70, 0x65, 0x72, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x08, 0x44, 0x61,
	0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x2a, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x62, 0x0a, 0x12, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x64,
	0x12, 0x20, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x32, 0xf9, 0x02, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x1b, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e,
	0x53, 0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x61, 0x76,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x2e,
	0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x67, 0x6f, 0x70, 0x68,
	0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x28, 0x01, 0x42, 0x12, 0x5a, 0x10, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_proto_rawDescOnce sync.Once
	file_data_proto_rawDescData = file_data_proto_rawDesc
)

func file_data_proto_rawDescGZIP() []byte {
	file_data_proto_rawDescOnce.Do(func() {
		file_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_proto_rawDescData)
	})
	return file_data_proto_rawDescData
}

var file_data_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_data_proto_goTypes = []any{
	(*Data)(nil),               // 0: gophkeeper.Data
	(*DataList)(nil),           // 1: gophkeeper.DataList
	(*SaveDataRequest)(nil),    // 2: gophkeeper.SaveDataRequest
	(*GetDataRequest)(nil),     // 3: gophkeeper.GetDataRequest
	(*DeleteDataRequest)(nil),  // 4: gophkeeper.DeleteDataRequest
	(*UploadFileRequest)(nil),  // 5: gophkeeper.UploadFileRequest
	(*GetDataResponse)(nil),    // 6: gophkeeper.GetDataResponse
	(*SaveDataResponse)(nil),   // 7: gophkeeper.SaveDataResponse
	(*DataListResponse)(nil),   // 8: gophkeeper.DataListResponse
	(*DeleteDataResponse)(nil), // 9: gophkeeper.DeleteDataResponse
	(*FileUploadResponse)(nil), // 10: gophkeeper.FileUploadResponse
	(*emptypb.Empty)(nil),      // 11: google.protobuf.Empty
}
var file_data_proto_depIdxs = []int32{
	0,  // 0: gophkeeper.SaveDataRequest.Data:type_name -> gophkeeper.Data
	0,  // 1: gophkeeper.GetDataResponse.Data:type_name -> gophkeeper.Data
	1,  // 2: gophkeeper.DataListResponse.DataList:type_name -> gophkeeper.DataList
	2,  // 3: gophkeeper.DataService.SaveData:input_type -> gophkeeper.SaveDataRequest
	11, // 4: gophkeeper.DataService.GetDataList:input_type -> google.protobuf.Empty
	3,  // 5: gophkeeper.DataService.GetData:input_type -> gophkeeper.GetDataRequest
	4,  // 6: gophkeeper.DataService.DeleteData:input_type -> gophkeeper.DeleteDataRequest
	5,  // 7: gophkeeper.DataService.UploadFile:input_type -> gophkeeper.UploadFileRequest
	7,  // 8: gophkeeper.DataService.SaveData:output_type -> gophkeeper.SaveDataResponse
	8,  // 9: gophkeeper.DataService.GetDataList:output_type -> gophkeeper.DataListResponse
	6,  // 10: gophkeeper.DataService.GetData:output_type -> gophkeeper.GetDataResponse
	9,  // 11: gophkeeper.DataService.DeleteData:output_type -> gophkeeper.DeleteDataResponse
	10, // 12: gophkeeper.DataService.UploadFile:output_type -> gophkeeper.FileUploadResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_data_proto_init() }
func file_data_proto_init() {
	if File_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Data); i {
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
		file_data_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DataList); i {
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
		file_data_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SaveDataRequest); i {
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
		file_data_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetDataRequest); i {
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
		file_data_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteDataRequest); i {
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
		file_data_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UploadFileRequest); i {
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
		file_data_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetDataResponse); i {
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
		file_data_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SaveDataResponse); i {
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
		file_data_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*DataListResponse); i {
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
		file_data_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteDataResponse); i {
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
		file_data_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*FileUploadResponse); i {
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
			RawDescriptor: file_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_proto_goTypes,
		DependencyIndexes: file_data_proto_depIdxs,
		MessageInfos:      file_data_proto_msgTypes,
	}.Build()
	File_data_proto = out.File
	file_data_proto_rawDesc = nil
	file_data_proto_goTypes = nil
	file_data_proto_depIdxs = nil
}
