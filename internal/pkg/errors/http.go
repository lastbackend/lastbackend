//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"encoding/json"
	"fmt"
	"net/http"
)

var HTTP Http

type Http struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (Http) Unauthorized(w http.ResponseWriter, msg ...string) {
	HTTP.getUnauthorized(msg...).send(w)
}

func (Http) Forbidden(w http.ResponseWriter, msg ...string) {
	HTTP.getForbidden(msg...).send(w)
}

func (Http) NotAllowed(w http.ResponseWriter, msg ...string) {
	HTTP.getNotAllowed(msg...).send(w)
}

func (Http) BadRequest(w http.ResponseWriter, msg ...string) {
	HTTP.getBadRequest(msg...).send(w)
}

func (Http) NotFound(w http.ResponseWriter, args ...string) {
	HTTP.getNotFound(args...).send(w)
}

func (Http) InternalServerError(w http.ResponseWriter, msg ...string) {
	HTTP.getInternalServerError(msg...).send(w)
}

func (Http) BadGateway(w http.ResponseWriter) {
	HTTP.getBadGateway().send(w)
}

func (Http) PaymentRequired(w http.ResponseWriter, msg ...string) {
	HTTP.getPaymentRequired(msg...).send(w)
}

func (Http) NotImplemented(w http.ResponseWriter, msg ...string) {
	HTTP.getPaymentRequired(msg...).send(w)
}

func (Http) BadParameter(w http.ResponseWriter, args ...string) {
	HTTP.getBadParameter(args...).send(w)
}

func (Http) InvalidJSON(w http.ResponseWriter, msg ...string) {
	HTTP.getIncorrectJSON(msg...).send(w)
}

func (Http) InvalidXML(w http.ResponseWriter, msg ...string) {
	HTTP.getIncorrectXML(msg...).send(w)
}

func (h Http) send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(h.Code)
	response, _ := json.Marshal(h)
	w.Write(response)
}

// ===================================================================================================================
// ============================================= INTERNAL HELPER METHODS =============================================
// ===================================================================================================================

func (Http) getUnauthorized(msg ...string) *Http {
	return getHttpError(http.StatusUnauthorized, msg...)
}

func (Http) getForbidden(msg ...string) *Http {
	return getHttpError(http.StatusForbidden, msg...)
}

func (Http) getNotAllowed(msg ...string) *Http {
	return getHttpError(http.StatusMethodNotAllowed, msg...)
}

func (Http) getPaymentRequired(msg ...string) *Http {
	return getHttpError(http.StatusPaymentRequired, msg...)
}

func (Http) getUnknown(msg ...string) *Http {
	return getHttpError(http.StatusInternalServerError, msg...)
}

func (Http) getInternalServerError(msg ...string) *Http {
	return getHttpError(http.StatusInternalServerError, msg...)
}

func (Http) getBadGateway() *Http {
	return getHttpError(http.StatusBadGateway)
}

func (Http) getNotImplemented(msg ...string) *Http {
	return getHttpError(http.StatusNotImplemented, msg...)
}

func (Http) getBadRequest(msg ...string) *Http {
	return getHttpError(http.StatusBadRequest, msg...)
}

func (Http) getNotFound(args ...string) *Http {
	message := "Not Found"
	for i, a := range args {
		switch i {
		case 0:
			message = fmt.Sprintf("%s not found", toUpperFirstChar(a))
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    http.StatusNotFound,
		Status:  http.StatusText(http.StatusNotFound),
		Message: message,
	}
}

func (Http) getNotUnique(name string) *Http {
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusNotUnique,
		Message: fmt.Sprintf("%s is already in use", toUpperFirstChar(name)),
	}
}

func (Http) getIncorrectJSON(msg ...string) *Http {
	message := "Incorrect json"
	for i, m := range msg {
		switch i {
		case 0:
			message = m
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusIncorrectJson,
		Message: message,
	}
}

func (Http) getIncorrectXML(msg ...string) *Http {
	message := "Incorrect json"
	for i, m := range msg {
		switch i {
		case 0:
			message = m
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusIncorrectXml,
		Message: message,
	}
}

func (Http) getAllocatedParameter(args ...string) *Http {
	message := "Value is in use"
	for i, a := range args {
		switch i {
		case 0:
			message = fmt.Sprintf("%s is already in use", toUpperFirstChar(a))
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusBadParameter,
		Message: message,
	}
}

func (Http) getBadParameter(args ...string) *Http {
	message := "Bad parameter"
	for i, a := range args {
		switch i {
		case 0:
			message = fmt.Sprintf("Bad %s parameter", a)
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusBadParameter,
		Message: message,
	}
}

func getHttpError(code int, msg ...string) *Http {
	status := http.StatusText(code)
	message := http.StatusText(code)

	for i, m := range msg {
		switch i {
		case 0:
			message = m
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return &Http{
		Code:    code,
		Status:  status,
		Message: message,
	}
}
