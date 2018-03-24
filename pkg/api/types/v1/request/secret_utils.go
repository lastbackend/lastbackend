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

type SecretRequest struct{}

func (SecretRequest) CreateOptions() *SecretCreateOptions {
	return new(SecretCreateOptions)
}

func (s *SecretCreateOptions) Validate() *errors.Err {
	switch true {
	case s.Data == nil:
		return errors.New("secret").BadParameter("data")
	}
	return nil
}

func (s *SecretCreateOptions) DecodeAndValidate(reader io.Reader) (*types.SecretCreateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("secret").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("secret").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return nil, errors.New("secret").IncorrectJSON(err)
	}

	opts := new(types.SecretCreateOptions)
	opts.Data = s.Data

	return opts, nil
}

func (s *SecretCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (SecretRequest) UpdateOptions() *SecretUpdateOptions {
	return new(SecretUpdateOptions)
}

func (s *SecretUpdateOptions) Validate() *errors.Err {
	switch true {
	case s.Data == nil:
		return errors.New("secret").BadParameter("data")
	}
	return nil
}

func (s *SecretUpdateOptions) DecodeAndValidate(reader io.Reader) (*types.SecretUpdateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("secret").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("secret").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return nil, errors.New("secret").IncorrectJSON(err)
	}

	if s.Data == nil {

	}

	opts := new(types.SecretUpdateOptions)
	opts.Data = s.Data

	return opts, nil
}

func (s *SecretUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (SecretRequest) RemoveOptions() *SecretRemoveOptions {
	return new(SecretRemoveOptions)
}

func (s *SecretRemoveOptions) Validate() *errors.Err {
	return nil
}
