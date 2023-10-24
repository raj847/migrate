package constans

const (
	// SUCCESS_CODE = "101"
	// PENDING_CODE = "102"
	// FAILED_CODE  = "103"

	// Error berhubungan dengan data
	DATA_NOT_FOUND_CODE = "201"
	WORKER_SUCCESS_CODE = "200"

	// SYSTEM_ERROR_CODE    = "501"
	// UNDEFINED_ERROR_CODE = "502"

	//EROR CODE
	SUCCESS_CODE            = "0000" //SUCCESS CODE
	DATA_ERROR_CODE         = "0001" //DATA ERROR/Data tidak ditemukan
	CARD_ALREADY_USED_CODE  = "0002" //Sesi Kartu Masih Ada
	FAILED_TRX_CODE         = "0003" //Transaksi Gagal
	MALFUNCTION_SYSTEM_CODE = "0004" //Terjadi Kesalahan Teknis
	UNDEFINED_ERROR_CODE    = "0005" //Terjadi Kesalahan Teknis
	VALIDATE_ERROR_CODE     = "0006" //Data Tidak ditemukan

	EMPTY_VALUE     = ""
	MANUAL          = "MANUAL"
	NONE_QRCODE     = "xxxxxxxxxxxxxxxxxxxxxxxxxx"
	NULL_LONG_VALUE = -99
	GET_KEY         = "testjwt"
	TRUE_VALUE      = true
	FALSE_VALUE     = false
	PREFIX_MKP      = "MKP"

	PRODUCT_COLLECTION                           = "product_col"
	TRX_LOG_COLLECTION                           = "trx_items_prepaid"
	YES                                          = "Y"
	NO                                           = "N"
	ACTIVE                                       = "ACTIVE"
	PRE_ACTIVE                                   = "PRE ACTIVE"
	EXPIRED                                      = "EXPIRED"
	TRX_OUTSTANDING_COLLECTION                   = "TrxOutstanding"
	TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION   = "TrxDepositCounterOutstanding"
	TRX_DEPOSIT_COUNTER_COLLECTION               = "TrxDepositCounter"
	LOST_TICKET_COLLECTION                       = "LostTicket"
	TRX_OUTSTANDING_FOR_CLEAR_SESSION_COLLECTION = "TrxOutstandingForClearSession"
	TRX_ADD_INFO_COLLECTION                      = "TrxAddInfo"
	TRX_PRODUCT_CUSTOM_COLLECTION                = "TrxProductCustom"
	TRX_COLLECTION                               = "Trx"
	MEMBER_COLLECTIONS                           = "Member"
	TRX_MEMBER_COLLECTIONS                       = "TrxMember"
	TRX_INVOICE_DETAIL_ITEM_COLLECTIONS          = "TrxInvoiceDetailItem"
	OU_STRUCTURE_COLLECTIONS                     = "OuStructure"
	TABLE_PARTNER                                = "DATATYPE_PARTNER"
	TYPE_PARTNER_FREE_PASS                       = "FREEPASS"
	TYPE_PARTNER_SPECIAL_MEMBER                  = "SPECIAL MEMBER"
	TYPE_PARTNER_ONE_TIME                        = "MEMBER ONE TIME"
	MEMBER                                       = "MEMBER"
	ROLE_TYPE_GENERAL                            = "GENERAL"
	ASCENDING                                    = "ASC"
	DESCENDING                                   = "DESC"

	TABLE_POLICY_OU_PRODUCT                        = "PolicyOUProduct"
	TABLE_POLICY_OU_PRODUCT_DEPOSIT_COUNTER        = "PolicyOUProductDepositCoutner"
	TABLE_DEVICE_OU                                = "DeviceOU"
	TABLE_USER_LOGIN                               = "UserLogin"
	EMPTY_VALUE_INT                                = 0
	ONE_HOUR_IN_MINUTE                             = 60
	TWELVE_HOUR_OVERNIGHT                          = 12
	ONE_DAY_OVERNIGHT                              = 1
	MILLISECOND                                    = 60000
	DATATYPE_TRX_LOCAL                             = "TRX_LOCAL"
	DATATYPE_TRX_QRIS                              = "TRX_LOCAL_QRIS"
	CHANNEL_REDIS_PG_PARKING                       = "PG_PARKING_CHECKIN"
	CHANNEL_REDIS_PG_PARKING_CHECKOUT              = "PG_PARKING_CHECKOUT"
	PG_PARKING_CHECKOUT_DEPOSIT_COUNTER            = "PG_PARKING_CHECKOUT_DEPOSIT"
	CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF   = "PG_PARKING_CHECKOUT_EXCLUDE_SF"
	PG_PARKING_CHECKOUT_DEPOSIT_COUNTER_EXCLUDE_SF = "PG_PARKING_CHECKOUT_DEPOSIT_COUNTER_EXCLUDE_SF"
	CHANNEL_PG_INQUIRY_PAYMENT                     = "PG_INQUIRY_PAYMENT"
	CHANNEL_PG_TRX_MEMBER                          = "PG_TRX_MEMBER"
	CHANNEL_PG_SYNC_PRODUCT                        = "PG_SYNC_PRODUCT"
	NON_CALLBACK_PAYMENT                           = "NON_CALLBACK_PAYMENT"
	CHANNEL_CALLBACK_REDIS                         = "CALLBACK_PAYMENT"
	PAYMENT_CATEGORY                               = "QRIS"
	EMPTY_NONE                                     = "-"
	VALIDATE_MEMBER_NOPOL                          = "VEHICLE"
	VALIDATE_MEMBER_CARD                           = "CARD"
	VALIDATE_MEMBER_MIX                            = "MIX"
	PAYMENT_METHOD_QRIS                            = "QRIS"
	CHARGE_WITH_RATES                              = "CHARGE_RATES"
	MAX_AMOUNT                                     = "MAX_AMOUNT"

	MILISECOND                      = 60000
	QUEUE_PG_SCHEDULING_PROGRESSIVE = "pg.scheduling-progressive-rate"
	DATE_TO_FREEPASS_MEMBER         = "2030-12-31"

	SETTLEMENT_CODE_QRIS          = "07"
	SETTLEMENT_CODE_CASH          = "08"
	DOCUMENT_SUCCESS_CODE         = "01"
	HOUR_24                       = 1440
	TABLE_COMBO_CONSTANT          = "ComboConstant"
	VEHICLE_NUMBER                = "VEHICLE_NUMBER"
	VEHICLE                       = "VEHICLE"
	MIX                           = "MIX"
	CARD_NUMBER                   = "CARD_NUMBER"
	ACTIVE_MEMBER                 = "ACTIVE"
	EXPIRED_MEMBER                = "EXPIRED"
	PRE_ACTIVE_MEMBER             = "PRE ACTIVE"
	VALIDATE_EXPIRED_MEMBER       = 7
	TABLE_LIST_PAYMENT_METHOD     = "ListPaymentMethod"
	TYPE_PAYMENT_ONLINE_MICROSITE = "TRX_PAYMENT_MICROSITE"
	INQUIRY                       = "INQUIRY"
	CONFIRM                       = "CONFIRM"
	SUCCESS                       = "SUCCESS"
	TRX_MEMBER_TEMPORARY          = "MemberTemporary"
	AMOUNT_TYPE                   = "AMT"
	PERCENTAGE_TYPE               = "PCT"
)
