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

package request

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/util/validator"
	"io"
	"io/ioutil"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
)

type DeploymentRequest struct{}

func (DeploymentRequest) Manifest() *DeploymentManifest {
	return new(DeploymentManifest)
}

func (s *DeploymentManifest) Validate() *errors.Err {
	switch true {
	case s.Meta.Name != nil && !validator.IsServiceName(*s.Meta.Name):
		return errors.New("deployment").BadParameter("name")
	case s.Meta.Description != nil && len(*s.Meta.Description) > DefaultDescriptionLimit:
		return errors.New("deployment").BadParameter("description")
	case len(s.Spec.Template.Containers) == 0:
		return errors.New("deployment").BadParameter("spec")
	case len(s.Spec.Template.Containers) != 0:
		for _, container := range s.Spec.Template.Containers {
			if len(container.Image.Name) == 0 {
				return errors.New("deployment").BadParameter("image")
			}
		}
	}

	return nil
}

func (s *DeploymentManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("deployment").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("deployment").Unknown(err)
	}

	fmt.Println(string(body))

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("deployment").IncorrectJSON(err)
	}

	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}

func (DeploymentRequest) UpdateOptions() *DeploymentUpdateOptions {
	return new(DeploymentUpdateOptions)
}

func (d *DeploymentUpdateOptions) Validate() *errors.Err {
	switch true {
	case d.Replicas == nil:
		return errors.New("deployment").BadParameter("replicas")
	case *d.Replicas < 1:
		return errors.New("deployment").BadParameter("replicas")
	}
	return nil
}

func (d *DeploymentUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("deployment").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("deployment").Unknown(err)
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return errors.New("deployment").IncorrectJSON(err)
	}

	return d.Validate()
}

func (d *DeploymentUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(d)
}

func (DeploymentRequest) RemoveOptions() *DeploymentRemoveOptions {
	return new(DeploymentRemoveOptions)
}
