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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"io/ioutil"
)

const (
	CentralUSRegions = "CU"
	WestEuropeRegion = "WE"
	EastAsiaRegion   = "EA"
)

type ClusterList []*Cluster

type Cluster struct {
	Meta  ClusterMeta  `json:"meta"`
	State ClusterState `json:"state"`
}

type ClusterMeta struct {
	Meta

	Region   string `json:"region"`
	Token    string `json:"token"`
	Provider string `json:"provider"`
	Shared   bool   `json:"shared"`
	Main     bool   `json:"main"`
}

type ClusterState struct {
	Nodes struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
	} `json:"nodes"`
	Capacity  ClusterResources `json:"capacity"`
	Allocated ClusterResources `json:"allocated"`
	Deleted   bool             `json:"deleted"`
}

type ClusterResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}

type ClusterCreateOptions struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *ClusterCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Cluster: decode and validate data for creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Cluster: decode and validate data for creating err: %s", err)
		return errors.New("cluster").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Cluster: convert struct from json err: %s", err)
		return errors.New("cluster").IncorrectJSON(err)
	}

	if s.Name == "" {
		log.V(logLevel).Error("Request: Cluster: parameter name can not be empty")
		return errors.New("cluster").BadParameter("name")
	}

	if len(s.Name) < 4 && len(s.Name) > 64 {
		log.V(logLevel).Error("Request: Cluster: parameter name not valid")
		return errors.New("cluster").BadParameter("name")
	}

	return nil
}

type ClusterUpdateOptions struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (s *ClusterUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Cluster: decode and validate data for updating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Cluster: decode and validate data for updating err: %s", err)
		return errors.New("cluster").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Cluster: convert struct from json err: %s", err)
		return errors.New("cluster").IncorrectJSON(err)
	}

	if s.Name != nil && *s.Name == "" {
		log.V(logLevel).Error("Request: Cluster: parameter name can not be empty")
		return errors.New("cluster").BadParameter("name")
	}

	if s.Name != nil {
		if len(*s.Name) < 4 && len(*s.Name) > 64 {
			log.V(logLevel).Error("Request: Cluster: parameter name not valid")
			return errors.New("cluster").BadParameter("name")
		}
	}

	return nil
}
