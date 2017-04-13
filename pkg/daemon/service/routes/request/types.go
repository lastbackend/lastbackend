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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"strings"
)

type RequestServiceCreateS struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Registry    string                 `json:"registry"`
	Region      string                 `json:"region"`
	Template    string                 `json:"template"`
	Image       string                 `json:"image"`
	Url         string                 `json:"url"`
	Config      map[string]interface{} `json:"config"`
	Source      types.ServiceSource
}

type resources struct {
	Region string `json:"region"`
	Memory int    `json:"memory"`
}

func (s *RequestServiceCreateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	s.Config = make(map[string]interface{})
	log := context.Get().GetLogger()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return errors.New("user").Unknown(err)
	}

	log.Debug(string(body))

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	if s.Template == "" && s.Image == "" && s.Url == "" {
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
			return errors.New("service").BadParameter("image")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}
	}

	if s.Url != "" {
		if !validator.IsGitUrl(s.Url) {
			return errors.New("service").BadParameter("url")
		}

		source, err := converter.GitUrlParse(s.Url)
		if err != nil {
			return errors.New("service").BadParameter("url")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		s.Source = types.ServiceSource{
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: "master",
		}
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name == "" {
		return errors.New("service").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
		return errors.New("service").BadParameter("name")
	}

	return nil
}

type RequestServiceUpdateS struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Domains     []string               `json:"domains"`
}

func (s *RequestServiceUpdateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	s.Config = make(map[string]interface{})
	log := context.Get().GetLogger()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name != "" {
		s.Name = strings.ToLower(s.Name)

		if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
			return errors.New("service").BadParameter("name")
		}
	}

	return nil
}
