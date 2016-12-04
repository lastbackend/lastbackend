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

	StatusForbidden = "FORBIDDEN"

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

func BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError(attr+": bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
