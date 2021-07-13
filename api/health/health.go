package health

import (
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

func status(w http.ResponseWriter, r *http.Request) {
	resp := &utils.Response{
		Code:    "Success",
		Success: true,
		Data:    nil,
		Errors:  nil,
		Message: "OK",
	}

	resp.ServeJSON(w)
	return
}
