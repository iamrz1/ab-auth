package response

import "github.com/iamrz1/ab-auth/model"

type RegistrationSuccessRes struct {
}

// CustomerSuccessRes example
type CustomerSuccessRes struct {
	Success   bool           `json:"success" example:"true"`
	Message   string         `json:"message" example:"success message"`
	Timestamp string         `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.Customer `json:"data"`
}

// CustomerErrorRes example
type CustomerErrorRes struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"failure message"`
	Timestamp string            `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.EmptyObject `json:"data"`
}

// CustomerListSuccessRes example
type CustomerListSuccessRes struct {
	Success   bool             `json:"success" example:"true"`
	Message   string           `json:"message" example:"success message"`
	Timestamp string           `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []model.Customer `json:"data"`
}

// CustomerListErrorRes example
type CustomerListErrorRes struct {
	Success   bool                `json:"success" example:"false"`
	Message   string              `json:"message" example:"failure message"`
	Timestamp string              `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []model.EmptyObject `json:"data"`
}

type CustomerResShort struct {
	Success   bool                `json:"success" example:"false"`
	Message   string              `json:"message" example:"failure message"`
	Timestamp string              `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.CustomerShort `json:"data"`
}
