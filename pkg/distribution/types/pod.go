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
	Meta PodMeta `json:"meta"`
	// Pod state
	State PodState `json:"state"`
	// Pod Spec
	Spec SpecTemplate `json:"spec"`
	// Containers status info
	Status PodStatus `json:"status"`
}

type PodMeta struct {
	Meta
	// Pod SelfLink
	SelfLink string `json:"self_link"`
	// Pod deployment name
	Deployment string `json:"deployment"`
	// Pod service name
	Namespace string `json:"namespace"`
	// Pod service name
	Service string `json:"service"`
	// Pod node
	Node string `json:"node"`
	// Pod status
	Status string `json:"status"`
	// Pod endpoint
	Endpoint string `json:"endpoint"`
}

type PodState struct {
	// Pod state scheduled
	Scheduled bool `json:"scheduled"`
	// Pod state provision
	Provision bool `json:"provision"`
	// Pod state error
	Error bool `json:"error"`
	// Pod state created
	Created bool `json:"created"`
	// Pod state created
	Pulling bool `json:"pulling"`
	// Pod state started
	Running bool `json:"started"`
	// Pod state stopped
	Stopped bool `json:"stopped"`
	// Pod state destroy
	Destroy bool `json:"destroy"`
}

type PodStatus struct {
	// Pod stage
	Stage string `json:"stage"`
	// Pod state message
	Message string `json:"message"`
	// Pod steps
	Steps PodSteps `json:"steps"`
	// Pod network
	Network PodNetwork `json:"network"`
	// Pod containers
	Containers map[string]*PodContainer `json:"containers"`
}

type PodSteps map[string]PodStep

type PodStep struct {
	// Pod step ready
	Ready bool `json:"ready"`
	// Pod step timestamp
	Timestamp time.Time `json:"timestamp"`
}

type PodNetwork struct {
	// Pod host IP
	HostIP string `json:"host_ip"`
	// Pod IP
	PodIP string `json:"pod_ip"`
}

type PodContainer struct {
	// Pod container ID
	ID string `json:"id"`
	// Pod ID
	Pod string `json:"pod"`
	// Pod container name
	Name string `json:"name"`
	// Pod container state
	State PodContainerState `json:"state"`
	// Pod container ready
	Ready bool `json:"ready"`
	// Pod container restart count
	Restart int `json:"restared"`
	// Pod container image meta
	Image PodContainerImage `json:"image"`
}

type PodContainerImage struct {
	// Pod container image ID
	ID string `json:"id"`
	// Pod container image name
	Name string `json:"name"`
}

type PodContainerState struct {
	// Container create state
	Created PodContainerStateCreated `json:"created"`

	// Container started state
	Started PodContainerStateStarted `json:"started"`

	// Container stopped state
	Stopped PodContainerStateStopped `json:"stopped"`

	// Container error state
	Error PodContainerStateError `json:"error"`
}

type PodContainerStateCreated struct {
	Created time.Time `json:"created"`
}

type PodContainerStateStarted struct {
	Started   bool      `json:"started"`
	Timestamp time.Time `json:"timestamp"`
}

type PodContainerStateStopped struct {
	Stopped bool                  `json:"stopped"`
	Exit    PodContainerStateExit `json:"exit"`
}

type PodContainerStateError struct {
	Error   bool                  `json:"error"`
	Message string                `json:"message"`
	Exit    PodContainerStateExit `json:"exit"`
}

type PodContainerStateExit struct {
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

func NewPod() *Pod {
	pod := new(Pod)
	pod.Status.Steps = make(PodSteps, 0)
	pod.Status.Containers = make(map[string]*PodContainer, 0)
	return pod
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
