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
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"strings"
)

const logLevel = 3

type RequestServiceCreateS struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Registry    string               `json:"registry"`
	Region      string               `json:"region"`
	Template    string               `json:"template"`
	Image       string               `json:"image"`
	Url         string               `json:"url"`
	Replicas    *int                 `json:"replicas,omitempty"`
	Spec        *RequestServiceSpecS `json:"spec"`
	Source      types.ServiceSource
}

type ServiceSpec struct {
	ID         *int64    `json:"id,omitempty"`
	Memory     *int64    `json:"memory,omitempty"`
	Entrypoint *string   `json:"entrypoint,omitempty"`
	Command    *string   `json:"command,omitempty"`
	Image      *string   `json:"image,omitempty"`
	EnvVars    *[]string `json:"env,omitempty"`
	Ports      *[]Port   `json:"ports,omitempty"`
}

type Port struct {
	Protocol  string `json:"protocol"`
	External  int    `json:"external"`
	Internal  int    `json:"internal"`
	Published bool   `json:"published"`
}

type resources struct {
	Region string `json:"region"`
	Memory int    `json:"memory"`
}

func (s *RequestServiceCreateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log := context.Get().GetLogger()

	log.V(logLevel).Debug("Request: Service: decode and validate data for creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: decode and validate data for creating err: %s", err.Error())
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: convert struct from json err: %s", err.Error())
		return errors.New("service").IncorrectJSON(err)
	}

	if s.Template == "" && s.Image == "" && s.Url == "" {
		log.V(logLevel).Error("Request: Service: One of the following parameters(template, image, url) is required")
		return errors.New("service").BadParameter("template,image,url")
	}

	if s.Template != "" {
		if s.Name == "" {
			s.Name = s.Template
		}
	}

	if s.Image != "" && s.Url == "" {
		source, err := converter.DockerNamespaceParse(s.Image)
		if err != nil {
			log.V(logLevel).Errorf("Request: Service: parameter image not valid err: %s", err.Error())
			return errors.New("service").BadParameter("image")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}
	}

	if s.Url != "" {

		match := strings.Split(s.Url, "#")

		if !validator.IsGitUrl(match[0]) {
			log.V(logLevel).Error("Request: Service: parameter url not valid")
			return errors.New("service").BadParameter("url")
		}

		source, err := converter.GitUrlParse(match[0])
		if err != nil {
			log.V(logLevel).Error("Request: Service: parameter url not valid err: %s", err.Error())
			return errors.New("service").BadParameter("url")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		if len(match) > 1 {
			source.Branch = match[len(match)-1]
		}

		s.Source = types.ServiceSource{
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: source.Branch,
		}
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name == "" {
		log.V(logLevel).Error("Request: Service: parameter name can not be empty")
		return errors.New("service").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
		log.V(logLevel).Error("Request: Namespace: parameter name not valid")
		return errors.New("service").BadParameter("name")
	}

	if s.Spec == nil {
		s.Spec = new(RequestServiceSpecS)
	}

	if s.Replicas != nil && *s.Replicas < 1 {
		*s.Replicas = 1
	}

	// TODO: Need validate data format in config

	return nil
}

type RequestServiceUpdateS struct {
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Replicas    *int         `json:"replicas"`
	Spec        *ServiceSpec `json:"spec"`
}

func (s *RequestServiceUpdateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log := context.Get().GetLogger()

	log.V(logLevel).Debug("Request: Service: decode and validate data for updating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: decode and validate data for creating err: %s", err.Error())
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: convert struct from json err: %s", err.Error())
		return errors.New("service").IncorrectJSON(err)
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name != "" {
		s.Name = strings.ToLower(s.Name)

		if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
			log.V(logLevel).Error("Request: Namespace: parameter name not valid")
			return errors.New("service").BadParameter("name")
		}
	}

	if s.Replicas != nil && *s.Replicas < 1 {
		*s.Replicas = 1
	}

	// TODO: Need validate data format in config

	return nil
}

type RequestServiceSpecS struct {
	ServiceSpec
}

func (s *RequestServiceSpecS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log := context.Get().GetLogger()

	log.V(logLevel).Debug("Request: Service: decode and validate data for service spec")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: decode and validate data for service spec err: %s", err.Error())
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: convert struct from json err: %s", err.Error())
		return errors.New("service").IncorrectJSON(err)
	}

	return nil
}
