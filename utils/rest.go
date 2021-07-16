package utils

import (
	"encoding/json"
	rest_error "github.com/iamrz1/ab-auth/error"
	"log"
	"net/http"
)

type Meta struct {
	Count    int64 `json:"count"`
	PageSize int64 `json:"page_size"`
}

// Response serializer util
type Response struct {
	Code    string      `json:"code,omitempty"`
	Status  int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
	Meta    interface{} `json:"meta,omitempty"`
	Data    interface{} `json:"data"`
	User    interface{} `json:"user,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ServeJSON serves json to http client
func (r *Response) ServeJSON(w http.ResponseWriter) error {
	resp := &Response{
		Code:    r.Code,
		Status:  r.Status,
		Message: r.Message,
		Data:    r.Data,
		Errors:  r.Errors,
		Success: r.Success,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func ServeJSONObject(w http.ResponseWriter, code string, status int, message string, data interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	var obj interface{}
	type EmptyObject struct{}
	if data == nil {
		obj = map[string]interface{}{
			"object": EmptyObject{},
		}
	} else {
		obj = map[string]interface{}{
			"object": data,
		}
	}

	resp = &Response{
		Code:    code,
		Message: message,
		Success: success,
		Data:    obj,
		Meta:    meta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func ServeJSONList(w http.ResponseWriter, code string, status int, message string, list interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	var objs interface{}
	type EmptyObject struct{}
	if list == nil {
		objs = []EmptyObject{}
	} else {
		objs = list
	}

	resp = &Response{
		Code:    code,
		Message: message,
		Success: success,
		Data:    objs,
		Meta:    meta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func HandleObjectError(w http.ResponseWriter, err error) {
	log.Println("service error: ", err)
	switch v := err.(type) {
	case rest_error.ValidationError:
		ServeJSONObject(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, v.ErrorMessage(), nil, nil, false)
		return
	case rest_error.GenericHttpError:
		ServeJSONObject(w, http.StatusText(v.Code()), v.Code(), v.Error(), nil, nil, false)
		return
	default:
		ServeJSONObject(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "Something went wrong", nil, nil, false)
		return
	}
}

func HandleListError(w http.ResponseWriter, err error) {
	log.Println("service error: ", err)
	switch v := err.(type) {
	case rest_error.ValidationError:
		ServeJSONList(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, v.ErrorMessage(), nil, nil, false)
		return
	case rest_error.GenericHttpError:
		ServeJSONList(w, http.StatusText(v.Code()), v.Code(), v.Error(), nil, nil, false)
		return
	default:
		ServeJSONList(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "Something went wrong", nil, nil, false)
		return
	}
}
