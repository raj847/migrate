package models

import "time"

type Response struct {
	StatusCode       string      `json:"statusCode"`
	Success          bool        `json:"success"`
	ResponseDatetime time.Time   `json:"responseDatetime"`
	Result           interface{} `json:"result"`
	Message          string      `json:"message"`
}

type ResponseTrxCloud struct {
	StatusCode       string    `json:"statusCode"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           struct {
		Trx            Trx                      `json:"trxLocal"`
		TrxInvoiceItem ResultFindTrxOutstanding `json:"trxInvoice"`
	} `json:"result"`
	Message string `json:"message"`
}

type ResponseTrxPaymentOnline struct {
	StatusCode       string    `json:"statusCode"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           struct {
		TrxPaymentOnline []TrxOnlinePayment `json:"trxPaymentOnline"`
	} `json:"result"`
	Message string `json:"message"`
}

type ResponseInquiryPayment struct {
	StatusCode       string    `json:"statusCode"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           struct {
		QrCode          string `json:"qrCode"`
		Type            string `json:"type"`
		PaymentRefDocNo string `json:"paymentRefDocNo"`
	} `json:"result"`
	Message string `json:"message"`
}

type ResponseCheckTrx struct {
	StatusCode       string    `json:"statusCode"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           bool      `json:"result"`
	Message          string    `json:"message"`
}

type ResponseListPaymentMethod struct {
	Status           string    `json:"status"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           []struct {
		CorporateCID                 string      `json:"corporateCID"`
		CorporateID                  int         `json:"corporateID"`
		CorporateName                string      `json:"corporateName"`
		CorporatePaymentCategoryID   int         `json:"corporatePaymentCategoryID"`
		CorporatePaymentCategoryName string      `json:"corporatePaymentCategoryName"`
		CorporatePaymentMethodAlias  string      `json:"corporatePaymentMethodAlias"`
		CorporatePaymentMethodID     int         `json:"corporatePaymentMethodID"`
		CorporatePaymentMethodName   string      `json:"corporatePaymentMethodName"`
		CorporateServiceFee          int         `json:"corporateServiceFee"`
		ID                           int         `json:"id"`
		IsPercentage                 bool        `json:"isPercentage"`
		MdrAmount                    int         `json:"mdrAmount"`
		MdrIsExclude                 bool        `json:"mdrIsExclude"`
		MdrPercentageStatus          bool        `json:"mdrPercentageStatus"`
		MerchantKey                  string      `json:"merchantKey"`
		PaymentCode                  string      `json:"paymentCode"`
		SamNum                       int         `json:"samNum"`
		SettlementKey                interface{} `json:"settlementKey"`
		UserID                       int         `json:"userID"`
		UserUsername                 string      `json:"userUsername"`
	} `json:"result"`
	Message string `json:"message"`
}

type CallbackPayment struct {
	Kodestatus             string  `json:"kodestatus"`
	Statuspayment          string  `json:"statuspayment"`
	Description            string  `json:"description"`
	Merchantnoref          string  `json:"merchantnoref"`
	Noheader               string  `json:"noheader"`
	Banknoref              string  `json:"banknoref"`
	Cardtype               string  `json:"cardtype"`
	Cardpan                string  `json:"cardpan"`
	Mid                    string  `json:"mid"`
	Tid                    string  `json:"tid"`
	Lastbalance            float64 `json:"lastbalance"`
	Currentbalance         float64 `json:"currentbalance"`
	Namakategoripembayaran string  `json:"namakategoripembayaran"`
	SettlementDatetime     string  `json:"settlement_datetime"`
	DeductDatetime         string  `json:"deduct_datetime"`
	Mdr                    float64 `json:"mdr"`
	Potongan               float64 `json:"potongan"`
	Kodepromo              string  `json:"kodepromo"`
	Promoissuer            string  `json:"promoissuer"`
	CreatedAt              string  `json:"created_at"`
}

type TrxOnlinePayment struct {
	SessionID                string           `json:"sessionId" bson:"sessionId"`
	ExtDocNo                 string           `json:"extDocNo" bson:"extDocNo"`
	DocNo                    string           `json:"docNo" bson:"docNo"`
	PaymentRefDocNo          string           `json:"paymentRefDocNo" bson:"paymentRefDocNo"`
	BankNoRef                string           `json:"bankNoRef" bson:"bankNoRef"`
	DocDate                  string           `json:"docDate" bson:"docDate"`
	CheckInDatetime          string           `json:"checkInDatetime" bson:"checkInDatetime"`
	CheckOutDatetime         string           `json:"checkOutDatetime" bson:"checkOutDatetime"`
	ExpiredTimeToCheckOut    string           `json:"expiredTimeToCheckOut" bson:"expiredTimeToCheckOut"`
	PaymentDatetime          string           `json:"paymentDatetime" bson:"paymentDatetime"`
	CheckInTime              int64            `json:"checkInTime" bson:"checkInTime"`
	CheckOutTime             int64            `json:"checkOutTime" bson:"checkOutTime"`
	VehicleNumberIn          string           `json:"vehicleNumberIn" bson:"vehicleNumberIn"`
	VehicleNumberOut         string           `json:"vehicleNumberOut" bson:"vehicleNumberOut"`
	CardNumberUUID           string           `json:"cardNumberUuid" bson:"cardNumberUuid"`
	CardNumber               string           `json:"cardNumber" bson:"cardNumber"`
	CardNumberUUIDIn         string           `json:"cardNumberUuidIn" bson:"cardNumberUuidIn"`
	CardNumberIn             string           `json:"cardNumberIn" bson:"cardNumberIn"`
	TypeCard                 string           `json:"typeCard" bson:"typeCard"`
	BeginningBalance         float64          `json:"beginningBalance" bson:"beginningBalance"`
	CurrentBalance           float64          `json:"currentBalance" bson:"currentBalance"`
	ExtLocalDatetime         string           `json:"extLocalDatetime" bson:"extLocalDatetime"`
	GrandTotal               float64          `json:"grandTotal" bson:"grandTotal"`
	ProductId                int64            `json:"productId" bson:"productId"`
	ProductCode              string           `json:"productCode" bson:"productCode"`
	ProductName              string           `json:"productName" bson:"productName"`
	ProductData              string           `json:"productData" bson:"productData"`
	IsPctServiceFee          string           `json:"isPctServiceFee" bson:"isPctServiceFee"`
	ServiceFee               float64          `json:"serviceFee" bson:"serviceFee"`
	ServiceFeeMember         float64          `json:"serviceFeeMember" bson:"serviceFeeMember"`
	Price                    float64          `json:"price" bson:"price"`
	IncludeSfee              string           `json:"includeSfee" bson:"includeSfee"`
	NormalSfee               string           `json:"normalSfe" bson:"normalSfee"`
	RequestData              string           `json:"requestData" bson:"requestData"`
	RequestOutData           string           `json:"requestOutData" bson:"requestOutData"`
	OuPartnerName            string           `json:"ouPartnerName" bson:"ouPartnerName"`
	OuId                     int64            `json:"ouId" bson:"ouId"`
	OuName                   string           `json:"ouName" bson:"ouName"`
	OuCode                   string           `json:"ouCode" bson:"ouCode"`
	OuSubBranchId            int64            `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName          string           `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode          string           `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId                 int64            `json:"mainOuId" bson:"mainOuId"`
	MainOuCode               string           `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName               string           `json:"mainOuName" bson:"mainOuName"`
	DurationTime             int64            `json:"durationTime" bson:"durationTime"`
	PaymentMethod            string           `json:"paymentMethod" bson:"paymentMethod"`
	Mdr                      float64          `json:"mdr" bson:"mdr"`
	Mid                      string           `json:"mid" bson:"mid"`
	Tid                      string           `json:"tid" bson:"tid"`
	CardPan                  string           `json:"cardPan" bson:"cardPan"`
	CardType                 string           `json:"cardType" bson:"cardType"`
	SettlementDatetime       string           `json:"settlementDatetime" bson:"settlementDatetime"`
	LogTrans                 string           `json:"logTrans" bson:"logTrans"`
	CustomerId               string           `json:"customerId" bson:"customerId"`
	CustomerName             string           `json:"customerName" bson:"customerName"`
	MerchantKey              string           `json:"merchantKey" bson:"merchantKey"`
	MKey                     string           `json:"mKey" bson:"mKey"`
	QrText                   string           `json:"qrText" bson:"qrText"`
	TrxInvoiceItem           []TrxInvoiceItem `json:"trxInvoiceItem" bson:"trxInvoiceItem"`
	PaymentCallback          *CallbackPayment `json:"paymentCallback" bson:"paymentCallback"`
	MemberData               *TrxMember       `json:"memberData" bson:"memberData"`
	FlagSyncData             bool             `json:"flagSyncData" bson:"flagSyncData"`
	FlagErrorCallAPI         bool             `json:"flagErrorCallAPI" bson:"flagErrorCallAPI"`
	ResponseSyncTrx          *string          `json:"responseSyncTrx" bson:"responseSyncTrx"`
	ResponsePaymentBiller    string           `json:"responsePaymentBiller" bson:"responsePaymentBiller"`
	ResponseStatusCodeBiller string           `json:"responseStatusCodeBiller" bson:"responseStatusCodeBiller"`
	Status                   string           `json:"status" bson:"status"`
	StatusDesc               string           `json:"statusDesc" bson:"statusDesc"`
	Msg                      string           `json:"msg" bson:"msg"`
	TrxAddInfo               interface{}      `json:"trxAddInfo" bson:"trxAddInfo"`
}

type ResponseA2PConfirmation struct {
	StatusCode       string    `json:"statusCode"`
	Success          bool      `json:"success"`
	ResponseDatetime time.Time `json:"responseDatetime"`
	Result           struct {
		AcquiringID       int         `json:"acquiringID"`
		Amount            int         `json:"amount"`
		BankNoRef         string      `json:"bankNoRef"`
		CardPan           string      `json:"cardPan"`
		CardType          string      `json:"cardType"`
		CorporateName     string      `json:"corporateName"`
		CreatedAt         string      `json:"createdAt"`
		CurrentBalance    int         `json:"currentBalance"`
		Description       string      `json:"description"`
		DeviceID          string      `json:"deviceID"`
		Discount          int         `json:"discount"`
		ErrorCode         string      `json:"errorCode"`
		ErrorDescription  string      `json:"errorDescription"`
		InquiryWorkerResp interface{} `json:"inquiryWorkerResp"`
		LastBalance       int         `json:"lastBalance"`
		Mdr               int         `json:"mdr"`
		MerchantNoRef     string      `json:"merchantNoRef"`
		Mid               string      `json:"mid"`
		NoHeader          string      `json:"noHeader"`
		PaymentCategory   string      `json:"paymentCategory"`
		PaymentCode       string      `json:"paymentCode"`
		PaymentFee        int         `json:"paymentFee"`
		PromoCode         string      `json:"promoCode"`
		PromoIssuer       string      `json:"promoIssuer"`
		ServiceFee        int         `json:"serviceFee"`
		SettleAt          string      `json:"settleAt"`
		StatusCode        string      `json:"statusCode"`
		StatusPayment     string      `json:"statusPayment"`
		Tid               string      `json:"tid"`
		VendorFee         int         `json:"vendorFee"`
	} `json:"result"`
	Message string `json:"message"`
}

type ProductMembership struct {
	ID                    int64   `json:"id"`
	ProductMembershipCode string  `json:"productMembershipCode"`
	ProductMembershipName string  `json:"productMembershipName"`
	OuId                  int64   `json:"ouId"`
	ProductId             int64   `json:"productId"`
	ProductCode           string  `json:"productCode"`
	ProductName           string  `json:"productName"`
	DueDate               int64   `json:"dueDate"`
	DiscType              string  `json:"discType"`
	DiscAmount            float64 `json:"discAmount"`
	DiscPct               int64   `json:"discPct"`
	GracePeriodDate       int64   `json:"gracePeriodDate"`
	Price                 float64 `json:"price"`
	Active                string  `json:"active"`
	ServiceFee            float64 `json:"serviceFee"`
	IsPct                 string  `json:"isPctSfee"`
	CreateUsername        string  `json:"createUsername"`
	CreatedAt             string  `json:"createdAt"`
	UpdateUsername        string  `json:"updateUsername"`
	UpdatedAt             string  `json:"updatedAt"`
}
