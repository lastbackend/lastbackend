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

type IngressRequest struct{}

func (IngressRequest) IngressConnectOptions() *IngressConnectOptions {
	cp := new(IngressConnectOptions)
	return cp
}

func (i *IngressConnectOptions) Validate() *errors.Err {
	return nil
}

func (i *IngressConnectOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("ingress").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("ingress").Unknown(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return errors.New("ingress").IncorrectJSON(err)
	}

	return i.Validate()
}

func (i *IngressConnectOptions) ToJson() ([]byte, error) {
	return json.Marshal(i)
}

func (IngressRequest) IngressStatusOptions() *IngressStatusOptions {
	ns := new(IngressStatusOptions)
	ns.Routes = make(map[string]*IngressRouteStatusOptions)
	return ns
}

func (i *IngressStatusOptions) Validate() *errors.Err {
	return nil
}

func (i *IngressStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("ingress").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("ingress").Unknown(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return errors.New("ingress").IncorrectJSON(err)
	}

	return i.Validate()
}

func (i *IngressStatusOptions) ToJson() ([]byte, error) {
	return json.Marshal(i)
}

func (IngressRequest) IngressRouteStatusOptions() *IngressRouteStatusOptions {
	return new(IngressRouteStatusOptions)
}

func (i *IngressRouteStatusOptions) Validate() *errors.Err {
	return nil
}

func (i *IngressRouteStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("ingress").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("ingress").Unknown(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return errors.New("ingress").IncorrectJSON(err)
	}

	return i.Validate()
}

func (i *IngressRouteStatusOptions) ToJson() ([]byte, error) {
	return json.Marshal(i)
}

func (IngressRequest) UpdateOptions() *IngressMetaOptions {
	return new(IngressMetaOptions)
}

func (i *IngressMetaOptions) ToJson() ([]byte, error) {
	return json.Marshal(i)
}

func (i *IngressMetaOptions) Validate() *errors.Err {
	return nil
}

func (i *IngressMetaOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("ingress").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("ingress").Unknown(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return errors.New("ingress").IncorrectJSON(err)
	}

	return i.Validate()
}

func (IngressRequest) RemoveOptions() *IngressRemoveOptions {
	return new(IngressRemoveOptions)
}

func (i *IngressRemoveOptions) ToJson() string {
	buf, _ := json.Marshal(i)
	return string(buf)
}

func (i *IngressRemoveOptions) Validate() *errors.Err {
	return nil
}
