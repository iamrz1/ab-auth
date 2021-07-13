package http_error

import (
	"fmt"
)

// ValidationError represents Validation Error
type ValidationError struct {
	message string
	err     error
}

// NewValidationError returns a validation error
func NewValidationError(message string, err error) ValidationError {
	return ValidationError{
		message: message,
		err:     err,
	}
}

func (ve ValidationError) Error() string {
	if ve.err != nil {
		return fmt.Sprintf("%s", ve.err)
	}

	return ve.message
}

func (ve ValidationError) ErrorMessage() string {
	if ve.err != nil {
		if ve.message != "" {
			return fmt.Sprintf("%s:%v", ve.message, ve.err)
		}
		return ve.err.Error()
	}

	return ve.message
}

// GetMessage returns error message
func (ve ValidationError) GetMessage() string {
	return ve.message
}

// GetError returns wrapped error
func (ve ValidationError) GetError() error {
	return ve.err
}
