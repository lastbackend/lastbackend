//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"time"
)

// swagger:ignore
// swagger:model types_pod
type Pod struct {
	System
	// Pod Meta
	Meta PodMeta `json:"meta" yaml:"meta"`
	// Pod Spec
	Spec PodSpec `json:"spec" yaml:"spec"`
	// Containers status info
	Status PodStatus `json:"status" yaml:"status"`
}

type PodList struct {
	System
	Items []*Pod
}

type PodMap struct {
	System
	Items map[string]*Pod
}

// swagger:ignore
// PodMeta is a meta of pod
// swagger:model types_pod_meta
type PodMeta struct {
	Meta `yaml:",inline"`
	// Pod SelfLink
	SelfLink PodSelfLink `json:"self_link" yaml:"self_link"`
	// Pod service id
	Namespace string `json:"namespace" yaml:"namespace"`
	// Pod node hostname
	Node string `json:"node" yaml:"node"`
	// Pod status
	Status string `json:"status" yaml:"status"`
	// Upstream
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type PodMetaParent struct {
	Kind     string `json:"kind" yaml:"kind"`
	SelfLink string `json:"self_link", yaml:"self_link"`
}

// PodSpec is a spec of pod
// swagger:model types_pod_spec
type PodSpec struct {
	Local    bool         `json:"local,omitempty"`
	State    SpecState    `json:"state"`
	Runtime  SpecRuntime  `json:"runtime"`
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
	// Pod runtime
	Runtime PodStatusRuntime `json:"runtime" yaml:"runtime"`
	// Pod volumes
	Volumes map[string]*VolumeClaim `json:"volumes" yaml:"volumes"`
}

type PodStatusRuntime struct {
	Services map[string]*PodContainer          `json:"containers" yaml:"containers"`
	Pipeline map[string]*PodStatusPipelineStep `json:"pipeline" yaml:"pipeline"`
}

type PodStatusPipelineStep struct {
	Status   string          `json:"status" yaml:"status"`
	Error    bool            `json:"error" yaml:"error"`
	Message  string          `json:"message" yaml:"message"`
	Commands []*PodContainer `json:"commands" yaml:"commands"`
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
type VolumeClaim struct {
	// Pod name
	Name string `json:"name" yaml:"name"`
	// Pod volume name
	Volume string `json:"volume" yaml:"volume"`
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
	Count     int       `json:"count" yaml:"count"`
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

func (s *PodStatus) SetExited() {
	s.State = StateExited
	s.Status = StateExited
	s.Running = false
	s.Message = EmptyString
}

func (s *PodStatus) SetDestroy() {
	s.State = StateDestroy
	s.Steps[StepDestroy] = PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
}

func (s *PodStatus) SetDestroyed() {
	s.State = StateDestroyed
	s.Running = false
}

func (s *PodStatus) SetPull() {
	s.State = StateProvision
	s.Status = StatusPull
	s.Running = false
	s.Steps[StepPull] = PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
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
	s.Steps[StepInitialized] = PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
}

func (s *PodStatus) SetStarting() {
	s.State = StateProvision
	s.Status = StatusStarting
	s.Running = false
	s.Message = EmptyString
	s.Steps[StepStarted] = PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
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
	s.Steps[StepStarted] = PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
}

func (s *PodStatus) SetError(err error) {
	s.State = StateError
	s.Status = StateError
	s.Message = err.Error()
}

func (s *PodSpec) SetSpecTemplate(selflink string, template SpecTemplate) {
	for _, c := range template.Containers {
		c.Labels = make(map[string]string)
		c.Labels[ContainerTypeLBC] = selflink
		c.DNS = SpecTemplateContainerDNS{}
		s.Template.Containers = append(s.Template.Containers, c)
	}

	for _, v := range template.Volumes {
		s.Template.Volumes = append(s.Template.Volumes, v)
	}
}

func (s *PodSpec) SetSpecSelector(selector SpecSelector) {
	s.Selector = selector
}

func (s *PodSpec) SetSpecRuntime(runtime SpecRuntime) {
	s.Runtime = runtime
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
		Steps: make(PodSteps, 0),
		Runtime: PodStatusRuntime{
			Services: make(map[string]*PodContainer, 0),
			Pipeline: make(map[string]*PodStatusPipelineStep, 0),
		},
		Volumes: make(map[string]*VolumeClaim, 0),
	}
	return &status
}

func (s *PodStatus) AddTask(name string) *PodStatusPipelineStep {

	pst := PodStatusPipelineStep{
		Status:   StateCreated,
		Error:    false,
		Message:  EmptyString,
		Commands: make([]*PodContainer, 0),
	}

	s.Runtime.Pipeline[name] = &pst
	return &pst
}

func (s *PodStatusPipelineStep) SetCreated() {
	s.Status = StateCreated
	s.Error = false
	s.Message = EmptyString
}

func (s *PodStatusPipelineStep) SetStarted() {
	s.Status = StateStarted
	s.Error = false
	s.Message = EmptyString
}

func (s *PodStatusPipelineStep) SetExited(error bool, message string) {
	s.Status = StateExited
	s.Error = error
	s.Message = message
}

func (s *PodStatusPipelineStep) AddTaskCommandContainer(c *PodContainer) {
	s.Commands = append(s.Commands, c)
}

func (p *Pod) SelfLink() *PodSelfLink {
	return &p.Meta.SelfLink
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
