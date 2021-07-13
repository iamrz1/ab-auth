package utils

import (
	"net/http"
	"strconv"
)

func GetPageLimit(r *http.Request) (int64, int64) {
	var page, limit = 1, 10
	pageQ := r.URL.Query().Get("page")
	limitQ := r.URL.Query().Get("limit")

	if len(pageQ) > 0 {
		n, err := strconv.Atoi(pageQ)
		if err == nil {
			page = n
		}
	}

	if len(limitQ) > 0 {
		n, err := strconv.Atoi(limitQ)
		if err == nil {
			limit = n
		}
	}

	return int64(page), int64(limit)
}
