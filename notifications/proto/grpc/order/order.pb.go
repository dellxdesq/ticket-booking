// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/order.proto

package order

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

type GetAvailableSeatsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	EventId       int64                  `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAvailableSeatsRequest) Reset() {
	*x = GetAvailableSeatsRequest{}
	mi := &file_proto_order_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAvailableSeatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableSeatsRequest) ProtoMessage() {}

func (x *GetAvailableSeatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableSeatsRequest.ProtoReflect.Descriptor instead.
func (*GetAvailableSeatsRequest) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{0}
}

func (x *GetAvailableSeatsRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type GetAvailableSeatsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	EventId       int64                  `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	Zones         []*Zone                `protobuf:"bytes,2,rep,name=zones,proto3" json:"zones,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAvailableSeatsResponse) Reset() {
	*x = GetAvailableSeatsResponse{}
	mi := &file_proto_order_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAvailableSeatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableSeatsResponse) ProtoMessage() {}

func (x *GetAvailableSeatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableSeatsResponse.ProtoReflect.Descriptor instead.
func (*GetAvailableSeatsResponse) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{1}
}

func (x *GetAvailableSeatsResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *GetAvailableSeatsResponse) GetZones() []*Zone {
	if x != nil {
		return x.Zones
	}
	return nil
}

type Zone struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // Название зоны (например, "A", "B", "C" или пусто, если зон нет)
	Rows          []*Row                 `protobuf:"bytes,2,rep,name=rows,proto3" json:"rows,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Zone) Reset() {
	*x = Zone{}
	mi := &file_proto_order_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Zone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Zone) ProtoMessage() {}

func (x *Zone) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Zone.ProtoReflect.Descriptor instead.
func (*Zone) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{2}
}

func (x *Zone) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Zone) GetRows() []*Row {
	if x != nil {
		return x.Rows
	}
	return nil
}

type Row struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Number        int64                  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`      // Номер ряда (если рядов нет, то 0)
	Seats         []int64                `protobuf:"varint,2,rep,packed,name=seats,proto3" json:"seats,omitempty"` // Список свободных мест
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Row) Reset() {
	*x = Row{}
	mi := &file_proto_order_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Row) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Row) ProtoMessage() {}

func (x *Row) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Row.ProtoReflect.Descriptor instead.
func (*Row) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{3}
}

func (x *Row) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Row) GetSeats() []int64 {
	if x != nil {
		return x.Seats
	}
	return nil
}

// Для создания заказа
type CreateOrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	EventId       int64                  `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	Zone          string                 `protobuf:"bytes,2,opt,name=zone,proto3" json:"zone,omitempty"`   // Название зоны
	Row           int64                  `protobuf:"varint,3,opt,name=row,proto3" json:"row,omitempty"`    // Номер ряда
	Seat          int64                  `protobuf:"varint,4,opt,name=seat,proto3" json:"seat,omitempty"`  // Номер места
	Email         string                 `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"` // Email для уведомлений
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	mi := &file_proto_order_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{4}
}

func (x *CreateOrderRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *CreateOrderRequest) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

func (x *CreateOrderRequest) GetRow() int64 {
	if x != nil {
		return x.Row
	}
	return 0
}

func (x *CreateOrderRequest) GetSeat() int64 {
	if x != nil {
		return x.Seat
	}
	return 0
}

func (x *CreateOrderRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type CreateOrderResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"` // Статус создания заказа
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrderResponse) Reset() {
	*x = CreateOrderResponse{}
	mi := &file_proto_order_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderResponse) ProtoMessage() {}

func (x *CreateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderResponse.ProtoReflect.Descriptor instead.
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{5}
}

func (x *CreateOrderResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_order_proto protoreflect.FileDescriptor

const file_proto_order_proto_rawDesc = "" +
	"\n" +
	"\x11proto/order.proto\x12\x05order\"5\n" +
	"\x18GetAvailableSeatsRequest\x12\x19\n" +
	"\bevent_id\x18\x01 \x01(\x03R\aeventId\"Y\n" +
	"\x19GetAvailableSeatsResponse\x12\x19\n" +
	"\bevent_id\x18\x01 \x01(\x03R\aeventId\x12!\n" +
	"\x05zones\x18\x02 \x03(\v2\v.order.ZoneR\x05zones\":\n" +
	"\x04Zone\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1e\n" +
	"\x04rows\x18\x02 \x03(\v2\n" +
	".order.RowR\x04rows\"3\n" +
	"\x03Row\x12\x16\n" +
	"\x06number\x18\x01 \x01(\x03R\x06number\x12\x14\n" +
	"\x05seats\x18\x02 \x03(\x03R\x05seats\"\x7f\n" +
	"\x12CreateOrderRequest\x12\x19\n" +
	"\bevent_id\x18\x01 \x01(\x03R\aeventId\x12\x12\n" +
	"\x04zone\x18\x02 \x01(\tR\x04zone\x12\x10\n" +
	"\x03row\x18\x03 \x01(\x03R\x03row\x12\x12\n" +
	"\x04seat\x18\x04 \x01(\x03R\x04seat\x12\x14\n" +
	"\x05email\x18\x05 \x01(\tR\x05email\"-\n" +
	"\x13CreateOrderResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status2\xac\x01\n" +
	"\fOrderService\x12V\n" +
	"\x11GetAvailableSeats\x12\x1f.order.GetAvailableSeatsRequest\x1a .order.GetAvailableSeatsResponse\x12D\n" +
	"\vCreateOrder\x12\x19.order.CreateOrderRequest\x1a\x1a.order.CreateOrderResponseB\x14Z\x12./proto/grpc/orderb\x06proto3"

var (
	file_proto_order_proto_rawDescOnce sync.Once
	file_proto_order_proto_rawDescData []byte
)

func file_proto_order_proto_rawDescGZIP() []byte {
	file_proto_order_proto_rawDescOnce.Do(func() {
		file_proto_order_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_order_proto_rawDesc), len(file_proto_order_proto_rawDesc)))
	})
	return file_proto_order_proto_rawDescData
}

var file_proto_order_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_order_proto_goTypes = []any{
	(*GetAvailableSeatsRequest)(nil),  // 0: order.GetAvailableSeatsRequest
	(*GetAvailableSeatsResponse)(nil), // 1: order.GetAvailableSeatsResponse
	(*Zone)(nil),                      // 2: order.Zone
	(*Row)(nil),                       // 3: order.Row
	(*CreateOrderRequest)(nil),        // 4: order.CreateOrderRequest
	(*CreateOrderResponse)(nil),       // 5: order.CreateOrderResponse
}
var file_proto_order_proto_depIdxs = []int32{
	2, // 0: order.GetAvailableSeatsResponse.zones:type_name -> order.Zone
	3, // 1: order.Zone.rows:type_name -> order.Row
	0, // 2: order.OrderService.GetAvailableSeats:input_type -> order.GetAvailableSeatsRequest
	4, // 3: order.OrderService.CreateOrder:input_type -> order.CreateOrderRequest
	1, // 4: order.OrderService.GetAvailableSeats:output_type -> order.GetAvailableSeatsResponse
	5, // 5: order.OrderService.CreateOrder:output_type -> order.CreateOrderResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_order_proto_init() }
func file_proto_order_proto_init() {
	if File_proto_order_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_order_proto_rawDesc), len(file_proto_order_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_order_proto_goTypes,
		DependencyIndexes: file_proto_order_proto_depIdxs,
		MessageInfos:      file_proto_order_proto_msgTypes,
	}.Build()
	File_proto_order_proto = out.File
	file_proto_order_proto_goTypes = nil
	file_proto_order_proto_depIdxs = nil
}
