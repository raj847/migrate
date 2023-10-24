package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestTrx struct {
	DateFrom string `json:"dateFrom" bson:"dateTo"`
	DateTo   string `json:"dateTo" bson:"dateTo"`
}

type RequestTrafficParking struct {
	Date string `json:"date" bson:"date"`
}

type ResponseTrx struct {
	TotalOmzet float64 `json:"totalOmzet" bson:"totalOmzet"`
	TotalTrx   int64   `json:"totalTrx" bson:"totalTrx"`
}

type ResponseTrxOutstanding struct {
	TotalTrxOutstanding int64 `json:"totalTrxOutstanding" bson:"totalTrxOutstanding"`
}

type ResponseSummaryMember struct {
	TotalMember   int64 `json:"totalMember"`
	TotalFreePass int64 `json:"totalFreePass"`
}

type SummaryMember struct {
	Member int64 `json:"member"`
}

type ResponseSummaryTrx struct {
	TotalOmzet          float64 `json:"totalOmzet" bson:"totalOmzet"`
	TotalTrx            int64   `json:"totalTrx" bson:"totalTrx"`
	TotalTrxOutstanding int64   `json:"totalTrxOutstanding" bson:"totalTrxOutstanding"`
}

type TrxInvoiceDetail struct {
	CreatedDate string `json:"createdDate" bson:"createdDate"`
}

type ResponseOvernight struct {
	OneNight       int64 `json:"oneNight"`
	TwoNight       int64 `json:"twoNight"`
	ThreeNight     int64 `json:"threeNight"`
	FourNight      int64 `json:"fourNight"`
	FiveNight      int64 `json:"fiveNight"`
	SixNight       int64 `json:"sixNight"`
	MoreSevenNight int64 `json:"moreSevenNight"`
}

type ChartDailyIntelligenceTrxList struct {
	Id         Date    `json:"id" bson:"_id"`
	TotalOmzet float64 `json:"totalOmzet" bson:"totalOmzet"`
	TotalTrx   int64   `json:"totalTrx" bson:"totalTrx"`
}

type Date struct {
	Days  int64 `json:"days" bson:"Days"`
	Month int64 `json:"month" bson:"Month"`
	Year  int64 `json:"year" bson:"Year"`
}

type ResponseChartDailyIntelligenceTrxList struct {
	Date       string  `json:"date"`
	TotalOmzet float64 `json:"totalOmzet"`
	TotalTrx   int64   `json:"totalTrx"`
}

type TrafficParking struct {
	Hour     string `json:"hour" bson:"hour"`
	TotalTrx int64  `json:"totalTrx" bson:"totalTrx"`
}

type TotalTrx struct {
	TotalTrx int64 `json:"totalTrx" bson:"totalTrx"`
}

type TrxItems struct {
	Keyword         string `json:"keyword" bson:"keyword"`
	TypeCard        string `json:"typeCard" bson:"typeCard"`
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type RequestExportTrxItems struct {
	Keyword         string `json:"keyword" bson:"keyword"`
	TypeCard        string `json:"typeCard" bson:"typeCard"`
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type RequestTrxItems struct {
	Keyword         string `json:"keyword" bson:"keyword"`
	TypeCard        string `json:"typeCard" bson:"typeCard"`
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	Draw            int64  `json:"draw" bson:"draw"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type GetTrxList struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	DocNo          string             `json:"docNo" bson:"docNo"`
	DocDate        string             `json:"docDate" bson:"docDate"`
	MemberCode     string             `json:"memberCode" bson:"memberCode"`
	MemberName     string             `json:"memberName" bson:"memberName"`
	CardNumberUUID string             `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber     string             `json:"cardNumber" bson:"cardNumber"`
	ProductData    string             `json:"productData" bson:"productData"`
	GrandTotal     float64            `json:"grandTotal" bson:"grandTotal"`
	TypeCard       string             `json:"typeCard" bson:"typeCard"`
}

type TrxList struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	DocNo          string             `json:"docNo" bson:"docNo"`
	DocDate        string             `json:"docDate" bson:"docDate"`
	MemberCode     string             `json:"memberCode" bson:"memberCode"`
	MemberName     string             `json:"memberName" bson:"memberName"`
	CardNumberUUID string             `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber     string             `json:"cardNumber" bson:"cardNumber"`
	ProductCode    string             `json:"productCode" bson:"productCode"`
	ProductName    string             `json:"productName" bson:"productName"`
	GrandTotal     float64            `json:"grandTotal" bson:"grandTotal"`
	TypeCard       string             `json:"typeCard" bson:"typeCard"`
}

type ResponseGetTrxList struct {
	CountTrxList int64                     `json:"countTrxList" bson:"countTrxList"`
	Data         []ResultFindTrxForPostBox `json:"data" bson:"data"`
}

type GetTrxOrderByVehicle struct {
	Id         Product `json:"id" bson:"_id"`
	TotalOmzet float64 `json:"totalOmzet" bson:"totalOmzet"`
	TotalTrx   int64   `json:"totalTrx" bson:"totalTrx"`
}

type Product struct {
	ProductCode string `json:"productCode" bson:"productCode"`
	ProductName string `json:"productName" bson:"productName"`
}

type ResponseGetTrxOrderByVehicle struct {
	ProductCode string  `json:"productCode" bson:"productCode"`
	ProductName string  `json:"productName" bson:"productName"`
	TotalOmzet  float64 `json:"totalOmzet" bson:"totalOmzet"`
	TotalTrx    int64   `json:"totalTrx" bson:"totalTrx"`
}

type RequestGetTrxOrderByVehicle struct {
	Date        string `json:"date" bson:"date"`
	ProductCode string `json:"productCode" bson:"productCode"`
}

type ResultFindTrxForPostBox struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	DocNo            string             `json:"docNo" bson:"docNo"`
	GrandTotal       float64            `json:"grandTotal" bson:"grandTotal"`
	CheckInDatetime  string             `json:"checkInDatetime" bson:"checkInDatetime"`
	CardNumber       string             `json:"cardNumber" bson:"cardNumber"`
	CardNumberUUID   string             `json:"cardNumberUuid" bson:"cardNumberUuid"`
	TypeCard         string             `json:"typeCard" bson:"typeCard"`
	OuCode           string             `json:"ouCode" bson:"ouCode"`
	VehicleNumberIn  string             `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	GateIn           string             `json:"gateIn" bson:"gateIn"`
	QrText           string             `json:"qrText" bson:"qrText"`
	QrTextLostTicket *string            `json:"qrTextLostTicket" bson:"qrTextLostTicket"`
	StatusLostTicket *bool              `json:"statusLostTicket" bson:"statusLostTicket"`
	TrxInvoiceItem   []TrxInvoiceItem   `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
}

type RequestFindTrxOutstanding struct {
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Keyword         string `json:"keyword" bson:"keyword"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	Draw            int64  `json:"draw" bson:"draw"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type ResponseTrxOutstandingForPostBox struct {
	TotalTrxOutstanding *ResponseTrxOutstanding   `json:"totalTrxOutstanding" bson:"totalTrxOutstanding"`
	DataOutstanding     []ResultFindTrxForPostBox `json:"dataOutstanding" bson:"dataOutstanding"`
}

type RequestAddTrxLostTicket struct {
	IdTrx                     string  `json:"idTrx" bson:"idTrx"`
	OfficerName               string  `json:"officerName" bson:"officerName"`
	VehicleNumber             string  `json:"vehicleNumber" bson:"vehicleNumber" validate:"required"`
	VehicleType               string  `json:"vehicleType" bson:"vehicleType"`
	VehicleColor              string  `json:"vehicleColor" bson:"vehicleColor"`
	ProductCode               string  `json:"productCode" bson:"productCode"`
	ProductName               string  `json:"productName" bson:"productName"`
	CheckInDatetime           string  `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime          string  `json:"checkOutDatetime" bson:"checkOutDatetime"`
	CustName                  string  `json:"custName" bson:"custName"`
	CustAddress               string  `json:"custAddress" bson:"custAddress"`
	CustIdentifier            string  `json:"custIdentifier" bson:"custIdentifier" validate:"required"`
	VehicleRegistrationName   string  `json:"vehicleRegistrationName" bson:"vehicleRegistrationName"`
	VehicleRegistrationNumber string  `json:"vehicleRegistrationNumber" bson:"vehicleRegistrationNumber"`
	Remark                    string  `json:"remark" bson:"remark"`
	GrandTotal                float64 `json:"grandTotal" bson:"grandTotal"`
	ChargeAmount              float64 `json:"chargeAmount" bson:"chargeAmount"`
	GenerateTrxOustanding     bool    `json:"generateTrxOutstanding" bson:"generateTrxOutstanding"`
	IsTrxChargeMaxAmount      bool    `json:"isTrxChargeMaxAmount" bson:"isTrxChargeMaxAmount"`
	CreatedAt                 string  `json:"createdAt" bson:"createdAt"`
}

type RequestInquiryLostTicket struct {
	ID                 string `json:"id" bson:"_id" validate:"required"`
	DocNo              string `json:"docNo" bson:"docNo"`
	InquiryDatetime    string `json:"inquiryDatetime" bson:"inquiryDatetime"`
	ProductCode        string `json:"productCode" bson:"productCode"`
	CustomerName       string `json:"customerName" bson:"customerName"`
	CustomerIdentifier string `json:"customerIdentifier" bson:"customerIdentifier"`
	Username           string `json:"username" bson:"username"`
	ShiftCode          string `json:"shiftCode" bson:"shiftCode"`
	TerminalId         string `json:"terminalId" bson:"terminalId"`
}

type RequestConfirmLostTicket struct {
	ID                    string  `json:"id" bson:"_id" validate:"required"`
	DeviceId              string  `json:"deviceId" bson:"deviceId"`
	IpTerminal            string  `json:"ipTerminal" bson:"ipTerminal"`
	CardNumber            string  `json:"cardNumber" bson:"cardNumber"`
	CardType              string  `json:"cardType" bson:"cardType"`
	LastBalance           int64   `json:"lastBalance" bson:"lastBalance"`
	CurrentBalance        int64   `json:"currentBalance" bson:"currentBalance"`
	UUIDCard              string  `json:"uuidCard" bson:"uuidCard"`
	ProductCode           string  `json:"productCode" validate:"required"`
	ProductName           string  `json:"productName" validate:"required"`
	LostTicketOutDatetime string  `json:"lostTicketOutDatetime" bson:"lostTicketOutDatetime" validate:"required"`
	ChargeAmount          float64 `json:"chargeAmount" bson:"chargeAmount"`
	GrandTotal            float64 `json:"grandTotal" bson:"grandTotal" validate:"min=0"`
	PaymentMethod         string  `json:"paymentMethod" bson:"paymentMethod"`
	LogTrans              string  `json:"logTrans" bson:"logTrans"`
	VehicleNumber         string  `json:"vehicleNumber" bson:"vehicleNumber"`
	Username              string  `json:"username" bson:"username"`
	ShiftCode             string  `json:"shiftCode" bson:"shiftCode"`
	ExcludeSf             *bool   `json:"excludeSf" bson:"excludeSf"`
}

type FindLostTicket struct {
	ID                        primitive.ObjectID `json:"_id" bson:"_id"`
	DocNo                     string             `json:"docNo" bson:"docNo"`
	PaymentRefDocNo           string             `json:"paymentRefDocNo" bson:"paymentRefDocNo"`
	ProductCode               string             `json:"productCode" bson:"productCode"`
	ProductName               string             `json:"productName" bson:"productName"`
	OfficerName               string             `json:"officerName" bson:"officerName"`
	VehicleNumber             string             `json:"vehicleNumber" bson:"vehicleNumber" validate:"required"`
	VehicleType               string             `json:"vehicleType" bson:"vehicleType"`
	VehicleColor              string             `json:"vehicleColor" bson:"vehicleColor"`
	CheckInDatetime           string             `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime          string             `json:"checkOutDatetime" bson:"checkOutDatetime"`
	CustName                  string             `json:"custName" bson:"custName"`
	CustAddress               string             `json:"custAddress" bson:"custAddress"`
	CustIdentifier            string             `json:"custIdentifier" bson:"custIdentifier" validate:"required"`
	VehicleRegistrationName   string             `json:"vehicleRegistrationName" bson:"vehicleRegistrationName"`
	VehicleRegistrationNumber string             `json:"vehicleRegistrationNumber" bson:"vehicleRegistrationNumber"`
	Remark                    string             `json:"remark" bson:"remark"`
	QRText                    string             `json:"qrText" bson:"qrText"`
	QrCodePayment             string             `json:"qrCodePayment" bson:"qrCodePayment"`
	QrCodeLostTicket          string             `json:"qrCodeLostTicket" bson:"qrCodeLostTicket"`
	GrandTotal                float64            `json:"grandTotal" bson:"grandTotal"`
	ChargeAmount              float64            `json:"chargeAmount" bson:"chargeAmount"`
	TrxOutstanding            *Trx               `json:"trxOutstanding" bson:"trxOutstanding"`
	LostTicketInDatetime      string             `json:"lostTicketInDatetime" bson:"lostTicketInDatetime"`
	LostTicketOutDatetime     string             `json:"lostTicketOutDatetime" bson:"lostTicketOutDatetime"`
	CreatedAt                 string             `json:"createdAt" bson:"createdAt"`
}

type LostTicket struct {
	DocNo                     string  `json:"docNo" bson:"docNo"`
	PaymentRefDocNo           string  `json:"paymentRefDocNo" bson:"paymentRefDocNo"`
	ProductCode               string  `json:"productCode" bson:"productCode"`
	ProductName               string  `json:"productName" bson:"productName"`
	OfficerName               string  `json:"officerName" bson:"officerName"`
	VehicleNumber             string  `json:"vehicleNumber" bson:"vehicleNumber" validate:"required"`
	VehicleType               string  `json:"vehicleType" bson:"vehicleType"`
	VehicleColor              string  `json:"vehicleColor" bson:"vehicleColor"`
	CheckInDatetime           string  `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime          string  `json:"checkOutDatetime" bson:"checkOutDatetime"`
	CustName                  string  `json:"custName" bson:"custName"`
	CustAddress               string  `json:"custAddress" bson:"custAddress"`
	CustIdentifier            string  `json:"custIdentifier" bson:"custIdentifier" validate:"required"`
	VehicleRegistrationName   string  `json:"vehicleRegistrationName" bson:"vehicleRegistrationName"`
	VehicleRegistrationNumber string  `json:"vehicleRegistrationNumber" bson:"vehicleRegistrationNumber"`
	Remark                    string  `json:"remark" bson:"remark"`
	InquiryDatetimeLostTicket string  `json:"inquiryDatetimeLostTicket" bson:"inquiryDatetimeLostTicket"`
	Username                  string  `json:"username" bson:"username"`
	ShiftCode                 string  `json:"shiftCode" bson:"shiftCode"`
	QRText                    string  `json:"qrText" bson:"qrText"`
	QrPayment                 string  `json:"qrPayment" bson:"qrPayment"`
	GrandTotal                float64 `json:"grandTotal" bson:"grandTotal"`
	ChargeAmount              float64 `json:"chargeAmount" bson:"chargeAmount"`
	TrxOutstanding            *Trx    `json:"trxOutstanding" bson:"trxOutstanding"`
	LostTicketInDatetime      string  `json:"lostTicketInDatetime" bson:"lostTicketInDatetime"`
	LostTicketOutDatetime     string  `json:"lostTicketOutDatetime" bson:"lostTicketOutDatetime"`
	Type                      string  `json:"type" bson:"type"`
	CreatedAt                 string  `json:"createdAt" bson:"createdAt"`
}

type RequestFindLostTicket struct {
	Draw            int64  `json:"draw"`
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Keyword         string `json:"keyword" bson:"keyword"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type ResponseFindLostTicket struct {
	TotalData          int64            `json:"totalData"`
	ListLostTicketData []FindLostTicket `json:"listLostTicketData"`
}

type RequestFindOneTimeMember struct {
	DateFrom        string `json:"dateFrom" bson:"dateFrom"`
	DateTo          string `json:"dateTo" bson:"dateTo"`
	Keyword         string `json:"keyword" bson:"keyword"`
	Limit           int64  `json:"limit" bson:"limit"`
	Offset          int64  `json:"offset" bson:"offset"`
	Draw            int64  `json:"draw" bson:"draw"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type TrxOneTimeMember struct {
	DocNo            string `json:"docNo" bson:"docNo"`
	DocDate          string `json:"docDate" bson:"docDate"`
	ProductName      string `json:"productName" bson:"productName"`
	MemberName       string `json:"memberName" bson:"memberName"`
	CardNumberUUID   string `json:"cardNumberUuidIn" bson:"cardNumberUuidIn"`
	CheckinDatetime  string `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckoutDatetime string `json:"checkOutDatetime" bson:"checkOutDatetime"`
	VehicleNumber    string `json:"vehicleNumber" bson:vehicleNumber""`
	PathInImage      string `json:"pathInImage" bson:"pathInImage"`
	PathOutImage     string `json:"pathOutImage" bson:"pathOutImage"`
}

type ResponseTrxOneTimeMember struct {
	TotalTrxMemberOneTime int64              `json:"totalTrxMemberOneTime" bson:"totalTrxMemberOneTime"`
	DataTrxMember         []TrxOneTimeMember `json:"dataTrxMember" bson:"dataTrxMember"`
}

type InquiryTrxLostTicket struct {
	ID                        string  `json:"_id" bson:"_id"`
	DocNo                     string  `json:"docNo" bson:"docNo"`
	InquiryDatetimeLostTicket string  `json:"inquiryDatetimeLostTicket" bson:"inquiryDatetimeLostTicket"`
	ProductCode               string  `json:"productCode" bson:"productCode"`
	ProductName               string  `json:"productName" bson:"productName"`
	Username                  string  `json:"username" bson:"username"`
	ShiftCode                 string  `json:"shiftCode" bson:"shiftCode"`
	CustName                  string  `json:"custName" bson:"custName"`
	VehicleNumber             string  `json:"vehicleNumber" bson:"vehicleNumber"`
	VehicleType               string  `json:"vehicleType" bson:"vehicleType"`
	CustIdentifier            string  `json:"custIdentifier" bson:"custIdentifier"`
	GrandTotal                float64 `json:"grandTotal" bson:"grandTotal"`
	ChargeAmount              float64 `json:"chargeAmount" bson:"chargeAmount"`
	PaymentMethod             string  `json:"paymentMethod" bson:"paymentMethod"`
	QrCode                    string  `json:"qrCode" bson:"qrCode"`
	Type                      string  `json:"type" bson:"type"`
	Remark                    string  `json:"remark" bson:"remark"`
}

type RequestUpdateTrxInquiryQris struct {
	DocNo           string  `json:"docNo" bson:"docNo"`
	PaymentRefDocNo string  `json:"paymentRefDocNo" bson:"paymentRefDocNo"`
	QrPayment       string  `json:"qrPayment" bson:"qrPayment"`
	CustName        string  `json:"custName" bson:"custName"`
	CustIdentifier  string  `json:"custIdentifier" bson:"custIdentifier"`
	GrandTotal      float64 `json:"grandTotal" bson:"grandTotal"`
}

type ResponseFindQrLostTicket struct {
	QrLostTicket string `json:"qrCodeLostTicket" bson:"qrCodeLostTicket"`
}
