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

func (Http) AccessDenied(w http.ResponseWriter) {
	Http{Code: 401, Status: StatusAccessDenied, Message: "Access denied"}.send(w)
}

func (Http) InvalidJSON(w http.ResponseWriter) {
	Http{Code: 400, Status: StatusIncorrectJson, Message: "Invalid josn"}.send(w)
}

func (Http) BadRequest(w http.ResponseWriter) {
	Http{Code: 400, Status: StatusBadRequest, Message: "Bad request"}.send(w)
}

func (Http) InternalServerError(w http.ResponseWriter) {
	Http{Code: 500, Status: StatusInternalServerError, Message: "internal server error"}.send(w)
}

func (Http) NotImplemented(w http.ResponseWriter) {
	Http{Code: 501, Status: StatusNotImplemented, Message: "not implemented"}.send(w)
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
		Code:    404,
		Status:  fmt.Sprintf("%s_NOT_FOUND", strings.ToUpper(name)),
		Message: fmt.Sprintf("%s not found", strings.ToLower(name)),
	}
}

func (Http) getBadParameter(name string) *Http {
	return &Http{
		Code:    406,
		Status:  fmt.Sprintf("BAD_PARAMETER_%s", strings.ToUpper(name)),
		Message: fmt.Sprintf("bad %s parameter", strings.ToLower(name)),
	}
}

func (Http) getNotUnique(name string) *Http {
	return &Http{
		Code:    400,
		Status:  fmt.Sprintf("%s_NOT_UNIQUE", strings.ToUpper(name)),
		Message: fmt.Sprintf("%s is already in use", strings.ToLower(name)),
	}
}

func (Http) getIncorrectJSON() *Http {
	return &Http{
		Code:    400,
		Status:  StatusIncorrectJson,
		Message: "incorrect json",
	}
}

func (Http) getAccessDenied() *Http {
	return &Http{
		Code:    401,
		Status:  StatusAccessDenied,
		Message: "Access denied",
	}
}

func (Http) getUnknown() *Http {
	return &Http{
		Code:    500,
		Status:  StatusInternalServerError,
		Message: "internal server error",
	}
}
