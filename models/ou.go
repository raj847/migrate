package models

type Ou struct {
	ID              int64  `json:"id" bson:"ID"`
	OuId            int64  `json:"ouId" bson:"ouId"`
	OuName          string `json:"ouName" bson:"ouName"`
	OuCode          string `json:"ouCode" bson:"ouCode"`
	OuSubBranchId   int64  `json:"ouSubBranchId" bson:"ouSubBranchId"`
	OuSubBranchName string `json:"ouSubBranchName" bson:"ouSubBranchName"`
	OuSubBranchCode string `json:"ouSubBranchCode" bson:"ouSubBranchCode"`
	MainOuId        int64  `json:"mainOuId" bson:"mainOuId"`
	MainOuCode      string `json:"mainOuCode" bson:"mainOuCode"`
	MainOuName      string `json:"mainOuName" bson:"mainOuId"`
}
