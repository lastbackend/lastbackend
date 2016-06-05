package errors

import (
	"errors"
	"strings"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

func InternalServerError() Error {
	return StatusError{500, errors.New("INTERNAL_SERVER_ERROR")}
}

func InvalidIncomingJSON() Error {
	return StatusError{500, errors.New("INVALID_INCOMING_JSON")}
}

func ParamInvalid(param string) Error {
	return StatusError{400, errors.New("PARAM_INVALID_" + strings.ToUpper(param))}
}

func ParamNotUnique(param string) Error {
	return StatusError{400, errors.New("PARAM_NOT_UNIQUE_" + strings.ToUpper(param))}
}

func Custom(code int, param string) Error {
	return StatusError{code, errors.New(param)}
}
