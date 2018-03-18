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

package v1

import (
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type PodView struct {
	// PodView meta id
	ID string `json:"id"`
	// PodView Meta
	Meta PodMeta `json:"meta"`
	// PodView state
	State PodState `json:"state"`
	// PodView Spec
	Spec PodSpec `json:"spec"`
	// PodView containers
	Status PodStatus `json:"status"`
}

type PodMeta struct {
	// Meta name
	Name string `json:"name"`
	// Meta name
	Description string `json:"description"`
	// PodView SelfLink
	SelfLink string `json:"self_link"`
	// PodView deployment id
	Deployment string `json:"deployment"`
	// PodView service id
	Namespace string `json:"namespace"`
	// PodView node id
	Node string `json:"node"`
	// PodView status
	Status string `json:"status"`
	// Meta labels
	Labels map[string]string `json:"labels"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type PodState struct {
	// PodView state scheduled
	Scheduled bool `json:"scheduled"`
	// PodView state provision
	Provision bool `json:"provision"`
	// PodView state error
	Error bool `json:"error"`
	// PodView state created
	Created bool `json:"created"`
	// PodView state created
	Pulling bool `json:"pulling"`
	// PodView state started
	Running bool `json:"started"`
	// PodView state stopped
	Stopped bool `json:"stopped"`
	// PodView state destroy
	Destroy bool `json:"destroy"`
}

type PodSpec struct {
	State    PodSpecState    `json:"state"`
	Template PodSpecTemplate `json:"template"`
}

type PodSpecState struct {
	Destroy     bool `json:"destroy"`
	Maintenance bool `json:"maintenance"`
}

type PodSpecTemplate struct {
	// Template Volume
	Volumes types.SpecTemplateVolumes `json:"volumes"`
	// Template main container
	Containers types.SpecTemplateContainers `json:"container"`
	// Termination period
	Termination int `json:"termination"`
}

type PodStatus struct {
	// PodView stage
	Stage string `json:"stage"`
	// PodView state message
	Message string `json:"message"`
	// PodView steps
	Steps PodSteps `json:"steps"`
	// PodView network
	Network PodNetwork `json:"network"`
	// PodView containers
	Containers PodContainers `json:"containers"`
}

type PodContainers []PodContainer

type PodContainer struct {
	// PodView container ID
	ID string `json:"id"`
	// PodView ID
	Pod string `json:"pod"`
	// PodView container name
	Name string `json:"name"`
	// PodView container state
	State PodContainerState `json:"state"`
	// PodView container ready
	Ready bool `json:"ready"`
	// PodView container restart count
	Restart int `json:"restared"`
	// PodView container image meta
	Image PodContainerImage `json:"image"`
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

type PodContainerImage struct {
	// PodView container image ID
	ID string `json:"id"`
	// PodView container image name
	Name string `json:"name"`
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
