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

type VolumeList []Volume

type Volume struct {
	// Volume meta
	Meta VolumeMeta `json:"meta" yaml:"meta"`
	// Volume stat
	State VolumeState `json:"stat" yaml:"stat"`
	// Volume spec
	Spec VolumeSpec `json:"spec" yaml:"spec"`
}

type VolumeMeta struct {
	Meta
	Namespace string `json:"namespace"`
}

type VolumeState struct {
	Provision bool `json:"provision"`
	Ready bool `json:"ready"`
}

type VolumeSpec struct {

}

func (v *Volume) SelfLink() string {
	if v.Meta.SelfLink == "" {
		v.Meta.SelfLink = fmt.Sprintf("%s:%s", v.Meta.Namespace, v.Meta.Name)
	}
	return v.Meta.SelfLink
}