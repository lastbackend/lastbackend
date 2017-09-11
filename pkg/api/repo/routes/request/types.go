//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"strings"
)

const logLevel = 3

type RequestRepoCreateS struct {
	Name string `json:"name"`
}

func (s *RequestRepoCreateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for creating err: %s", err.Error())
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err.Error())
		return errors.New("repo").IncorrectJSON(err)
	}

	if s.Name == "" {
		log.V(logLevel).Error("Request: Repo: parameter name can not be empty")
		return errors.New("repo").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if !validator.IsRepoName(s.Name) {
		return errors.New("repo").BadParameter("name")
	}

	return nil
}

type RequestRepoUpdateS struct {
	Technology  *string `json:"technology,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *RequestRepoUpdateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for updating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for creating err: %s", err.Error())
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err.Error())
		return errors.New("repo").IncorrectJSON(err)
	}

	return nil
}

type RequestRepoDeployTempalteCreateS struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	EnvVars     []string `json:"environments"`
	Ports       []struct {
		Protocol  string `json:"protocol"`
		Container int    `json:"internal"`
		Host      int    `json:"external"`
		Published bool   `json:"published"`
	} `json:"ports"`
	Memory int64 `json:"memory"`
	Shared bool  `json:"shared"`
}

func (s *RequestRepoDeployTempalteCreateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for deploy template creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for deploy template creating err: %s", err.Error())
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err.Error())
		return errors.New("repo").IncorrectJSON(err)
	}

	return nil
}

type RequestRepoDeployTempalteUpdateS struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	EnvVars     []string `json:"environments"`
	Ports       []struct {
		Protocol  string `json:"protocol"`
		Container int    `json:"internal"`
		Host      int    `json:"external"`
		Published bool   `json:"published"`
	} `json:"ports"`
	Memory int64 `json:"memory"`
	Shared bool  `json:"shared"`
}

func (s *RequestRepoDeployTempalteUpdateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for deploy template updating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for deploy template updating err: %s", err.Error())
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err.Error())
		return errors.New("repo").IncorrectJSON(err)
	}

	return nil
}
