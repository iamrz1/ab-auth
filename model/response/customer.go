package response

import "github.com/iamrz1/ab-auth/model"

type RegistrationSuccessRes struct {
}

// CustomerSuccessRes example
type CustomerSuccessRes struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"success message"`
	Data    CustomerData `json:"data"`
}

type CustomerData struct {
	Customer model.Customer `json:"object"`
}

// CustomerErrorRes example
type CustomerErrorRes struct {
	Success bool            `json:"success" example:"false"`
	Message string          `json:"message" example:"failure message"`
	Data    CustomerErrData `json:"data"`
}

type CustomerErrData struct {
	Customer model.EmptyObject `json:"object"`
}

// CustomerListSuccessRes example
type CustomerListSuccessRes struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"success message"`
	Data    CustomerListData `json:"data"`
}

type CustomerListData struct {
	Customers []model.Customer `json:"objects"`
	Count     int64            `json:"count"`
}

// CustomerListErrorRes example
type CustomerListErrorRes struct {
	Success bool                  `json:"success" example:"false"`
	Message string                `json:"message" example:"failure message"`
	Data    CustomerListErrorData `json:"data"`
}

type CustomerListErrorData struct {
	Customer []model.EmptyObject `json:"objects"`
	Count    int64               `json:"count"`
}
