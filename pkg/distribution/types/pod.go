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
	"sync"
	"time"
	"fmt"
)

const PodStepInitialized = "initialized"
const PodStepScheduled = "scheduled"
const PodStepPull = "pull"
const PodStepDestroyed = "destroyed"
const PodStepReady = "ready"

const PodStageInitialized = PodStepInitialized
const PodStageScheduled = PodStepScheduled
const PodStagePull = PodStepPull

const PodStageStarting = "starting"
const PodStageRunning = "running"
const PodStageStopped = "stopped"
const PodStageError = "error"
const PodStageDestroy = "destroy"
const PodStageDestroyed = "destroyed"

type Pod struct {
	// Lock map
	lock sync.RWMutex
	// Pod Meta
	Meta PodMeta `json:"meta" yaml:"meta"`
	// Pod state
	State PodState `json:"state" yaml:"state"`
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
	Template SpecTemplate `json:"spec" yaml:"spec"`
}

type PodState struct {
	// Pod state ready
	Ready bool `json:"ready" yaml:"ready"`
	// Pod state scheduled
	Scheduled bool `json:"scheduled" yaml:"scheduled"`
	// Pod state provision
	Provision bool `json:"provision" yaml:"provision"`
	// Pod state error
	Error bool `json:"error" yaml:"error"`
	// Pod state created
	Created bool `json:"created" yaml:"created"`
	// Pod state created
	Pulling bool `json:"pulling" yaml:"pulling"`
	// Pod state started
	Running bool `json:"started" yaml:"started"`
	// Pod state stopped
	Stopped bool `json:"stopped" yaml:"stopped"`
	// Pod state destroy
	Destroy bool `json:"destroy" yaml:"destroy"`
}

type PodStatus struct {
	// Pod stage
	Stage string `json:"stage" yaml:"stage"`
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

func NewPod() *Pod {
	pod := new(Pod)
	pod.Status.Steps = make(PodSteps, 0)
	pod.Status.Containers = make(map[string]*PodContainer, 0)
	return pod
}

func (p *Pod) SelfLink() string {
	if p.Meta.SelfLink == "" {
		p.Meta.SelfLink = fmt.Sprintf("%s:%s:%s:%s", p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment, p.Meta.Name)
	}
	return p.Meta.SelfLink
}

func (p *Pod) MarkAsInitialized() {
	p.State.Pulling = false
	p.State.Created = true
	p.State.Running = false
	p.State.Stopped = true
	p.State.Provision = true
	p.State.Destroy = false
	p.State.Error = false
	p.Status.Stage = PodStageInitialized
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsPull() {
	p.State.Pulling = true
	p.State.Created = false
	p.State.Running = false
	p.State.Stopped = true
	p.State.Provision = true
	p.State.Destroy = false
	p.State.Error = false
	p.Status.Stage = PodStepPull
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsDestroyed() {
	p.State.Pulling = false
	p.State.Created = false
	p.State.Running = false
	p.State.Stopped = true
	p.State.Provision = false
	p.State.Destroy = true
	p.State.Error = false
	p.Status.Stage = PodStageDestroyed
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsStarting() {
	p.State.Pulling = false
	p.State.Created = false
	p.State.Running = false
	p.State.Stopped = false
	p.State.Provision = true
	p.State.Error = false
	p.Status.Stage = PodStageStarting
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsRunning() {
	p.State.Pulling = false
	p.State.Created = false
	p.State.Running = true
	p.State.Stopped = false
	p.State.Provision = false
	p.State.Error = false
	p.Status.Stage = PodStageRunning
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsStopped() {
	p.State.Pulling = false
	p.State.Created = false
	p.State.Running = false
	p.State.Stopped = true
	p.State.Provision = false
	p.State.Error = false
	p.Status.Stage = PodStageStopped
	p.Status.Message = EmptyString
}

func (p *Pod) MarkAsError(err error) {
	p.State.Pulling = false
	p.State.Created = false
	p.State.Running = false
	p.State.Stopped = true
	p.State.Provision = false
	p.State.Destroy = false
	p.State.Error = true
	p.Status.Stage = PodStageError
	p.Status.Message = err.Error()
}
