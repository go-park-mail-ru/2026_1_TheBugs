package repository

import "errors"

var NotFoundRecord error = errors.New("not found record")
var AlreadyExist error = errors.New("record alredy exists")
