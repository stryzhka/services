package word

import "errors"

var (
	ErrWordNotFound      = errors.New("word not found")
	ErrWordAlreadyExists = errors.New("word already exists")
	ErrValidation        = errors.New("validation error")
)
