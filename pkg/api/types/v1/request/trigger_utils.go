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

type TriggerRequest struct{}

func (TriggerRequest) CreateOptions() *TriggerCreateOptions {
	return new(TriggerCreateOptions)
}

func (t *TriggerCreateOptions) Validate() *errors.Err {
	return nil
}

func (t *TriggerCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("trigger").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("trigger").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("trigger").IncorrectJSON(err)
	}

	return t.Validate()
}

func (t *TriggerCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (TriggerRequest) UpdateOptions() *TriggerUpdateOptions {
	return new(TriggerUpdateOptions)
}

func (t *TriggerUpdateOptions) Validate() *errors.Err {
	return nil
}

func (t *TriggerUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("trigger").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("trigger").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("trigger").IncorrectJSON(err)
	}

	return t.Validate()
}

func (t *TriggerUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (TriggerRequest) RemoveOptions() *TriggerRemoveOptions {
	return new(TriggerRemoveOptions)
}

func (t *TriggerRemoveOptions) Validate() *errors.Err {
	return nil
}

func (t *TriggerRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("trigger").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("trigger").IncorrectJSON(err)
	}

	return t.Validate()
}
