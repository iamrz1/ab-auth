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
