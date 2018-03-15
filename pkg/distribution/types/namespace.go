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
	"io"
	"io/ioutil"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type NamespaceList []*Namespace

type Namespace struct {
	Meta      NamespaceMeta      `json:"meta"`
	Env       NamespaceEnvs      `json:"env"`
	Resources NamespaceResources `json:"resources"`
	Quotas    NamespaceQuotas    `json:"quotas,omitempty"`
	Labels    map[string]string  `json:"labels"`
}

type NamespaceEnvs []NamespaceEnv

type NamespaceEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NamespaceMeta struct {
	Meta            `yaml:",inline"`
	Endpoint string `json:"endpoint"`
	Type     string `json:"type"`
}

type NamespaceResources struct {
	RAM    int64 `json:"ram"`
	Routes int   `json:"routes"`
}

type NamespaceQuotas struct {
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
	Disabled bool  `json:"disabled"`
}

func (n *Namespace) ToJson() ([]byte, error) {
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (n *NamespaceList) ToJson() ([]byte, error) {
	if n == nil {
		return []byte("[]"), nil
	}
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type NamespaceCreateOptions struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

func (s *NamespaceCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Namespace: decode and validate data")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Namespace: decode and validate data err: %s", err)
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Namespace: convert struct from json err: %s", err)
		return errors.New("namespace").IncorrectJSON(err)
	}

	if s.Name == "" {
		log.V(logLevel).Error("Request: Namespace: parameter name can not be empty")
		return errors.New("namespace").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 || len(s.Name) > 64 || !validator.IsNamespaceName(s.Name) {
		log.V(logLevel).Error("Request: Namespace: parameter name not valid")
		return errors.New("namespace").BadParameter("name")
	}

	return nil
}

type NamespaceUpdateOptions struct {
	Description *string                 `json:"description"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

type NamespaceQuotasOptions struct {
	Disabled bool  `json:"disabled"`
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
}

func (s *NamespaceUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Namespace: decode and validate data")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Namespace: decode and validate data err: %s", err)
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Namespace: convert struct from json err: %s", err)
		return errors.New("namespace").IncorrectJSON(err)
	}

	return nil
}
