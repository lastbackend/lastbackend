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

import "time"

type NodeManifest struct {
	Secrets   map[string]*SecretManifest   `json:"secrets"`
	Endpoints map[string]*EndpointManifest `json:"endpoint"`
	Network   map[string]*SubnetManifest   `json:"network"`
	Pods      map[string]*PodManifest      `json:"pods"`
	Volumes   map[string]*VolumeManifest   `json:"volumes"`
}

type PodManifest PodSpec

type PodManifestList struct {
	Runtime
	Items []*PodManifest
}

type PodManifestMap struct {
	Runtime
	Items map[string]*PodManifest
}

type VolumeManifest VolumeSpec

type VolumeManifestList struct {
	Runtime
	Items []*VolumeManifest
}

type VolumeManifestMap struct {
	Runtime
	Items map[string]*VolumeManifest
}

type SubnetManifest struct {
	Runtime
	SubnetSpec
}

type SubnetManifestList struct {
	Runtime
	Items []*SubnetManifest
}

type SubnetManifestMap struct {
	Runtime
	Items map[string]*SubnetManifest
}

type EndpointManifest struct {
	Runtime
	EndpointSpec `json:",inline"`
	Upstreams    []string `json:"upstreams"`
}

type EndpointManifestList struct {
	Runtime
	Items []*EndpointManifest
}

type EndpointManifestMap struct {
	Runtime
	Items map[string]*EndpointManifest
}

type SecretManifest struct {
	Runtime
	State   string    `json:"state"`
	Kind    string    `json:"kind"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type SecretManifestList struct {
	Runtime
	Items []*SecretManifest
}

type SecretManifestMap struct {
	Runtime
	Items map[string]*SecretManifest
}

func NewPodManifestList() *PodManifestList {
	dm := new(PodManifestList)
	dm.Items = make([]*PodManifest, 0)
	return dm
}

func NewPodManifestMap() *PodManifestMap {
	dm := new(PodManifestMap)
	dm.Items = make(map[string]*PodManifest)
	return dm
}

func NewVolumeManifestList() *VolumeManifestList {
	dm := new(VolumeManifestList)
	dm.Items = make([]*VolumeManifest, 0)
	return dm
}

func NewVolumeManifestMap() *VolumeManifestMap {
	dm := new(VolumeManifestMap)
	dm.Items = make(map[string]*VolumeManifest)
	return dm
}

func NewSubnetManifestList() *SubnetManifestList {
	dm := new(SubnetManifestList)
	dm.Items = make([]*SubnetManifest, 0)
	return dm
}

func NewSubnetManifestMap() *SubnetManifestMap {
	dm := new(SubnetManifestMap)
	dm.Items = make(map[string]*SubnetManifest)
	return dm
}

func NewEndpointManifestList() *EndpointManifestList {
	dm := new(EndpointManifestList)
	dm.Items = make([]*EndpointManifest, 0)
	return dm
}

func NewEndpointManifestMap() *EndpointManifestMap {
	dm := new(EndpointManifestMap)
	dm.Items = make(map[string]*EndpointManifest)
	return dm
}

func NewSecretManifestList() *SecretManifestList {
	dm := new(SecretManifestList)
	dm.Items = make([]*SecretManifest, 0)
	return dm
}

func NewSecretManifestMap() *SecretManifestMap {
	dm := new(SecretManifestMap)
	dm.Items = make(map[string]*SecretManifest)
	return dm
}
