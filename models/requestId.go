package models

type RequestID struct {
	ID int `json:"id" validate:"required"`
}

type RequestConfirmA2P struct {
	Merchantkey     string `json:"merchantkey"`
	MerchantNoRef   string `json:"merchantNoRef"`
	PaymentCategory string `json:"paymentCategory"`
}
