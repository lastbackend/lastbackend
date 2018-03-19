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
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type NamespaceRequest struct{}

func (NamespaceRequest) CreateOptions() *NamespaceCreateOptions {
	return new(NamespaceCreateOptions)
}

func (s *NamespaceCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("namespace").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	switch true {
	case len(s.Name) == 0:
		return errors.New("namespace").BadParameter("name")
	case len(s.Name) < 4 || len(s.Name) > 64:
		return errors.New("namespace").BadParameter("name")
	case !validator.IsNamespaceName(strings.ToLower(s.Name)):
		return errors.New("namespace").BadParameter("name")
	case len(s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("namespace").BadParameter("description")
	case s.Quotas != nil:
		// TODO: check quotas data
	}

	return nil
}

func (NamespaceRequest) UpdateOptions() *NamespaceUpdateOptions {
	return new(NamespaceUpdateOptions)
}

func (s *NamespaceUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("namespace").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	switch true {
	case s.Description != nil && len(*s.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("namespace").BadParameter("description")
	case s.Quotas != nil:
		// TODO: check quotas data
	}

	return nil
}

func (NamespaceRequest) RemoveOptions() *NamespaceRemoveOptions {
	return new(NamespaceRemoveOptions)
}

func (s *NamespaceRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	return nil
}
