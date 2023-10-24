package models

//-- INPUT DECRYPT UTILS --
type RequestDecryptMerchKey struct {
	CryptoText string `json:"cryptotext" validate:"required"`
}

//-- RESPONSE ENCRYPT UTILS --
type ResultEncryptMerchKey struct {
	Result string `json:"result"`
}

//-- RESPONSE DECRYPT & INPUT ENCRYPT UTILS --
type MerchantKey struct {
	ID				int64  `json:"ID"`
	OuId            int64  `json:"ouId" validate:"required"`
	OuName          string `json:"ouName" validate:"required"`
	OuCode          string `json:"ouCode" validate:"required"`
	OuSubBranchId   int64  `json:"ouSubBranchId" validate:"required"`
	OuSubBranchName string `json:"ouSubBranchName"`
	OuSubBranchCode string `json:"ouSubBranchCode"`
	MainOuId        int64  `json:"mainOuId"`
	MainOuCode      string `json:"mainOuCode"`
	MainOuName      string `json:"mainOuName"`
}
