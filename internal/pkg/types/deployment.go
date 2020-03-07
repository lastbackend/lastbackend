//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

type DeploymentMap struct {
	System
	Items map[string]*Deployment
}

type DeploymentList struct {
	System
	Items []*Deployment
}

type Deployment struct {
	System
	// Deployment Meta
	Meta DeploymentMeta `json:"meta"`
	// Deployment status
	Status DeploymentStatus `json:"status"`
	// Deployment spec
	Spec DeploymentSpec `json:"spec"`
}

type DeploymentMeta struct {
	Meta
	// Version
	Version int `json:"version"`
	// Environment id
	Namespace string `json:"namespace"`
	// Applications id
	Service string `json:"service"`
	// Upstream
	Endpoint string `json:"endpoint"`
	// Self Link
	SelfLink DeploymentSelfLink `json:"self_link"`
}

type DeploymentSpec struct {
	Replicas int          `json:"replicas"`
	State    SpecState    `json:"state"`
	Selector SpecSelector `json:"selector"`
	Template SpecTemplate `json:"template"`
}

type DeploymentStatus struct {
	State        string             `json:"state"`
	Message      string             `json:"message"`
	Dependencies StatusDependencies `json:"dependencies"`
}

type StatusDependencies struct {
	Volumes map[string]StatusDependency `json:"volumes"`
	Secrets map[string]StatusDependency `json:"secrets"`
	Configs map[string]StatusDependency `json:"configs"`
}

type StatusDependency struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Status string `json:"status"`
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

type DeploymentManifest struct {
	Meta DeploymentMeta `json:"meta"`
}

func (d *Deployment) SelfLink() *DeploymentSelfLink {
	return &d.Meta.SelfLink
}

func (d *Deployment) ServiceLink() *ServiceSelfLink {
	return d.Meta.SelfLink.parent.SelfLink.(*ServiceSelfLink)
}

func (ds *DeploymentStatus) CheckDeps() bool {

	for _, d := range ds.Dependencies.Volumes {
		if d.Status != StateReady {
			return false
		}
	}

	for _, d := range ds.Dependencies.Secrets {
		if d.Status != StateReady {
			return false
		}
	}

	for _, d := range ds.Dependencies.Configs {
		if d.Status != StateReady {
			return false
		}
	}

	return true
}

func (ds *DeploymentStatus) SetCreated() {
	ds.State = StateCreated
	ds.Message = ""
}

func (ds *DeploymentStatus) SetProvision() {
	ds.State = StateProvision
	ds.Message = ""
}

func (ds *DeploymentStatus) SetReady() {
	ds.State = StateReady
	ds.Message = ""
}

func (ds *DeploymentStatus) SetCancel() {
	ds.State = StateCancel
	ds.Message = ""
}

func (ds *DeploymentStatus) SetDestroy() {
	ds.State = StateDestroy
	ds.Message = ""
}

func (s *DeploymentSpec) GetResourceRequest() ResourceRequest {

	rr := ResourceRequest{}

	var (
		limitsRAM int64
		limitsCPU int64

		requestRAM int64
		requestCPU int64
	)

	for _, c := range s.Template.Containers {

		limitsCPU += c.Resources.Limits.CPU
		limitsRAM += c.Resources.Limits.RAM

		requestCPU += c.Resources.Request.CPU
		requestRAM += c.Resources.Request.RAM
	}

	if requestRAM > 0 {
		requestRAM = int64(s.Replicas) * requestRAM
		rr.Request.RAM = requestRAM
	}

	if requestCPU > 0 {
		requestCPU = int64(s.Replicas) * requestCPU
		rr.Request.CPU = requestCPU
	}

	if limitsRAM > 0 {
		limitsRAM = int64(s.Replicas) * limitsRAM
		rr.Limits.RAM = limitsRAM
	}

	if limitsCPU > 0 {
		limitsCPU = int64(s.Replicas) * limitsCPU
		rr.Limits.CPU = limitsCPU
	}

	return rr
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

func NewDeploymentList() *DeploymentList {
	dm := new(DeploymentList)
	dm.Items = make([]*Deployment, 0)
	return dm
}

func NewDeploymentMap() *DeploymentMap {
	dm := new(DeploymentMap)
	dm.Items = make(map[string]*Deployment)
	return dm
}
