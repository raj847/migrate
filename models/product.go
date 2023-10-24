package models

type PolicyOuProductWithRules struct {
	OuId                  int64               `json:"ouId" bson:"ouId"`
	OuCode                string              `json:"ouCode" bson:"ouCode"`
	OuName                string              `json:"ouName" bson:"ouName"`
	ProductId             int64               `json:"productId" bson:"productId"`
	ProductCode           string              `json:"productCode" bson:"productCode"`
	ProductName           string              `json:"productName" bson:"productName"`
	ServiceFee            float64             `json:"serviceFee" bson:"serviceFee"`
	IsPctServiceFee       string              `json:"isPctServiceFee" bson:"isPctServiceFee"`
	IsPctServiceFeeMember string              `json:"isPctServiceFeeMember" bson:"isPctServiceFeeMember"`
	ServiceFeeMember      float64             `json:"serviceFeeMember" bson:"serviceFeeMember"`
	ProductOuWithRules    *ProductOuWithRules `json:"productRules" bson:"productRules"`
}

type ProductOuWithRules struct {
	OuId             int64   `json:"ouId" bson:"ouId"`
	ProductId        int64   `json:"productId" bson:"productId"`
	Price            float64 `json:"price" bson:"price"`
	BaseTime         int64   `json:"baseTime" bson:"baseTime"`
	ProgressiveTime  int64   `json:"progressiveTime" bson:"progressiveTime"`
	ProgressivePrice float64 `json:"progressivePrice" bson:"progressivePrice"`
	IsPct            string  `json:"isPct" bson:"isPct"`
	ProgressivePct   int64   `json:"progressivePct" bson:"progressivePct"`
	MaxPrice         float64 `json:"maxPrice" bson:"maxPrice"`
	Is24H            string  `json:"is24h" bson:"is24h"`
	OvernightTime    string  `json:"overnightTime" bson:"overnightTime"`
	OvernightPrice   float64 `json:"overnightPrice" bson:"overnightPrice"`
	GracePeriod      int64   `json:"gracePeriod" bson:"gracePeriod"`
	FlgRepeat        string  `json:"flgRepeat" bson:"flgRepeat"`
}

type RequestGetPolicyOuProductList struct {
	Keyword         string `json:"keyword" bson:"keyword"`
	AscDesc         string `json:"ascDesc" bson:"ascDesc"`
	ColumnOrderName string `json:"columnOrderName" bson:"columnOrderName"`
}

type ResponseGetPolicyOuProductList struct {
	CountProductList int64                      `json:"countProductList" bson:"countProductList"`
	Data             []PolicyOuProductWithRules `json:"data" bson:"data"`
}

type ResponseTrxProduct struct {
	ProductCode string `json:"productCode" bson:"productCode"`
	ProductName string `json:"productName" bson:"productName"`
}

type PolicyOuProductDepositCounterWithRules struct {
	OuId                             int64                            `json:"ouId" bson:"ouId,omitempty"`
	OuCode                           string                           `json:"ouCode" bson:"ouCode,omitempty"`
	OuName                           string                           `json:"ouName" bson:"ouName,omitempty"`
	ProductId                        int64                            `json:"productId" bson:"productId,omitempty"`
	ProductCode                      string                           `json:"productCode" bson:"productCode,omitempty"`
	ProductName                      string                           `json:"productName" bson:"productName,omitempty"`
	ServiceFee                       float64                          `json:"serviceFee" bson:"serviceFee,omitempty"`
	IsPctServiceFee                  string                           `json:"isPctServiceFee" bson:"isPctServiceFee,omitempty"`
	ServiceFeeMember                 float64                          `json:"serviceFeeMember" bson:"serviceFeeMember,omitempty"`
	IsPctServiceFeeMember            string                           `json:"isPctServiceFeeMember" bson:"isPctServiceFeeMember,omitempty"`
	ProductOuDepositCounterWithRules ProductOuDepositCounterWithRules `json:"productDepositCounterRules" bson:"productDepositCounterRules,omitempty"`
}

type ProductOuDepositCounterWithRules struct {
	OuId      int64   `json:"ouId" bson:"ouId,omitempty"`
	ProductId int64   `json:"productId" bson:"productId,omitempty"`
	Price     float64 `json:"price" bson:"price,omitempty"`
	IsPct     string  `json:"isPct" bson:"isPct"`
}
