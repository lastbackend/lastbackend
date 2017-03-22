//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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

func IncorrectXML(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectXml,
		origin: getError("incorrect xml", e...),
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
		origin: getError(joinNameAndMessage(self.name, "access denied"), e...),
		http:   HTTP.getUnauthorized(),
	}
}

func (self *err) NotFound(e ...error) *Err {
	return &Err{
		Code:   http.StatusText(http.StatusNotFound),
		origin: getError(joinNameAndMessage(self.name, "not found"), e...),
		http:   HTTP.getNotFound(self.name),
	}
}

func (self *err) NotUnique(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusNotUnique,
		origin: getError(joinNameAndMessage(self.name, strings.ToLower(attr)+" not unique"), e...),
		http:   HTTP.getNotUnique(strings.ToLower(attr)),
	}
}

func (self *err) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError(joinNameAndMessage(self.name, "bad parameter"), e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (self *err) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError(joinNameAndMessage(self.name, "incorrect json"), e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (self *err) IncorrectXML(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError(joinNameAndMessage(self.name, "incorrect xml"), e...),
		http:   HTTP.getIncorrectXML(),
	}
}

func (self *err) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError(joinNameAndMessage(self.name, "unknow error"), e...),
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

func joinNameAndMessage(name, message string) string {
	return toUpperFirstChar(name) + ": " + message
}

func toUpperFirstChar(srt string) string {
	return strings.ToUpper(srt[0:1]) + srt[1:]
}
