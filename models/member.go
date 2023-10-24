package models

type RequestMember struct {
	PartnerCode   string `json:"partnerCode" bson:"partnerCode"`
	FirstName     string `json:"firstName" bson:"firstName" validate:"required"`
	LastName      string `json:"lastName" bson:"lastName"`
	RoleType      string `json:"roleType" bson:"roleType"`
	PhoneNumber   string `json:"phoneNumber" bson:"phoneNumber"`
	Email         string `json:"email" bson:"email"`
	TypePartner   string `json:"typePartner" bson:"typePartner" validate:"required"`
	CardNumber    string `json:"cardNumber" bson:"cardNumber"`
	VehicleNumber string `json:"vehicleNumber" bson:"vehicleNumber"`
	DateFrom      string `json:"dateFrom" bson:"dateFrom" validate:"required"`
	DateTo        string `json:"dateTo" bson:"dateTo"`
	ProductId     int64  `json:"productId" bson:"productId"`
}

type RequestMemberPayment struct {
	PartnerId           int64  `json:"partnerId"`
	PartnerCode         string `json:"partnerCode"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	RoleType            string `json:"roleType"`
	PhoneNumber         string `json:"phoneNumber"`
	Email               string `json:"email"`
	TypePartner         string `json:"typePartner" validate:"required"`
	CardNumber          string `json:"cardNumber"`
	VehicleNumber       string `json:"vehicleNumber"`
	DateFrom            string `json:"dateFrom" validate:"required"`
	DateTo              string `json:"dateTo"`
	ProductId           int64  `json:"productId"`
	ProductMembershipId int64  `json:"productMembershipId"`
	IsExtend            bool   `json:"isExtend" bson:"isExtend"`
}

type Member struct {
	ID                  int64   `json:"id"`
	PartnerCode         string  `json:"partnerCode" bson:"partnerCode"`
	FirstName           string  `json:"firstName" bson:"firstName"`
	LastName            string  `json:"lastName" bson:"lastName"`
	RoleType            string  `json:"roleType" bson:"roleType"`
	PhoneNumber         string  `json:"phoneNumber" bson:"phoneNumber"`
	Email               string  `json:"email" bson:"email"`
	Active              string  `json:"active" bson:"active"`
	ActiveAt            string  `json:"activeAt" bson:"activeAt"`
	NonActiveAt         *string `json:"nonActiveAt" bson:"nonActiveAt"`
	OuId                int64   `json:"ouId" bson:"ouId"`
	TypePartner         string  `json:"typePartner" bson:"typePartner"`
	CardNumber          string  `json:"cardNumber" bson:"cardNumber"`
	VehicleNumber       string  `json:"vehicleNumber" bson:"vehicleNumber"`
	RegisteredDatetime  string  `json:"registeredDatetime" bson:"registeredDatetime"`
	DateFrom            string  `json:"dateFrom" bson:"dateFrom"`
	DateTo              string  `json:"dateTo" bson:"dateTo"`
	ProductId           int64   `json:"productId" bson:"productId"`
	ProductCode         string  `json:"productCode" bson:"productCode"`
	ProductMembershipId int64   `json:"productMembershipId" bson:"productMembershipId"`
	Price               float64 `json:"price" bson:"price"`
	ServiceFee          float64 `json:"serviceFee" bson:"serviceFee"`
	IsPctSfee           string  `json:"isPctSfee" bson:"isPctSfee"`
	DueDate             int64   `json:"dueDate" bson:"dueDate"`
	DiscType            string  `json:"discType" bson:"discType"`
	DiscAmount          float64 `json:"discAmount" bson:"discAmount"`
	DiscPct             int64   `json:"discPct" bson:"discPct"`
	GracePeriodDate     int64   `json:"gracePeriodDate" bson:"gracePeriodDate"`
	Username            string  `json:"username" bson:"username"`
	IsExtendMember      bool    `json:"isExtendMember" bson:"isExtendMember"`
	CreatedAt           string  `json:"createdAt" bson:"createdAt"`
	CreatedBy           string  `json:"createdBy" bson:"createdBy"`
	UpdatedAt           string  `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy           string  `json:"updatedBy" bson:"updatedBy"`
}

type MemberCustom struct {
	ID                 int64   `json:"id"`
	PartnerCode        string  `json:"partnerCode" bson:"partnerCode"`
	FirstName          string  `json:"firstName" bson:"firstName"`
	LastName           string  `json:"lastName" bson:"lastName"`
	RoleType           string  `json:"roleType" bson:"roleType"`
	PhoneNumber        string  `json:"phoneNumber" bson:"phoneNumber"`
	Status             string  `json:"status" bson:"status"`
	Email              string  `json:"email" bson:"email"`
	Active             string  `json:"active" bson:"active"`
	ActiveAt           string  `json:"activeAt" bson:"activeAt"`
	NonActiveAt        *string `json:"nonActiveAt" bson:"nonActiveAt"`
	OuId               int64   `json:"ouId" bson:"ouId"`
	TypePartner        string  `json:"typePartner" bson:"typePartner"`
	CardNumber         string  `json:"cardNumber" bson:"cardNumber"`
	VehicleNumber      string  `json:"vehicleNumber" bson:"vehicleNumber"`
	RegisteredDatetime string  `json:"registeredDatetime" bson:"registeredDatetime"`
	DateFrom           string  `json:"dateFrom" bson:"dateFrom"`
	DateTo             string  `json:"dateTo" bson:"dateTo"`
	ProductId          int64   `json:"productId" bson:"productId"`
	ProductCode        string  `json:"productCode" bson:"productCode"`
	CreatedAt          string  `json:"createdAt" bson:"createdAt"`
	CreatedBy          string  `json:"createdBy" bson:"createdBy"`
	UpdatedAt          string  `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy          string  `json:"updatedBy" bson:"updatedBy"`
}

type ValidateRegis struct {
	RequestDateFrom string `json:"requestDateFrom"`
	RequestDateTo   string `json:"requestDateTo"`
}

type ValidateMember struct {
	ReqDateFrom    string `json:"reqDateFrom"`
	ReqDateTo      string `json:"reqDateTo"`
	DateFrom       string `json:"dateFrom"`
	DateTo         string `json:"dateTo"`
	ReqCardNumber  string `json:"reqCardNumber"`
	CardNumber     string `json:"cardNumber"`
	ReqProductId   int64  `json:"reqProductId"`
	ProductId      int64  `json:"productId"`
	ReqPartnerCode string `json:"reqPartnerCode"`
	PartnerCode    string `json:"partnerCode"`
}

type RequestEditPartner struct {
	ID           int64 `json:"id" validate:"required"`
	ActiveMember bool  `json:"activeMember"`
}

type ResponseIsPartnerExistsByID struct {
	PartnerCode string `json:"partner_code"`
	DateFrom    string `json:"date_from"`
	DateTo      string `json:"date_to"`
	ActiveAt    string `json:"active_at"`
}

type ValidateTime struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type EditPartner struct {
	ID          int64   `json:"id"`
	Active      string  `json:"active" bson:"active"`
	ActiveAt    string  `json:"activeAt" bson:"activeAt"`
	NonActiveAt *string `json:"nonActiveAt" bson:"nonActiveAt"`
}

type RequestRemoveMember struct {
	ID int64 `json:"id" validate:"required"`
}

type RequestSearchTrxLog struct {
	DateFrom  string `json:"dateFrom"`
	DateTo    string `json:"dateTo"`
	OuId      int64  `json:"ouId"`
	StatusTrx string `json:"statusTrx"`
}

type WorkerRequestAddMemberExt struct {
	KodeMember      string `json:"kode_member" validate:"required"`
	NamaDepan       string `json:"nama_depan" validate:"required"`
	NamaBelakang    string `json:"nama_belakang" validate:"required"`
	NoHandphone     string `json:"no_handphone" valdiate:"required"`
	RoleType        string `json:"roleType"`
	Email           string `json:"email"`
	TanggalMulai    string `json:"tanggal_mulai"`
	TanggalBerakhir string `json:"tanggal_berakhir"`
	MerchantKey     string `json:"merchant_key"`
	NomorKartu      string `json:"nomor_kartu"`
	PlatNomor       string `json:"plat_nomor" validate:"required"`
	TypePartner     string `json:"typePartner"`
	KodeProdukMKP   string `json:"kode_produk_mkp"`
}

type WorkerResponseAddMemberExt struct {
	StatusCode      string      `json:"statusCode" bson:"statusCode"`
	StatusDesc      string      `json:"statusDesc" bson:"statusDesc"`
	SnapshotData    interface{} `json:"snapshotData" bson:"snapshotData"`
	URL             string      `json:"url" bson:"URL"`
	RequestDatetime string      `json:"requestDatetime"`
	ResponseBody    interface{} `json:"responseBody" bson:"responseBody"`
}

type WorkerRequestActiveMemberExt struct {
	KodeMember  string `json:"kode_member"`
	AktifMember bool   `json:"aktif_member"`
	MerchantKey string `json:"merchant_key"`
}

type TrxMember struct {
	DocNo              string  `json:"docNo" bson:"docNo"`
	PartnerCode        string  `json:"partnerCode" bson:"partnerCode"`
	FirstName          string  `json:"firstName" bson:"firstName"`
	LastName           string  `json:"lastName" bson:"lastName"`
	RoleType           string  `json:"roleType" bson:"roleType"`
	PhoneNumber        string  `json:"phoneNumber" bson:"phoneNumber"`
	Email              string  `json:"email" bson:"email"`
	Active             string  `json:"active" bson:"active"`
	ActiveAt           string  `json:"activeAt" bson:"activeAt"`
	NonActiveAt        *string `json:"nonActiveAt" bson:"nonActiveAt"`
	OuId               int64   `json:"ouId" bson:"ouId"`
	TypePartner        string  `json:"typePartner" bson:"typePartner"`
	CardNumber         string  `json:"cardNumber" bson:"cardNumber"`
	VehicleNumber      string  `json:"vehicleNumber" bson:"vehicleNumber"`
	RegisteredDatetime string  `json:"registeredDatetime" bson:"registeredDatetime"`
	DateFrom           string  `json:"dateFrom" bson:"dateFrom"`
	DateTo             string  `json:"dateTo" bson:"dateTo"`
	ProductId          int64   `json:"productId" bson:"productId"`
	ProductCode        string  `json:"productCode" bson:"productCode"`
}

type ResponseFindPartner struct {
	PartnerId   int64  `json:"partnerId"`
	DateFrom    string `json:"dateFrom"`
	DateTo      string `json:"dateTo"`
	CardNumber  string `json:"cardNumber"`
	OuId        int64  `json:"ouId"`
	ProductId   int64  `json:"productId"`
	PartnerCode string `json:"partnerCode"`
}

type RequestFindPartnerAdvance struct {
	Keyword         string `json:"keyword"`
	OuList          string `json:"ouList"`
	StatusMember    string `json:"statusMember"`
	AscDesc         string `json:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName"`
	Limit           int64  `json:"limit"`
	Offset          int64  `json:"offset"`
}

type ResponseFindMemberAdvance struct {
	CountMemberList int64    `json:"countMemberList"`
	Data            []Member `json:"data"`
}

type UpdateMember struct {
	PartnerCode string `json:"partnerCode" bson:"partnerCode"`
	CardNumber  string `json:"cardNumber" bson:"cardNumber"`
	DateFrom    string `json:"dateFrom" bson:"dateFrom"`
	DateTo      string `json:"dateTo" bson:"dateTo"`
	UpdatedAt   string `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy   string `json:"updatedBy" bson:"updatedBy"`
}

type ValPartnerAdvance struct {
	PartnerCode    string `json:"partnerCode"`
	OuId           int64  `json:"ouId"`
	CardNumber     string `json:"cardNumber"`
	VehicleNumber  string `json:"vehicleNumber"`
	RegisteredType string `json:"registeredType"`
	DateNow        string `json:"dateNow"`
}

type RequestPartnerInternal struct {
	PartnerCode    string `json:"partnerCode"`
	FirstName      string `json:"firstName" validate:"required"`
	LastName       string `json:"lastName"`
	RoleType       string `json:"roleType" validate:"required"`
	TypePartner    string `json:"typePartner" validate:"required"`
	PhoneNumber    string `json:"phoneNumber" validate:"required"`
	RegisteredType string `json:"registeredType" validate:"required"`
	Email          string `json:"email"`
	StartDate      string `json:"startDate" validate:"required"`
	EndDate        string `json:"endDate" validate:"required"`
	OuID           int64  `json:"ouId" validate:"required"`
	CardNumber     string `json:"cardNumber"`
	VehicleNumber  string `json:"vehicleNumber"`
	ProductID      int64  `json:"productId" validate:"required"`
	Remark         string `json:"remark"`
}

type FindPolicyOuPartnerByIndex struct {
	ID                 int64  `json:"id"`
	PartnerCode        string `json:"partnerCode"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	OuId               int64  `json:"ouId"`
	OuName             string `json:"ouName"`
	ProductId          int64  `json:"productId"`
	ProductName        string `json:"productName"`
	Remark             string `json:"remark"`
	RegisteredType     string `json:"registeredType"`
	RegisteredDatetime string `json:"registeredDatetime"`
	DateFrom           string `json:"dateFrom"`
	DateTo             string `json:"dateTo"`
	Status             string `json:"status"`
	TypePartner        string `json:"typePartner"`
	RoleType           string `json:"roleType"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phoneNumber"`
	VehicleNumber      string `json:"vehicleNumber"`
	CardNumber         string `json:"cardNumber"`
	RefPartnerId       *int64 `json:"refPartnerId"`
	CreatedAt          string `json:"createdAt"`
	CreatedBy          string `json:"createdBy"`
	UpdatedAt          string `json:"updatedAt"`
	UpdatedBy          string `json:"updatedBy"`
}

type RequestExtendMember struct {
	MemberName    string `json:"memberName"`
	DateFrom      string `json:"dateFrom" validate:"required"`
	DateTo        string `json:"dateTo" validate:"required"`
	VehicleNumber string `json:"vehicleNumber"`
	CardNumber    string `json:"cardNumber"`
}

type PolicyOuPartner struct {
	ID                 int64  `json:"id"`
	PartnerId          int64  `json:"partnerId"`
	OuId               int64  `json:"ouId" validate:"required"`
	TypePartner        string `json:"typePartner" validate:"required"`
	RoleType           string `json:"roleType"`
	CardNumber         string `json:"cardNumber" validate:"required"`
	VehicleNumber      string `json:"vehicleNumber" validate:"required"`
	RegisteredDatetime string `json:"registeredDatetime"`
	RegisteredType     string `json:"registeredType"`
	DateFrom           string `json:"dateFrom"`
	DateTo             string `json:"dateTo"`
	ProductId          int64  `json:"productId" validate:"required"`
	Remark             string `json:"remark"`
	CreatedAt          string `json:"createdAt"`
	CreatedBy          string `json:"createdBy"`
	UpdatedAt          string `json:"updatedAt"`
	UpdatedBy          string `json:"updatedBy"`
}

type EditPartnerInternal struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	DateFrom      string `json:"dateFrom"`
	DateTo        string `json:"dateTo"`
	ProductId     int64  `json:"productId"`
	ProductCode   string `json:"productCode"`
	RoleType      string `json:"roleType"`
	PhoneNumber   string `json:"phoneNumber"`
	CardNumber    string `json:"cardNumber"`
	VehicleNumber string `json:"vehicleNumber"`
	ActiveAt      string `json:"activeAt"`
	UpdatedBy     string `json:"updatedBy"`
}

type ResponseExtendMember struct {
	Message       string `json:"message"`
	MemberName    string `json:"memberName"`
	VehicleNumber string `json:"vehicleNumber"`
}

type ResponseGetMemberName struct {
	FirstName   string `json:"firstName"`
	OuId        int64  `json:"ouId"`
	OuName      string `json:"ouName"`
	RoleType    string `json:"roleType"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type GetMemberNameByIndex struct {
	MemberName  string `json:"memberName"`
	OuId        int64  `json:"ouId"`
	OuName      string `json:"ouName"`
	RoleType    string `json:"roleType"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type EditPolicyOuPartner struct {
	ID             int64  `json:"id"`
	TypePartner    string `json:"typePartner"`
	RoleType       string `json:"roleType"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	RegisteredType string `json:"registeredType"`
	VehicleNumber  string `json:"vehicleNumber"`
	CardNumber     string `json:"cardNumber"`
	ProductId      int64  `json:"productId"`
	UpdatedAt      string `json:"updatedAt"`
	UpdatedBy      string `json:"updatedBy"`
}
