// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: driver.proto

package driver

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type DriverInfo struct {
	DriverID             string   `protobuf:"bytes,1,opt,name=DriverID,proto3" json:"DriverID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Location             string   `protobuf:"bytes,3,opt,name=Location,proto3" json:"Location,omitempty"`
	ImageURL             string   `protobuf:"bytes,4,opt,name=ImageURL,proto3" json:"ImageURL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DriverInfo) Reset()         { *m = DriverInfo{} }
func (m *DriverInfo) String() string { return proto.CompactTextString(m) }
func (*DriverInfo) ProtoMessage()    {}
func (*DriverInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_521003751d596b5e, []int{0}
}
func (m *DriverInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DriverInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DriverInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DriverInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DriverInfo.Merge(m, src)
}
func (m *DriverInfo) XXX_Size() int {
	return m.Size()
}
func (m *DriverInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DriverInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DriverInfo proto.InternalMessageInfo

func (m *DriverInfo) GetDriverID() string {
	if m != nil {
		return m.DriverID
	}
	return ""
}

func (m *DriverInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DriverInfo) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *DriverInfo) GetImageURL() string {
	if m != nil {
		return m.ImageURL
	}
	return ""
}

type DriverLocationRequest struct {
	Location             string   `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DriverLocationRequest) Reset()         { *m = DriverLocationRequest{} }
func (m *DriverLocationRequest) String() string { return proto.CompactTextString(m) }
func (*DriverLocationRequest) ProtoMessage()    {}
func (*DriverLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_521003751d596b5e, []int{1}
}
func (m *DriverLocationRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DriverLocationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DriverLocationRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DriverLocationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DriverLocationRequest.Merge(m, src)
}
func (m *DriverLocationRequest) XXX_Size() int {
	return m.Size()
}
func (m *DriverLocationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DriverLocationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DriverLocationRequest proto.InternalMessageInfo

func (m *DriverLocationRequest) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

type DriverLocation struct {
	DriverID             string      `protobuf:"bytes,1,opt,name=driverID,proto3" json:"driverID,omitempty"`
	Location             string      `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Driver               *DriverInfo `protobuf:"bytes,3,opt,name=driver,proto3" json:"driver,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *DriverLocation) Reset()         { *m = DriverLocation{} }
func (m *DriverLocation) String() string { return proto.CompactTextString(m) }
func (*DriverLocation) ProtoMessage()    {}
func (*DriverLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_521003751d596b5e, []int{2}
}
func (m *DriverLocation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DriverLocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DriverLocation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DriverLocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DriverLocation.Merge(m, src)
}
func (m *DriverLocation) XXX_Size() int {
	return m.Size()
}
func (m *DriverLocation) XXX_DiscardUnknown() {
	xxx_messageInfo_DriverLocation.DiscardUnknown(m)
}

var xxx_messageInfo_DriverLocation proto.InternalMessageInfo

func (m *DriverLocation) GetDriverID() string {
	if m != nil {
		return m.DriverID
	}
	return ""
}

func (m *DriverLocation) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *DriverLocation) GetDriver() *DriverInfo {
	if m != nil {
		return m.Driver
	}
	return nil
}

type DriverLocationResponse struct {
	Locations            []*DriverLocation `protobuf:"bytes,1,rep,name=locations,proto3" json:"locations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DriverLocationResponse) Reset()         { *m = DriverLocationResponse{} }
func (m *DriverLocationResponse) String() string { return proto.CompactTextString(m) }
func (*DriverLocationResponse) ProtoMessage()    {}
func (*DriverLocationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_521003751d596b5e, []int{3}
}
func (m *DriverLocationResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DriverLocationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DriverLocationResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DriverLocationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DriverLocationResponse.Merge(m, src)
}
func (m *DriverLocationResponse) XXX_Size() int {
	return m.Size()
}
func (m *DriverLocationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DriverLocationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DriverLocationResponse proto.InternalMessageInfo

func (m *DriverLocationResponse) GetLocations() []*DriverLocation {
	if m != nil {
		return m.Locations
	}
	return nil
}

func init() {
	proto.RegisterType((*DriverInfo)(nil), "driver.DriverInfo")
	proto.RegisterType((*DriverLocationRequest)(nil), "driver.DriverLocationRequest")
	proto.RegisterType((*DriverLocation)(nil), "driver.DriverLocation")
	proto.RegisterType((*DriverLocationResponse)(nil), "driver.DriverLocationResponse")
}

func init() { proto.RegisterFile("driver.proto", fileDescriptor_521003751d596b5e) }

var fileDescriptor_521003751d596b5e = []byte{
	// 271 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x29, 0xca, 0x2c,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x4a, 0xb8, 0xb8,
	0x5c, 0xc0, 0x2c, 0xcf, 0xbc, 0xb4, 0x7c, 0x21, 0x29, 0x2e, 0x0e, 0x28, 0xcf, 0x45, 0x82, 0x51,
	0x81, 0x51, 0x83, 0x33, 0x08, 0xce, 0x17, 0x12, 0xe2, 0x62, 0xf1, 0x4b, 0xcc, 0x4d, 0x95, 0x60,
	0x02, 0x8b, 0x83, 0xd9, 0x20, 0xf5, 0x3e, 0xf9, 0xc9, 0x89, 0x25, 0x99, 0xf9, 0x79, 0x12, 0xcc,
	0x10, 0xf5, 0x30, 0x3e, 0x48, 0xce, 0x33, 0x37, 0x31, 0x3d, 0x35, 0x34, 0xc8, 0x47, 0x82, 0x05,
	0x22, 0x07, 0xe3, 0x2b, 0x19, 0x73, 0x89, 0x42, 0xcc, 0x85, 0xa9, 0x0e, 0x4a, 0x2d, 0x2c, 0x4d,
	0x2d, 0x2e, 0x01, 0x69, 0xca, 0x81, 0x19, 0x08, 0x75, 0x00, 0x8c, 0xaf, 0x54, 0xc2, 0xc5, 0x87,
	0xaa, 0x09, 0xa4, 0x3a, 0x05, 0xcd, 0xb9, 0x30, 0x3e, 0x8a, 0x49, 0x4c, 0xa8, 0x26, 0x09, 0x69,
	0x71, 0x41, 0xbd, 0x0f, 0x76, 0x34, 0xb7, 0x91, 0x90, 0x1e, 0x34, 0x6c, 0x10, 0x41, 0x11, 0x04,
	0x0b, 0x20, 0x3f, 0x2e, 0x31, 0x74, 0xa7, 0x16, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x99, 0x70,
	0x71, 0xc2, 0x4c, 0x2c, 0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36, 0x12, 0x43, 0x35, 0x08, 0xae,
	0x05, 0xa1, 0xd0, 0x28, 0x96, 0x8b, 0x17, 0x22, 0x19, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a,
	0xe4, 0xc3, 0xc5, 0xed, 0x96, 0x99, 0x97, 0xe2, 0x97, 0x9a, 0x58, 0x04, 0x0a, 0x01, 0x59, 0x1c,
	0x46, 0x40, 0x02, 0x48, 0x4a, 0x0e, 0x97, 0x34, 0xc4, 0x51, 0x4e, 0x22, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0x63, 0x14, 0xd4, 0x13, 0x49, 0x6c, 0xe0, 0x48,
	0x37, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x6b, 0x69, 0xc1, 0x04, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DriverServiceClient is the client API for DriverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DriverServiceClient interface {
	FindNearest(ctx context.Context, in *DriverLocationRequest, opts ...grpc.CallOption) (*DriverLocationResponse, error)
}

type driverServiceClient struct {
	cc *grpc.ClientConn
}

func NewDriverServiceClient(cc *grpc.ClientConn) DriverServiceClient {
	return &driverServiceClient{cc}
}

func (c *driverServiceClient) FindNearest(ctx context.Context, in *DriverLocationRequest, opts ...grpc.CallOption) (*DriverLocationResponse, error) {
	out := new(DriverLocationResponse)
	err := c.cc.Invoke(ctx, "/driver.DriverService/FindNearest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverServiceServer is the server API for DriverService service.
type DriverServiceServer interface {
	FindNearest(context.Context, *DriverLocationRequest) (*DriverLocationResponse, error)
}

// UnimplementedDriverServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDriverServiceServer struct {
}

func (*UnimplementedDriverServiceServer) FindNearest(ctx context.Context, req *DriverLocationRequest) (*DriverLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindNearest not implemented")
}

func RegisterDriverServiceServer(s *grpc.Server, srv DriverServiceServer) {
	s.RegisterService(&_DriverService_serviceDesc, srv)
}

func _DriverService_FindNearest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DriverLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).FindNearest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driver.DriverService/FindNearest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).FindNearest(ctx, req.(*DriverLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DriverService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "driver.DriverService",
	HandlerType: (*DriverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindNearest",
			Handler:    _DriverService_FindNearest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "driver.proto",
}

func (m *DriverInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DriverInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DriverInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ImageURL) > 0 {
		i -= len(m.ImageURL)
		copy(dAtA[i:], m.ImageURL)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.ImageURL)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Location) > 0 {
		i -= len(m.Location)
		copy(dAtA[i:], m.Location)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Location)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DriverID) > 0 {
		i -= len(m.DriverID)
		copy(dAtA[i:], m.DriverID)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.DriverID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DriverLocationRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DriverLocationRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DriverLocationRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Location) > 0 {
		i -= len(m.Location)
		copy(dAtA[i:], m.Location)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Location)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DriverLocation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DriverLocation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DriverLocation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Driver != nil {
		{
			size, err := m.Driver.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintDriver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Location) > 0 {
		i -= len(m.Location)
		copy(dAtA[i:], m.Location)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Location)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DriverID) > 0 {
		i -= len(m.DriverID)
		copy(dAtA[i:], m.DriverID)
		i = encodeVarintDriver(dAtA, i, uint64(len(m.DriverID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DriverLocationResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DriverLocationResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DriverLocationResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Locations) > 0 {
		for iNdEx := len(m.Locations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Locations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDriver(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintDriver(dAtA []byte, offset int, v uint64) int {
	offset -= sovDriver(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DriverInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DriverID)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.ImageURL)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DriverLocationRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DriverLocation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DriverID)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	if m.Driver != nil {
		l = m.Driver.Size()
		n += 1 + l + sovDriver(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DriverLocationResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Locations) > 0 {
		for _, e := range m.Locations {
			l = e.Size()
			n += 1 + l + sovDriver(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDriver(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDriver(x uint64) (n int) {
	return sovDriver(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DriverInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
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
			return fmt.Errorf("proto: DriverInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DriverInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DriverID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DriverID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageURL", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageURL = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DriverLocationRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
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
			return fmt.Errorf("proto: DriverLocationRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DriverLocationRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DriverLocation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
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
			return fmt.Errorf("proto: DriverLocation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DriverLocation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DriverID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DriverID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Driver", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Driver == nil {
				m.Driver = &DriverInfo{}
			}
			if err := m.Driver.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DriverLocationResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
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
			return fmt.Errorf("proto: DriverLocationResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DriverLocationResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Locations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
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
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDriver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Locations = append(m.Locations, &DriverLocation{})
			if err := m.Locations[len(m.Locations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDriver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDriver
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
					return 0, ErrIntOverflowDriver
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
					return 0, ErrIntOverflowDriver
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
				return 0, ErrInvalidLengthDriver
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDriver
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDriver
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDriver        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDriver          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDriver = fmt.Errorf("proto: unexpected end of group")
)
