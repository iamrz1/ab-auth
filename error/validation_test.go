package http_error

import (
	"fmt"
	"testing"
)

const (
	ValidationMessage = "A message about bad requests"
)

func TestNewValidationError(t *testing.T) {
	possibleErrorMessage := "missing fields"
	err := NewValidationError(ValidationMessage, fmt.Errorf(possibleErrorMessage))
	if err.Error() != possibleErrorMessage || err.GetMessage() != ValidationMessage {
		t.Fail()
	}

	if err.ErrorMessage() != fmt.Sprintf("%s:%s", err.GetMessage(), err.Error()) {
		t.Fail()
	}
}
