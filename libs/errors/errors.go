package errors

import (
	"errors"
	"net/http"
)

const (
	StatusIncorrectXml      = "INCORRECT_XML"
	StatusIncorrectJson     = "INCORRECT_JSON"
	StatusIncorrectName     = "INCORRECT_NAME"
	StatusIncorrectEmail    = "INCORRECT_EMAIL"
	StatusIncorrectUsename  = "INCORRECT_USERNAME"
	StatusIncorrectPassword = "INCORRECT_PASSWORD"
	StatusIncorrectAuth     = "INCORRECT_AUTH"
	StatusIncorrectPayload  = "INCORRECT_PAYLOAD"

	StatusBadRequest   = "BAD_REQUEST"
	StatusBadGateway   = "BAD_GATEWAY"
	StatusBadParameter = "BAD_PARAMETER"

	StatusNotFound       = "NOT_FOUND"
	StatusNotUnique      = "NOT_UNIQUE"
	StatusNotSupported   = "NOT_SUPPORTED"
	StatusNotAcceptable  = "NOT_ACCEPTABLE"
	StatusNotImplemented = "NOT_IMPLEMENTED"

	StatusPaymentRequired = "PAYMENT_REQUIRED"
	StatusAccessDenied    = "ACCESS_DENIED"

	StatusForbidden        = "FORBIDDEN"
	StatusMethodNotAllowed = "METHOD_NOT_ALLOWED"

	StatusInternalServerError = "INTERNAL_SERVER_ERROR"

	StatusUnknown = "UNKNOWN"
)

type Err struct {
	Code   string
	Attr   string
	origin error
	http   *Http
}

func (e *Err) Err() error {
	return e.origin
}

func (e *Err) Http(w http.ResponseWriter) {
	e.http.send(w)
}

func getError(msg string, e ...error) error {
	if len(e) == 0 {
		return errors.New(msg)
	} else {
		return e[0]
	}
}
