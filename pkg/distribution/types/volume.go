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

import "fmt"

// swagger:ignore
// swagger:model types_volume
type Volume struct {
	Runtime
	// Volume meta
	Meta VolumeMeta `json:"meta" yaml:"meta"`
	// Volume spec
	Spec VolumeSpec `json:"spec" yaml:"spec"`
	// Volume status
	Status VolumeStatus `json:"status" yaml:"status"`
}

// swagger:ignore
// swagger:model types_volume_map
type VolumeMap struct {
	Runtime
	Items map[string]*Volume
}

// swagger:ignore
// swagger:model types_volume_list
type VolumeList struct {
	Runtime
	Items []*Volume
}

// swagger:ignore
// swagger:model types_volume_meta
type VolumeMeta struct {
	Meta
	Namespace string `json:"namespace"`
}

// swagger:model types_volume_spec
type VolumeSpec struct {
	State VolumeSpecState `json:"state"`
}

// swagger:model types_volume_spec_state
type VolumeSpecState struct {
	Destroy bool `json:"destroy"`
}

// swagger:ignore
// swagger:model types_volume_create
type VolumeCreateOptions struct {
}

// swagger:ignore
// swagger:model types_volume_update
type VolumeUpdateOptions struct {
}

// swagger:ignore
// swagger:model types_volume_status
type VolumeStatus struct {
	// volume status state
	State string `json:"state" yaml:"state"`
	// volume status message
	Message string `json:"message" yaml:"message"`
}

func (v *Volume) SelfLink() string {
	if v.Meta.SelfLink == "" {
		v.Meta.SelfLink = v.CreateSelfLink(v.Meta.Namespace, v.Meta.Name)
	}
	return v.Meta.SelfLink
}

func (v *Volume) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}
