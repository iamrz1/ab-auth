package infra

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// List of errors
var (
	ErrUnsupportedType = errors.New("infra: unsupported type")
	ErrNotFound        = mongo.ErrNoDocuments
	ErrDuplicateKey    = errors.New("infra: duplicate key")
	ErrInvalidData     = errors.New("infra: invalid data")
)
