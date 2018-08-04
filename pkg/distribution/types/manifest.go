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

type NodeManifest struct {
	Endpoints map[string]*EndpointManifest `json:"endpoint"`
	Network   map[string]*NetworkManifest  `json:"network"`
	Pods      map[string]*PodManifest      `json:"pods"`
	Volumes   map[string]*VolumeManifest   `json:"volumes"`
}

type PodManifest PodSpec

type PodManifestList struct {
	Runtime
	Items  []*PodManifest
}

type PodManifestMap struct {
	Runtime
	Items  map[string]*PodManifest
}

type VolumeManifest VolumeSpec

type VolumeManifestList struct {
	Runtime
	Items  []*VolumeManifest
}

type VolumeManifestMap struct {
	Runtime
	Items  map[string]*VolumeManifest
}

type NetworkManifest struct {
	NetworkSpec
}

type EndpointManifest struct {
	EndpointSpec
}


func NewPodManifestList () *PodManifestList {
	dm := new(PodManifestList)
	dm.Items = make([]*PodManifest, 0)
	return dm
}

func NewPodManifestMap () *PodManifestMap {
	dm := new(PodManifestMap)
	dm.Items = make(map[string]*PodManifest)
	return dm
}

func NewVolumeManifestList () *VolumeManifestList {
	dm := new(VolumeManifestList)
	dm.Items = make([]*VolumeManifest, 0)
	return dm
}

func NewVolumeManifestMap () *VolumeManifestMap {
	dm := new(VolumeManifestMap)
	dm.Items = make(map[string]*VolumeManifest)
	return dm
}