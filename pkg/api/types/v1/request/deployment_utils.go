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

type DeploymentRequest struct{}

func (DeploymentRequest) UpdateOptions() *DeploymentUpdateOptions {
	return new(DeploymentUpdateOptions)
}

func (s *DeploymentUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *DeploymentUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("deployment").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("deployment").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("deployment").IncorrectJSON(err)
	}

	switch true {
	case s.Replicas == nil:
		return errors.New("deployment").BadParameter("replicas")
	case *s.Replicas < 1:
		return errors.New("deployment").BadParameter("replicas")
	}

	return nil
}
