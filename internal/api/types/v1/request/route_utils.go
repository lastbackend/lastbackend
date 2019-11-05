//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
)

type RouteRequest struct{}

func (RouteRequest) Manifest() *RouteManifest {
	return new(RouteManifest)
}

func (r *RouteManifest) Validate() *errors.Err {
	return nil
}

func (r *RouteManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("route").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return errors.New("route").IncorrectJSON(err)
	}

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (RouteRequest) RemoveOptions() *RouteRemoveOptions {
	return new(RouteRemoveOptions)
}

func (r *RouteRemoveOptions) Validate() *errors.Err {
	return nil
}
