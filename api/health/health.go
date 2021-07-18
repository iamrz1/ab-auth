package health

import (
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

func status(w http.ResponseWriter, r *http.Request) {
	utils.ServeJSONObject(w, http.StatusOK, "OK", nil, nil, true)
	return
}
