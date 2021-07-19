package response

import "github.com/iamrz1/ab-auth/model"

// EmptyErrorRes example
type EmptyErrorRes struct {
	Success   bool              `json:"success" example:"false"`
	Status    string            `json:"status" example:"Status string corresponding to the error"`
	Message   string            `json:"message" example:"failure message"`
	Timestamp string            `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.EmptyObject `json:"data"`
}

// EmptyListErrorRes example
type EmptyListErrorRes struct {
	Success   bool                `json:"success" example:"false"`
	Status    string              `json:"status" example:"Status string corresponding to the error"`
	Message   string              `json:"message" example:"failure message"`
	Timestamp string              `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []model.EmptyObject `json:"data"`
}

// EmptySuccessRes example
type EmptySuccessRes struct {
	Success   bool              `json:"success" example:"false"`
	Status    string            `json:"status" example:"OK"`
	Message   string            `json:"message" example:"success message"`
	Timestamp string            `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.EmptyObject `json:"data"`
}

type TokenSuccessRes struct {
	Success   bool        `json:"success" example:"false"`
	Status    string      `json:"status" example:"OK"`
	Message   string      `json:"message" example:"success message"`
	Timestamp string      `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      model.Token `json:"data"`
}

type ListMeta struct {
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Limit int `json:"Limit"`
	Count int `json:"count"`
}
