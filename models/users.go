package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Active    string             `json:"active"`
	CreatedBy string             `json:"createdBy"`
	CreatedAt string             `json:"createdAt"`
	UpdatedBy string             `json:"updatedBy"`
	UpdatedAt string             `json:"updatedAt"`
}

type RequestAuthLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	DeviceId string `json:"deviceId"`
}

type UserLogin struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Username        string  `json:"username"`
	Pin             int64   `json:"pin"`
	TypeUser        string  `json:"type"`
	Email           string  `json:"email"`
	RolesId         int64   `json:"rolesId"`
	Password        string  `json:"password"`
	Active          string  `json:"active"`
	IsAdmin         string  `json:"isAdmin"`
	IsInternal      string  `json:"isInternal"`
	RolesName       *string `json:"rolesName"`
	OuDefaultId     int64   `json:"ouDefaultId"`
	OuCode          string  `json:"ouCode"`
	OuName          string  `json:"ouName"`
	PolicyDefaultId int64   `json:"policyDefaultId"`
}

type UserLoginLocal struct {
	User               UserSync     `json:"user"`
	OuList             []OuSync     `json:"ouList"`
	MerchantKeyParking string       `json:"merchantKeyParking"`
	TaskList           string       `json:"taskList"`
	DeviceList         []DeviceSync `json:"deviceList"`
}

type OuSync struct {
	Id         int64  `json:"id"`
	OuCode     string `json:"ouCode"`
	OuName     string `json:"ouName"`
	OuType     string `json:"ouType"`
	OuParentId int64  `json:"ouParentId"`
}

type DeviceSync struct {
	DeviceId       string `json:"deviceId"`
	FlgProgressive string `json:"flgProgressive"`
	MerchantKey    string `json:"merchantKey"`
}

type UserSync struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TypeUser        string `json:"type"`
	Email           string `json:"email"`
	RolesId         int64  `json:"rolesId"`
	RolesName       string `json:"rolesName"`
	IsAdmin         string `json:"isAdmin"`
	IsInternal      string `json:"isInternal"`
	Active          string `json:"active"`
	OuList          []Ou   `json:"ouList"`
	TaskList        string `json:"taskList"`
	OuDefaultId     int64  `json:"ouDefaultId"`
	PolicyDefaultId int64  `json:"policyDefaultId"`
	OuCode          string `json:"ouCode"`
	OuName          string `json:"ouName"`
	PinUser         int64  `json:"pinUser"`
}

type RequestAuthDashboard struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Type     string `json:"type"`
	DeviceId string `json:"deviceId"`
}

type UserData struct {
	Id              int64     `json:"id"`
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Pin             int64     `json:"pin"`
	TypeUser        string    `json:"type"`
	Email           string    `json:"email"`
	RolesId         int64     `json:"rolesId"`
	Password        string    `json:"password"`
	IsAdmin         string    `json:"isAdmin"`
	IsInternal      string    `json:"isInternal"`
	Active          string    `json:"active"`
	RolesName       string    `json:"rolesName"`
	OuList          *[]OuSync `json:"ouList"`
	TaskList        *string   `json:"taskList"`
	OuDefaultId     int64     `json:"ouDefaultId"`
	PolicyDefaultId int64     `json:"policyDefaultId"`
	OuCode          string    `json:"ouCode"`
	OuName          string    `json:"ouName"`
}

type ResponseAuth struct {
	Token           string   `json:"token"`
	RefreshToken    string   `json:"refreshToken"`
	User            UserSync `json:"user"`
	MKey            string   `json:"mKey"`
	MerchantKey     string   `json:"merchantKey"`
	AdditionalInfo  *OuInfo  `json:"additionalInfo"`
	FlagProgressive string   `json:"flagProgressive"`
}

type OuInfo struct {
	OuId        int64  `json:"ouId" bson:"ouId"`
	ImageLogo   string `json:"imageLogo" bson:"imageLogo"`
	ImageFooter string `json:"imageFooter" bson:"imageFooter"`
	Desc1       string `json:"desc1" bson:"desc1"`
	Desc2       string `json:"desc2" bson:"desc2"`
	Desc3       string `json:"desc3" bson:"desc3"`
}
