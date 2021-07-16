package response

import "github.com/iamrz1/ab-auth/model"

// EmptySuccessRes example
type EmptySuccessRes struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"success message"`
	Data    model.EmptyObject `json:"data"`
}

type TokenSuccessRes struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"success message"`
	Data    model.Token `json:"data"`
}
