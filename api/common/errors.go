package common

import (
	err "github.com/micro/go-micro/errors"
)

// CleanErrResponse to cancel circle Error
func CleanErrResponse(id string, e error, fn func(id, format string, a ...interface{}) error) error {
	if e == nil {
		return nil
	}
	E, ok := e.(*err.Error)
	if ok {
		E.Id = id
		return E
	}
	return fn(id, e.Error())
}
