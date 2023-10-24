package models

import (
	"fmt"
	"time"
)

type ComboConstant struct {
	ComboId   string `json:"comboId" bson:"comboId"`
	ComboCode string `json:"comboCode" bson:"comboCode"`
	ComboName string `json:"comboName" bson:"comboName"`
}

type RequestAddComboConstan struct {
	ComboId   string `json:"comboId" bson:"comboId" validate:"required"`
	ComboCode string `json:"comboCode" bson:"comboCode"`
	ComboName string `json:"comboName" bson:"comboName"`
	Sort      int64  `json:"sort" bson:"sort"`
}

type RequestAddRouteIPWhitelist struct {
	OuCode    string `json:"ouCode" bson:"ouCode"`
	IPAddress string `json:"ipAddress" bson:"ipAddress" validate:"required"`
	Sort      int64  `json:"sort" bson:"sort"`
}

type RequestListPaymentMethod struct {
	MerchantKey string `json:"merchantkey"`
}

type ListPaymentMethod struct {
	OuId              int64  `json:"ouId"`
	OuCode            string `json:"ouCode"`
	PaymentMethodList []struct {
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
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
