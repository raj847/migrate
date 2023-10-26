package repositories

import (
	"database/sql"
	"github.com/raj847/togrpc/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceRepository interface {
	FindDeviceByDeviceID(deviceID string) (models.DeviceDB, error)
}

type OuMongoRepository interface {
	FindOuByOuId(OuId int64) (models.Ou, error)
}

type ProductMongoRepository interface {
	FindProductByProductCode(productCode string, ouId int64) (models.PolicyOuProductWithRules, error)
	FindPolicyOuProductByAdvance(OuID int64, ProductId int64) (models.PolicyOuProductWithRules, error)
	GetPolicyOuProductList() ([]models.PolicyOuProductWithRules, error)
	GetPolicyOuProductListAdvance(product models.RequestGetPolicyOuProductList) ([]models.PolicyOuProductWithRules, error)
	CountGetPolicyOuProductListAdvance(keyword string) (int64, error)
	GetPolicyOuProductByOuId(OuID int64) (*[]models.PolicyOuProductWithRules, error)
	FindProductDepositByProductCode(productCode string, ouId int64) (models.PolicyOuProductDepositCounterWithRules, error)
	GetPolicyOuProductDepositByOuId(OuID int64) (*[]models.PolicyOuProductDepositCounterWithRules, error)
	GetPolicyOuProductDepositListAdvance(product models.RequestGetPolicyOuProductList) ([]models.PolicyOuProductDepositCounterWithRules, error)
}

type TrxMongoRepository interface {
	AddTrxCheckIn(trx models.Trx) (*primitive.ObjectID, error)
	AddTrx(trx models.Trx) (string, error)
	ValTrxExistByUUIDCard(uuidCard string) (string, error)
	FindTrxOutstandingByDocNo(docNo string, productCode string) (models.ResultFindTrxOutstanding, error)
	FindTrxOutstandingByID(ID primitive.ObjectID) (models.Trx, error)
	FindTrxOutstandingByUUID(uuid string, productCode string) (models.ResultFindTrxOutstanding, error)
	UpdateProgressivePriceByDocNo(docNo string, progressivePrice float64) error
	RemoveTrxByID(ID primitive.ObjectID) error
	LogActivity(logActivity models.LogActivityTrx) error
	InvoiceTrx(invoiceTrx models.InvoiceTrx) error
	UpdateTrxInvoiceItemForTrxOutstanding(docNo string, productCode string, trxInvoiceItems models.TrxInvoiceItem) error
	IsTrxOutstandingExistByUUIDCard(uuidCard string) (models.Trx, bool)
	GetTrxMemberByPeriodList(docNo string, inquiryDate string) ([]models.TrxMember, error)
	GetTrxMemberByPeriodListByProductCode(docNo, productCode, inquiryDate string) ([]models.TrxMember, error)
	GetTrxInvoiceDetailsItemsLists(trxMemberList []models.TrxMember) (float64, error)
	AddTrxAddInfoInterfaces(filter interface{}, updateSet interface{}) error
	IsTrxAddInfoInterfacesExistsByDocNo(docNo string) (map[string]interface{}, bool)
	GetTrxInvoiceDetailsItemsList(docNo, productCode string) ([]models.TrxInvoiceDetailItem, error)
	IsTrxOutstandingByUUID(uuidCard string) (*models.Trx, bool, error)
	RemoveTrxByDocNo(docNo string) error
	IsTrxOutstandingByCardNumber(cardNumber string) (*models.Trx, bool, error)
	IsTrxOutstandingByDocNoForCustom(docNo string) (trx *models.Trx, exists bool, err error)
	IsTrxOutstandingByDocNo(docNo string, productCode string) (*models.ResultFindTrxOutstanding, bool, error)
	GetTrxListForSyncDataFailed() (trxList []models.Trx, err error)
	UpdateTrxByInterface(filter interface{}, updateSet interface{}) error
	IsHistoryTrxListByCheckInAndCheckoutDate(cardNumberUuid string, productCode, checkInDate, vehicleNumber string) (trx *models.Trx, exists bool, err error)
	GetTrxListByDocDate(docDate string) (trxList []models.Trx, err error)
	AddTrxOutstandingForClearSession(trx models.TrxOutstandingForClearSession) (ID *primitive.ObjectID, err error)
	IsTrxOutstandingForClearSession(docNo string) (trx *models.TrxOutstandingForClearSession, exists bool, err error)
	GetListTrxByRangeDate(dateFrom, dateTo, ouCode string, limit int64) (*[]models.Trx, error)
	FindTrxOutstandingByDocNoCustom(docNo string) (*models.ResultFindTrxOutstanding, error)
	FindTrxOutstandingByUUIDCustom(uuidCard string) (*models.ResultFindTrxOutstanding, error)
	UpdateProductById(keyword, productCode, productName string) error
	IsTrxProductCustomExistsByKeyword(keyword string) (*models.TrxProductCustom, bool)
	RemoveTrxProductCustom(keyword string) error
	UpdateProgressivePriceAndChargeAmountByDocNo(docNo string, progressivePrice, chargeAmount float64) error
	FindTrxFlgSyncFalse(flgSyncFalse bool, checkInDateFrom, checkInDateTo string) (*[]models.Trx, error)
	IsTrxOutstandingByDocNoForLostTicket(docNo string) (trx *models.TrxWithId, exists bool, err error)
	IsTrxOutstandingByDocNoNew(DocNo string) (models.Trx, bool)
	AddTrxCheckInv2(docNo string, trx models.Trx) error
	AddTrxInvoiceDetailItem(trxInvoice models.TrxInvoiceDetailItem) error
	FindIDTrxOutstandingByCard(uuidCard string) (models.TrxWithId, error)
}

type UserMongoRepository interface {
	FindUserByIndex(username string) (*models.UserLoginLocal, error)
	FindDeviceIdByIndex(username, deviceId string) (*models.DeviceSync, error)
}

type MemberMongoRepository interface {
	AddTrxMember(trxMember models.TrxMember) error
	AddMemberTemp(memberData models.Member) (*primitive.ObjectID, error)
}

type MemberRepository interface {
	AddMember(requestMember models.Member, tx *sql.Tx) (int64, error)
	IsPartnerExistsByCardNumber(cardNumber string) (models.ResponseFindPartner, bool)
	IsCardNumberExists(CardNumber string) bool
	IsMemberActiveExistsByUUIDCard(cardNumber string, date string) (models.Member, bool)
	IsPartnerExistsById(id int64) (models.ResponseIsPartnerExistsByID, bool)
	ActivationPartner(updatePartner models.EditPartner, tx *sql.Tx) error
	GetListPartnerAdvance(requestPartner models.RequestFindPartnerAdvance) ([]models.Member, error)
	RemoveMember(id int64, tx *sql.Tx) error
	FindMemberActiveByDate(Date string, RoleType string) (models.SummaryMember, error)
	CountGetListPartnerAdvance(requestPartner models.RequestFindPartnerAdvance) (int64, error)
	IsMemberFreePassByIndex(uuidCard string) (*models.Member, bool, error)
	IsMemberByAdvanceIndex(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.Member, bool, error)
	GetMemberActiveListByPeriod(uuidCard, vehicleNumber, checkinDate string, inquiryDate string, memberBy string, isFreePass bool) ([]models.Member, error)
	IsMemberExistsByPartnerCode(partnerCode, startDate, keyword string) (*models.Member, bool, error)
	IsMemberByAdvanceIndexCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.Member, bool, error)
	GetMemberActiveListByPeriodCustom(uuidCard, vehicleNumber, checkinDate string, inquiryDate string, memberBy string, isFreePass bool, isHotelMember bool) ([]models.Member, error)
	UpdateMemberByPartnerCode(updateMember models.UpdateMember, tx *sql.Tx) error
	IsMemberActiveByCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.MemberCustom, bool, error)
	GetListMemberActiveByCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*[]models.MemberCustom, bool, error)
	GetPolicyOuPartnerByIndex(partner models.RequestFindPartnerAdvance) (*[]models.FindPolicyOuPartnerByIndex, error)
	CountPolicyOuPartnerByIndex(partner models.RequestFindPartnerAdvance) (*int64, error)
	ValPartnerActiveByAdvance(valPartnerAdvance models.ValPartnerAdvance, productCode, typePartner string) error
	FindMemberByIndex(member models.RequestExtendMember) (*[]models.FindPolicyOuPartnerByIndex, error)
	EditPartnerInternal(partner models.EditPartnerInternal, tx *sql.Tx) (bool, error)
	GetMemberNameList(partner models.RequestFindPartnerAdvance) (*[]models.ResponseGetMemberName, error)
	EditPolicyOuPartner(policyOuPartner models.EditPolicyOuPartner, tx *sql.Tx) (bool, error)
}

type ReportTrxMongoRepository interface {
	GetTrxOmzetByDate(dateFrom string, dateTo string) (models.ResponseTrx, error)
	GetTrxOutstandingByDate(dateFrom string, dateTo string) (models.ResponseTrxOutstanding, error)
	GetTrxOvernightByDate() ([]models.TrxInvoiceDetail, error)
	GetRevenueForIntelligenceTrxList(dateFrom string, dateTo string) ([]models.ChartDailyIntelligenceTrxList, error)
	GetTrafficParkingFromCheckin(date string, listTime []string) ([]models.TrafficParking, error)
	GetTrxListAdvance(trx models.TrxItems, isExport bool) ([]models.ResultFindTrxForPostBox, error)
	CountGetTrxListAdvance(trx models.TrxItems) (int64, error)
	GetTrxVehicleByDocDate(requestTrx models.RequestGetTrxOrderByVehicle) ([]models.GetTrxOrderByVehicle, error)
	GetListTrxOutstanding(trx models.RequestFindTrxOutstanding) (*[]models.ResultFindTrxForPostBox, error)
	CountTrxOutstanding(dateFrom string, dateTo string, keyword string) (models.ResponseTrxOutstanding, error)
	GetListTrxLostTicket(lostTicket models.RequestFindLostTicket) (*[]models.FindLostTicket, error)
	AddTrxLostTicket(trxLostTicket models.LostTicket) (*primitive.ObjectID, error)
	CountListTrxLostTicket(lostTicket models.RequestFindLostTicket) (*int64, error)
	GetListTrxOneTimeMember(trx models.RequestFindOneTimeMember, isExport bool) (*[]models.TrxOneTimeMember, error)
	CountGetListTrxOneTimeMember(trx models.RequestFindOneTimeMember) (int64, error)
	FindQrTextByDocNo(docNo string) (*string, error)
	FindTrxLostTicketByDocNo(id primitive.ObjectID) (*models.FindLostTicket, error)
	FindTrxLostTicketByID(ID primitive.ObjectID) (*models.LostTicket, error)
	UpdateTypeTrxLostTicket(ID primitive.ObjectID, typeTrx string, trx models.LostTicket) error
	UpdateQrCodeTrxLostTicket(ID primitive.ObjectID, qrCode string) error
	UpdatePaymentRefDocNoCodeTrxLostTicket(ID primitive.ObjectID, inquiryTrx models.RequestUpdateTrxInquiryQris) error
	FindQrLostTicketByDocNo(docNo string) (result *models.ResponseFindQrLostTicket, exists bool)
}

type TrxMongoDepositCounterRepository interface {
	AddTrxCheckInDepositCenter(trx models.TrxDepositCounter) (*primitive.ObjectID, error)
	AddTrxDepositCounter(trx models.TrxDepositCounter) (string, error)
	IsTrxDepositCenterOutstandingByDocNo(docNo, productCode string) (*models.ResultFindTrxDepositCounterOutstanding, error)
	FindTrxDepositCenterOutstandingByID(ID primitive.ObjectID) (models.TrxDepositCounter, error)
	RemoveTrxDepositCounterByID(ID primitive.ObjectID) error
	GetListTrxDepositCounter(trxDepo models.TrxItemsDepositCounter) (*[]models.FindTrxDepositCounter, error)
	CountGetTrxDepositCounter(trxDepo models.TrxItemsDepositCounter) (int64, error)
}

type CommonRepository interface {
	GetComboConstantListByComboId(comboId string) ([]models.ComboConstant, error)
	AddComboConstan(comboId, comboCode string, comboConstan models.RequestAddComboConstan) error
	AddPaymentMethodList(ouId int64, ouCode string, listPaymentMethod models.ListPaymentMethod) error
	FindListPaymentMethod(ouId int64) (models.ListPaymentMethod, error)
}

type ProductMembershipRepository interface {
	FindProductMembershipById(id int64) (*models.ProductMembership, error)
}
