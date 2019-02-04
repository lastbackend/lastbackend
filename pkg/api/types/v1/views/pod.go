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

package views

import (
	"time"
)

// swagger:model views_pod
type Pod struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod Spec
	Spec PodSpec `json:"spec"`
	// Pod containers
	Status PodStatus `json:"status"`
}

// swagger:ignore
// PodList is a map of pods
// swagger:model views_pod_list
type PodList map[string]Pod

// PodMeta is a meta of pod
// swagger:model views_pod_meta
type PodMeta struct {
	// Meta name
	Name string `json:"name"`
	// Meta description
	Description string `json:"description"`
	// Pod SelfLink
	SelfLink string `json:"self_link"`
	// Pod parent
	Parent PodMetaParent `json:"parent"`
	// Pod namespace
	Namespace string `json:"namespace"`
	// Pod node id
	Node string `json:"node"`
	// Pod status
	Status string `json:"status"`
	// Meta labels
	Labels map[string]string `json:"labels"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type PodMetaParent struct {
	Kind     string `json:"kind"`
	SelfLink string `json:"self_link"`
}

// PodSpec is a spec of pod
// swagger:model views_pod_spec
type PodSpec struct {
	State    PodSpecState         `json:"state"`
	Template ManifestSpecTemplate `json:"template"`
}

// PodSpecState is a state of pod spec
// swagger:model views_pod_spec_state
type PodSpecState struct {
	Destroy     bool `json:"destroy"`
	Maintenance bool `json:"maintenance"`
}

// PodStatus is a status of pod
// swagger:model views_pod_status
type PodStatus struct {
	// Pod state
	State string `json:"state"`
	// Pod state message
	Message string `json:"message"`
	// Pod steps
	Steps PodSteps `json:"steps"`
	// Pod network
	Network PodNetwork `json:"network"`
	// Pod containers
	Containers PodContainers `json:"containers"`
}

// PodContainers is a list of pod containers
// swagger:model views_pod_container_list
type PodContainers []PodContainer

// PodContainer is a container of the pod
// swagger:model views_pod_container
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

// PodContainerState is a state of pod container
// swagger:model views_pod_container_state
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

// swagger:ignore
// PodContainerStateCreated represents creation time of the pod container
// swagger:model views_pod_container_state_created
type PodContainerStateCreated struct {
	Created time.Time `json:"created"`
}

// swagger:ignore
// PodContainerStateStarted represents time when pod container was started
// swagger:model views_pod_container_state_started
type PodContainerStateStarted struct {
	// was the container started
	Started   bool      `json:"started"`
	Timestamp time.Time `json:"timestamp"`
}

// swagger:ignore
// PodContainerStateStopped shows if pod container was stopped
// swagger:model views_pod_container_state_stopped
type PodContainerStateStopped struct {
	// was the container stopped
	Stopped bool                  `json:"stopped"`
	Exit    PodContainerStateExit `json:"exit"`
}

// swagger:ignore
// PodContainerStateError shows if pod container got error
// swagger:model views_pod_container_state_error
type PodContainerStateError struct {
	// was error happened
	Error   bool                  `json:"error"`
	Message string                `json:"message"`
	Exit    PodContainerStateExit `json:"exit"`
}

// swagger:ignore
// PodContainerStateExit represents an exit status of pod container after stop or error
// swagger:model views_pod_container_state_exit
type PodContainerStateExit struct {
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// PodContainerImage is an image of pod container
// swagger:model views_pod_container_state
type PodContainerImage struct {
	// Pod container image ID
	ID string `json:"id"`
	// Pod container image name
	Name string `json:"name"`
}

// PodSteps is a map of pod steps
// swagger:model views_pod_step_map
type PodSteps map[string]PodStep

// swagger:model views_pod_step
type PodStep struct {
	// Pod step ready
	Ready bool `json:"ready"`
	// Pod step timestamp
	Timestamp time.Time `json:"timestamp"`
}

// swagger:model views_pod_network
type PodNetwork struct {
	// Pod host IP
	HostIP string `json:"host_ip"`
	// Pod IP
	PodIP string `json:"pod_ip"`
}
