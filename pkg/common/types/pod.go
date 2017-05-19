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
	"sync"
	"time"
)

type Pod struct {
	lock sync.RWMutex

	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod state
	State PodState `json:"state"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Containers status info
	Containers map[string]*Container `json:"containers"`
	// Secrets
	Secrets map[string]*PodSecret `json:"secrets"`
}

type PodNodeSpec struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod state
	State PodState `json:"state"`
	// Pod spec
	Spec PodSpec `json:"spec"`
}

type PodNodeState struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Containers status info
	Containers map[string]*Container `json:"containers"`
}

type PodMeta struct {
	Meta
	// Pod hostname
	Hostname string `json:"hostname"`
}

type PodSpec struct {
	// Provision ID
	ID string `json:"id"`
	// Provision state
	State string `json:"state"`
	// Provision status
	Status string `json:"status,omitempty"`

	// Containers spec for pod
	Containers map[string]*ContainerSpec `json:"containers"`

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
	// Pod Provision state
	Provision bool `json:"provision"`
	// Pod ready state
	Ready bool `json:"created"`
}

type PodSecret struct{}

type PodCRIMeta struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type PodNetwork struct {
	Interface string   `json:"interface,omitempty"`
	IP        []string `json:"ip,omitempty"`
}

func (p *Pod) AddContainer(c *Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
}

func (p *Pod) SetContainer(c *Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
}

func (p *Pod) DelContainer(ID string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.Containers, ID)
}

func (p *Pod) UpdateState() {

	p.State.Status = ""

	for _, c := range p.Containers {
		switch c.State {
		case StateStarted:
			switch p.State.State {
			case "":
				p.State.State = StateStarted
			case StateStarted:
				p.State.State = StateStarted
			case StateError:
				p.State.State = StateError
			case StateDestroy:
				p.State.State = StateDestroy
			default:
				p.State.State = StateWarning
			}
		case StateRunning:
			switch p.State.State {
			case "":
				p.State.State = StateStarted
			case StateStarted:
				p.State.State = StateStarted
			case StateError:
				p.State.State = StateError
			case StateDestroy:
				p.State.State = StateDestroy
			default:
				p.State.State = StateWarning
			}
		case StateStopped:
			switch p.State.State {
			case "":
				p.State.State = StateStopped
			case StateStopped:
				p.State.State = StateStopped
			case StateDestroy:
				p.State.State = StateDestroy
			case StateError:
				p.State.State = StateError
			default:
				p.State.State = StateWarning
			}
		case StateRestarted:
			switch p.State.State {
			case "":
				p.State.State = StateStarted
			case StateStarted:
				p.State.State = StateStarted
			case StateDestroy:
				p.State.State = StateDestroy
			case StateError:
				p.State.State = StateError
			default:
				p.State.State = StateWarning
			}
		case StateDestroyed:
			p.State.State = StateDestroy
		case StateExited:
			switch p.State.State {
			case "":
				p.State.State = StateStopped
			case StateStopped:
				p.State.State = StateStopped
			case StateDestroy:
				p.State.State = StateDestroy
			case StateError:
				p.State.State = StateError
			default:
				p.State.State = StateWarning
			}
		case StateError:
			p.State.State = StateError
			p.State.Status = c.Status
		}
	}

	if len(p.Containers) == 0 && p.Spec.State == StateDestroy {
		p.State.State = StateDestroyed
	}

}

func NewPod() *Pod {
	return &Pod{
		Containers: make(map[string]*Container),
		Secrets:    make(map[string]*PodSecret),
	}
}
