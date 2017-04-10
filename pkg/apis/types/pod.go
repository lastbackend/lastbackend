//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"sync"
	"time"
)

type PodList []*Pod

func (pl *PodList) ToJson() []byte {
	j, _ := json.Marshal(pl)
	return j
}

type PodMap struct {
	Items map[string]*Pod `json:"pods"`
}

type Pod struct {
	lock sync.RWMutex

	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Pod state
	State PodState `json:"state"`
	// Containers status info
	Containers map[string]*Container `json:"containers"`
	// Secrets
	Secrets map[string]*PodSecret `json:"secrets"`
}

type PodMeta struct {
	Meta

	// Pod owner
	Owner string `json:"owner"`
	// Pod project
	Project string `json:"project"`
	// Pod service
	Service string `json:"service"`
	// Current Spec ID
	Spec string `json:"spec"`
}

type PodSpec struct {
	// Provision ID
	ID string `json:"id"`
	// Provision state
	State string `json:"state"`
	// Provision status
	Status string `json:"status"`

	// Containers spec for pod
	Containers []*ContainerSpec

	// Provision create time
	Created time.Time `json:"created"`
	// Provision update time
	Updated time.Time `json:"updated"`
}

type PodState struct {
	// Pod current state
	State string `json:"state"`
	// Pod current status
	Status string `json:"status"`
	// Pod provision
}

type PodSecret struct {
}

func (p *Pod) AddContainer(c *Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
	p.CalculateState()
}

func (p *Pod) SetContainer(c *Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
	p.CalculateState()
}

func (p *Pod) DelContainer(ID string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.Containers, ID)
	p.CalculateState()
}

func (p *Pod) CalculateState() {
	for _, c := range p.Containers {

		if c.State == "exited" && p.State.Status == "" {
			p.State.State = "stopped"
			continue
		}

		if c.State == "exited" && p.State.Status != "stopped" {
			p.State.State = "warning"
			continue
		}

		if c.State == "running" && p.State.Status != "running" {
			p.State.State = "warning"
			continue
		}

		if p.State.State == "" {
			p.State.State = c.State
			continue
		}

		if c.State == p.State.State {
			continue
		}

		if p.State.State == "error" {
			continue
		}

		if c.State == "error" {
			p.State.State = c.State
			p.State.Status = c.Status
			continue
		}

	}
}

const PodStateStarted = "started"
const PodStateRestarted = "restarted"
const PodStateStopped = "stopped"

func NewPod() *Pod {
	return &Pod{
		Containers: make(map[string]*Container),
		Secrets:    make(map[string]*PodSecret),
	}
}
