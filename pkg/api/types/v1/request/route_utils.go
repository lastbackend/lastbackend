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
)

type RouteRequest struct{}

func (RouteRequest) CreateOptions() *RouteCreateOptions {
	return new(RouteCreateOptions)
}

func (s *RouteCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("route").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("route").IncorrectJSON(err)
	}

	//TODO: need checking arguments

	return nil
}

func (RouteRequest) UpdateOptions() *RouteUpdateOptions {
	return new(RouteUpdateOptions)
}

func (s *RouteUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("route").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("route").IncorrectJSON(err)
	}

	//TODO: need checking arguments

	return nil
}

func (RouteRequest) RemoveOptions() *RouteRemoveOptions {
	return new(RouteRemoveOptions)
}

func (s *RouteRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("route").IncorrectJSON(err)
	}

	return nil
}
