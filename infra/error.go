package infra

import (
	"errors"
	"fmt"
)

// List of errors
var (
	ErrUnsupportedType = errors.New("infra: unsupported type")
	ErrNotFound        = fmt.Errorf("not found")
	ErrDuplicateKey    = errors.New("infra: duplicate key")
	ErrInvalidData     = errors.New("infra: invalid data")
)
