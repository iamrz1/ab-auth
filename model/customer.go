package model

import (
	"github.com/iamrz1/ab-auth/utils"
	"time"
)

type Customer struct {
	Username            string    `json:"username,omitempty" bson:"username,omitempty"`
	FullName            string    `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Password            string    `json:"-" bson:"password,omitempty"`
	RecoveryPhoneNumber string    `json:"recovery_phone_number,omitempty" bson:"recovery_phone_number,omitempty"`
	Gender              string    `json:"gender,omitempty" bson:"gender,omitempty"`
	BirthDate           time.Time `json:"-" bson:"birth_date,omitempty"`
	BirthDateString     string    `json:"birth_date,omitempty" bson:"-"`
	Status              string    `json:"status,omitempty" bson:"status,omitempty"`
	IsVerified          *bool     `json:"is_verified,omitempty"`
	IsDeleted           *bool     `json:"is_deleted,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (d *Customer) ToCustomerResponse() *Customer {
	if !d.BirthDate.IsZero() {
		d.BirthDateString = d.BirthDate.Format(utils.ISOLayout)
	}

	// todo: process d if needed
	return d
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

type CustomerUpdateReq struct {
	Username            string `json:"-" validate:"nonzero"` //username will come from either token or url param
	FullName            string `json:"full_name,omitempty"`
	Gender              string `json:"gender,omitempty"`
	BirthDate           string `json:"birth_date,omitempty"`
	RecoveryPhoneNumber string `json:"recovery_phone_number,omitempty"`
	IsVerified          *bool  `json:"is_verified,omitempty"`
	IsDeleted           *bool  `json:"is_deleted,omitempty"`
}

type CustomerDeleteReq struct {
	Username string `json:"-" validate:"nonzero"`
}
