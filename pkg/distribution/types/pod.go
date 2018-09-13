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
	"fmt"
	"time"
)

// swagger:ignore
// swagger:model types_pod
type Pod struct {
	Runtime
	// Pod Meta
	Meta PodMeta `json:"meta" yaml:"meta"`
	// Pod Spec
	Spec PodSpec `json:"spec" yaml:"spec"`
	// Containers status info
	Status PodStatus `json:"status" yaml:"status"`
}

type PodList struct {
	Runtime
	Items []*Pod
}

type PodMap struct {
	Runtime
	Items map[string]*Pod
}

// swagger:ignore
// PodMeta is a meta of pod
// swagger:model types_pod_meta
type PodMeta struct {
	Meta `yaml:",inline"`
	// Pod SelfLink
	SelfLink string `json:"self_link" yaml:"self_link"`
	// Pod deployment
	Deployment string `json:"deployment" yaml:"deployment"`
	// Pod service
	Service string `json:"service" yaml:"service"`
	// Pod service id
	Namespace string `json:"namespace" yaml:"namespace"`
	// Pod node hostname
	Node string `json:"node" yaml:"node"`
	// Pod status
	Status string `json:"status" yaml:"status"`
	// Endpoint
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// PodSpec is a spec of pod
// swagger:model types_pod_spec
type PodSpec struct {
	Local    bool         `json:"local,omitempty"`
	State    SpecState    `json:"state"`
	Selector SpecSelector `json:"selector"`
	Template SpecTemplate `json:"template" yaml:"template"`
}

// swagger:ignore
// PodSpecStatus is a status of pod
// swagger:model types_pod_status
type PodStatus struct {
	// Pod type
	Local bool `json:"local" yaml:"local"`
	// Pod state
	State string `json:"state" yaml:"state"`
	// Pod status
	Status string `json:"status" yaml:"status"`
	// Pod state
	Running bool `json:"running" yaml:"state"`
	// Pod state message
	Message string `json:"message" yaml:"message"`
	// Pod steps
	Steps PodSteps `json:"steps" yaml:"steps"`
	// Pod network
	Network PodNetwork `json:"network" yaml:"network"`
	// Pod containers
	Containers map[string]*PodContainer `json:"containers" yaml:"containers"`
	// Pod volumes
	Volumes map[string]*PodVolume `json:"volumes" yaml:"volumes"`
}

// PodSteps is a map of pod steps
// swagger:model types_pod_step_map
type PodSteps map[string]PodStep

// swagger:model types_pod_step
type PodStep struct {
	// Pod step ready
	Ready bool `json:"ready" yaml:"ready"`
	// Pod step timestamp
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

// swagger:model types_pod_network
type PodNetwork struct {
	// Pod host IP
	HostIP string `json:"host_ip" yaml:"host_ip"`
	// Pod IP
	PodIP string `json:"pod_ip" yaml:"pod_ip"`
}

// PodContainer is a container of the pod
// swagger:model types_pod_container
type PodContainer struct {
	// Pod container ID
	ID string `json:"id" yaml:"id"`
	// Pod ID
	Pod string `json:"pod" yaml:"pod"`
	// Pod container name
	Name string `json:"name" yaml:"name"`
	// Pod container exec
	Exec SpecTemplateContainerExec `json:"exec" yaml:"exec"`
	// Pod container state
	State PodContainerState `json:"state" yaml:"state"`
	// Pod container ready
	Ready bool `json:"ready" yaml:"ready"`
	// Pod container restart count
	Restart struct {
		Policy  string `json:"policy"`
		Attempt int    `json:"count"`
	} `json:"restart" yaml:"restart"`
	// Pod container image meta
	Image PodContainerImage `json:"image" yaml:"image"`
	// Pod container envs
	Envs []string `json:"-"`
	// Pod container binds
	Binds []string `json:"-"`
	// Pod container ports
	Ports []*SpecTemplateContainerPort `json:"ports"`
}

// PodContainer is a container of the pod
// swagger:model types_pod_container
type PodVolume struct {
	// Pod name
	Pod string `json:"pod" yaml:"pod"`
	// Pod volume name
	Name string `json:"name" yaml:"name"`
	// Pod volume ready flag
	Ready bool `json:"ready" yaml:"ready"`
	// Pod volume string
	Type string `json:"type" yaml:"type"`
	// Pod volume Path
	Path string `json:"path" yaml:"path"`
}

// swagger:model types_pod_container_image
type PodContainerImage struct {
	// Pod container image ID
	ID string `json:"id" yaml:"id"`
	// Pod container image name
	Name string `json:"name" yaml:"name"`
}

// swagger:model types_pod_container_state
type PodContainerState struct {
	// Container restart state
	Restarted PodContainerStateRestarted `json:"restarted" yaml:"restarted"`
	// Container create state
	Created PodContainerStateCreated `json:"created" yaml:"created"`

	// Container started state
	Started PodContainerStateStarted `json:"started" yaml:"started"`

	// Container stopped state
	Stopped PodContainerStateStopped `json:"stopped" yaml:"stopped"`

	// Container error state
	Error PodContainerStateError `json:"error" yaml:"error"`
}

// swagger:model types_pod_container_state_restarted
type PodContainerStateRestarted struct {
	Count int `json:"count" yaml:"count"`
	Restarted time.Time `json:"restarted" yaml:"restarted"`
}

// swagger:model types_pod_container_state_created
type PodContainerStateCreated struct {
	Created time.Time `json:"created" yaml:"created"`
}

// swagger:model types_pod_container_state_started
type PodContainerStateStarted struct {
	Started   bool      `json:"started" yaml:"started"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

// swagger:model types_pod_container_state_stopped
type PodContainerStateStopped struct {
	Stopped bool                  `json:"stopped" yaml:"stopped"`
	Exit    PodContainerStateExit `json:"exit" yaml:"exit"`
}

// swagger:model types_pod_container_state_error
type PodContainerStateError struct {
	Error   bool                  `json:"error" yaml:"error"`
	Message string                `json:"message" yaml:"message"`
	Exit    PodContainerStateExit `json:"exit" yaml:"exit"`
}

// swagger:model types_pod_container_state_exit
type PodContainerStateExit struct {
	Code      int       `json:"code" yaml:"code"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

func (s *PodStatus) SetInitialized() {
	s.State = StateProvision
	s.Status = StatusInitialized
	s.Running = false
	s.Message = EmptyString
}

func (s *PodStatus) SetDestroy() {
	s.State = StateDestroy
}

func (s *PodStatus) SetDestroyed() {
	s.State = StateDestroyed
	s.Running = false
}

func (s *PodStatus) SetPull() {
	s.State = StateProvision
	s.Status = StatusPull
	s.Running = false
}

func (s *PodStatus) SetProvision() {
	s.State = StateProvision
	s.Running = false
}

func (s *PodStatus) SetCreated() {
	s.State = StateProvision
	s.Status = StateCreated
	s.Running = false
	s.Message = EmptyString
}

func (s *PodStatus) SetStarting() {
	s.State = StateProvision
	s.Status = StatusStarting
	s.Running = false
	s.Message = EmptyString
}

func (s *PodStatus) SetRunning() {
	s.State = StateReady
	s.Status = StatusRunning
	s.Running = true
	s.Message = EmptyString
}

func (s *PodStatus) SetStopped() {
	s.State = StateReady
	s.Status = StatusStopped
	s.Running = false
	s.Message = EmptyString
}

func (s *PodStatus) SetError(err error) {
	s.State = StateError
	s.Message = err.Error()
}

func NewPod() *Pod {
	pod := new(Pod)
	pod.Status = *NewPodStatus()
	return pod
}

func NewPodList() *PodList {
	dm := new(PodList)
	dm.Items = make([]*Pod, 0)
	return dm
}

func NewPodMap() *PodMap {
	dm := new(PodMap)
	dm.Items = make(map[string]*Pod)
	return dm
}

func NewPodStatus() *PodStatus {
	status := PodStatus{
		Steps:      make(PodSteps, 0),
		Containers: make(map[string]*PodContainer, 0),
		Volumes:    make(map[string]*PodVolume, 0),
	}
	return &status
}

func (p *Pod) SelfLink() string {
	if p.Meta.SelfLink == "" {
		p.Meta.SelfLink = p.CreateSelfLink(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment, p.Meta.Name)
	}
	return p.Meta.SelfLink
}

func (p *Pod) ServiceLink() string {
	return new(Service).CreateSelfLink(p.Meta.Namespace, p.Meta.Service)
}

func (p *Pod) DeploymentLink() string {
	return new(Deployment).CreateSelfLink(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
}

func (p *Pod) CreateSelfLink(namespace, service, deployment, name string) string {
	return fmt.Sprintf("%s:%s:%s:%s", namespace, service, deployment, name)
}

func (c *PodContainer) GetManifest() *ContainerManifest {
	var manifest = new(ContainerManifest)

	manifest.Name = c.Name
	manifest.Image = c.Image.Name
	manifest.Binds = c.Binds
	manifest.Envs = c.Envs
	manifest.Ports = c.Ports
	manifest.Exec = c.Exec

	manifest.RestartPolicy.Policy = c.Restart.Policy
	manifest.RestartPolicy.Attempt = c.Restart.Attempt

	return manifest
}
