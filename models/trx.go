package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestTrxCheckIn struct {
	CheckInDatetime  string  `json:"checkInDatetime"`
	ProductCode      string  `json:"productCode"`
	DeviceId         string  `json:"deviceId" validate:"required"`
	IpTerminal       string  `json:"ipTerminal" validate:"required"`
	CardNumber       string  `json:"cardNumber" validate:"required"`
	UUIDCard         string  `json:"uuidCard" validate:"required,min=14,max=14"`
	TypeCard         string  `json:"typeCard" validate:"required"`
	BeginningBalance float64 `json:"beginningBalance"`
}

type RequestTrxCheckInWithoutCard struct {
	CheckInDatetime string `json:"checkInDatetime" validate:"required"`
	ProductCode     string `json:"productCode"`
	IpTerminal      string `json:"ipTerminal"`
	RefId           string `json:"refId"`
}

type Trx struct {
	DocNo                          string                 `json:"docNo" bson:"docNo"`
	DocDate                        string                 `json:"docDate" bson:"docDate"`
	PaymentRefDocNo                string                 `json:"paymentRefDocNo" bson:"paymentRefDocNo"`
	CheckInDatetime                string                 `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime               string                 `json:"checkOutDatetime" bson:"checkOutDatetime"`
	DeviceIdIn                     string                 `json:"deviceIdIn" bson:"deviceIdIn"`
	DeviceId                       string                 `json:"device_id" bson:"deviceId"`
	GateIn                         string                 `json:"gateIn" bson:"gateIn"`
	GateOut                        string                 `json:"gateOut" bson:"gateOut"`
	CardNumberUUIDIn               string                 `json:"cardNumberUuidIn" bson:"cardNumberUuidIn"`
	CardNumberIn                   string                 `json:"cardNumberIn" bson:"cardNumberIn"`
	CardNumberUUID                 string                 `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber                     string                 `json:"cardNumber" bson:"cardNumber"`
	TypeCard                       string                 `json:"typeCard" bson:"typeCard"`
	BeginningBalance               float64                `json:"beginningBalance" bson:"beginningBalance"`
	ExtLocalDatetime               string                 `json:"extLocalDatetime" bson:"extLocalDatetime"`
	ChargeAmount                   float64                `json:"chargeAmount" bson:"chargeAmount"`
	GrandTotal                     float64                `json:"grandTotal" bson:"grandTotal"`
	ProductCode                    string                 `json:"productCode" bson:"productCode"`
	ProductName                    string                 `json:"productName" bson:"productName"`
	ProductData                    string                 `json:"productData" bson:"productData"`
	RequestData                    string                 `json:"requestData" bson:"requestData"`
	RequestOutData                 string                 `json:"requestOutData" bson:"requestOutData"`
	OuId                           int64                  `json:"ouId" bson:"ouId"`
	OuName                         string                 `json:"ouName" bson:"ouName"`
	OuCode                         string                 `json:"ouCode" bson:"ouCode"`
	OuSubBranchId                  int64                  `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName                string                 `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode                string                 `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId                       int64                  `json:"mainOuId" bson:"mainOuId"`
	MainOuCode                     string                 `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName                     string                 `json:"mainOuName" bson:"mainOuName"`
	MemberCode                     string                 `json:"memberCode" bson:"memberCode"`
	MemberName                     string                 `json:"memberName" bson:"memberName"`
	MemberType                     string                 `json:"memberType" bson:"memberType"`
	MemberStatus                   string                 `json:"memberStatus" bson:"memberStatus"`
	MemberExpiredDate              string                 `json:"memberExpiredDate" bson:"memberExpiredDate"`
	CheckInTime                    int64                  `json:"checkInTime" bson:"checkInTime"`
	CheckOutTime                   int64                  `json:"checkOutTime" bson:"checkOutTime"`
	DurationTime                   int64                  `json:"durationTime" bson:"durationTime"`
	VehicleNumberIn                string                 `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	VehicleNumberOut               string                 `json:"vehicleNumberOut" bson:"vehicleNumberOut"`
	LogTrans                       string                 `json:"logTrans" bson:"logTrans"`
	MerchantKey                    string                 `json:"merchantKey" bson:"merchantKey"`
	QrText                         string                 `json:"qrText" bson:"qrText"`
	QrA2P                          string                 `json:"qrA2P" bson:"qrA2P"`
	QrTextPaymentOnline            string                 `json:"qrTextPaymentOnline" bson:"qrTextPaymentOnline"`
	TrxInvoiceItem                 []TrxInvoiceItem       `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
	FlagSyncData                   bool                   `json:"flagSyncData" bson:"flagSyncData"`
	MemberData                     *TrxMember             `json:"memberData" bson:"memberData"`
	TrxAddInfo                     map[string]interface{} `json:"trxAddInfo" bson:"trxAddInfo"`
	FlagTrxFromCloud               bool                   `json:"flagTrxFromCloud" bson:"flagTrxFromCloud"`
	IsRsyncDataTrx                 bool                   `json:"isRsyncDataTrx" bson:"isRsyncDataTrx"`
	ExcludeSf                      bool                   `json:"excludeSf" bson:"excludeSf"`
	FlagCharge                     bool                   `json:"flagCharge" bson:"flagCharge"`
	ChargeType                     string                 `json:"chargeType" bson:"chargeType"`
	RequestAddTrxInvoiceDetailItem *TrxInvoiceDetailItem  `json:"requestAddTrxInvoiceDetailItem" bson:"requestAddTrxInvoiceDetailItem"`
	LastUpdatedAt                  string                 `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
}

type ResponseTrxTicket struct {
	ID                   *primitive.ObjectID `json:"_id"`
	CheckInDatetime      string              `json:"checkInDatetime"`
	DocNo                string              `json:"docNo"`
	ProductName          string              `json:"productName"`
	VehicleNumberIn      string              `json:"vehicleNumberIn"`
	QRCode               string              `json:"qrCode"`
	QRCodePaymentOnlinet string              `json:"qrCodePaymentOnline"`
	OuCode               string              `json:"ouCode"`
	OuName               string              `json:"ouName"`
	Address              string              `json:"address"`
}

type RequestInquiryWithoutCard struct {
	QRCode          string `json:"qrCode"`
	ProductCode     string `json:"productCode"`
	InquiryDatetime string `json:"inquiryDatetime"`
	TerminalId      string `json:"terminalId"`
}

type RequestCustomInquiryWithoutCard struct {
	QRCode          string `json:"qrCode"`
	ProductCode     string `json:"productCode"`
	InquiryDatetime string `json:"inquiryDatetime"`
	TerminalId      string `json:"terminalId"`
	VehicleNumber   string `json:"vehicleNumber" validate:"required"`
}

type RequestInquiryTrxDepoWithoutCard struct {
	DocNoDepo       string `json:"docNoDepo"`
	ProductCode     string `json:"productCode"`
	InquiryDatetime string `json:"inquiryDatetime"`
	TerminalId      string `json:"terminalId"`
}

type RequestInquiryWithCard struct {
	UUIDCard        string `json:"uuidCard" validate:"required"`
	ProductCode     string `json:"productCode"`
	InquiryDatetime string `json:"inquiryDatetime"`
	TerminalId      string `json:"terminalId"`
}

type RequestInquiryWithCardP3 struct {
	CardNumber       string  `json:"cardNumber"`
	UUIDCard         string  `json:"uuidCard" validate:"required"`
	TypeCard         string  `json:"typeCard"`
	BeginningBalance float64 `json:"beginningBalance"`
	ProductCode      string  `json:"productCode"`
	CheckInDatetime  string  `json:"checkInDatetime"`
	DeviceId         string  `json:"deviceId"`
	TerminalId       string  `json:"terminalId"`
}

type RequestInquiryWithCardCustom struct {
	UUIDCard        string `json:"uuidCard" validate:"required"`
	VehicleNumber   string `json:"vehicleNumber"`
	ProductCode     string `json:"productCode"`
	InquiryDatetime string `json:"inquiryDatetime"`
	TerminalId      string `json:"terminalId"`
}

type ResultFindTrxOutstanding struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	DocNo           string             `json:"docNo" bson:"docNo"`
	GrandTotal      float64            `json:"grandTotal" bson:"grandTotal"`
	CheckInDatetime string             `json:"checkInDatetime" bson:"checkInDatetime"`
	OverNightPrice  float64            `json:"overnightPrice" bson:"overnightPrice"`
	Is24H           string             `json:"is24H" bson:"is24H"`
	CardNumber      string             `json:"cardNumber" bson:"cardNumber"`
	CardNumberUUID  string             `json:"cardNumberUuid" bson:"cardNumberUuid"`
	OuCode          string             `json:"ouCode" bson:"ouCode"`
	VehicleNumberIn string             `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	TrxInvoiceItem  []TrxInvoiceItem   `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
}

type ResultInquiryTrx struct {
	ID              string      `json:"_id" bson:"_id"`
	DocNo           string      `json:"docNo" bson:"docNo"`
	Nominal         float64     `json:"nominal" bson:"nominal"`
	ProductCode     string      `json:"productCode" bson:"productCode"`
	ProductName     string      `json:"productName" bson:"productName"`
	VehicleNumberIn string      `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	QRCode          string      `json:"qrCode" bson:"qrCode"`
	ExcludeSf       bool        `json:"excludeSf" bson:"excludeSf"`
	Duration        interface{} `json:"duration" bson:"duration"`
	InvoiceDetail   interface{} `json:"invoiceDetail" bson:"invoiceDetail"`
	OuCode          string      `json:"ouCode"`
}

type ResultInquiryTrxCustomV2 struct {
	ID                string      `json:"_id" bson:"_id"`
	DocNo             string      `json:"docNo" bson:"docNo"`
	Nominal           float64     `json:"nominal" bson:"nominal"`
	ProductCode       string      `json:"productCode" bson:"productCode"`
	ProductName       string      `json:"productName" bson:"productName"`
	VehicleNumberIn   string      `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	MemberCode        string      `json:"memberCode" bson:"memberCode"`
	MemberName        string      `json:"memberName" bson:"memberName"`
	MemberExpiredDate string      `json:"memberExpiredDate" bson:"memberExpiredDate"`
	MemberType        string      `json:"memberType" bson:"memberType"`
	MemberStatus      string      `json:"memberStatus" bson:"memberStatus"`
	QRCode            string      `json:"qrCode" bson:"qrCode"`
	ExcludeSf         bool        `json:"excludeSf" bson:"excludeSf"`
	Duration          interface{} `json:"duration" bson:"duration"`
	InvoiceDetail     interface{} `json:"invoiceDetail" bson:"invoiceDetail"`
	OuCode            string      `json:"ouCode"`
}

type ResultInquiryTrxCustom struct {
	ID              string      `json:"_id" bson:"_id"`
	DocNo           string      `json:"docNo" bson:"docNo"`
	Nominal         float64     `json:"nominal" bson:"nominal"`
	ProductCode     string      `json:"productCode" bson:"productCode"`
	ProductName     string      `json:"productName" bson:"productName"`
	VehicleNumberIn string      `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	QRCode          string      `json:"qrCode" bson:"qrCode"`
	ExcludeSf       bool        `json:"excludeSf" bson:"excludeSf"`
	Duration        interface{} `json:"duration" bson:"duration"`
	MemberCode      string      `json:"memberCode"`
	MemberName      string      `json:"memberName"`
	MemberType      string      `json:"memberType"`
	OuCode          string      `json:"ouCode"`
}

type ResultInquiryTrxWithCard struct {
	ID                string      `json:"_id" bson:"_id"`
	DocNo             string      `json:"docNo" bson:"docNo"`
	CardNumberUUID    string      `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber        string      `json:"cardNumber" bson:"cardNumber"`
	Nominal           float64     `json:"nominal" bson:"nominal"`
	ProductCode       string      `json:"productCode" bson:"productCode"`
	ProductName       string      `json:"productName" bson:"productName"`
	VehicleNumberIn   string      `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	QRCode            string      `json:"qrCode" bson:"qrCode"`
	ExcludeSf         bool        `json:"excludeSf" bson:"excludeSf"`
	Duration          interface{} `json:"duration" bson:"duration"`
	Duration24H       interface{} `json:"duration24h" bson:"duration24h"`
	MemberCode        string      `json:"memberCode"`
	MemberName        string      `json:"memberName"`
	MemberType        string      `json:"memberType"`
	MemberExpiredDate string      `json:"memberExpiredDate" bson:"memberExpiredDate"`
	MemberStatus      string      `json:"memberStatus" bson:"memberStatus"`
	InvoiceDetail     interface{} `json:"invoiceDetail" bson:"invoiceDetail"`
	OuCode            string      `json:"ouCode"`
}

type RequestConfirmTrx struct {
	ID               string  `json:"id" bson:"_id" validate:"required"`
	DeviceId         string  `json:"deviceId" bson:"deviceId"`
	IpTerminal       string  `json:"ipTerminal" bson:"ipTerminal"`
	CardNumber       string  `json:"cardNumber" bson:"cardNumber"`
	CardType         string  `json:"cardType" bson:"cardType"`
	LastBalance      int64   `json:"lastBalance" bson:"lastBalance"`
	CurrentBalance   int64   `json:"currentBalance" bson:"currentBalance"`
	UUIDCard         string  `json:"uuidCard" bson:"uuidCard"`
	ProductCode      string  `json:"productCode" validate:"required"`
	ProductName      string  `json:"productName" validate:"required"`
	CheckOutDatetime string  `json:"checkOutDatetime" bson:"checkOutDatetime" validate:"required"`
	GrandTotal       float64 `json:"grandTotal" bson:"grandTotal" validate:"min=0"`
	LogTrans         string  `json:"logTrans" bson:"logTrans"`
	VehicleNumber    string  `json:"vehicleNumber" bson:"vehicleNumber"`
	Username         string  `json:"username" bson:"username"`
	ShiftCode        string  `json:"shiftCode" bson:"shiftCode"`
	ExcludeSf        *bool   `json:"excludeSf" bson:"excludeSf"`
}

type RequestConfirmTrxDepositCounter struct {
	ID               string  `json:"id" bson:"_id" validate:"required"`
	DeviceId         string  `json:"deviceId" bson:"deviceId"`
	IpTerminal       string  `json:"ipTerminal" bson:"ipTerminal"`
	CardNumber       string  `json:"cardNumber" bson:"cardNumber"`
	CardType         string  `json:"cardType" bson:"cardType"`
	LastBalance      int64   `json:"lastBalance" bson:"lastBalance"`
	CurrentBalance   int64   `json:"currentBalance" bson:"currentBalance"`
	UUIDCard         string  `json:"uuidCard" bson:"uuidCard"`
	ProductCode      string  `json:"productCode" validate:"required"`
	ProductName      string  `json:"productName" validate:"required"`
	CheckOutDatetime string  `json:"checkOutDatetime" bson:"checkOutDatetime" validate:"required"`
	GrandTotal       float64 `json:"grandTotal" bson:"grandTotal" validate:"min=0"`
	LogTrans         string  `json:"logTrans" bson:"logTrans"`
	Username         string  `json:"username" bson:"username"`
	ShiftCode        string  `json:"shiftCode" bson:"shiftCode"`
	ExcludeSf        *bool   `json:"excludeSf" bson:"excludeSf"`
}

type TrxCheckOut struct {
	TrxOutStanding Trx `json:"trxOutStanding" bson:"trxOutStanding"`
	Trx            Trx `json:"trx" bson:"trx"`
}

type LogActivityTrx struct {
	DocNo     string `json:"docNo" bson:"docNo"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	Remark    string `json:"remark" bson:"remark"`
}

type InvoiceTrx struct {
	DocNo       string  `json:"docNo" bson:"docNo"`
	CreatedBy   string  `json:"createdBy" bson:"createdBy"`
	CreatedDate string  `json:"createdDate" bson:"createdDate"`
	TotalAmount float64 `json:"totalAmount" bson:"totalAmount"`
	TypeInvoice string  `json:"typeInvoice" bson:"typeInvoice"`
}

type TrxInvoiceItem struct {
	DocNo                  string  `json:"docNo" bson:"docNo"`
	ProductId              int64   `json:"productId" bson:"productId"`
	ProductCode            string  `json:"productCode" bson:"productCode"`
	ProductName            string  `json:"productName" bson:"productName"`
	IsPctServiceFee        string  `json:"isPctServiceFee" bson:"isPctServiceFee"`
	ServiceFee             float64 `json:"serviceFee" bson:"serviceFee"`
	ServiceFeeMember       float64 `json:"serviceFeeMember" bson:"serviceFeeMember"`
	Price                  float64 `json:"price" bson:"price"`
	BaseTime               int64   `json:"baseTime" bson:"baseTime"`
	ProgressiveTime        int64   `json:"progressiveTime" bson:"progressiveTime"`
	ProgressivePrice       float64 `json:"progressivePrice" bson:"progressivePrice"`
	IsPct                  string  `json:"isPct" bson:"isPct"`
	ProgressivePct         int64   `json:"progressivePct" bson:"progressivePct"`
	MaxPrice               float64 `json:"maxPrice" bson:"maxPrice"`
	Is24H                  string  `json:"is24H" bson:"is24H"`
	OvernightTime          string  `json:"overnightTime" bson:"overnightTime"`
	OvernightPrice         float64 `json:"overnightPrice" bson:"overnightPrice"`
	GracePeriod            int64   `json:"gracePeriod" bson:"gracePeriod"`
	FlgRepeat              string  `json:"flgRepeat" bson:"flgRepeat"`
	TotalAmount            float64 `json:"totalAmount" bson:"totalAmount"`
	TotalProgressiveAmount float64 `json:"totalProgressiveAmount" bson:"totalProgressiveAmount"`
}

type TrxInvoiceDetailItem struct {
	DocNo         string  `json:"docNo" bson:"docNo"`
	ProductCode   string  `json:"productCode" bson:"productCode"`
	InvoiceAmount float64 `json:"invoiceAmount" bson:"invoiceAmount"`
	CreatedAt     string  `json:"createdAt" bson:"createdAt"`
	CreatedDate   string  `json:"createdDate" bson:"createdDate"`
}

type ResponseTrxInvoiceDetailsItemList struct {
	TotalAmount float64 `json:"invoiceAmount" bson:"invoiceAmount"`
}

type RequestInquiryRedis struct {
	DocNo           string  `json:"docNo" bson:"docNo"`
	ProductCode     string  `json:"productCode" bson:"productCode"`
	ProductName     string  `json:"productName" bson:"productName"`
	GrandTotal      float64 `json:"grandTotal" bson:"grandTotal"`
	OuCode          string  `json:"ouCode" bson:"ouCode"`
	MKey            string  `json:"mKey" bson:"mKey"`
	PaymentCategory string  `json:"paymentCategory" bson:"paymentCategory"`
	DeviceId        string  `json:"deviceId" bson:"deviceId"`
	ChannelCallback string  `json:"channelCallback" bson:"channelCallback"`
}

type InvoiceFeeDetailParking struct {
	Amount       float64 `json:"amount"`
	OvernightFee float64 `json:"overnightFee"`
}

type ResponseConfirm struct {
	DocNo            string      `json:"docNo"`
	ProductData      string      `json:"productData"`
	ProductName      string      `json:"productName"`
	CardNumber       string      `json:"cardNumber"`
	CardType         string      `json:"cardType"`
	CheckInDatetime  string      `json:"checkInDatetime"`
	CheckOutDatetime string      `json:"checkOutDatetime"`
	VehicleNumberIn  string      `json:"vehicleNumberIn"`
	VehicleNumberOut string      `json:"vehicleNumberOut"`
	UUIDCard         string      `json:"uuidCard"`
	ShowQRISArea     string      `json:"showQRISArea"`
	CurrentBalance   int64       `json:"currentBalance"`
	GrandTotal       float64     `json:"grandTotal"`
	InvoiceDetail    interface{} `json:"invoiceDetail"`
	OuCode           string      `json:"ouCode"`
	OuName           string      `json:"ouName"`
	Address          string      `json:"address"`
	IPAddr           string      `json:"ipAddr"`
}

type ResponseConfirmTrxDeposit struct {
	DocNo            string  `json:"docNo"`
	ProductData      string  `json:"productData"`
	ProductName      string  `json:"productName"`
	CardNumber       string  `json:"cardNumber"`
	CardType         string  `json:"cardType"`
	CheckInDatetime  string  `json:"checkInDatetime"`
	CheckOutDatetime string  `json:"checkOutDatetime"`
	UUIDCard         string  `json:"uuidCard"`
	ShowQRISArea     string  `json:"showQRISArea"`
	CurrentBalance   int64   `json:"currentBalance"`
	GrandTotal       float64 `json:"grandTotal"`
	OuCode           string  `json:"ouCode"`
	OuName           string  `json:"ouName"`
	Address          string  `json:"address"`
	IPAddr           string  `json:"ipAddr"`
}

type ConfirmTrxByPass struct {
	CardType    string  `json:"cardType"`
	DeviceId    string  `json:"deviceId" validate:"required"`
	CardNumber  string  `json:"cardNumber" validate:"required"`
	ProductCode string  `json:"productCode" validate:"required"`
	LogTrans    string  `json:"logTrans" validate:"required"`
	IpTerminal  string  `json:"ip_terminal"`
	GrandTotal  float64 `json:"grandTotal" validate:"min=0"`
}

type RequestInquiryPayment struct {
	DocNo         string  `json:"docNo" validate:"required"`
	ProductCode   string  `json:"productCode" validate:"required"`
	ProductName   string  `json:"productName" validate:"required"`
	PaymentMethod string  `json:"paymentMethod" validate:"required"`
	GrandTotal    float64 `json:"grandTotal" validate:"required"`
}

type RequestInquiryPaymentP3 struct {
	ProductCode     string  `json:"productCode" validate:"required"`
	PaymentMethod   string  `json:"paymentMethod"`
	InquiryDatetime string  `json:"inquiryDatetime" bson:"inquiryDatetime"`
	GrandTotal      float64 `json:"grandTotal"`
	TerminalId      string  `json:"terminalId" bson:"terminalId"`
}

type TrxOutstandingForClearSession struct {
	RefDocNo               string `json:"refDocNo" bson:"refDocNo"`
	TappingDate            string `json:"tappingDate" bson:"tappingDate"`
	TappingDatetime        string `json:"tappingDatetime" bson:"tappingDatetime"`
	CardNumberUuid         string `json:"cardNumberUuid" bson:"cardNumberUuid"`
	FlagClearSession       bool   `json:"flagClearSession" bson:"flagClearSession"`
	ClearDatetime          string `json:"clearDatetime" bson:"clearDatetime"`
	TrxOutstandingSnapshot Trx    `json:"trxOutstandingSnapshot" bson:"trxOutstandingSnapshot"`
}

type RequestSyncTrxToCLoud struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
	OuCode   string `json:"ouCode"`
	Limit    int64  `json:"limit"`
}

type RequestUpdateProductPrice struct {
	ID          string `json:"_id" bson:"_id" validate:"required"`
	Username    string `json:"username" bson:"username" validate:"required"`
	Pin         int64  `json:"pin" bson:"pin"`
	ProductCode string `json:"productCode" bson:"productCode"`
	ProductId   int64  `json:"productid" bson:"productId"`
}

type UpdateProductCloud struct {
	OuId        int64  `json:"ouId" bson:"ouId"`
	DocNo       string `json:"docNo" bson:"docNo"`
	ProductId   int64  `json:"productId" bson:"productId"`
	ProductCode string `json:"productCode" bson:"productCode"`
	ProductName string `json:"productName" bson:"productName"`
}

type RequestRegistrationMemberLocal struct {
	Keyword            string `json:"keyword" bson:"keyword" validate:"required"`
	FirstName          string `json:"firstName" bson:"firstName" validate:"required"`
	LastName           string `json:"lastName" bson:"lastName"`
	PhoneNumber        string `json:"phoneNumber" bson:"phoneNumber"`
	VehicleNumber      string `json:"vehicleNumber"`
	StartDate          string `json:"startDate" bson:"startDate" validate:"required"`
	EndDate            string `json:"endDate" bson:"endDate" validate:"required"`
	IdentifierCustomer string `json:"identifierCustomer" bson:"identifierCustomer"`
}

type TrxProductCustom struct {
	Keyword     string `json:"keyword" bson:"keyword"`
	ProductName string `json:"productName" bson:"productName"`
	ProductCode string `json:"productCode" bson:"productCode"`
}

type TrxInvoiceDepositCounterItem struct {
	DocNo            string  `json:"docNoDepo" bson:"docNoDepo"`
	ProductId        int64   `json:"productId" bson:"productId"`
	ProductCode      string  `json:"productCode" bson:"productCode"`
	ProductName      string  `json:"productName" bson:"productName"`
	IsPctServiceFee  string  `json:"isPctServiceFee" bson:"isPctServiceFee"`
	ServiceFee       float64 `json:"serviceFee" bson:"serviceFee"`
	ServiceFeeMember float64 `json:"serviceFeeMember" bson:"serviceFeeMember"`
	Price            float64 `json:"price" bson:"price"`
	TotalAmount      float64 `json:"totalAmount" bson:"totalAmount"`
}

type RequestTrxDepositCounter struct {
	CheckInDatetime string `json:"checkInDatetime" validate:"required"`
	ProductCode     string `json:"productCode"`
	DepositorName   string `json:"depositorName"`
	Merk            string `json:"merk"`
	Username        string `json:"username" bson:"username"`
	ShiftCode       string `json:"shiftCode" bson:"shiftCode"`
}

type TrxDepositCounter struct {
	DocNoDepo        string                         `json:"docNoDepo" bson:"docNoDepo"`
	DocDateDepo      string                         `json:"docDateDepo" bson:"docDate"`
	ProductData      string                         `json:"productData" bson:"productData"`
	ProductCode      string                         `json:"productCode" bson:"productCode"`
	ProductName      string                         `json:"productName" bson:"productName"`
	Merk             string                         `json:"merk" bson:"merk"`
	DeviceId         string                         `json:"device_id" bson:"deviceId"`
	DepositorName    string                         `json:"depositorName" bson:"depositorName"`
	CheckInDatetime  string                         `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime string                         `json:"checkOutDatetime" bson:"checkOutDatetime"`
	OuId             int64                          `json:"ouId" bson:"ouId"`
	OuName           string                         `json:"ouName" bson:"ouName"`
	OuCode           string                         `json:"ouCode" bson:"ouCode"`
	OuSubBranchId    int64                          `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName  string                         `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode  string                         `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId         int64                          `json:"mainOuId" bson:"mainOuId"`
	MainOuCode       string                         `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName       string                         `json:"mainOuName" bson:"mainOuName"`
	MemberCode       string                         `json:"memberCode" bson:"memberCode"`
	MemberName       string                         `json:"memberName" bson:"memberName"`
	MemberType       string                         `json:"memberType" bson:"memberType"`
	RequestData      string                         `json:"requestData" bson:"requestData"`
	RequestOutData   string                         `json:"requestOutData" bson:"requestOutData"`
	CardNumberUUID   string                         `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber       string                         `json:"cardNumber" bson:"cardNumber"`
	TypeCard         string                         `json:"typeCard" bson:"typeCard"`
	BeginningBalance float64                        `json:"beginningBalance" bson:"beginningBalance"`
	ExtLocalDatetime string                         `json:"extLocalDatetime" bson:"extLocalDatetime"`
	GrandTotal       float64                        `json:"grandTotal" bson:"grandTotal"`
	LogTrans         string                         `json:"logTrans" bson:"logTrans"`
	MerchantKey      string                         `json:"merchantKey" bson:"merchantKey"`
	QrText           string                         `json:"qrText" bson:"qrText"`
	TrxInvoiceItem   []TrxInvoiceDepositCounterItem `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
	FlagSyncData     bool                           `json:"flagSyncData" bson:"flagSyncData"`
	Username         string                         `json:"username" bson:"username"`
	ShiftCode        string                         `json:"shiftCode" bson:"shiftCode"`
	UsernameOut      string                         `json:"usernameOut" bson:"usernameOut"`
	ShiftCodeOut     string                         `json:"shiftCodeOut" bson:"shiftCodeOut"`
}

type ResponseTrxDepositCounter struct {
	ID              *primitive.ObjectID `json:"_id"`
	CheckInDatetime string              `json:"checkInDatetime"`
	DocNoDepo       string              `json:"docNoDepo"`
	ProductName     string              `json:"productName"`
	DepositorName   string              `json:"depositorName"`
	Merk            string              `json:"merk"`
	QRCode          string              `json:"qrCode"`
	OuCode          string              `json:"ouCode"`
	OuName          string              `json:"ouName"`
	Address         string              `json:"address"`
}

type ResultFindTrxDepositCounterOutstanding struct {
	ID              primitive.ObjectID             `json:"_id" bson:"_id"`
	DocNoDepo       string                         `json:"docNoDepo" bson:"docNoDepo"`
	GrandTotal      float64                        `json:"grandTotal" bson:"grandTotal"`
	DepositorName   string                         `json:"depositorName" bson:"depositorName"`
	Merk            string                         `json:"merk" bson:"merk"`
	ProductName     string                         `json:"productName" bson:"productName"`
	CheckInDatetime string                         `json:"checkInDatetime" bson:"checkInDatetime"`
	OuCode          string                         `json:"ouCode" bson:"ouCode"`
	TrxInvoiceItem  []TrxInvoiceDepositCounterItem `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
}

type TrxItemsDepositCounter struct {
	Draw            int64  `json:"draw" bson:"draw"`
	Keyword         string `json:"keyword" bson:"keyword"`
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type ResponseFindTrxDepositCounter struct {
	TotalData                 int64                   `json:"totalData"`
	ListTrxDepositCounterData []FindTrxDepositCounter `json:"listTrxDepositCounterData"`
}

type FindTrxDepositCounter struct {
	ID              primitive.ObjectID             `json:"_id" bson:"_id"`
	DocNoDepo       string                         `json:"docNoDepo" bson:"docNoDepo"`
	ProductData     string                         `json:"productData" bson:"productData"`
	ProductCode     string                         `json:"productCode" bson:"productCode"`
	ProductName     string                         `json:"productName" bson:"productName"`
	Merk            string                         `json:"merk" bson:"merk"`
	DepositorName   string                         `json:"depositorName" bson:"depositorName"`
	CheckInDatetime string                         `json:"checkInDatetime" bson:"checkInDatetime"`
	OuId            int64                          `json:"ouId" bson:"ouId"`
	OuName          string                         `json:"ouName" bson:"ouName"`
	OuCode          string                         `json:"ouCode" bson:"ouCode"`
	OuSubBranchId   int64                          `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName string                         `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode string                         `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId        int64                          `json:"mainOuId" bson:"mainOuId"`
	MainOuCode      string                         `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName      string                         `json:"mainOuName" bson:"mainOuName"`
	QrText          string                         `json:"qrText" bson:"qrText"`
	TrxInvoiceItem  []TrxInvoiceDepositCounterItem `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
	Username        string                         `json:"username" bson:"username"`
	ShiftCode       string                         `json:"shiftCode" bson:"shiftCode"`
}

type ResultInquiryTrxCustomDepositCounter struct {
	ID            string  `json:"_id" bson:"_id"`
	DocNo         string  `json:"docNo" bson:"docNo"`
	Nominal       float64 `json:"nominal" bson:"nominal"`
	ProductCode   string  `json:"productCode" bson:"productCode"`
	ProductName   string  `json:"productName" bson:"productName"`
	DepositorName string  `json:"depositorName" bson:"depositorName"`
	Merk          string  `json:"merk" bson:"merk"`
	QRCode        string  `json:"qrCode" bson:"qrCode"`
	Type          string  `json:"type" bson:"type"`
	ExcludeSf     bool    `json:"excludeSf" bson:"excludeSf"`
	OuCode        string  `json:"ouCode"`
}

type Decrypt struct {
	Keyword string `json:"keyword"`
}

type TrxCustom struct {
	DocNo                        string                         `json:"docNo" bson:"docNo"`
	DocDate                      string                         `json:"docDate" bson:"docDate"`
	ExtDocNo                     string                         `json:"extDocNo" bson:"extDocNo"`
	CheckInDatetime              string                         `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime             string                         `json:"checkOutDatetime" bson:"checkOutDatetime"`
	DeviceIdIn                   string                         `json:"deviceIdIn" bson:"deviceIdIn"`
	DeviceId                     string                         `json:"device_id" bson:"deviceId"`
	GateIn                       string                         `json:"gateIn" bson:"gateIn"`
	GateOut                      string                         `json:"gateOut" bson:"gateOut"`
	CardNumberUUIDIn             string                         `json:"cardNumberUuidIn" bson:"cardNumberUuidIn"`
	CardNumberIn                 string                         `json:"cardNumberIn" bson:"cardNumberIn"`
	CardNumberUUID               string                         `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber                   string                         `json:"cardNumber" bson:"cardNumber"`
	TypeCard                     string                         `json:"typeCard" bson:"typeCard"`
	BeginningBalance             float64                        `json:"beginningBalance" bson:"beginningBalance"`
	ExtLocalDatetime             string                         `json:"extLocalDatetime" bson:"extLocalDatetime"`
	GrandTotal                   float64                        `json:"grandTotal" bson:"grandTotal"`
	ChargeAmount                 float64                        `json:"chargeAmount" bson:"chargeAmount"`
	ProductId                    int64                          `json:"productId" bson:"productId"`
	ProductCode                  string                         `json:"productCode" bson:"productCode"`
	ProductName                  string                         `json:"productName" bson:"productName"`
	ProductMembershipId          int64                          `json:"productMembershipId" bson:"productMembershipId"`
	ProductMembershipCode        string                         `json:"productMembershipCode" bson:"productMembershipCode"`
	ProductMembershipName        string                         `json:"productMembershipName" bson:"productMembershipName"`
	IsPctServiceFee              string                         `json:"isPctServiceFee" bson:"isPctServiceFee"`
	ServiceFee                   float64                        `json:"serviceFee" bson:"serviceFee"`
	ServiceFeeMember             float64                        `json:"serviceFeeMember" bson:"serviceFeeMember"`
	DueDate                      int64                          `json:"dueDate" bson:"dueDate"`
	DiscType                     string                         `json:"discType" bson:"discType"`
	DiscPct                      int64                          `json:"discPct" bson:"discPct"`
	DiscAmount                   float64                        `json:"discAmount" bson:"discAmount"`
	GracePeriodDate              int64                          `json:"gracePeriodDate" bson:"gracePeriodDate"`
	ProductData                  string                         `json:"productData" bson:"productData"`
	RequestData                  string                         `json:"requestData" bson:"requestData"`
	RequestOutData               string                         `json:"requestOutData" bson:"requestOutData"`
	OuId                         int64                          `json:"ouId" bson:"ouId"`
	OuName                       string                         `json:"ouName" bson:"ouName"`
	OuCode                       string                         `json:"ouCode" bson:"ouCode"`
	OuSubBranchId                int64                          `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName              string                         `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode              string                         `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId                     int64                          `json:"mainOuId" bson:"mainOuId"`
	MainOuCode                   string                         `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName                   string                         `json:"mainOuName" bson:"mainOuName"`
	Price                        float64                        `json:"price" bson:"price"`
	LogTrans                     string                         `json:"logTrans" bson:"logTrans"`
	QrText                       string                         `json:"qrText" bson:"qrText"`
	MerchantKey                  string                         `json:"merchantKey" bson:"merchantKey"`
	TrxInvoiceItemDepositCounter []TrxInvoiceDepositCounterItem `json:"trxInvoiceItemDepositCounter" bson:"trxInvoiceItemDepositCounter"`
	TrxInvoiceItemMemberDeposit  []TrxInvoiceMemberDeposit      `json:"trxInvoiceItemMemberDeposit" bson:"trxInvoiceItemMemberDeposit"`
	FlagSyncData                 bool                           `json:"flagSyncData" bson:"flagSyncData"`
	MemberData                   *TrxMember                     `json:"memberData" bson:"memberData"`
	TrxAddInfo                   map[string]interface{}         `json:"trxAddInfo" bson:"trxAddInfo"`
	FlagTrxFromCloud             bool                           `json:"flagTrxFromCloud" bson:"flagTrxFromCloud"`
	IsRsyncDataTrx               bool                           `json:"isRsyncDataTrx" bson:"isRsyncDataTrx"`
	ExcludeSf                    bool                           `json:"excludeSf" bson:"excludeSf"`
	FlagCharge                   bool                           `json:"flagCharge" bson:"flagCharge"`
	FlgDepositCounter            bool                           `json:"flgDepositCounter" bson:"flgDepositCounter"`
}

type TrxInvoiceMemberDeposit struct {
	DocNo                 string  `json:"docNoDepo" bson:"docNoDepo"`
	PartnerCode           string  `json:"partnerCode" bson:"partnerCode"`
	ProductId             int64   `json:"productId" bson:"productId"`
	ProductCode           string  `json:"productCode" bson:"productCode"`
	ProductName           string  `json:"productName" bson:"productName"`
	ProductMembershipId   int64   `json:"productMembershipId" bson:"productMembershipId"`
	ProductMembershipCode string  `json:"productMembershipCode" bson:"productMembershipCode"`
	ProductMembershipName string  `json:"productMembershipName" bson:"productMembershipName"`
	IsPctServiceFee       string  `json:"isPctServiceFee" bson:"isPctServiceFee"`
	ServiceFee            float64 `json:"serviceFee" bson:"serviceFee"`
	ServiceFeeMember      float64 `json:"serviceFeeMember" bson:"serviceFeeMember"`
	DueDate               int64   `json:"dueDate" bson:"dueDate"`
	DiscType              string  `json:"discType" bson:"discType"`
	DiscPct               int64   `json:"discPct" bson:"discPct"`
	DiscAmount            float64 `json:"discAmount" bson:"discAmount"`
	GracePeriodDate       int64   `json:"gracePeriodDate" bson:"gracePeriodDate"`
	Price                 float64 `json:"price" bson:"price"`
	TotalAmount           float64 `json:"totalAmount" bson:"totalAmount"`
}

type RequestInquiryPaymentQris struct {
	DocNo         string  `json:"docNo" validate:"required"`
	ProductCode   string  `json:"productCode" validate:"required"`
	ProductName   string  `json:"productName" validate:"required"`
	PaymentMethod string  `json:"paymentMethod" validate:"required"`
	GrandTotal    float64 `json:"grandTotal" validate:"required"`
	MKey          string  `json:"mKey" validate:"required"`
}

type ResponseInquiryQris struct {
	Type            string `json:"type"`
	QrCode          string `json:"qrCode"`
	PaymentRefDocNo string `json:"paymentRefDocNo"`
}

type ResponseConfirmLostTicket struct {
	DocNo                 string  `json:"docNo"`
	PaymentMethod         string  `json:"paymentMethod"`
	ProductName           string  `json:"productName"`
	CardNumber            string  `json:"cardNumber"`
	CardType              string  `json:"cardType"`
	LostTicketInDatetime  string  `json:"lostTicketInDatetime"`
	LostTicketOutDatetime string  `json:"lostTicketOutDatetime"`
	VehicleNumberIn       string  `json:"vehicleNumberIn"`
	VehicleNumberOut      string  `json:"vehicleNumberOut"`
	UUIDCard              string  `json:"uuidCard"`
	QrCodeLostTicket      string  `json:"qrCodeLostTicket"`
	CurrentBalance        int64   `json:"currentBalance"`
	ChargeAmount          float64 `json:"chargeAmount"`
	GrandTotal            float64 `json:"grandTotal"`
	OuCode                string  `json:"ouCode"`
	OuName                string  `json:"ouName"`
	Address               string  `json:"address"`
	IPAddr                string  `json:"ipAddr"`
}

type RequestCheckStatusPaymentQris struct {
	DocNo       string `json:"docNo" validate:"required"`
	MerchantKey string `json:"merchantKey" validate:"required"`
}

type ResponseCheckStatusA2P struct {
	AcquiringID     float64 `json:"acquiringID"`
	Amount          float64 `json:"amount"`
	BankNoRef       string  `json:"bankNoRef"`
	CardPan         string  `json:"cardPan"`
	CardType        string  `json:"cardType"`
	CorporateName   string  `json:"corporateName"`
	CreatedAt       string  `json:"createdAt"`
	CurrentBalance  float64 `json:"currentBalance"`
	DeviceID        string  `json:"deviceID"`
	Discount        float64 `json:"discount"`
	LastBalance     float64 `json:"lastBalance"`
	Mdr             float64 `json:"mdr"`
	Mid             string  `json:"mid"`
	NoHeader        string  `json:"noHeader"`
	PaymentCategory string  `json:"paymentCategory"`
	PaymentFee      float64 `json:"paymentFee"`
	PromoCode       string  `json:"promoCode"`
	PromoIssuer     string  `json:"promoIssuer"`
	ServiceFee      float64 `json:"serviceFee"`
	SettleAt        string  `json:"settleAt"`
	StatusCode      string  `json:"statusCode"`
	StatusPayment   string  `json:"statusPayment"`
	Tid             string  `json:"tid"`
}

type ResponseConfirmTrxVip struct {
	DocNo            string  `json:"docNo"`
	ProductData      string  `json:"productData"`
	ProductName      string  `json:"productName"`
	CardNumber       string  `json:"cardNumber"`
	CardType         string  `json:"cardType"`
	CheckInDatetime  string  `json:"checkInDatetime"`
	CheckOutDatetime string  `json:"checkOutDatetime"`
	VehicleNumberIn  string  `json:"vehicleNumberIn"`
	VehicleNumberOut string  `json:"vehicleNumberOut"`
	UUIDCard         string  `json:"uuidCard"`
	ShowQRISArea     string  `json:"showQRISArea"`
	CurrentBalance   int64   `json:"currentBalance"`
	ChargeAmount     float64 `json:"chargeAmount" bson:"chargeAmount"`
	GrandTotal       float64 `json:"grandTotal"`
	OuCode           string  `json:"ouCode"`
	OuName           string  `json:"ouName"`
	Address          string  `json:"address"`
	IPAddr           string  `json:"ipAddr"`
}

type TrxWithId struct {
	ID               string                 `json:"_id" bson:"_id"`
	DocNo            string                 `json:"docNo" bson:"docNo"`
	DocDate          string                 `json:"docDate" bson:"docDate"`
	CheckInDatetime  string                 `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime string                 `json:"checkOutDatetime" bson:"checkOutDatetime"`
	DeviceIdIn       string                 `json:"deviceIdIn" bson:"deviceIdIn"`
	DeviceId         string                 `json:"device_id" bson:"deviceId"`
	GateIn           string                 `json:"gateIn" bson:"gateIn"`
	GateOut          string                 `json:"gateOut" bson:"gateOut"`
	CardNumberUUIDIn string                 `json:"cardNumberUuidIn" bson:"cardNumberUuidIn"`
	CardNumberIn     string                 `json:"cardNumberIn" bson:"cardNumberIn"`
	CardNumberUUID   string                 `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber       string                 `json:"cardNumber" bson:"cardNumber"`
	TypeCard         string                 `json:"typeCard" bson:"typeCard"`
	BeginningBalance float64                `json:"beginningBalance" bson:"beginningBalance"`
	ExtLocalDatetime string                 `json:"extLocalDatetime" bson:"extLocalDatetime"`
	ChargeAmount     float64                `json:"chargeAmount" bson:"chargeAmount"`
	GrandTotal       float64                `json:"grandTotal" bson:"grandTotal"`
	ProductCode      string                 `json:"productCode" bson:"productCode"`
	ProductName      string                 `json:"productName" bson:"productName"`
	ProductData      string                 `json:"productData" bson:"productData"`
	RequestData      string                 `json:"requestData" bson:"requestData"`
	RequestOutData   string                 `json:"requestOutData" bson:"requestOutData"`
	OuId             int64                  `json:"ouId" bson:"ouId"`
	OuName           string                 `json:"ouName" bson:"ouName"`
	OuCode           string                 `json:"ouCode" bson:"ouCode"`
	OuSubBranchId    int64                  `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName  string                 `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode  string                 `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId         int64                  `json:"mainOuId" bson:"mainOuId"`
	MainOuCode       string                 `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName       string                 `json:"mainOuName" bson:"mainOuName"`
	MemberCode       string                 `json:"memberCode" bson:"memberCode"`
	MemberName       string                 `json:"memberName" bson:"memberName"`
	MemberType       string                 `json:"memberType" bson:"memberType"`
	CheckInTime      int64                  `json:"checkInTime" bson:"checkInTime"`
	CheckOutTime     int64                  `json:"checkOutTime" bson:"checkOutTime"`
	DurationTime     int64                  `json:"durationTime" bson:"durationTime"`
	VehicleNumberIn  string                 `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	VehicleNumberOut string                 `json:"vehicleNumberOut" bson:"vehicleNumberOut"`
	LogTrans         string                 `json:"logTrans" bson:"logTrans"`
	MerchantKey      string                 `json:"merchantKey" bson:"merchantKey"`
	QrText           string                 `json:"qrText" bson:"qrText"`
	TrxInvoiceItem   []TrxInvoiceItem       `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
	FlagSyncData     bool                   `json:"flagSyncData" bson:"flagSyncData"`
	MemberData       *TrxMember             `json:"memberData" bson:"memberData"`
	TrxAddInfo       map[string]interface{} `json:"trxAddInfo" bson:"trxAddInfo"`
	FlagTrxFromCloud bool                   `json:"flagTrxFromCloud" bson:"flagTrxFromCloud"`
	IsRsyncDataTrx   bool                   `json:"isRsyncDataTrx" bson:"isRsyncDataTrx"`
	ExcludeSf        bool                   `json:"excludeSf" bson:"excludeSf"`
	FlagCharge       bool                   `json:"flagCharge" bson:"flagCharge"`
	ChargeType       string                 `json:"chargeType" bson:"chargeType"`
	QrTextLostTicket *string                `json:"qrTextLostTicket" bson:"qrTextLostTicket"`
	StatusLostTicket *bool                  `json:"statusLostTicket" bson:"statusLostTicket"`
}

type ResponseIDTrxOutstanding struct {
	ID *primitive.ObjectID `json:"_id"`
}

type ExtendMember struct {
	OuId          int64  `json:"ouId"`
	RegisteredBy  string `json:"registeredBy"`
	TypePartner   string `json:"typePartner"`
	DateFrom      string `json:"dateFrom" validate:"required"`
	DateTo        string `json:"dateTo" validate:"required"`
	VehicleNumber string `json:"vehicleNumber"`
	CardNumber    string `json:"cardNumber"`
	ProductId     int64  `json:"productId"`
	Username      string `json:"username"`
	UpdatedAt     string `json:"updatedAt"`
	UpdatedBy     string `json:"updatedBy"`
}
