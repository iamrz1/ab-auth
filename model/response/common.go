package response

import "github.com/iamrz1/ab-auth/model"

// EmptySuccessRes example
type EmptySuccessRes struct {
	Success bool             `json:"success" example:"false"`
	Message string           `json:"message" example:"success message"`
	Data    emptySuccessData `json:"data"`
}

type emptySuccessData struct {
	Customer model.EmptyObject `json:"object"`
}

type TokenSuccessRes struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"success message"`
	Data    token  `json:"data"`
}

type token struct {
	TokenObject model.Token `json:"object"`
}
