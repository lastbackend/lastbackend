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
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type ServiceRequest struct{}

func (ServiceRequest) CreateOptions() *ServiceCreateOptions {
	return new(ServiceCreateOptions)
}

func (s *ServiceCreateOptions) Validate() *errors.Err {
	switch true {
	case s.Name != nil && !validator.IsServiceName(*s.Name):
		return errors.New("service").BadParameter("name")
	case s.Image == nil:
		return errors.New("service").BadParameter("image")
	case s.Description != nil && len(*s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("service").BadParameter("description")
	case s.Spec != nil:
		if s.Spec.Replicas != nil && *s.Spec.Replicas < DEFAULT_REPLICAS_MIN {
			return errors.New("service").BadParameter("replicas")
		}

		if s.Spec.Memory != nil && *s.Spec.Memory < DEFAULT_MEMORY_MIN {
			return errors.New("service").BadParameter("memory")
		}
	}
	return nil
}

func (s *ServiceCreateOptions) DecodeAndValidate(reader io.Reader) (*types.ServiceCreateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return nil, errors.New("service").IncorrectJSON(err)
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	opts := new(types.ServiceCreateOptions)
	opts.Name = s.Name
	opts.Description = s.Description
	opts.Image = s.Image

	if s.Spec != nil {
		opts.Spec = new(types.ServiceOptionsSpec)
		opts.Spec.Replicas = s.Spec.Replicas
		opts.Spec.Memory = s.Spec.Memory
		opts.Spec.EnvVars = s.Spec.EnvVars
		opts.Spec.Entrypoint = s.Spec.Entrypoint
		opts.Spec.Command = s.Spec.Command

		if s.Spec.Ports != nil {
			opts.Spec.Ports = s.Spec.Ports
		}

	}

	return opts, nil
}

func (s *ServiceCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (ServiceRequest) UpdateOptions() *ServiceUpdateOptions {
	return new(ServiceUpdateOptions)
}

func (s *ServiceUpdateOptions) Validate() *errors.Err {
	switch true {
	case s.Description != nil && len(*s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("service").BadParameter("description")
	case s.Spec != nil:
		if s.Spec.Memory != nil && *s.Spec.Memory < DEFAULT_MEMORY_MIN {
			return errors.New("service").BadParameter("memory")
		}
	}
	return nil
}

func (s *ServiceUpdateOptions) DecodeAndValidate(reader io.Reader) (*types.ServiceUpdateOptions, *errors.Err) {

	if reader == nil {
		err := errors.New("data body can not be null")
		return nil, errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return nil, errors.New("service").IncorrectJSON(err)
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	opts := new(types.ServiceUpdateOptions)
	opts.Description = s.Description

	if s.Spec != nil {
		opts.Spec = new(types.ServiceOptionsSpec)
		opts.Spec.Memory = s.Spec.Memory
		opts.Spec.EnvVars = s.Spec.EnvVars
		opts.Spec.Entrypoint = s.Spec.Entrypoint
		opts.Spec.Command = s.Spec.Command

		if len(s.Spec.Ports) != 0 {
			opts.Spec.Ports = s.Spec.Ports
		}
	}

	return opts, nil
}

func (s *ServiceUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (ServiceRequest) RemoveOptions() *ServiceRemoveOptions {
	return new(ServiceRemoveOptions)
}

func (s *ServiceRemoveOptions) Validate() *errors.Err {
	return nil
}
