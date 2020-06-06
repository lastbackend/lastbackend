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

package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrUnknown                = errors.New("unknown error")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrNotFound               = errors.New("not found")
	ErrAlreadyExists          = errors.New("already exists")
	ErrFailedPrecondition     = errors.New("failed precondition")
	ErrUnavailable            = errors.New("unavailable")
	ErrNotImplemented         = errors.New("not implemented")
	ErrTooManyRequests        = errors.New("too many requests")
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrNoErrorsInBody         = errors.New("no error details found in HTTP response body")
	ErrInvalidAuthorization   = errors.New("authorization failed")
	ErrInvalid                = errors.New("invalid reference")
	ErrObjectRequired         = errors.New("object required")
	ErrHostnameRequired       = errors.New("hostname required")
	ErrNoToken                = errors.New("authorization server did not include a token in the response")
)

var (
	errorCodeToDescriptors = map[ErrorCode]ErrorDescriptor{}
	idToDescriptors        = map[string]ErrorDescriptor{}
	groupToDescriptors     = map[string][]ErrorDescriptor{}
)

var (
	// ErrorCodeUnknown is a generic error that can be used as a last
	// resort if there is no situation-specific error message that can be used
	ErrorCodeUnknown = Register("errcode", ErrorDescriptor{
		Value:   "UNKNOWN",
		Message: "unknown error",
		Description: `Generic error returned when the error does not have an
			                                            API classification.`,
		HTTPStatusCode: http.StatusInternalServerError,
	})

	// ErrorCodeUnsupported is returned when an operation is not supported.
	ErrorCodeUnsupported = Register("errcode", ErrorDescriptor{
		Value:   "UNSUPPORTED",
		Message: "The operation is unsupported.",
		Description: `The operation was unsupported due to a missing
		implementation or invalid set of parameters.`,
		HTTPStatusCode: http.StatusMethodNotAllowed,
	})

	// ErrorCodeUnauthorized is returned if a request requires
	// authentication.
	ErrorCodeUnauthorized = Register("errcode", ErrorDescriptor{
		Value:   "UNAUTHORIZED",
		Message: "authentication required",
		Description: `The access controller was unable to authenticate
		the client. Often this will be accompanied by a
		Www-Authenticate HTTP response header indicating how to
		authenticate.`,
		HTTPStatusCode: http.StatusUnauthorized,
	})

	// ErrorCodeDenied is returned if a client does not have sufficient
	// permission to perform an action.
	ErrorCodeDenied = Register("errcode", ErrorDescriptor{
		Value:   "DENIED",
		Message: "requested access to the resource is denied",
		Description: `The access controller denied access for the
		operation on a resource.`,
		HTTPStatusCode: http.StatusForbidden,
	})

	// ErrorCodeUnavailable provides a common error to report unavailability
	// of a service or endpoint.
	ErrorCodeUnavailable = Register("errcode", ErrorDescriptor{
		Value:          "UNAVAILABLE",
		Message:        "service unavailable",
		Description:    "Returned when a service is not available",
		HTTPStatusCode: http.StatusServiceUnavailable,
	})

	// ErrorCodeTooManyRequests is returned if a client attempts too many
	// times to contact a service endpoint.
	ErrorCodeTooManyRequests = Register("errcode", ErrorDescriptor{
		Value:   "TOOMANYREQUESTS",
		Message: "too many requests",
		Description: `Returned when a client attempts to contact a
		service too many times`,
		HTTPStatusCode: http.StatusTooManyRequests,
	})
)

var nextCode = 1000
var registerLock sync.Mutex

// Register will make the passed-in error known to the environment and
// return a new ErrorCode
func Register(group string, descriptor ErrorDescriptor) ErrorCode {
	registerLock.Lock()
	defer registerLock.Unlock()

	descriptor.Code = ErrorCode(nextCode)

	if _, ok := idToDescriptors[descriptor.Value]; ok {
		panic(fmt.Sprintf("ErrorValue %q is already registered", descriptor.Value))
	}
	if _, ok := errorCodeToDescriptors[descriptor.Code]; ok {
		panic(fmt.Sprintf("ErrorCode %v is already registered", descriptor.Code))
	}

	groupToDescriptors[group] = append(groupToDescriptors[group], descriptor)
	errorCodeToDescriptors[descriptor.Code] = descriptor
	idToDescriptors[descriptor.Value] = descriptor

	nextCode++
	return descriptor.Code
}

type byValue []ErrorDescriptor

func (a byValue) Len() int           { return len(a) }
func (a byValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

// GetGroupNames returns the list of Error group names that are registered
func GetGroupNames() []string {
	keys := []string{}

	for k := range groupToDescriptors {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GetErrorCodeGroup returns the named group of error descriptors
func GetErrorCodeGroup(name string) []ErrorDescriptor {
	desc := groupToDescriptors[name]
	sort.Sort(byValue(desc))
	return desc
}

// GetErrorAllDescriptors returns a slice of all ErrorDescriptors that are
// registered, irrespective of what group they're in
func GetErrorAllDescriptors() []ErrorDescriptor {
	result := []ErrorDescriptor{}

	for _, group := range GetGroupNames() {
		result = append(result, GetErrorCodeGroup(group)...)
	}
	sort.Sort(byValue(result))
	return result
}

// UnexpectedHTTPStatusError is returned when an unexpected HTTP status is
// returned when making a RegistryHost api call.
type UnexpectedHTTPStatusError struct {
	Status string
}

func (e *UnexpectedHTTPStatusError) Error() string {
	return fmt.Sprintf("received unexpected HTTP status: %s", e.Status)
}

// UnexpectedHTTPResponseError is returned when an expected HTTP status code
// is returned, but the content was unexpected and failed to be parsed.
type UnexpectedHTTPResponseError struct {
	ParseErr   error
	StatusCode int
	Response   []byte
}

func (e *UnexpectedHTTPResponseError) Error() string {
	return fmt.Sprintf("error parsing HTTP %d response body: %s: %q", e.StatusCode, e.ParseErr.Error(), string(e.Response))
}

// IsInvalidArgument returns true if the error is due to an invalid argument
func IsInvalidArgument(err error) bool {
	return errors.Cause(err) == ErrInvalidArgument
}

// IsNotFound returns true if the error is due to a missing object
func IsNotFound(err error) bool {
	return errors.Cause(err) == ErrNotFound
}

// IsAlreadyExists returns true if the error is due to an already existing
// metadata item
func IsAlreadyExists(err error) bool {
	return errors.Cause(err) == ErrAlreadyExists
}

// IsFailedPrecondition returns true if an operation could not proceed to the
// lack of a particular condition
func IsFailedPrecondition(err error) bool {
	return errors.Cause(err) == ErrFailedPrecondition
}

// IsUnavailable returns true if the error is due to a resource being unavailable
func IsUnavailable(err error) bool {
	return errors.Cause(err) == ErrUnavailable
}

// IsNotImplemented returns true if the error is due to not being implemented
func IsNotImplemented(err error) bool {
	return errors.Cause(err) == ErrNotImplemented
}

// IsCanceled returns true if the error is due to `context.Canceled`.
func IsCanceled(err error) bool {
	return errors.Cause(err) == context.Canceled
}

// IsDeadlineExceeded returns true if the error is due to
// `context.DeadlineExceeded`.
func IsDeadlineExceeded(err error) bool {
	return errors.Cause(err) == context.DeadlineExceeded
}

// ErrorCoder is the base interface for ErrorCode and Error allowing
// users of each to just call ErrorCode to get the real ID of each
type ErrorCoder interface {
	ErrorCode() ErrorCode
}

// ErrorCode represents the error type. The errors are serialized via strings
// and the integer format may change and should *never* be exported.
type ErrorCode int

var _ error = ErrorCode(0)

// ErrorCode just returns itself
func (ec ErrorCode) ErrorCode() ErrorCode {
	return ec
}

// Error returns the ID/Value
func (ec ErrorCode) Error() string {
	// NOTE(stevvooe): Cannot use message here since it may have unpopulated args.
	return strings.ToLower(strings.Replace(ec.String(), "_", " ", -1))
}

// Descriptor returns the descriptor for the error code.
func (ec ErrorCode) Descriptor() ErrorDescriptor {
	d, ok := errorCodeToDescriptors[ec]

	if !ok {
		return ErrorCodeUnknown.Descriptor()
	}

	return d
}

// String returns the canonical identifier for this error code.
func (ec ErrorCode) String() string {
	return ec.Descriptor().Value
}

// Message returned the human-readable error message for this error code.
func (ec ErrorCode) Message() string {
	return ec.Descriptor().Message
}

// MarshalText encodes the receiver into UTF-8-encoded text and returns the
// result.
func (ec ErrorCode) MarshalText() (text []byte, err error) {
	return []byte(ec.String()), nil
}

// UnmarshalText decodes the form generated by MarshalText.
func (ec *ErrorCode) UnmarshalText(text []byte) error {
	desc, ok := idToDescriptors[string(text)]

	if !ok {
		desc = ErrorCodeUnknown.Descriptor()
	}

	*ec = desc.Code

	return nil
}

// WithMessage creates a new Error struct based on the passed-in info and
// overrides the Message property.
func (ec ErrorCode) WithMessage(message string) Error {
	return Error{
		Code:    ec,
		Message: message,
	}
}

// WithDetail creates a new Error struct based on the passed-in info and
// set the Detail property appropriately
func (ec ErrorCode) WithDetail(detail interface{}) Error {
	return Error{
		Code:    ec,
		Message: ec.Message(),
	}.WithDetail(detail)
}

// WithArgs creates a new Error struct and sets the Args slice
func (ec ErrorCode) WithArgs(args ...interface{}) Error {
	return Error{
		Code:    ec,
		Message: ec.Message(),
	}.WithArgs(args...)
}

// Error provides a wrapper around ErrorCode with extra Details provided.
type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

var _ error = Error{}

// ErrorCode returns the ID/Value of this Error
func (e Error) ErrorCode() ErrorCode {
	return e.Code
}

// Error returns a human readable representation of the error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code.Error(), e.Message)
}

// WithDetail will return a new Error, based on the current one, but with
// some Detail info added
func (e Error) WithDetail(detail interface{}) Error {
	return Error{
		Code:    e.Code,
		Message: e.Message,
		Detail:  detail,
	}
}

// WithArgs uses the passed-in list of interface{} as the substitution
// variables in the Error's Message string, but returns a new Error
func (e Error) WithArgs(args ...interface{}) Error {
	return Error{
		Code:    e.Code,
		Message: fmt.Sprintf(e.Code.Message(), args...),
		Detail:  e.Detail,
	}
}

// ErrorDescriptor provides relevant information about a given error code.
type ErrorDescriptor struct {
	// Code is the error code that this descriptor describes.
	Code ErrorCode

	// Value provides a unique, string key, often captilized with
	// underscores, to identify the error code. This value is used as the
	// keyed value when serializing api errors.
	Value string

	// Message is a short, human readable decription of the error condition
	// included in API responses.
	Message string

	// Description provides a complete account of the errors purpose, suitable
	// for use in documentation.
	Description string

	// HTTPStatusCode provides the http status code that is associated with
	// this error condition.
	HTTPStatusCode int
}

// ParseErrorCode returns the value by the string error code.
// `ErrorCodeUnknown` will be returned if the error is not known.
func ParseErrorCode(value string) ErrorCode {
	ed, ok := idToDescriptors[value]
	if ok {
		return ed.Code
	}

	return ErrorCodeUnknown
}

// Errors provides the envelope for multiple errors and a few sugar methods
// for use within the application.
type Errors []error

var _ error = Errors{}

func (errs Errors) Error() string {
	switch len(errs) {
	case 0:
		return "<nil>"
	case 1:
		return errs[0].Error()
	default:
		msg := "errors:\n"
		for _, err := range errs {
			msg += err.Error() + "\n"
		}
		return msg
	}
}

// Len returns the current number of errors.
func (errs Errors) Len() int {
	return len(errs)
}

// MarshalJSON converts slice of error, ErrorCode or Error into a
// slice of Error - then serializes
func (errs Errors) MarshalJSON() ([]byte, error) {
	var tmpErrs struct {
		Errors []Error `json:"errors,omitempty"`
	}

	for _, daErr := range errs {
		var err Error

		switch daErr.(type) {
		case ErrorCode:
			err = daErr.(ErrorCode).WithDetail(nil)
		case Error:
			err = daErr.(Error)
		default:
			err = ErrorCodeUnknown.WithDetail(daErr)

		}

		// If the Error struct was setup and they forgot to set the
		// Message field (meaning its "") then grab it from the ErrCode
		msg := err.Message
		if msg == "" {
			msg = err.Code.Message()
		}

		tmpErrs.Errors = append(tmpErrs.Errors, Error{
			Code:    err.Code,
			Message: msg,
			Detail:  err.Detail,
		})
	}

	return json.Marshal(tmpErrs)
}

// UnmarshalJSON deserializes []Error and then converts it into slice of
// Error or ErrorCode
func (errs *Errors) UnmarshalJSON(data []byte) error {
	var tmpErrs struct {
		Errors []Error
	}

	if err := json.Unmarshal(data, &tmpErrs); err != nil {
		return err
	}

	var newErrs Errors
	for _, daErr := range tmpErrs.Errors {
		// If Message is empty or exactly matches the Code's message string
		// then just use the Code, no need for a full Error struct
		if daErr.Detail == nil && (daErr.Message == "" || daErr.Message == daErr.Code.Message()) {
			// Error's w/o details get converted to ErrorCode
			newErrs = append(newErrs, daErr.Code)
		} else {
			// Error's w/ details are untouched
			newErrs = append(newErrs, Error{
				Code:    daErr.Code,
				Message: daErr.Message,
				Detail:  daErr.Detail,
			})
		}
	}

	*errs = newErrs
	return nil
}
