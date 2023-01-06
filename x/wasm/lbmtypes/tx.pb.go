// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/wasm/v1/tx.proto

package lbmtypes

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	github_com_line_lbm_sdk_types "github.com/line/lbm-sdk/types"
	types1 "github.com/line/lbm-sdk/types"
	github_com_line_wasmd_x_wasm_types "github.com/line/wasmd/x/wasm/types"
	types "github.com/line/wasmd/x/wasm/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgStoreCodeAndInstantiateContract submit Wasm code to the system and instantiate a contract using it.
type MsgStoreCodeAndInstantiateContract struct {
	// Sender is the that actor that signed the messages
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode          []byte              `protobuf:"bytes,2,opt,name=wasm_byte_code,json=wasmByteCode,proto3" json:"wasm_byte_code,omitempty"`
	InstantiatePermission *types.AccessConfig `protobuf:"bytes,5,opt,name=instantiate_permission,json=instantiatePermission,proto3" json:"instantiate_permission,omitempty"`
	// Admin is an optional address that can execute migrations
	Admin string `protobuf:"bytes,6,opt,name=admin,proto3" json:"admin,omitempty"`
	// Label is optional metadata to be stored with a contract instance.
	Label string `protobuf:"bytes,7,opt,name=label,proto3" json:"label,omitempty"`
	// Msg json encoded message to be passed to the contract on instantiation
	Msg github_com_line_wasmd_x_wasm_types.RawContractMessage `protobuf:"bytes,8,opt,name=msg,proto3,casttype=github.com/line/wasmd/x/wasm/types.RawContractMessage" json:"msg,omitempty"`
	// Funds coins that are transferred to the contract on instantiation
	Funds github_com_line_lbm_sdk_types.Coins `protobuf:"bytes,9,rep,name=funds,proto3,castrepeated=github.com/line/lbm-sdk/types.Coins" json:"funds"`
}

func (m *MsgStoreCodeAndInstantiateContract) Reset()         { *m = MsgStoreCodeAndInstantiateContract{} }
func (m *MsgStoreCodeAndInstantiateContract) String() string { return proto.CompactTextString(m) }
func (*MsgStoreCodeAndInstantiateContract) ProtoMessage()    {}
func (*MsgStoreCodeAndInstantiateContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_751e1d2b9f9bf9e8, []int{0}
}
func (m *MsgStoreCodeAndInstantiateContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgStoreCodeAndInstantiateContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgStoreCodeAndInstantiateContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgStoreCodeAndInstantiateContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgStoreCodeAndInstantiateContract.Merge(m, src)
}
func (m *MsgStoreCodeAndInstantiateContract) XXX_Size() int {
	return m.Size()
}
func (m *MsgStoreCodeAndInstantiateContract) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgStoreCodeAndInstantiateContract.DiscardUnknown(m)
}

var xxx_messageInfo_MsgStoreCodeAndInstantiateContract proto.InternalMessageInfo

// MsgStoreCodeAndInstantiateContractResponse returns store and instantiate result data.
type MsgStoreCodeAndInstantiateContractResponse struct {
	// CodeID is the reference to the stored WASM code
	CodeID uint64 `protobuf:"varint,1,opt,name=code_id,json=codeId,proto3" json:"code_id,omitempty"`
	// Address is the bech32 address of the new contract instance
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// Data contains base64-encoded bytes to returned from the contract
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *MsgStoreCodeAndInstantiateContractResponse) Reset() {
	*m = MsgStoreCodeAndInstantiateContractResponse{}
}
func (m *MsgStoreCodeAndInstantiateContractResponse) String() string {
	return proto.CompactTextString(m)
}
func (*MsgStoreCodeAndInstantiateContractResponse) ProtoMessage() {}
func (*MsgStoreCodeAndInstantiateContractResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_751e1d2b9f9bf9e8, []int{1}
}
func (m *MsgStoreCodeAndInstantiateContractResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgStoreCodeAndInstantiateContractResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgStoreCodeAndInstantiateContractResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgStoreCodeAndInstantiateContractResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgStoreCodeAndInstantiateContractResponse.Merge(m, src)
}
func (m *MsgStoreCodeAndInstantiateContractResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgStoreCodeAndInstantiateContractResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgStoreCodeAndInstantiateContractResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgStoreCodeAndInstantiateContractResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgStoreCodeAndInstantiateContract)(nil), "lbm.wasm.v1.MsgStoreCodeAndInstantiateContract")
	proto.RegisterType((*MsgStoreCodeAndInstantiateContractResponse)(nil), "lbm.wasm.v1.MsgStoreCodeAndInstantiateContractResponse")
}

func init() { proto.RegisterFile("lbm/wasm/v1/tx.proto", fileDescriptor_751e1d2b9f9bf9e8) }

var fileDescriptor_751e1d2b9f9bf9e8 = []byte{
	// 540 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x4f, 0x8b, 0xd3, 0x4e,
	0x18, 0x6e, 0x7e, 0x49, 0xff, 0x4d, 0xcb, 0x8f, 0x32, 0xd4, 0x25, 0x16, 0x49, 0x4a, 0x17, 0xa1,
	0x28, 0x4e, 0x68, 0x45, 0xc5, 0x63, 0x5b, 0x2f, 0x5d, 0x29, 0x48, 0x16, 0x11, 0x44, 0x28, 0x93,
	0xcc, 0x6c, 0x0c, 0x26, 0x33, 0x25, 0x33, 0xdb, 0xdd, 0x9e, 0xfc, 0x00, 0x5e, 0x04, 0xbf, 0x85,
	0x07, 0x3f, 0x47, 0x8f, 0x7b, 0xf4, 0x54, 0xb5, 0xfd, 0x16, 0x9e, 0x64, 0x26, 0xcd, 0xba, 0x78,
	0xb0, 0x78, 0x9b, 0xe7, 0xfd, 0xf3, 0xbc, 0xef, 0x3c, 0xcf, 0x0c, 0x68, 0x27, 0x41, 0xea, 0x5d,
	0x60, 0x91, 0x7a, 0xcb, 0x81, 0x27, 0x2f, 0xd1, 0x22, 0xe3, 0x92, 0xc3, 0x46, 0x12, 0xa4, 0x48,
	0x45, 0xd1, 0x72, 0xd0, 0x69, 0x47, 0x3c, 0xe2, 0x3a, 0xee, 0xa9, 0x53, 0x5e, 0xd2, 0x71, 0x42,
	0x2e, 0x52, 0x2e, 0xbc, 0x00, 0x0b, 0xea, 0x2d, 0x07, 0x01, 0x95, 0x78, 0xe0, 0x85, 0x3c, 0x66,
	0xfb, 0xfc, 0x1d, 0x95, 0xd7, 0xc4, 0xd7, 0xec, 0xab, 0x05, 0x15, 0x79, 0xb6, 0xf7, 0xc5, 0x04,
	0xbd, 0x99, 0x88, 0x4e, 0x25, 0xcf, 0xe8, 0x84, 0x13, 0x3a, 0x62, 0x64, 0xca, 0x84, 0xc4, 0x4c,
	0xc6, 0x58, 0xd2, 0x09, 0x67, 0x32, 0xc3, 0xa1, 0x84, 0x47, 0xa0, 0x22, 0x28, 0x23, 0x34, 0xb3,
	0x8d, 0xae, 0xd1, 0xaf, 0xfb, 0x7b, 0x04, 0x1f, 0x83, 0xff, 0x15, 0xeb, 0x3c, 0x58, 0x49, 0x3a,
	0x0f, 0x39, 0xa1, 0xf6, 0x7f, 0x5d, 0xa3, 0xdf, 0x1c, 0xb7, 0xb6, 0x1b, 0xb7, 0xf9, 0x6a, 0x74,
	0x3a, 0x1b, 0xaf, 0xa4, 0xe6, 0xf5, 0x9b, 0xaa, 0xae, 0x40, 0xf0, 0x25, 0x38, 0x8a, 0x7f, 0x8f,
	0x99, 0x2f, 0x68, 0x96, 0xc6, 0x42, 0xc4, 0x9c, 0xd9, 0xe5, 0xae, 0xd1, 0x6f, 0x0c, 0x1d, 0x54,
	0x6c, 0x5d, 0xdc, 0x1e, 0x8d, 0xc2, 0x90, 0x0a, 0x31, 0xe1, 0xec, 0x2c, 0x8e, 0xfc, 0x5b, 0x37,
	0xba, 0x5f, 0x5c, 0x37, 0xc3, 0x36, 0x28, 0x63, 0x92, 0xc6, 0xcc, 0xae, 0xe8, 0x2d, 0x73, 0xa0,
	0xa2, 0x09, 0x0e, 0x68, 0x62, 0x57, 0xf3, 0xa8, 0x06, 0xf0, 0x39, 0x30, 0x53, 0x11, 0xd9, 0x35,
	0xbd, 0xef, 0xd3, 0x9f, 0x1b, 0xf7, 0x51, 0x14, 0xcb, 0xb7, 0xe7, 0x01, 0x0a, 0x79, 0xea, 0x25,
	0x31, 0xa3, 0x5a, 0x2f, 0xe2, 0x5d, 0xe6, 0xba, 0xe5, 0xa2, 0xf9, 0xf8, 0xa2, 0xd0, 0x64, 0x46,
	0x85, 0xc0, 0x11, 0xf5, 0x15, 0x0b, 0x7c, 0x03, 0xca, 0x67, 0xe7, 0x8c, 0x08, 0xbb, 0xde, 0x35,
	0xfb, 0x8d, 0xe1, 0x6d, 0x94, 0x9b, 0x82, 0x94, 0x29, 0x68, 0x6f, 0x0a, 0x9a, 0xf0, 0x98, 0x8d,
	0xef, 0xaf, 0x37, 0x6e, 0xe9, 0xf3, 0x37, 0xf7, 0xf8, 0xcf, 0x69, 0x49, 0x90, 0x3e, 0x10, 0xe4,
	0xdd, 0x7e, 0x90, 0xaa, 0x15, 0x7e, 0x4e, 0x7a, 0x62, 0xd5, 0xcc, 0x96, 0x75, 0x62, 0xd5, 0xac,
	0x56, 0xb9, 0xf7, 0x1e, 0xdc, 0x3b, 0xec, 0x97, 0x4f, 0xc5, 0x82, 0x33, 0x41, 0xe1, 0x31, 0xa8,
	0x2a, 0x57, 0xe6, 0x31, 0xd1, 0xc6, 0x59, 0x63, 0xb0, 0xdd, 0xb8, 0x15, 0xd5, 0x38, 0x7d, 0xe6,
	0x57, 0x54, 0x6a, 0x4a, 0xa0, 0x0d, 0xaa, 0x98, 0x90, 0x8c, 0x0a, 0xa1, 0xdd, 0xab, 0xfb, 0x05,
	0x84, 0x10, 0x58, 0x04, 0x4b, 0x6c, 0x9b, 0x4a, 0x24, 0x5f, 0x9f, 0x87, 0x9f, 0x0c, 0x60, 0xce,
	0x44, 0x04, 0x3f, 0x18, 0xc0, 0x3d, 0xf4, 0x6c, 0x3c, 0x74, 0xe3, 0xfd, 0xa2, 0xc3, 0x7b, 0x77,
	0x9e, 0xfc, 0x63, 0x43, 0x71, 0xd1, 0xf1, 0x64, 0xfd, 0xc3, 0x29, 0xad, 0xb7, 0x8e, 0x71, 0xb5,
	0x75, 0x8c, 0xef, 0x5b, 0xc7, 0xf8, 0xb8, 0x73, 0x4a, 0x57, 0x3b, 0xa7, 0xf4, 0x75, 0xe7, 0x94,
	0x5e, 0xdf, 0xfd, 0xab, 0xb5, 0x49, 0x90, 0x6a, 0xd1, 0x83, 0x8a, 0xfe, 0x13, 0x0f, 0x7f, 0x05,
	0x00, 0x00, 0xff, 0xff, 0xca, 0xf8, 0xad, 0x18, 0x8c, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// StoreCodeAndInstantiateContract upload code and instantiate a contract using it
	StoreCodeAndInstantiateContract(ctx context.Context, in *MsgStoreCodeAndInstantiateContract, opts ...grpc.CallOption) (*MsgStoreCodeAndInstantiateContractResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) StoreCodeAndInstantiateContract(ctx context.Context, in *MsgStoreCodeAndInstantiateContract, opts ...grpc.CallOption) (*MsgStoreCodeAndInstantiateContractResponse, error) {
	out := new(MsgStoreCodeAndInstantiateContractResponse)
	err := c.cc.Invoke(ctx, "/lbm.wasm.v1.Msg/StoreCodeAndInstantiateContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// StoreCodeAndInstantiateContract upload code and instantiate a contract using it
	StoreCodeAndInstantiateContract(context.Context, *MsgStoreCodeAndInstantiateContract) (*MsgStoreCodeAndInstantiateContractResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) StoreCodeAndInstantiateContract(ctx context.Context, req *MsgStoreCodeAndInstantiateContract) (*MsgStoreCodeAndInstantiateContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCodeAndInstantiateContract not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_StoreCodeAndInstantiateContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgStoreCodeAndInstantiateContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).StoreCodeAndInstantiateContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbm.wasm.v1.Msg/StoreCodeAndInstantiateContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).StoreCodeAndInstantiateContract(ctx, req.(*MsgStoreCodeAndInstantiateContract))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lbm.wasm.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreCodeAndInstantiateContract",
			Handler:    _Msg_StoreCodeAndInstantiateContract_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbm/wasm/v1/tx.proto",
}

func (m *MsgStoreCodeAndInstantiateContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgStoreCodeAndInstantiateContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgStoreCodeAndInstantiateContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Funds) > 0 {
		for iNdEx := len(m.Funds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Funds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.Label) > 0 {
		i -= len(m.Label)
		copy(dAtA[i:], m.Label)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Label)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Admin) > 0 {
		i -= len(m.Admin)
		copy(dAtA[i:], m.Admin)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Admin)))
		i--
		dAtA[i] = 0x32
	}
	if m.InstantiatePermission != nil {
		{
			size, err := m.InstantiatePermission.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if len(m.WASMByteCode) > 0 {
		i -= len(m.WASMByteCode)
		copy(dAtA[i:], m.WASMByteCode)
		i = encodeVarintTx(dAtA, i, uint64(len(m.WASMByteCode)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgStoreCodeAndInstantiateContractResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgStoreCodeAndInstantiateContractResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgStoreCodeAndInstantiateContractResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.CodeID != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.CodeID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgStoreCodeAndInstantiateContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.WASMByteCode)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.InstantiatePermission != nil {
		l = m.InstantiatePermission.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Admin)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.Funds) > 0 {
		for _, e := range m.Funds {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	return n
}

func (m *MsgStoreCodeAndInstantiateContractResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CodeID != 0 {
		n += 1 + sovTx(uint64(m.CodeID))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgStoreCodeAndInstantiateContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgStoreCodeAndInstantiateContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgStoreCodeAndInstantiateContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WASMByteCode", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WASMByteCode = append(m.WASMByteCode[:0], dAtA[iNdEx:postIndex]...)
			if m.WASMByteCode == nil {
				m.WASMByteCode = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InstantiatePermission", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.InstantiatePermission == nil {
				m.InstantiatePermission = &types.AccessConfig{}
			}
			if err := m.InstantiatePermission.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Admin", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Admin = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Label", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Label = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Funds = append(m.Funds, types1.Coin{})
			if err := m.Funds[len(m.Funds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgStoreCodeAndInstantiateContractResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgStoreCodeAndInstantiateContractResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgStoreCodeAndInstantiateContractResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeID", wireType)
			}
			m.CodeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CodeID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)