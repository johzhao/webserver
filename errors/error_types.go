package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

type ErrorCode string

const (
	ErrUnknown       = ErrorCode("ERR-10000000")
	ErrRouteNotFound = ErrorCode("ERR-10000404")
	ErrorBadRequest  = ErrorCode("ERR-10000400")
)

func (e ErrorCode) New(msg string) error {
	return customError{
		code:          e,
		originalError: errors.New(msg),
	}
}

func (e ErrorCode) Newf(msg string, args ...interface{}) error {
	return customError{
		code:          e,
		originalError: fmt.Errorf(msg, args...),
	}
}

func (e ErrorCode) Wrap(err error, msg string) error {
	return e.Wrapf(err, msg)
}

func (e ErrorCode) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)
	return customError{
		code:          e,
		originalError: newErr,
		contextInfo:   nil,
	}
}
