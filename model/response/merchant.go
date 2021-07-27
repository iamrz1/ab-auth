package response

import "github.com/iamrz1/ab-auth/model"

// MerchantSuccessRes example
type MerchantSuccessRes struct {
	Success   bool           `json:"success" example:"true"`
	Status    string         `json:"status" example:"OK"`
	Message   string         `json:"message" example:"success message"`
	Timestamp string         `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.Merchant `json:"data"`
}

// MerchantListSuccessRes example
type MerchantListSuccessRes struct {
	Success   bool             `json:"success" example:"true"`
	Status    string           `json:"status" example:"OK"`
	Message   string           `json:"message" example:"success message"`
	Timestamp string           `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []model.Merchant `json:"data"`
	ListMeta  ListMeta         `json:"meta"`
}

type MerchantResShort struct {
	Success   bool                `json:"success" example:"true"`
	Status    string              `json:"status" example:"OK"`
	Message   string              `json:"message" example:"failure message"`
	Timestamp string              `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.MerchantShort `json:"data"`
}
