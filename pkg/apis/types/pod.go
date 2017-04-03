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
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type PodList []*Pod

func (pl *PodList) ToJson() []byte {
	j, _ := json.Marshal(pl)
	return j
}

type PodMap struct {
	Items map[PodID]*Pod `json:"pods"`
}

type Pod struct {
	lock sync.RWMutex

	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod provision flag
	Policy PodPolicy `json:"provision"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Containers status info
	Containers map[ContainerID]Container `json:"containers"`
	// Secrets
	Secrets map[string]PodSecret `json:"secrets"`
	// Container created time
	Created time.Time `json:"created"`
	// Container updated time
	Updated time.Time `json:"updated"`
}

func (p *Pod) ID() PodID {
	return p.Meta.ID
}

func (p *Pod) AddContainer(c Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
}

func (p *Pod) SetContainer(c Container) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Containers[c.ID] = c
}

func (p *Pod) DelContainer(ID ContainerID) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.Containers, ID)
}

type PodSpec struct {
	Containers []ContainerSpec
	Volumes    []VolumesSpec
	Images     map[string]ImageSpec
}

type PodMeta struct {
	// Pod ID
	ID PodID
	// Pod owner
	Owner string
	// Pod project
	Project string
	// Pod service
	Service string
}

type PodPolicy struct {
	// Pull image flag
	PullImage bool
	// Restart containers flag
	Restart bool
}

type PodSecret struct {
}

type PodID uuid.UUID

func (s *PodSpec) Equal(spec PodSpec) bool {

	ohash, _ := json.Marshal(s)
	nhash, _ := json.Marshal(spec)

	if string(ohash) == string(nhash) {
		return true
	}

	return false
}

func (s *PodSpec) NotEqual(spec PodSpec) bool {
	return !s.Equal(spec)
}
