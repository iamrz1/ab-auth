package response

import "github.com/iamrz1/ab-auth/model"

type RegistrationSuccessRes struct {
}

// CustomerSuccessRes example
type CustomerSuccessRes struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"success message"`
	Data    customerData `json:"data"`
}

type customerData struct {
	Customer model.Customer `json:"object"`
}

// CustomerErrorRes example
type CustomerErrorRes struct {
	Success bool            `json:"success" example:"false"`
	Message string          `json:"message" example:"failure message"`
	Data    customerErrData `json:"data"`
}

type customerErrData struct {
	Customer model.EmptyObject `json:"object"`
}

// CustomerListSuccessRes example
type CustomerListSuccessRes struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"success message"`
	Data    customerListData `json:"data"`
}

type customerListData struct {
	Customers []model.Customer `json:"objects"`
	Count     int64            `json:"count"`
}

// CustomerListErrorRes example
type CustomerListErrorRes struct {
	Success bool                  `json:"success" example:"false"`
	Message string                `json:"message" example:"failure message"`
	Data    customerListErrorData `json:"data"`
}

type customerListErrorData struct {
	Customer []model.EmptyObject `json:"objects"`
	Count    int64               `json:"count"`
}

type CustomerResShort struct {
	Success bool          `json:"success" example:"false"`
	Message string        `json:"message" example:"failure message"`
	Data    customerShort `json:"data"`
}

type customerShort struct {
	Customer model.CustomerShort `json:"object"`
}
