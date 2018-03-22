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

import "fmt"

type Deployment struct {
	Meta DeploymentMeta `json:"meta"`
	// Deployment spec
	Spec DeploymentSpec `json:"spec"`
	// Deployment status
	Status DeploymentStatus `json:"status"`
	// Deployment replicas
	Replicas DeploymentReplicas `json:"replicas"`
}

type DeploymentMeta struct {
	Meta
	// Version
	Version int `json:"version"`
	// Namespace id
	Namespace string `json:"namespace"`
	// Service id
	Service string `json:"service"`
	// Endpoint
	Endpoint string `json:"endpoint"`
	// Self Link
	Status string `json:"status"`
}

type DeploymentSpec struct {
	Meta     Meta         `json:"meta"`
	Replicas int          `json:"replicas"`
	State    SpecState    `json:"state"`
	Strategy SpecStrategy `json:"strategy"`
	Triggers SpecTriggers `json:"triggers"`
	Selector SpecSelector `json:"selector"`
	Template SpecTemplate `json:"template"`
}

type DeploymentStatus struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

type DeploymentReplicas struct {
	Total     int `json:"total"`
	Provision int `json:"provision"`
	Pulling   int `json:"pulling"`
	Created   int `json:"created"`
	Started   int `json:"started"`
	Stopped   int `json:"stopped"`
	Errored   int `json:"errored"`
}

type DeploymentOptions struct {
	Replicas int `json:"replicas"`
}

func (d *Deployment) SelfLink() string {
	if d.Meta.SelfLink == "" {
		d.Meta.SelfLink = fmt.Sprintf("%s:%s:%s", d.Meta.Namespace, d.Meta.Service, d.Meta.Name)
	}
	return d.Meta.SelfLink
}

func (d *DeploymentStatus) SetProvision() {
	d.Stage = StageProvision
	d.Message = ""
}

func (d *DeploymentStatus) SetReady() {
	d.Stage = StageReady
	d.Message = ""
}

func (d *DeploymentStatus) SetCancel() {
	d.Stage = StageCancel
	d.Message = ""
}

func (d *DeploymentStatus) SetDestroy() {
	d.Stage = StageDestroy
	d.Message = ""
}
