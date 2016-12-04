package errors

import (
	"errors"
	"net/http"
	"strings"
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

func (self *Err) Err() error {
	return self.origin
}

func (self *Err) Http(w http.ResponseWriter) {
	self.http.send(w)
}

type err struct {
	name string
}

func New(name string) *err {
	return &err{strings.ToLower(name)}
}

func (self *err) AccessDenied(e ...error) *Err {
	return &Err{
		Code:   StatusAccessDenied,
		origin: getError(toUpperFirstChar(self.name)+"access denied", e...),
		http:   HTTP.getAccessDenied(),
	}
}

func (self *err) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError(toUpperFirstChar(self.name)+": not found", e...),
		http:   HTTP.getNotFound(self.name),
	}
}

func (self *err) NotUnique(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusNotUnique,
		origin: getError(toUpperFirstChar(self.name)+":"+strings.ToLower(attr)+" not unique", e...),
		http:   HTTP.getNotUnique(strings.ToLower(attr)),
	}
}

func (self *err) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError(toUpperFirstChar(self.name)+": bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (self *err) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError(toUpperFirstChar(self.name)+": incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (self *err) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError(toUpperFirstChar(self.name)+": unknow error", e...),
		http:   HTTP.getUnknown(),
	}
}

func getError(msg string, e ...error) error {
	if len(e) == 0 {
		return errors.New(msg)
	} else {
		return e[0]
	}
}

func toUpperFirstChar(srt string) string {
	return strings.ToUpper(srt[0:1]) + srt[1:]
}
