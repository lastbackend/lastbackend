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

type NodeRequest struct{}

func (NodeRequest) NodeInfoOptions() *NodeInfoOptions {
	return new(NodeInfoOptions)
}

func (n *NodeInfoOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeInfoOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (NodeRequest) NodeStateOptions() *NodeStateOptions {
	return new(NodeStateOptions)
}

func (n *NodeStateOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeStateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (NodeRequest) NodePodStatusOptions() *NodePodStatusOptions {
	return new(NodePodStatusOptions)
}

func (n *NodePodStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *NodePodStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (NodeRequest) NodeVolumeStatusOptions() *NodeVolumeStatusOptions {
	return new(NodeVolumeStatusOptions)
}

func (n *NodeVolumeStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeVolumeStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (NodeRequest) NodeRouteStatusOptions() *NodeRouteStatusOptions {
	return new(NodeRouteStatusOptions)
}

func (n *NodeRouteStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeRouteStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (n *NodeRouteStatusOptions) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (NodeRequest) UpdateOptions() *NodeUpdateOptions {
	return new(NodeUpdateOptions)
}

func (n *NodeUpdateOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (n *NodeUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (NodeRequest) RemoveOptions() *NodeRemoveOptions {
	return new(NodeRemoveOptions)
}

func (n *NodeRemoveOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
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
