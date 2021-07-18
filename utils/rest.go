package utils

import (
	"encoding/json"
	rest_error "github.com/iamrz1/ab-auth/error"
	"log"
	"net/http"
	"os"
	"time"
)

type Meta struct {
	Count    int64 `json:"count"`
	PageSize int64 `json:"page_size"`
}

// Response serializer util
type Response struct {
	Status    string      `json:"status,omitempty"`
	Message   string      `json:"message,omitempty"`
	Success   bool        `json:"success"`
	Meta      interface{} `json:"meta,omitempty"`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

func ServeJSONObject(w http.ResponseWriter, code int, message string, data interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	type EmptyObject struct{}
	if data == nil {
		data = EmptyObject{}
	}

	resp = &Response{
		Status:    http.StatusText(code),
		Message:   message,
		Success:   success,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().Format(ISOLayout),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func ServeJSONList(w http.ResponseWriter, code int, message string, list interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	type EmptyObject struct{}
	if list == nil {
		list = []EmptyObject{}
	}

	resp = &Response{
		Status:    http.StatusText(code),
		Message:   message,
		Success:   success,
		Data:      list,
		Meta:      meta,
		Timestamp: time.Now().Format(ISOLayout),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func HandleObjectError(w http.ResponseWriter, err error) {
	log.Println("service error: ", err)

	errMeta := map[string]string{}
	if os.Getenv("ENV") != "prod" {
		errMeta["error"] = err.Error()
	}

	switch v := err.(type) {
	case rest_error.ValidationError:
		ServeJSONObject(w, http.StatusBadRequest, v.ErrorMessage(), nil, errMeta, false)
		return
	case rest_error.GenericHttpError:
		ServeJSONObject(w, v.Code(), v.Error(), nil, errMeta, false)
		return
	default:
		ServeJSONObject(w, http.StatusInternalServerError, "Something went wrong", nil, errMeta, false)
		return
	}
}

func HandleListError(w http.ResponseWriter, err error) {
	log.Println("service error: ", err)

	errMeta := map[string]string{}
	if os.Getenv("ENV") != "prod" {
		errMeta["error"] = err.Error()
	}

	switch v := err.(type) {
	case rest_error.ValidationError:
		ServeJSONList(w, http.StatusBadRequest, v.ErrorMessage(), nil, errMeta, false)
		return
	case rest_error.GenericHttpError:
		ServeJSONList(w, v.Code(), v.Error(), nil, errMeta, false)
		return
	default:
		ServeJSONList(w, http.StatusInternalServerError, "Something went wrong", nil, errMeta, false)
		return
	}
}
