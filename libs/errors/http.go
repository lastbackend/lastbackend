package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var HTTP Http

type Http struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (Http) Unauthorized(w http.ResponseWriter) {
	Http{Code: http.StatusUnauthorized, Status: http.StatusText(http.StatusUnauthorized), Message: "access denied"}.send(w)
}

func (Http) InvalidJSON(w http.ResponseWriter) {
	Http{Code: http.StatusBadRequest, Status: StatusIncorrectJson, Message: "invalid json"}.send(w)
}

func (Http) BadRequest(w http.ResponseWriter) {
	Http{Code: http.StatusBadRequest, Status: http.StatusText(http.StatusBadRequest), Message: "bad request"}.send(w)
}

func (Http) NotFound(w http.ResponseWriter) {
	Http{Code: http.StatusNotFound, Status: http.StatusText(http.StatusNotFound), Message: "not found"}.send(w)
}

func (Http) InternalServerError(w http.ResponseWriter) {
	Http{Code: http.StatusInternalServerError, Status: http.StatusText(http.StatusInternalServerError), Message: "internal server error"}.send(w)
}

func (Http) NotImplemented(w http.ResponseWriter) {
	Http{Code: http.StatusNotImplemented, Status: http.StatusText(http.StatusNotImplemented), Message: "not implemented"}.send(w)
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
		Message: fmt.Sprintf("%s not found", strings.ToLower(name)),
	}
}

func (Http) getBadParameter(name string) *Http {
	return &Http{
		Code:    http.StatusNotAcceptable,
		Status:  StatusBadParameter,
		Message: fmt.Sprintf("bad %s parameter", strings.ToLower(name)),
	}
}

func (Http) getNotUnique(name string) *Http {
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusNotUnique,
		Message: fmt.Sprintf("%s is already in use", strings.ToLower(name)),
	}
}

func (Http) getIncorrectJSON() *Http {
	return &Http{
		Code:    http.StatusBadRequest,
		Status:  StatusIncorrectJson,
		Message: "incorrect json",
	}
}

func (Http) getUnauthorized() *Http {
	return &Http{
		Code:    http.StatusUnauthorized,
		Status:  http.StatusText(http.StatusUnauthorized),
		Message: "access denied",
	}
}

func (Http) getUnknown() *Http {
	return &Http{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: "internal server error",
	}
}
