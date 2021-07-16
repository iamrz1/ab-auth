package response

import "github.com/iamrz1/ab-auth/model"

type RegistrationSuccessRes struct {
}

// CustomerSuccessRes example
type CustomerSuccessRes struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"success message"`
	Data    model.Customer `json:"data"`
}

// CustomerErrorRes example
type CustomerErrorRes struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"failure message"`
	Data    model.EmptyObject `json:"data"`
}

// CustomerListSuccessRes example
type CustomerListSuccessRes struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"success message"`
	Data    []model.Customer `json:"data"`
}

// CustomerListErrorRes example
type CustomerListErrorRes struct {
	Success bool                `json:"success" example:"false"`
	Message string              `json:"message" example:"failure message"`
	Data    []model.EmptyObject `json:"data"`
}

type CustomerResShort struct {
	Success bool                `json:"success" example:"false"`
	Message string              `json:"message" example:"failure message"`
	Data    model.CustomerShort `json:"data"`
}
