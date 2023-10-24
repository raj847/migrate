// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: product.proto

package product

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

type ProductOuWithRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OuId             int64   `protobuf:"varint,1,opt,name=ouId,proto3" json:"ouId,omitempty"`
	ProductId        int64   `protobuf:"varint,2,opt,name=productId,proto3" json:"productId,omitempty"`
	Price            float64 `protobuf:"fixed64,3,opt,name=price,proto3" json:"price,omitempty"`
	BaseTime         int64   `protobuf:"varint,4,opt,name=baseTime,proto3" json:"baseTime,omitempty"`
	ProgressiveTime  int64   `protobuf:"varint,5,opt,name=progressiveTime,proto3" json:"progressiveTime,omitempty"`
	ProgressivePrice float64 `protobuf:"fixed64,6,opt,name=progressivePrice,proto3" json:"progressivePrice,omitempty"`
	IsPct            string  `protobuf:"bytes,7,opt,name=isPct,proto3" json:"isPct,omitempty"`
	ProgressivePct   int64   `protobuf:"varint,8,opt,name=progressivePct,proto3" json:"progressivePct,omitempty"`
	MaxPrice         float64 `protobuf:"fixed64,9,opt,name=maxPrice,proto3" json:"maxPrice,omitempty"`
	Is24H            string  `protobuf:"bytes,10,opt,name=is24h,proto3" json:"is24h,omitempty"`
	OvernightTime    string  `protobuf:"bytes,11,opt,name=overnightTime,proto3" json:"overnightTime,omitempty"`
	OvernightPrice   float64 `protobuf:"fixed64,12,opt,name=overnightPrice,proto3" json:"overnightPrice,omitempty"`
	GracePeriod      int64   `protobuf:"varint,13,opt,name=gracePeriod,proto3" json:"gracePeriod,omitempty"`
	FlgRepeat        string  `protobuf:"bytes,14,opt,name=flgRepeat,proto3" json:"flgRepeat,omitempty"`
}

func (x *ProductOuWithRules) Reset() {
	*x = ProductOuWithRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductOuWithRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductOuWithRules) ProtoMessage() {}

func (x *ProductOuWithRules) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductOuWithRules.ProtoReflect.Descriptor instead.
func (*ProductOuWithRules) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{0}
}

func (x *ProductOuWithRules) GetOuId() int64 {
	if x != nil {
		return x.OuId
	}
	return 0
}

func (x *ProductOuWithRules) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *ProductOuWithRules) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ProductOuWithRules) GetBaseTime() int64 {
	if x != nil {
		return x.BaseTime
	}
	return 0
}

func (x *ProductOuWithRules) GetProgressiveTime() int64 {
	if x != nil {
		return x.ProgressiveTime
	}
	return 0
}

func (x *ProductOuWithRules) GetProgressivePrice() float64 {
	if x != nil {
		return x.ProgressivePrice
	}
	return 0
}

func (x *ProductOuWithRules) GetIsPct() string {
	if x != nil {
		return x.IsPct
	}
	return ""
}

func (x *ProductOuWithRules) GetProgressivePct() int64 {
	if x != nil {
		return x.ProgressivePct
	}
	return 0
}

func (x *ProductOuWithRules) GetMaxPrice() float64 {
	if x != nil {
		return x.MaxPrice
	}
	return 0
}

func (x *ProductOuWithRules) GetIs24H() string {
	if x != nil {
		return x.Is24H
	}
	return ""
}

func (x *ProductOuWithRules) GetOvernightTime() string {
	if x != nil {
		return x.OvernightTime
	}
	return ""
}

func (x *ProductOuWithRules) GetOvernightPrice() float64 {
	if x != nil {
		return x.OvernightPrice
	}
	return 0
}

func (x *ProductOuWithRules) GetGracePeriod() int64 {
	if x != nil {
		return x.GracePeriod
	}
	return 0
}

func (x *ProductOuWithRules) GetFlgRepeat() string {
	if x != nil {
		return x.FlgRepeat
	}
	return ""
}

type PolicyOuProductWithRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OuId                  int64               `protobuf:"varint,1,opt,name=ouId,proto3" json:"ouId,omitempty"`
	OuCode                string              `protobuf:"bytes,2,opt,name=ouCode,proto3" json:"ouCode,omitempty"`
	OuName                string              `protobuf:"bytes,3,opt,name=ouName,proto3" json:"ouName,omitempty"`
	ProductId             int64               `protobuf:"varint,4,opt,name=productId,proto3" json:"productId,omitempty"`
	ProductCode           string              `protobuf:"bytes,5,opt,name=productCode,proto3" json:"productCode,omitempty"`
	ProductName           string              `protobuf:"bytes,6,opt,name=productName,proto3" json:"productName,omitempty"`
	ServiceFee            float64             `protobuf:"fixed64,7,opt,name=serviceFee,proto3" json:"serviceFee,omitempty"`
	IsPctServiceFee       string              `protobuf:"bytes,8,opt,name=isPctServiceFee,proto3" json:"isPctServiceFee,omitempty"`
	IsPctServiceFeeMember string              `protobuf:"bytes,9,opt,name=isPctServiceFeeMember,proto3" json:"isPctServiceFeeMember,omitempty"`
	ServiceFeeMember      float64             `protobuf:"fixed64,10,opt,name=serviceFeeMember,proto3" json:"serviceFeeMember,omitempty"`
	ProductRules          *ProductOuWithRules `protobuf:"bytes,11,opt,name=productRules,proto3" json:"productRules,omitempty"`
}

func (x *PolicyOuProductWithRules) Reset() {
	*x = PolicyOuProductWithRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyOuProductWithRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyOuProductWithRules) ProtoMessage() {}

func (x *PolicyOuProductWithRules) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PolicyOuProductWithRules.ProtoReflect.Descriptor instead.
func (*PolicyOuProductWithRules) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{1}
}

func (x *PolicyOuProductWithRules) GetOuId() int64 {
	if x != nil {
		return x.OuId
	}
	return 0
}

func (x *PolicyOuProductWithRules) GetOuCode() string {
	if x != nil {
		return x.OuCode
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetOuName() string {
	if x != nil {
		return x.OuName
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *PolicyOuProductWithRules) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetServiceFee() float64 {
	if x != nil {
		return x.ServiceFee
	}
	return 0
}

func (x *PolicyOuProductWithRules) GetIsPctServiceFee() string {
	if x != nil {
		return x.IsPctServiceFee
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetIsPctServiceFeeMember() string {
	if x != nil {
		return x.IsPctServiceFeeMember
	}
	return ""
}

func (x *PolicyOuProductWithRules) GetServiceFeeMember() float64 {
	if x != nil {
		return x.ServiceFeeMember
	}
	return 0
}

func (x *PolicyOuProductWithRules) GetProductRules() *ProductOuWithRules {
	if x != nil {
		return x.ProductRules
	}
	return nil
}

type RequestGetPolicyOuProductList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keyword         string `protobuf:"bytes,1,opt,name=keyword,proto3" json:"keyword,omitempty"`
	AscDesc         string `protobuf:"bytes,2,opt,name=ascDesc,proto3" json:"ascDesc,omitempty"`
	ColumnOrderName string `protobuf:"bytes,3,opt,name=columnOrderName,proto3" json:"columnOrderName,omitempty"`
}

func (x *RequestGetPolicyOuProductList) Reset() {
	*x = RequestGetPolicyOuProductList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestGetPolicyOuProductList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestGetPolicyOuProductList) ProtoMessage() {}

func (x *RequestGetPolicyOuProductList) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestGetPolicyOuProductList.ProtoReflect.Descriptor instead.
func (*RequestGetPolicyOuProductList) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{2}
}

func (x *RequestGetPolicyOuProductList) GetKeyword() string {
	if x != nil {
		return x.Keyword
	}
	return ""
}

func (x *RequestGetPolicyOuProductList) GetAscDesc() string {
	if x != nil {
		return x.AscDesc
	}
	return ""
}

func (x *RequestGetPolicyOuProductList) GetColumnOrderName() string {
	if x != nil {
		return x.ColumnOrderName
	}
	return ""
}

type ResponseGetPolicyOuProductList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CountProductList int64                       `protobuf:"varint,1,opt,name=countProductList,proto3" json:"countProductList,omitempty"`
	Data             []*PolicyOuProductWithRules `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ResponseGetPolicyOuProductList) Reset() {
	*x = ResponseGetPolicyOuProductList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseGetPolicyOuProductList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseGetPolicyOuProductList) ProtoMessage() {}

func (x *ResponseGetPolicyOuProductList) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseGetPolicyOuProductList.ProtoReflect.Descriptor instead.
func (*ResponseGetPolicyOuProductList) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseGetPolicyOuProductList) GetCountProductList() int64 {
	if x != nil {
		return x.CountProductList
	}
	return 0
}

func (x *ResponseGetPolicyOuProductList) GetData() []*PolicyOuProductWithRules {
	if x != nil {
		return x.Data
	}
	return nil
}

type ResponseTrxProduct struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductCode string `protobuf:"bytes,1,opt,name=productCode,proto3" json:"productCode,omitempty"`
	ProductName string `protobuf:"bytes,2,opt,name=productName,proto3" json:"productName,omitempty"`
}

func (x *ResponseTrxProduct) Reset() {
	*x = ResponseTrxProduct{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseTrxProduct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseTrxProduct) ProtoMessage() {}

func (x *ResponseTrxProduct) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseTrxProduct.ProtoReflect.Descriptor instead.
func (*ResponseTrxProduct) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{4}
}

func (x *ResponseTrxProduct) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *ResponseTrxProduct) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

type ProductOuDepositCounterWithRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OuId      int64   `protobuf:"varint,1,opt,name=ouId,proto3" json:"ouId,omitempty"`
	ProductId int64   `protobuf:"varint,2,opt,name=productId,proto3" json:"productId,omitempty"`
	Price     float64 `protobuf:"fixed64,3,opt,name=price,proto3" json:"price,omitempty"`
	IsPct     string  `protobuf:"bytes,4,opt,name=isPct,proto3" json:"isPct,omitempty"`
}

func (x *ProductOuDepositCounterWithRules) Reset() {
	*x = ProductOuDepositCounterWithRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductOuDepositCounterWithRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductOuDepositCounterWithRules) ProtoMessage() {}

func (x *ProductOuDepositCounterWithRules) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductOuDepositCounterWithRules.ProtoReflect.Descriptor instead.
func (*ProductOuDepositCounterWithRules) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{5}
}

func (x *ProductOuDepositCounterWithRules) GetOuId() int64 {
	if x != nil {
		return x.OuId
	}
	return 0
}

func (x *ProductOuDepositCounterWithRules) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *ProductOuDepositCounterWithRules) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ProductOuDepositCounterWithRules) GetIsPct() string {
	if x != nil {
		return x.IsPct
	}
	return ""
}

type PolicyOuProductDepositCounterWithRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OuId                       int64                             `protobuf:"varint,1,opt,name=ouId,proto3" json:"ouId,omitempty"`
	OuCode                     string                            `protobuf:"bytes,2,opt,name=ouCode,proto3" json:"ouCode,omitempty"`
	OuName                     string                            `protobuf:"bytes,3,opt,name=ouName,proto3" json:"ouName,omitempty"`
	ProductId                  int64                             `protobuf:"varint,4,opt,name=productId,proto3" json:"productId,omitempty"`
	ProductCode                string                            `protobuf:"bytes,5,opt,name=productCode,proto3" json:"productCode,omitempty"`
	ProductName                string                            `protobuf:"bytes,6,opt,name=productName,proto3" json:"productName,omitempty"`
	ServiceFee                 float64                           `protobuf:"fixed64,7,opt,name=serviceFee,proto3" json:"serviceFee,omitempty"`
	IsPctServiceFee            string                            `protobuf:"bytes,8,opt,name=isPctServiceFee,proto3" json:"isPctServiceFee,omitempty"`
	ServiceFeeMember           float64                           `protobuf:"fixed64,9,opt,name=serviceFeeMember,proto3" json:"serviceFeeMember,omitempty"`
	IsPctServiceFeeMember      string                            `protobuf:"bytes,10,opt,name=isPctServiceFeeMember,proto3" json:"isPctServiceFeeMember,omitempty"`
	ProductDepositCounterRules *ProductOuDepositCounterWithRules `protobuf:"bytes,11,opt,name=productDepositCounterRules,proto3" json:"productDepositCounterRules,omitempty"`
}

func (x *PolicyOuProductDepositCounterWithRules) Reset() {
	*x = PolicyOuProductDepositCounterWithRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyOuProductDepositCounterWithRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyOuProductDepositCounterWithRules) ProtoMessage() {}

func (x *PolicyOuProductDepositCounterWithRules) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PolicyOuProductDepositCounterWithRules.ProtoReflect.Descriptor instead.
func (*PolicyOuProductDepositCounterWithRules) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{6}
}

func (x *PolicyOuProductDepositCounterWithRules) GetOuId() int64 {
	if x != nil {
		return x.OuId
	}
	return 0
}

func (x *PolicyOuProductDepositCounterWithRules) GetOuCode() string {
	if x != nil {
		return x.OuCode
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetOuName() string {
	if x != nil {
		return x.OuName
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *PolicyOuProductDepositCounterWithRules) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetServiceFee() float64 {
	if x != nil {
		return x.ServiceFee
	}
	return 0
}

func (x *PolicyOuProductDepositCounterWithRules) GetIsPctServiceFee() string {
	if x != nil {
		return x.IsPctServiceFee
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetServiceFeeMember() float64 {
	if x != nil {
		return x.ServiceFeeMember
	}
	return 0
}

func (x *PolicyOuProductDepositCounterWithRules) GetIsPctServiceFeeMember() string {
	if x != nil {
		return x.IsPctServiceFeeMember
	}
	return ""
}

func (x *PolicyOuProductDepositCounterWithRules) GetProductDepositCounterRules() *ProductOuDepositCounterWithRules {
	if x != nil {
		return x.ProductDepositCounterRules
	}
	return nil
}

var File_product_proto protoreflect.FileDescriptor

var file_product_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0xcc, 0x03, 0x0a, 0x12, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x4f, 0x75, 0x57, 0x69, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6f,
	0x75, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x70, 0x72,
	0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2a, 0x0a,
	0x10, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x76, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x50,
	0x63, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x73, 0x50, 0x63, 0x74, 0x12,
	0x26, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x50, 0x63,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x76, 0x65, 0x50, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6d, 0x61, 0x78, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x32, 0x34, 0x68, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x69, 0x73, 0x32, 0x34, 0x68, 0x12, 0x24, 0x0a, 0x0d, 0x6f, 0x76, 0x65,
	0x72, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x26, 0x0a, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x69, 0x67,
	0x68, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x67, 0x72, 0x61, 0x63, 0x65,
	0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x67, 0x72,
	0x61, 0x63, 0x65, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x6c, 0x67,
	0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x6c,
	0x67, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x22, 0xad, 0x03, 0x0a, 0x18, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x4f, 0x75, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x57, 0x69, 0x74, 0x68, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x43, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6f, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x73,
	0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x46, 0x65, 0x65, 0x12, 0x34, 0x0a, 0x15, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x15, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x46, 0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4f, 0x75,
	0x57, 0x69, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x22, 0x7d, 0x0a, 0x1d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4f, 0x75, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6b, 0x65, 0x79, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x73, 0x63, 0x44, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x73, 0x63, 0x44, 0x65, 0x73, 0x63, 0x12, 0x28, 0x0a, 0x0f,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x83, 0x01, 0x0a, 0x1e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4f, 0x75, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x10, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x4f, 0x75, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x57, 0x69, 0x74,
	0x68, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x58, 0x0a, 0x12,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x72, 0x78, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x80, 0x01, 0x0a, 0x20, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x4f, 0x75, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6f,
	0x75, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x50, 0x63, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x69, 0x73, 0x50, 0x63, 0x74, 0x22, 0xe5, 0x03, 0x0a, 0x26, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x4f, 0x75, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x44, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x6f, 0x75, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x43, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6f, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x73,
	0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x46, 0x65, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46,
	0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x34, 0x0a, 0x15, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x46, 0x65, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x15, 0x69, 0x73, 0x50, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x46, 0x65, 0x65,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x69, 0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4f, 0x75, 0x44, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x1a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x44, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x75, 0x6c, 0x65,
	0x73, 0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_product_proto_rawDescOnce sync.Once
	file_product_proto_rawDescData = file_product_proto_rawDesc
)

func file_product_proto_rawDescGZIP() []byte {
	file_product_proto_rawDescOnce.Do(func() {
		file_product_proto_rawDescData = protoimpl.X.CompressGZIP(file_product_proto_rawDescData)
	})
	return file_product_proto_rawDescData
}

var file_product_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_product_proto_goTypes = []interface{}{
	(*ProductOuWithRules)(nil),                     // 0: product.ProductOuWithRules
	(*PolicyOuProductWithRules)(nil),               // 1: product.PolicyOuProductWithRules
	(*RequestGetPolicyOuProductList)(nil),          // 2: product.RequestGetPolicyOuProductList
	(*ResponseGetPolicyOuProductList)(nil),         // 3: product.ResponseGetPolicyOuProductList
	(*ResponseTrxProduct)(nil),                     // 4: product.ResponseTrxProduct
	(*ProductOuDepositCounterWithRules)(nil),       // 5: product.ProductOuDepositCounterWithRules
	(*PolicyOuProductDepositCounterWithRules)(nil), // 6: product.PolicyOuProductDepositCounterWithRules
}
var file_product_proto_depIdxs = []int32{
	0, // 0: product.PolicyOuProductWithRules.productRules:type_name -> product.ProductOuWithRules
	1, // 1: product.ResponseGetPolicyOuProductList.data:type_name -> product.PolicyOuProductWithRules
	5, // 2: product.PolicyOuProductDepositCounterWithRules.productDepositCounterRules:type_name -> product.ProductOuDepositCounterWithRules
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_product_proto_init() }
func file_product_proto_init() {
	if File_product_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_product_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductOuWithRules); i {
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
		file_product_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PolicyOuProductWithRules); i {
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
		file_product_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestGetPolicyOuProductList); i {
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
		file_product_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseGetPolicyOuProductList); i {
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
		file_product_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseTrxProduct); i {
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
		file_product_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductOuDepositCounterWithRules); i {
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
		file_product_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PolicyOuProductDepositCounterWithRules); i {
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
			RawDescriptor: file_product_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_product_proto_goTypes,
		DependencyIndexes: file_product_proto_depIdxs,
		MessageInfos:      file_product_proto_msgTypes,
	}.Build()
	File_product_proto = out.File
	file_product_proto_rawDesc = nil
	file_product_proto_goTypes = nil
	file_product_proto_depIdxs = nil
}
