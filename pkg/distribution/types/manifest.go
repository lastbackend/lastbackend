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

package types

type NodeManifest struct {
	Meta      NodeManifestMeta             `json:"meta"`
	Resolvers map[string]*ResolverManifest `json:"resolvers"`
	Secrets   map[string]*SecretManifest   `json:"secrets"`
	Configs   map[string]*ConfigManifest   `json:"configs"`
	Endpoints map[string]*EndpointManifest `json:"endpoint"`
	Network   map[string]*SubnetManifest   `json:"network"`
	Pods      map[string]*PodManifest      `json:"pods"`
	Volumes   map[string]*VolumeManifest   `json:"volumes"`
}

type NodeManifestMeta struct {
	Initial bool `json:"initial"`
}

type ResolverManifest struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type IngressManifest struct {
	Meta      IngressManifestMeta          `json:"meta"`
	Resolvers map[string]*ResolverManifest `json:"resolvers"`
	Routes    map[string]*RouteManifest    `json:"routes"`
	Endpoints map[string]*EndpointManifest `json:"endpoints"`
	Network   map[string]*SubnetManifest   `json:"network"`
}

type IngressManifestMeta struct {
	Initial bool `json:"initial"`
}

type DiscoveryManifest struct {
	Meta    DiscoveryManifestMeta      `json:"meta"`
	Network map[string]*SubnetManifest `json:"network"`
}

type DiscoveryManifestMeta struct {
	Initial bool `json:"initial"`
}

type PodManifest PodSpec

type PodManifestList struct {
	System
	Items []*PodManifest
}

type PodManifestMap struct {
	System
	Items map[string]*PodManifest
}

type VolumeManifest VolumeSpec

type VolumeManifestList struct {
	System
	Items []*VolumeManifest
}

type VolumeManifestMap struct {
	System
	Items map[string]*VolumeManifest
}

type SubnetManifest struct {
	System
	SubnetSpec
}

type SubnetManifestList struct {
	System
	Items []*SubnetManifest
}

type SubnetManifestMap struct {
	System
	Items map[string]*SubnetManifest
}

type EndpointManifest struct {
	System
	EndpointSpec `json:",inline"`
	Upstreams    []string `json:"upstreams"`
}

type EndpointManifestList struct {
	System
	Items []*EndpointManifest
}

type EndpointManifestMap struct {
	System
	Items map[string]*EndpointManifest
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
