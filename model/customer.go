package model

import (
	"encoding/json"
	"github.com/iamrz1/ab-auth/utils"
	"strings"
	"time"
)

type Customer struct {
	Username            string    `json:"username,omitempty" bson:"username,omitempty"`
	FullName            string    `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Password            string    `json:"-" bson:"password,omitempty"`
	RecoveryPhoneNumber string    `json:"recovery_phone_number,omitempty" bson:"recovery_phone_number,omitempty"`
	Gender              string    `json:"gender,omitempty" bson:"gender,omitempty"`
	Email               string    `json:"email,omitempty" bson:"email,omitempty"`
	Occupation          string    `json:"occupation,omitempty" bson:"occupation,omitempty"`
	Organization        string    `json:"organization,omitempty" bson:"organization,omitempty"`
	BirthDate           time.Time `json:"-" bson:"birth_date,omitempty"`
	BirthDateString     string    `json:"birth_date,omitempty" bson:"-"`
	Status              string    `json:"status,omitempty" bson:"status,omitempty"`
	IsVerified          *bool     `json:"is_verified,omitempty" bson:"is_verified,omitempty"`
	ProfilePicURL       string    `json:"profile_pic_url,omitempty" bson:"profile_pic_url,omitempty"`
	IsDeleted           *bool     `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	LastResetAt         time.Time `json:"-" bson:"last_reset_at,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (d *Customer) ToResponse() *Customer {
	if !d.BirthDate.IsZero() {
		d.BirthDateString = d.BirthDate.Format(utils.ISOLayout)
	}

	d.Gender = strings.Title(d.Gender)

	return d
}

func (d *Customer) ToShortResponse() *CustomerShort {
	cs := CustomerShort{}

	b, err := json.Marshal(d)
	if err != nil {
		return &cs
	}

	err = json.Unmarshal(b, &cs)
	if err != nil {
		return &cs
	}
	return &cs
}

type CustomerSignupReq struct {
	Username     string `json:"username" validate:"nonzero"`
	FullName     string `json:"full_name" validate:"nonzero"`
	Password     string `json:"password" validate:"nonzero"`
	CaptchaID    string `json:"captcha_id" validate:"nonzero"`
	CaptchaValue string `json:"captcha_value" validate:"nonzero"`
}

type CustomerSignupVerificationReq struct {
	Username string `json:"username" validate:"nonzero"`
	OTP      string `json:"otp" validate:"nonzero"`
}

type CustomerListReq struct {
	Page   int64
	Limit  int64
	Sort   string
	Order  string
	Search string
}

type CustomerProfileUpdateReq struct {
	Username      string `json:"-" validate:"nonzero"` //username will come from either token or url param
	FullName      string `json:"full_name,omitempty"`
	Gender        string `json:"gender,omitempty" example:"male/female/other"`
	Email         string `json:"email,omitempty"`
	Occupation    string `json:"occupation,omitempty"`
	Organization  string `json:"organization,omitempty"`
	BirthDate     string `json:"birth_date,omitempty" example:"2006-01-02T15:04:05.000Z"`
	ProfilePicURL string `json:"profile_pic_url,omitempty" bson:"profile_pic_url,omitempty"`
	//RecoveryPhoneNumber string `json:"recovery_phone_number,omitempty"`
	IsVerified *bool `json:"-"`
	IsDeleted  *bool `json:"-"`
}

type CustomerDeleteReq struct {
	Username string `json:"-" validate:"nonzero"`
}

type CustomerShort struct {
	Username string `json:"username,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Status   string `json:"status,omitempty"`
}
