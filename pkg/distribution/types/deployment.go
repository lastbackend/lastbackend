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
type DeploymentMap map[string]*Deployment
type DeploymentList []*Deployment

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
	Replicas int          `json:"replicas"`
	State    SpecState    `json:"state"`
	Selector SpecSelector `json:"selector"`
	Template SpecTemplate `json:"template"`
}

type DeploymentStatus struct {
	State   string `json:"state"`
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
		d.Meta.SelfLink = d.CreateSelfLink(d.Meta.Namespace, d.Meta.Service, d.Meta.Name)
	}
	return d.Meta.SelfLink
}

func (d *Deployment) CreateSelfLink(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

func (d *DeploymentStatus) SetProvision() {
	d.State = StateProvision
	d.Message = ""
}

func (d *DeploymentStatus) SetReady() {
	d.State = StateReady
	d.Message = ""
}

func (d *DeploymentStatus) SetCancel() {
	d.State = StateCancel
	d.Message = ""
}

func (d *DeploymentStatus) SetDestroy() {
	d.State = StateDestroy
	d.Message = ""
}

type DeploymentUpdateOptions struct {
	// Number of replicas
	Replicas *int
	// Deployment status for update
	Status *struct {
		State   string
		Message string
	}
}
