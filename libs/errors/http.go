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

func (Http) Unauthorized(w http.ResponseWriter) {
	Http{Code: http.StatusUnauthorized, Status: http.StatusText(http.StatusUnauthorized), Message: "Access denied"}.send(w)
}

func (Http) InvalidJSON(w http.ResponseWriter) {
	Http{Code: http.StatusBadRequest, Status: StatusIncorrectJson, Message: "Invalid json"}.send(w)
}

func (Http) BadRequest(w http.ResponseWriter) {
	Http{Code: http.StatusBadRequest, Status: http.StatusText(http.StatusBadRequest), Message: "Bad request"}.send(w)
}

func (Http) NotFound(w http.ResponseWriter) {
	Http{Code: http.StatusNotFound, Status: http.StatusText(http.StatusNotFound), Message: "Not found"}.send(w)
}

func (Http) InternalServerError(w http.ResponseWriter) {
	Http{Code: http.StatusInternalServerError, Status: http.StatusText(http.StatusInternalServerError), Message: "Internal server error"}.send(w)
}

func (Http) NotImplemented(w http.ResponseWriter) {
	Http{Code: http.StatusNotImplemented, Status: http.StatusText(http.StatusNotImplemented), Message: "Not implemented"}.send(w)
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

func (Http) getNotFound(name string) *Http {
	return &Http{
		Code:    http.StatusNotFound,
		Status:  http.StatusText(http.StatusNotFound),
		Message: fmt.Sprintf("%s not found", toUpperFirstChar(name)),
	}
}

func (Http) getBadParameter(name string) *Http {
	return &Http{
		Code:    http.StatusNotAcceptable,
		Status:  StatusBadParameter,
		Message: fmt.Sprintf("Bad %s parameter", name),
	}
}

func (Http) getNotUnique(name string) *Http {
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusNotUnique,
		Message: fmt.Sprintf("%s is already in use", toUpperFirstChar(name)),
	}
}

func (Http) getIncorrectJSON() *Http {
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusIncorrectJson,
		Message: "Incorrect json",
	}
}

func (Http) getUnauthorized() *Http {
	return &Http{
		Code:    http.StatusUnauthorized,
		Status:  http.StatusText(http.StatusUnauthorized),
		Message: "Access denied",
	}
}

func (Http) getUnknown() *Http {
	return &Http{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: "Internal server error",
	}
}
