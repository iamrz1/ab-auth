package response

import "github.com/iamrz1/ab-auth/model"

// AddressSuccessRes example
type AddressSuccessRes struct {
	Success   bool          `json:"success" example:"true"`
	Status    string        `json:"status" example:"OK"`
	Message   string        `json:"message" example:"success message"`
	Timestamp string        `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.Address `json:"data"`
}

// AddressListSuccessRes example
type AddressListSuccessRes struct {
	Success   bool            `json:"success" example:"true"`
	Status    string          `json:"status" example:"OK"`
	Message   string          `json:"message" example:"success message"`
	Timestamp string          `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []model.Address `json:"data"`
}
