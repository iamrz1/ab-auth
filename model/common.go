package model

import "gopkg.in/validator.v2"

type ListOptions struct {
	Page  int64       `json:"page" bson:"page,omitempty"`
	Limit int64       `json:"limit" bson:"limit,omitempty"`
	Sort  interface{} `json:"sort" bson:"sort,omitempty"`
}

type EmptyObject struct{}

func Validate(s interface{}) error {
	return validator.Validate(s)
}

type LoginReq struct {
	Username string `json:"username" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordReq struct {
	Username     string `json:"username" validate:"nonzero"`
	CaptchaID    string `json:"captcha_id" validate:"nonzero"`
	CaptchaValue string `json:"captcha_value" validate:"nonzero"`
}
