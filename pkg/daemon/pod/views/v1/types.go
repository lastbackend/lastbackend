package v1

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/container/views/v1"
	"time"
)

type Pod struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Pod state
	State PodState `json:"state"`
}

type PodState struct {
	// Pod current state
	State string `json:"state"`
	// Pod current status
	Status string `json:"status"`
}

type PodMeta struct {
	// Meta id
	ID string `json:"id"`
	// Meta labels
	Labels map[string]string `json:"lables"`
	// Pod owner
	Owner string `json:"owner"`
	// Pod project
	Project string `json:"project"`
	// Pod service
	Service string `json:"service"`
	// Current Spec ID
	Spec string `json:"spec"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type PodSpec struct {
	// Provision ID
	ID string `json:"id"`
	// Provision state
	State string `json:"state"`
	// Provision status
	Status string `json:"status"`

	// Containers spec for pod
	Containers []v1.ContainerSpec `json:"containers"`

	// Provision create time
	Created time.Time `json:"created"`
	// Provision update time
	Updated time.Time `json:"updated"`
}
