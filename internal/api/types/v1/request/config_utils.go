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
// patents in process, and are protected by trade config or copyright law.
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

type ConfigRequest struct{}

func (ConfigRequest) Manifest() *ConfigManifest {
	return new(ConfigManifest)
}

func (v *ConfigManifest) Validate() *errors.Err {
	return nil
}

func (v *ConfigManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("config").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("config").Unknown(err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return errors.New("config").IncorrectJSON(err)
	}

	return v.Validate()
}

func (ConfigRequest) RemoveOptions() *ConfigRemoveOptions {
	return new(ConfigRemoveOptions)
}

func (v *ConfigRemoveOptions) Validate() *errors.Err {
	return nil
}
