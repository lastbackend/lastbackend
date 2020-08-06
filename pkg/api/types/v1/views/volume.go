//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

import "time"

type Volume struct {
	Meta   VolumeMeta   `json:"meta"`
	Spec   VolumeSpec   `json:"spec"`
	Status VolumeStatus `json:"status"`
}

type VolumeMeta struct {
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	Description string    `json:"description"`
	SelfLink    string    `json:"self_link"`
	Updated     time.Time `json:"updated"`
	Created     time.Time `json:"created"`
}

type VolumeSpec struct {
	Selector   VolumeSpecSelector `json:"selector"`
	State      VolumeSpecState    `json:"state"`
	Type       string             `json:"type"`
	HostPath   string             `json:"path"`
	AccessMode string             `json:"mode"`
	Capacity   VolumeSpecCapacity `json:"capacity"`
}

type VolumeSpecSelector struct {
	Node   string            `json:"node"`
	Labels map[string]string `json:"labels"`
}

type VolumeSpecCapacity struct {
	Storage string `json:"storage"`
}

type VolumeSpecState struct {
	Destroy bool `json:"destroy"`
}

type VolumeStatus struct {
	// volume state
	State string `json:"state" yaml:"state"`
	// volume status
	Status VolumeState `json:"status" yaml:"status"`
	// volume status message
	Message string `json:"message" yaml:"message"`
}

// swagger:ignore
// swagger:model types_volume_status
type VolumeState struct {
	Type string `json:"type" yaml:"type"`
	// Volume root path
	Path string `json:"path" yaml:"path"`
	// Volume state ready
	Ready bool `json:"ready" yaml:"ready"`
}

type VolumeList []*Volume
