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
	"sync"
	"time"
)

type Pod struct {
	// Lock map
	lock sync.RWMutex
	// Pod Meta
	Meta PodMeta `json:"meta" yaml:"meta"`
	// Pod Spec
	Spec PodSpec `json:"spec" yaml:"spec"`
	// Containers status info
	Status PodStatus `json:"status" yaml:"status"`
}

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

type PodSpec struct {
	State    SpecState    `json:"state"`
	Template SpecTemplate `json:"template" yaml:"template"`
}

type PodStatus struct {
	// Pod stage
	State string `json:"stage" yaml:"stage"`
	// Pod state message
	Message string `json:"message" yaml:"message"`
	// Pod steps
	Steps PodSteps `json:"steps" yaml:"steps"`
	// Pod network
	Network PodNetwork `json:"network" yaml:"network"`
	// Pod containers
	Containers map[string]*PodContainer `json:"containers" yaml:"containers"`
}

type PodSteps map[string]PodStep

type PodStep struct {
	// Pod step ready
	Ready bool `json:"ready" yaml:"ready"`
	// Pod step timestamp
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

type PodNetwork struct {
	// Pod host IP
	HostIP string `json:"host_ip" yaml:"host_ip"`
	// Pod IP
	PodIP string `json:"pod_ip" yaml:"pod_ip"`
}

type PodContainer struct {
	// Pod container ID
	ID string `json:"id" yaml:"id"`
	// Pod ID
	Pod string `json:"pod" yaml:"pod"`
	// Pod container name
	Name string `json:"name" yaml:"name"`
	// Pod container state
	State PodContainerState `json:"state" yaml:"state"`
	// Pod container ready
	Ready bool `json:"ready" yaml:"ready"`
	// Pod container restart count
	Restart int `json:"restared" yaml:"restared"`
	// Pod container image meta
	Image PodContainerImage `json:"image" yaml:"image"`
}

type PodContainerImage struct {
	// Pod container image ID
	ID string `json:"id" yaml:"id"`
	// Pod container image name
	Name string `json:"name" yaml:"name"`
}

type PodContainerState struct {
	// Container create state
	Created PodContainerStateCreated `json:"created" yaml:"created"`

	// Container started state
	Started PodContainerStateStarted `json:"started" yaml:"started"`

	// Container stopped state
	Stopped PodContainerStateStopped `json:"stopped" yaml:"stopped"`

	// Container error state
	Error PodContainerStateError `json:"error" yaml:"error"`
}

type PodContainerStateCreated struct {
	Created time.Time `json:"created" yaml:"created"`
}

type PodContainerStateStarted struct {
	Started   bool      `json:"started" yaml:"started"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

type PodContainerStateStopped struct {
	Stopped bool                  `json:"stopped" yaml:"stopped"`
	Exit    PodContainerStateExit `json:"exit" yaml:"exit"`
}

type PodContainerStateError struct {
	Error   bool                  `json:"error" yaml:"error"`
	Message string                `json:"message" yaml:"message"`
	Exit    PodContainerStateExit `json:"exit" yaml:"exit"`
}

type PodContainerStateExit struct {
	Code      int       `json:"code" yaml:"code"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

func (s *PodStatus) SetInitialized() {
	s.State = StateInitialized
	s.Message = EmptyString
}

func (s *PodStatus) SetDestroy () {
	s.State = StateDestroy
}

func (s *PodStatus) SetDestroyed () {
	s.State = StateDestroyed
}

func (s *PodStatus) SetPull () {
	s.State = StatePull
}

func (s *PodStatus) SetProvision () {
	s.State = StateProvision
}

func (s *PodStatus) SetCreated() {
	s.State = StateCreated
	s.Message = EmptyString
}

func (s *PodStatus) SetRunning() {
	s.State = StateRunning
	s.Message = EmptyString
}

func (s *PodStatus) SetStopped() {
	s.State = StateStopped
	s.Message = EmptyString
}

func (s *PodStatus) SetError (err error) {
	s.State = StateError
	s.Message = err.Error()
}

func NewPod() *Pod {
	pod := new(Pod)
	pod.Status = *NewPodStatus()
	return pod
}

func NewPodStatus () *PodStatus {
	status := PodStatus{
		Steps: make(PodSteps, 0),
		Containers: make(map[string]*PodContainer, 0),
	}
	return &status
}

func (p *Pod) SelfLink() string {
	if p.Meta.SelfLink == "" {
		p.Meta.SelfLink = fmt.Sprintf("%s:%s:%s:%s", p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment, p.Meta.Name)
	}
	return p.Meta.SelfLink
}
