package utils

import (
	goErrors "errors"

	"github.com/micro/go-micro/v2/errors"
)

// MicroErrorDetail of errors.Error - can't be nil
func MicroErrorDetail(e error) string {
	err, ok := e.(*errors.Error)
	if ok {
		return err.Detail
	}
	return e.Error()
}

// MicroError of errors.Error
func MicroError(e error) error {
	if e == nil {
		return nil
	}
	err, ok := e.(*errors.Error)
	if ok {
		return goErrors.New(err.Detail)
	}
	return e
}
