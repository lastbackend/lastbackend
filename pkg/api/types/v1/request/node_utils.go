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

func (s *NodeInfoOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (NodeRequest) NodeStatusOptions() *NodeStatusOptions {
	return new(NodeStatusOptions)
}

func (n *NodeStatusOptions) Validate() *errors.Err {
	return nil
}

func (n *NodeStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

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

func (s *NodeStatusOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
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

func (s *NodePodStatusOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
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

func (s *NodeVolumeStatusOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
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

func (n *NodeRouteStatusOptions) ToJson() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func (NodeRequest) UpdateOptions() *NodeUpdateOptions {
	return new(NodeUpdateOptions)
}

func (s *NodeUpdateOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
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

func (NodeRequest) RemoveOptions() *NodeRemoveOptions {
	return new(NodeRemoveOptions)
}

func (s *NodeRemoveOptions) ToJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (n *NodeRemoveOptions) Validate() *errors.Err {
	return nil
}
