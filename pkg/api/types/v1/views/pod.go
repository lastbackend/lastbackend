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

package views

import (
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type Pod struct {
	// Pod meta id
	ID string `json:"id"`
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod Spec
	Spec PodSpec `json:"spec"`
	// Pod containers
	Status PodStatus `json:"status"`
}

type PodList map[string]*Pod

type PodMeta struct {
	// Meta name
	Name string `json:"name"`
	// Meta name
	Description string `json:"description"`
	// Pod SelfLink
	SelfLink string `json:"self_link"`
	// Pod deployment id
	Deployment string `json:"deployment"`
	// Pod service id
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
	Volumes types.SpecTemplateVolumeList `json:"volumes"`
	// Template main container
	Containers types.SpecTemplateContainers `json:"container"`
	// Termination period
	Termination int `json:"termination"`
}

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

type PodContainers []PodContainer

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
	// Pod container image ID
	ID string `json:"id"`
	// Pod container image name
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
