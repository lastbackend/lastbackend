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
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/validator"
	"io"
	"io/ioutil"
)

type JobRequest struct{}

func (JobRequest) Manifest() *JobManifest {
	return new(JobManifest)
}

func (j *JobManifest) Validate() *errors.Err {
	switch true {
	case j.Meta.Name != nil && !validator.IsJobName(*j.Meta.Name):
		return errors.New("job").BadParameter("name")
	case j.Meta.Description != nil && len(*j.Meta.Description) > DefaultDescriptionLimit:
		return errors.New("job").BadParameter("description")
	case j.Spec.Task.Template != nil:
		if len(j.Spec.Task.Template.Containers) == 0 {
			return errors.New("job").BadParameter("spec")
		}
		if len(j.Spec.Task.Template.Containers) != 0 {
			for _, container := range j.Spec.Task.Template.Containers {
				if len(container.Image.Name) == 0 {
					return errors.New("job").BadParameter("image")
				}
			}
		}
	}

	return nil
}

func (j *JobManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("job").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("job").Unknown(err)
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return errors.New("job").IncorrectJSON(err)
	}

	if err := j.Validate(); err != nil {
		return err
	}

	return nil
}

func (JobRequest) RemoveOptions() *JobRemoveOptions {
	return new(JobRemoveOptions)
}

func (s *JobRemoveOptions) Validate() *errors.Err {
	return nil
}
