package http_error

import (
	"net/http"
	"testing"
)

const (
	Code    = http.StatusUnauthorized
	Message = "A message about unauthorized access"
)

func TestNewGenericError(t *testing.T) {
	err := NewGenericError(Code, Message)
	if err.Code() != Code || err.Error() != Message {
		t.Fail()
	}
}
