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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type EndpointRequest struct{}

func (EndpointRequest) CreateOptions() *EndpointCreateOptions {
	return new(EndpointCreateOptions)
}

func (t *EndpointCreateOptions) Validate() *errors.Err {
	return nil
}

func (t *EndpointCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("endpoint").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("endpoint").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("endpoint").IncorrectJSON(err)
	}

	return t.Validate()
}

func (t *EndpointCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (EndpointRequest) UpdateOptions() *EndpointUpdateOptions {
	return new(EndpointUpdateOptions)
}

func (t *EndpointUpdateOptions) Validate() *errors.Err {
	return nil
}

func (t *EndpointUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("endpoint").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("endpoint").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("endpoint").IncorrectJSON(err)
	}

	return t.Validate()
}

func (t *EndpointUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (EndpointRequest) RemoveOptions() *EndpointRemoveOptions {
	return new(EndpointRemoveOptions)
}

func (t *EndpointRemoveOptions) Validate() *errors.Err {
	return nil
}
