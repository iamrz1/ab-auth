package model

import (
	"encoding/json"
	"github.com/iamrz1/ab-auth/utils"
)

type Address struct {
	ID              string  `json:"id,omitempty" bson:"_id,omitempty"`
	Username        string  `json:"username,omitempty" bson:"username,omitempty"`
	PhoneNumber     string  `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	FullName        string  `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Division        string  `json:"division,omitempty" bson:"division,omitempty"`
	District        string  `json:"district,omitempty" bson:"district,omitempty"`
	SubDistrict     string  `json:"sub_district,omitempty" bson:"sub_district,omitempty"`
	Union           string  `json:"union,omitempty" bson:"union,omitempty"`
	DivisionSlug    string  `json:"division_slug,omitempty" bson:"division_slug,omitempty"`
	DistrictSlug    string  `json:"district_slug,omitempty" bson:"district_slug,omitempty"`
	SubDistrictSlug string  `json:"sub_district_slug,omitempty" bson:"sub_district_slug,omitempty"`
	UnionSlug       string  `json:"union_slug,omitempty" bson:"union_slug,omitempty"`
	Address         string  `json:"address,omitempty" bson:"address,omitempty"`
	IsPrimary       *bool   `json:"is_primary,omitempty" bson:"is_primary,omitempty"`
	Longitude       float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	IsDeleted       *bool   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

type AddressCreateReq struct {
	Username        string  `json:"-" validate:"nonzero"`
	PhoneNumber     string  `json:"phone_number,omitempty" validate:"nonzero"`
	FullName        string  `json:"full_name,omitempty" validate:"nonzero"`
	Division        string  `json:"division,omitempty" validate:"nonzero"`
	District        string  `json:"district,omitempty"`
	SubDistrict     string  `json:"sub_district,omitempty"`
	DivisionSlug    string  `json:"division_slug,omitempty" validate:"nonzero"`
	DistrictSlug    string  `json:"district_slug,omitempty" validate:"nonzero"`
	SubDistrictSlug string  `json:"sub_district_slug,omitempty"`
	UnionSlug       string  `json:"union_slug,omitempty"`
	Union           string  `json:"union,omitempty"`
	Address         string  `json:"address,omitempty"`
	IsPrimary       *bool   `json:"is_primary,omitempty"`
	Longitude       float64 `json:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty"`
}

func (acr *AddressCreateReq) ToAddress() *Address {
	res := &Address{}
	b, _ := json.Marshal(acr)
	json.Unmarshal(b, res)
	res.Username = acr.Username
	res.IsDeleted = utils.BoolP(false)

	return res
}

type AddressUpdateReq struct {
	ID              string  `json:"id,omitempty" validate:"nonzero"`
	Username        string  `json:"-" validate:"nonzero"`
	PhoneNumber     string  `json:"phone_number,omitempty"`
	FullName        string  `json:"full_name,omitempty"`
	Division        string  `json:"division,omitempty"`
	District        string  `json:"district,omitempty"`
	SubDistrict     string  `json:"sub_district,omitempty"`
	DivisionSlug    string  `json:"division_slug,omitempty"`
	DistrictSlug    string  `json:"district_slug,omitempty"`
	SubDistrictSlug string  `json:"sub_district_slug,omitempty"`
	UnionSlug       string  `json:"union_slug,omitempty"`
	Union           string  `json:"union,omitempty"`
	Address         string  `json:"address,omitempty"`
	Longitude       float64 `json:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty"`
}

func (acr *AddressUpdateReq) ToAddress() *Address {
	res := &Address{}
	b, _ := json.Marshal(acr)
	json.Unmarshal(b, res)
	res.Username = acr.Username
	res.IsDeleted = utils.BoolP(false)

	return res
}

// BD location preset models

type BDLocationType string

const (
	LocationTypeDivision    BDLocationType = "division"
	LocationTypeDistrict    BDLocationType = "district"
	LocationTypeSubDistrict BDLocationType = "sub_district"
	LocationTypeUnion       BDLocationType = "union"
)

type BDLocation struct {
	ID     int32          `json:"id" bson:"id"`
	Name   string         `json:"name" bson:"name"`
	NameBn string         `json:"name_bn" bson:"name_bn"`
	Slug   string         `json:"slug" bson:"slug"`
	Parent string         `json:"parent" json:"parent"`
	Type   BDLocationType `json:"type" bson:"type"`
}

type BDLocationReq struct {
	Parent string `json:"parent" bson:"parent"`
	Page   int64  `json:"-" bson:"-"`
	Limit  int64  `json:"-" bson:"-"`
}
