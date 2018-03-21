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
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type ServiceRequest struct{}

func (ServiceRequest) CreateOptions() *ServiceCreateOptions {
	return new(ServiceCreateOptions)
}

func (s *ServiceCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	switch true {
	case s.Name == nil:
		return errors.New("service").BadParameter("name")
	case !validator.IsServiceName(*s.Name):
		return errors.New("service").BadParameter("name")
	case s.Replicas == nil:
		return errors.New("service").BadParameter("replicas")
	case s.Replicas != nil && *s.Replicas < DEFAULT_REPLICAS_MIN:
		return errors.New("service").BadParameter("replicas")
	case s.Sources == nil:
		return errors.New("service").BadParameter("source")
	case s.Description != nil && len(*s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("service").BadParameter("description")
	case s.Spec != nil:
		if s.Spec.Memory != nil && *s.Spec.Memory < DEFAULT_MEMORY_MIN {
			return errors.New("service").BadParameter("memory")
		}

		// TODO: need checking anymore arguments
	}

	return nil
}

func (s *ServiceCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (ServiceRequest) UpdateOptions() *ServiceUpdateOptions {
	return new(ServiceUpdateOptions)
}

func (s *ServiceUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ServiceUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	switch true {
	case s.Sources == nil:
		return errors.New("service").BadParameter("source")
	case s.Description != nil && len(*s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("service").BadParameter("description")
	case s.Spec != nil:
		if s.Spec.Memory != nil && *s.Spec.Memory < DEFAULT_MEMORY_MIN {
			return errors.New("service").BadParameter("memory")
		}

		// TODO: need checking anymore arguments
	}

	return nil
}

func (ServiceRequest) RemoveOptions() *ServiceRemoveOptions {
	return new(ServiceRemoveOptions)
}

func (s *ServiceRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	return nil
}
