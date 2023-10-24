package models

type DeviceDB struct {
	DeviceID    string `json:"deviceID" bson:"deviceId"`
	MerchantKey string `json:"merchantKey" bson:"merchantKey"`
}
