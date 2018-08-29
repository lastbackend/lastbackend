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

type Image struct {
	Meta   ImageMeta
	Status ImageStatus
	Spec   ImageSpec
}

type ImageMeta struct {
	ID   string `json:"id"`
	Hash string `json:"hash"`
	Name string `json:"name"`
}

type ImageStatus struct {
	State       string `json:"state"`
	Size        int64  `json:"size"`
	VirtualSize int64  `json:"virtual_size"`
}

type ImageSpec struct {
	// Name full name
	Name string `json:"name"`
	// Secret name for pulling
	Secret string `json:"auth"`
}

type ImageManifest struct {
	Name   string `json:"name" yaml:"name"`
	Auth   string `json:"auth" yaml:"auth"`
	Policy string `json:"policy" yaml:"policy"`
}

func (i *Image) SelfLink() string {
	return fmt.Sprintf("%s:%s", i.Meta.Name)
}
