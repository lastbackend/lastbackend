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
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type request struct {
	method string
	path   string
	header http.Header
	host   RegistryHost
	body   func() (io.ReadCloser, error)
	size   int64
}

func (r *request) do(ctx context.Context) (*http.Response, error) {
	u := r.host.GetScheme() + "://" + r.host.GetHost() + r.path
	req, err := http.NewRequest(r.method, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	if r.body != nil {
		body, err := r.body()
		if err != nil {
			return nil, err
		}
		req.Body = body
		req.GetBody = r.body
		if r.size > 0 {
			req.ContentLength = r.size
		}
	}

	if err := r.authorize(ctx, req); err != nil {
		return nil, errors.Wrap(err, "failed to authorize")
	}

	resp, err := r.host.client.Do(req.WithContext(ctx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return nil, errors.Wrap(err, "failed to do request")
	}

	return resp, nil
}

func (r *request) authorize(ctx context.Context, req *http.Request) error {
	// Check if has header for host
	if r.host.authorizer != nil {
		if err := r.host.authorizer.Authorize(ctx, req); err != nil {
			return err
		}
	}

	return nil
}

func (r *request) doWithRetries(ctx context.Context, responses []*http.Response) (*http.Response, error) {
	resp, err := r.do(ctx)
	if err != nil {
		return nil, err
	}

	responses = append(responses, resp)
	retry, err := r.retryRequest(ctx, responses)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	if retry {
		resp.Body.Close()
		return r.doWithRetries(ctx, responses)
	}
	return resp, err
}

func (r *request) retryRequest(ctx context.Context, responses []*http.Response) (bool, error) {
	if len(responses) > 5 {
		return false, nil
	}
	last := responses[len(responses)-1]
	switch last.StatusCode {
	case http.StatusUnauthorized:
		if r.host.authorizer != nil {
			if err := r.host.authorizer.AddResponses(ctx, responses); err == nil {
				return true, nil
			} else if !IsNotImplemented(err) {
				return false, err
			}
		}

		return false, nil
	case http.StatusMethodNotAllowed:
		// Support hosts which have not properly implemented the HEAD method for
		// manifests endpoint
		if r.method == http.MethodHead && strings.Contains(r.path, "/manifests/") {
			r.method = http.MethodGet
			return true, nil
		}
	case http.StatusRequestTimeout, http.StatusTooManyRequests:
		return true, nil
	}

	// TODO: Handle 50x errors accounting for attempt history
	return false, nil
}

func (r *request) String() string {
	return r.host.GetScheme() + "://" + r.host.GetHost() + r.path
}

// handleErrorResponse returns error parsed from HTTP response for an
// unsuccessful HTTP response code (in the range 400 - 499 inclusive). An
// UnexpectedHTTPStatusError returned for response code outside of expected
// range.
func handleErrorResponse(resp *http.Response) error {
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		for _, c := range responseChallenges(resp) {
			if c.Scheme == bearerAuth {
				var err Error
				switch c.Parameters["error"] {
				case "invalid_token":
					err.Code = ErrorCodeUnauthorized
				case "insufficient_scope":
					err.Code = ErrorCodeDenied
				default:
					continue
				}
				if description := c.Parameters["error_description"]; description != "" {
					err.Message = description
				} else {
					err.Message = err.Code.Message()
				}

				return mergeErrors(err, parseHTTPErrorResponse(resp.StatusCode, resp.Body))
			}
		}
		err := parseHTTPErrorResponse(resp.StatusCode, resp.Body)
		if uErr, ok := err.(*UnexpectedHTTPResponseError); ok && resp.StatusCode == http.StatusUnauthorized {
			return ErrorCodeUnauthorized.WithDetail(uErr.Response)
		}
		return err
	}
	return &UnexpectedHTTPStatusError{Status: resp.Status}
}

func makeErrorList(err error) []error {
	if errs, ok := err.(Errors); ok {
		return errs
	}
	return []error{err}
}

func mergeErrors(err1, err2 error) error {
	return Errors(append(makeErrorList(err1), makeErrorList(err2)...))
}

func parseHTTPErrorResponse(statusCode int, r io.Reader) error {
	var errs Errors
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var detailsErr struct {
		Details string `json:"details"`
	}
	err = json.Unmarshal(body, &detailsErr)
	if err == nil && detailsErr.Details != "" {
		switch statusCode {
		case http.StatusUnauthorized:
			return ErrorCodeUnauthorized.WithMessage(detailsErr.Details)
		case http.StatusTooManyRequests:
			return ErrorCodeTooManyRequests.WithMessage(detailsErr.Details)
		default:
			return ErrorCodeUnknown.WithMessage(detailsErr.Details)
		}
	}

	if err := json.Unmarshal(body, &errs); err != nil {
		return errors.Wrapf(err, "error parsing HTTP %d response body: %q", statusCode, string(body))
	}

	if len(errs) == 0 {
		return errors.Wrapf(ErrNoErrorsInBody, "error parsing HTTP %d response body: %q", statusCode, string(body))
	}

	return errs
}

// responseChallenges returns a list of authorization challenges
// for the given http Response. Challenges are only checked if
// the response status code was a 401.
func responseChallenges(resp *http.Response) []Challenge {
	if resp.StatusCode == http.StatusUnauthorized {
		// Parse the WWW-Authenticate Header and store the challenges
		// on this endpoint object.
		return parseAuthHeader(resp.Header)
	}

	return nil
}
