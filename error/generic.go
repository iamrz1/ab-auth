package http_error

// GenericHttpError represents Validation Error
type GenericHttpError struct {
	code    int
	message string
}

// NewGenericError returns a Generics error
func NewGenericError(code int, message string) GenericHttpError {
	return GenericHttpError{
		code:    code,
		message: message,
	}
}

func (ge GenericHttpError) Error() string {
	return ge.message
}

func (ge GenericHttpError) Code() int {
	return ge.code
}
