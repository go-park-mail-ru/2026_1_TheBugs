package entity

import "errors"

var ServiceError error = errors.New("servicе error")

var NotFoundError error = errors.New("not found error")
var AlredyExitError error = errors.New("alredy exits error")
var BadCredentials error = errors.New("bad credential")
var InvalidInput error = errors.New("invalid data")
