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
	"io"
	"io/ioutil"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
)

type IngressRequest struct{}

func (IngressRequest) IngressConnectOptions() *IngressConnectOptions {
	cp := new(IngressConnectOptions)
	return cp
}

func (n *IngressConnectOptions) Validate() *errors.Err {
	return nil
}

func (n *IngressConnectOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, n)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return n.Validate()
}

func (n *IngressConnectOptions) ToJson() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func (IngressRequest) IngressStatusOptions() *IngressStatusOptions {
	ns := new(IngressStatusOptions)
	return ns
}

func (n *IngressStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *IngressStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, n)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return n.Validate()
}

func (s *IngressStatusOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (IngressRequest) RemoveOptions() *IngressRemoveOptions {
	return new(IngressRemoveOptions)
}

func (s *IngressRemoveOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (n *IngressRemoveOptions) Validate() *errors.Err {
	return nil
}
