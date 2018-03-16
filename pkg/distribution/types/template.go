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

package types

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"io/ioutil"
	"time"
)

type Template struct {
	ID          string    `json:"id"`
	RepoID      string    `json:"repo"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Main        bool      `json:"main"`
	Shared      bool      `json:"shared"`
	Deleted     bool      `json:"deleted"`
	Spec        PodSpec   `json:"spec"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type TemplateCreateOptions struct {
	TemplateOptions
}

func (s *TemplateCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for deploy template creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for deploy template creating err: %s", err)
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err)
		return errors.New("repo").IncorrectJSON(err)
	}

	if s.Name == nil {
		log.V(logLevel).Error("Request: Template: parameter name can not be empty")
		return errors.New("template").BadParameter("name")
	}

	return nil
}

type TemplateOptions TemplateOptionsSpec

func (s *TemplateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Repo: decode and validate data for deploy template creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: decode and validate data for deploy template creating err: %s", err)
		return errors.New("repo").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Repo: convert struct from json err: %s", err)
		return errors.New("repo").IncorrectJSON(err)
	}

	return nil
}

type TemplateOptionsSpec struct {
	ID          *string                    `json:"id,omitempty"`
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Main        *bool                      `json:"main"`
	Tag         *string                    `json:"tag"`
	Shared      *bool                      `json:"shared"`
	Memory      *int64                     `json:"memory,omitempty"`
	Entrypoint  *string                    `json:"entrypoint,omitempty"`
	Command     *string                    `json:"command,omitempty"`
	EnvVars     *[]string                  `json:"env,omitempty"`
	Ports       *[]TemplateOptionsSpecPort `json:"ports,omitempty"`
}

type TemplateOptionsSpecPort struct {
	Protocol  string `json:"protocol"`
	External  int    `json:"external"`
	Internal  int    `json:"internal"`
	Published bool   `json:"published"`
}
