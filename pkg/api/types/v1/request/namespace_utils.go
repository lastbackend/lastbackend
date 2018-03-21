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

func (n *NamespaceCreateOptions) Validate() *errors.Err {
	switch true {
	case len(n.Name) == 0:
		return errors.New("namespace").BadParameter("name")
	case len(n.Name) < 4 || len(n.Name) > 64:
		return errors.New("namespace").BadParameter("name")
	case !validator.IsNamespaceName(strings.ToLower(n.Name)):
		return errors.New("namespace").BadParameter("name")
	case len(n.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("namespace").BadParameter("description")
	default:
		// TODO: check quotas data
		return nil
	}
}

func (n *NamespaceCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("namespace").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, n)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	return n.Validate()
}

func (n *NamespaceCreateOptions) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (NamespaceRequest) UpdateOptions() *NamespaceUpdateOptions {
	return new(NamespaceUpdateOptions)
}

func (n *NamespaceUpdateOptions) Validate() *errors.Err {
	switch true {
	case n.Description != nil && len(*n.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("namespace").BadParameter("description")
	default:
		// TODO: check quotas data
		return nil
	}
}

func (n *NamespaceUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("namespace").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, n)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	return n.Validate()
}

func (n *NamespaceUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (NamespaceRequest) RemoveOptions() *NamespaceRemoveOptions {
	return new(NamespaceRemoveOptions)
}

func (n *NamespaceRemoveOptions) Validate() *errors.Err {
	return nil
}

func (n *NamespaceRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("namespace").Unknown(err)
	}

	err = json.Unmarshal(body, n)
	if err != nil {
		return errors.New("namespace").IncorrectJSON(err)
	}

	return n.Validate()
}

func (n *NamespaceRemoveOptions) ToJson() ([]byte, error) {
	return json.Marshal(n)
}
