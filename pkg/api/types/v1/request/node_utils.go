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

func (s *NodeInfoOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) NodeStateOptions() *NodeStateOptions {
	return new(NodeStateOptions)
}

func (s *NodeStateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) NodePodStatusOptions() *NodePodStatusOptions {
	return new(NodePodStatusOptions)
}

func (s *NodePodStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) NodeVolumeStatusOptions() *NodeVolumeStatusOptions {
	return new(NodeVolumeStatusOptions)
}

func (s *NodeVolumeStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) NodeRouteStatusOptions() *NodeRouteStatusOptions {
	return new(NodeRouteStatusOptions)
}

func (s *NodeRouteStatusOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) UpdateOptions() *NodeUpdateOptions {
	return new(NodeUpdateOptions)
}

func (s *NodeUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("node").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}

func (NodeRequest) RemoveOptions() *NodeRemoveOptions {
	return new(NodeRemoveOptions)
}

func (s *NodeRemoveOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		return nil
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("node").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}
