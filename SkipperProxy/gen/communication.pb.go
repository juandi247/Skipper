// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: communication.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Method        string                 `protobuf:"bytes,1,opt,name=Method,proto3" json:"Method,omitempty"`
	Proto         string                 `protobuf:"bytes,2,opt,name=Proto,proto3" json:"Proto,omitempty"`
	TargetUri     string                 `protobuf:"bytes,3,opt,name=TargetUri,proto3" json:"TargetUri,omitempty"`
	Path          string                 `protobuf:"bytes,4,opt,name=Path,proto3" json:"Path,omitempty"`
	Headers       map[string]string      `protobuf:"bytes,5,rep,name=Headers,proto3" json:"Headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Body          string                 `protobuf:"bytes,6,opt,name=Body,proto3" json:"Body,omitempty"`
	RequestId     string                 `protobuf:"bytes,7,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Request) Reset() {
	*x = Request{}
	mi := &file_communication_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_communication_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_communication_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Request) GetProto() string {
	if x != nil {
		return x.Proto
	}
	return ""
}

func (x *Request) GetTargetUri() string {
	if x != nil {
		return x.TargetUri
	}
	return ""
}

func (x *Request) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Request) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Request) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Request) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=Status,proto3" json:"Status,omitempty"`
	StatusCode    int32                  `protobuf:"varint,2,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	ProtoMajor    int32                  `protobuf:"varint,3,opt,name=ProtoMajor,proto3" json:"ProtoMajor,omitempty"`
	ProtoMinor    int32                  `protobuf:"varint,4,opt,name=ProtoMinor,proto3" json:"ProtoMinor,omitempty"`
	Proto         string                 `protobuf:"bytes,5,opt,name=Proto,proto3" json:"Proto,omitempty"`
	Headers       map[string]string      `protobuf:"bytes,6,rep,name=Headers,proto3" json:"Headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Body          string                 `protobuf:"bytes,7,opt,name=Body,proto3" json:"Body,omitempty"`
	RequestId     string                 `protobuf:"bytes,8,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_communication_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_communication_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_communication_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Response) GetProtoMajor() int32 {
	if x != nil {
		return x.ProtoMajor
	}
	return 0
}

func (x *Response) GetProtoMinor() int32 {
	if x != nil {
		return x.ProtoMinor
	}
	return 0
}

func (x *Response) GetProto() string {
	if x != nil {
		return x.Proto
	}
	return ""
}

func (x *Response) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Response) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Response) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

var File_communication_proto protoreflect.FileDescriptor

const file_communication_proto_rawDesc = "" +
	"\n" +
	"\x13communication.proto\x12\rcommunication\"\x96\x02\n" +
	"\aRequest\x12\x16\n" +
	"\x06Method\x18\x01 \x01(\tR\x06Method\x12\x14\n" +
	"\x05Proto\x18\x02 \x01(\tR\x05Proto\x12\x1c\n" +
	"\tTargetUri\x18\x03 \x01(\tR\tTargetUri\x12\x12\n" +
	"\x04Path\x18\x04 \x01(\tR\x04Path\x12=\n" +
	"\aHeaders\x18\x05 \x03(\v2#.communication.Request.HeadersEntryR\aHeaders\x12\x12\n" +
	"\x04Body\x18\x06 \x01(\tR\x04Body\x12\x1c\n" +
	"\tRequestId\x18\a \x01(\tR\tRequestId\x1a:\n" +
	"\fHeadersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\xc6\x02\n" +
	"\bResponse\x12\x16\n" +
	"\x06Status\x18\x01 \x01(\tR\x06Status\x12\x1e\n" +
	"\n" +
	"StatusCode\x18\x02 \x01(\x05R\n" +
	"StatusCode\x12\x1e\n" +
	"\n" +
	"ProtoMajor\x18\x03 \x01(\x05R\n" +
	"ProtoMajor\x12\x1e\n" +
	"\n" +
	"ProtoMinor\x18\x04 \x01(\x05R\n" +
	"ProtoMinor\x12\x14\n" +
	"\x05Proto\x18\x05 \x01(\tR\x05Proto\x12>\n" +
	"\aHeaders\x18\x06 \x03(\v2$.communication.Response.HeadersEntryR\aHeaders\x12\x12\n" +
	"\x04Body\x18\a \x01(\tR\x04Body\x12\x1c\n" +
	"\tRequestId\x18\b \x01(\tR\tRequestId\x1a:\n" +
	"\fHeadersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01B\x06Z\x04gen/b\x06proto3"

var (
	file_communication_proto_rawDescOnce sync.Once
	file_communication_proto_rawDescData []byte
)

func file_communication_proto_rawDescGZIP() []byte {
	file_communication_proto_rawDescOnce.Do(func() {
		file_communication_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_communication_proto_rawDesc), len(file_communication_proto_rawDesc)))
	})
	return file_communication_proto_rawDescData
}

var file_communication_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_communication_proto_goTypes = []any{
	(*Request)(nil),  // 0: communication.Request
	(*Response)(nil), // 1: communication.Response
	nil,              // 2: communication.Request.HeadersEntry
	nil,              // 3: communication.Response.HeadersEntry
}
var file_communication_proto_depIdxs = []int32{
	2, // 0: communication.Request.Headers:type_name -> communication.Request.HeadersEntry
	3, // 1: communication.Response.Headers:type_name -> communication.Response.HeadersEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_communication_proto_init() }
func file_communication_proto_init() {
	if File_communication_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_communication_proto_rawDesc), len(file_communication_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_communication_proto_goTypes,
		DependencyIndexes: file_communication_proto_depIdxs,
		MessageInfos:      file_communication_proto_msgTypes,
	}.Build()
	File_communication_proto = out.File
	file_communication_proto_goTypes = nil
	file_communication_proto_depIdxs = nil
}
