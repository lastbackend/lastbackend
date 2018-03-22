//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package request

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type RouteRequest struct{}

func (RouteRequest) CreateOptions() *RouteCreateOptions {
	return new(RouteCreateOptions)
}

func (r *RouteCreateOptions) Validate() *errors.Err {
	return nil
}

func (r *RouteCreateOptions) DecodeAndValidate(reader io.Reader) (*types.RouteCreateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("route").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, errors.New("route").IncorrectJSON(err)
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	opts := new(types.RouteCreateOptions)
	opts.Domain = r.Domain
	opts.Subdomain = r.Subdomain
	opts.Security = r.Security
	opts.Custom = r.Custom

	for _, rule := range r.Rules {
		opts.Rules = append(opts.Rules, types.RulesOption{
			Endpoint: rule.Endpoint,
			Path:     rule.Path,
			Port:     rule.Port,
		})
	}

	return opts, nil
}

func (r *RouteCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (RouteRequest) UpdateOptions() *RouteUpdateOptions {
	return new(RouteUpdateOptions)
}

func (r *RouteUpdateOptions) Validate() *errors.Err {
	return nil
}

func (r *RouteUpdateOptions) DecodeAndValidate(reader io.Reader) (*types.RouteUpdateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("route").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, errors.New("route").IncorrectJSON(err)
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	opts := new(types.RouteUpdateOptions)
	opts.Domain = r.Domain
	opts.Subdomain = r.Subdomain
	opts.Security = r.Security
	opts.Custom = r.Custom

	for _, rule := range r.Rules {
		opts.Rules = append(opts.Rules, types.RulesOption{
			Endpoint: rule.Endpoint,
			Path:     rule.Path,
			Port:     rule.Port,
		})
	}

	return opts, nil
}

func (r *RouteUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (RouteRequest) RemoveOptions() *RouteRemoveOptions {
	return new(RouteRemoveOptions)
}

func (r *RouteRemoveOptions) Validate() *errors.Err {
	return nil
}
