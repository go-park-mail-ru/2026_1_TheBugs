package entity

import (
	"errors"
	"fmt"
)

var ServiceError error = errors.New("servicе error")

var NotFoundError error = errors.New("not found error")
var AlredyExitError error = errors.New("alredy exits error")
var BadCredentials error = errors.New("bad credential")
var InvalidInput error = errors.New("invalid data")
var OffsetOutOfRange error = errors.New("offset out of range")
var JWTError error = errors.New("jwt error")
var CollectPostersErr error = errors.New("collect posters error")
var ToManyRequest error = errors.New("to many requests")

type ValidationError struct {
	Err     error
	Field   string
	Details string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %v", v.Field, v.Err)
}

func (v *ValidationError) Unwrap() error {
	return v.Err
}

func NewValidationError(field string) *ValidationError {
	return &ValidationError{Err: InvalidInput, Field: field}
}
