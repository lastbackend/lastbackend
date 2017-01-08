package errors

import (
	"errors"
	"net/http"
	"strings"
)

const (
	StatusBadParameter  = "Bad Parameter"
	StatusUnknown       = "Unknow"
	StatusIncorrectXml  = "Incorrect XML"
	StatusIncorrectJson = "Incorrect json"
	StatusNotUnique     = "Not Unique"
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

func (self *err) Unauthorized(e ...error) *Err {
	return &Err{
		Code:   http.StatusText(http.StatusUnauthorized),
		origin: getError(toUpperFirstChar(self.name)+"access denied", e...),
		http:   HTTP.getUnauthorized(),
	}
}

func (self *err) NotFound(e ...error) *Err {
	return &Err{
		Code:   http.StatusText(http.StatusNotFound),
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

func (self *err) NotImplemented(e ...error) *Err {
	return &Err{
		Code:   http.StatusText(http.StatusNotImplemented),
		origin: getError("not implemented", e...),
		http:   HTTP.getUnknown(),
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
