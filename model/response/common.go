package response

import "github.com/iamrz1/ab-auth/model"

// EmptySuccessRes example
type EmptySuccessRes struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"success message"`
	Timestamp string            `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.EmptyObject `json:"data"`
}

type TokenSuccessRes struct {
	Success   bool        `json:"success" example:"false"`
	Message   string      `json:"message" example:"success message"`
	Timestamp string      `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.Token `json:"data"`
}
