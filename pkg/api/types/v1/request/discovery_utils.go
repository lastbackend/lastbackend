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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type DiscoveryRequest struct{}

func (DiscoveryRequest) DiscoveryConnectOptions() *DiscoveryConnectOptions {
	cp := new(DiscoveryConnectOptions)
	return cp
}

func (n *DiscoveryConnectOptions) Validate() *errors.Err {
	return nil
}

func (n *DiscoveryConnectOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (s *DiscoveryConnectOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (DiscoveryRequest) DiscoveryStatusOptions() *DiscoveryStatusOptions {
	ns := new(DiscoveryStatusOptions)
	return ns
}

func (n *DiscoveryStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *DiscoveryStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (s *DiscoveryStatusOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (DiscoveryRequest) RemoveOptions() *DiscoveryRemoveOptions {
	return new(DiscoveryRemoveOptions)
}

func (s *DiscoveryRemoveOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (n *DiscoveryRemoveOptions) Validate() *errors.Err {
	return nil
}
