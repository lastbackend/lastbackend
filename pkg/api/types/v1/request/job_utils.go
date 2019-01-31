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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"io"
	"io/ioutil"
)

type JobRequest struct{}

func (JobRequest) Manifest() *JobManifest {
	return new(JobManifest)
}

func (j *JobManifest) Validate() *errors.Err {
	switch true {
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
