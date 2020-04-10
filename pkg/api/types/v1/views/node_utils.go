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

import (
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/resource"
)

type NodeView struct{}

func (nv *NodeView) New(obj *models.Node) *Node {
	n := Node{}
	n.Meta = nv.ToNodeMeta(obj.Meta)
	n.Status = nv.ToNodeStatus(obj.Status)
	n.Spec = nv.ToNodeSpec(obj.Spec)
	return &n
}

func (nv *NodeView) ToNodeMeta(meta models.NodeMeta) NodeMeta {
	nm := NodeMeta{}
	nm.Name = meta.Name
	nm.SelfLink = meta.SelfLink.String()
	nm.Created = meta.Created
	nm.Updated = meta.Updated
	nm.Hostname = meta.Hostname
	nm.OSName = meta.OSName
	nm.OSType = meta.OSType
	nm.Architecture = meta.Architecture
	nm.IP.External = meta.ExternalIP
	nm.IP.Internal = meta.InternalIP
	nm.CIDR = meta.CIDR
	return nm
}

func (nv *NodeView) ToNodeStatus(status models.NodeStatus) NodeStatus {
	ns := NodeStatus{}

	ns.Online = status.Online

	ns.Capacity.Containers = status.Capacity.Containers
	ns.Capacity.Pods = status.Capacity.Pods
	ns.Capacity.Memory = resource.EncodeMemoryResource(status.Capacity.RAM)
	ns.Capacity.CPU = resource.EncodeCpuResource(status.Capacity.CPU)
	ns.Capacity.Storage = resource.EncodeMemoryResource(status.Capacity.Storage)

	ns.Allocated.Containers = status.Allocated.Containers
	ns.Allocated.Pods = status.Allocated.Pods
	ns.Allocated.Memory = resource.EncodeMemoryResource(status.Allocated.RAM)
	ns.Allocated.CPU = resource.EncodeCpuResource(status.Allocated.CPU)
	ns.Allocated.Storage = resource.EncodeMemoryResource(status.Allocated.Storage)

	ns.State.CNI.Type = status.State.CNI.Type
	ns.State.CNI.State = status.State.CNI.State
	ns.State.CNI.Version = status.State.CNI.Version
	ns.State.CNI.Message = status.State.CNI.Message

	ns.State.CPI.Type = status.State.CPI.Type
	ns.State.CPI.State = status.State.CPI.State
	ns.State.CPI.Version = status.State.CPI.Version
	ns.State.CPI.Message = status.State.CPI.Message

	ns.State.CRI.Type = status.State.CRI.Type
	ns.State.CRI.State = status.State.CRI.State
	ns.State.CRI.Version = status.State.CRI.Version
	ns.State.CRI.Message = status.State.CRI.Message

	ns.State.CSI.Type = status.State.CSI.Type
	ns.State.CSI.State = status.State.CSI.State
	ns.State.CSI.Version = status.State.CSI.Version
	ns.State.CSI.Message = status.State.CSI.Message

	return ns
}

func (nv *NodeView) ToNodeSpec(spec models.NodeSpec) NodeSpec {
	ns := NodeSpec{}
	ns.Security.TLS = spec.Security.TLS
	return ns
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeManifest) Decode() *models.NodeManifest {

	manifest := models.NodeManifest{
		Secrets:   make(map[string]*models.SecretManifest, 0),
		Configs:   make(map[string]*models.ConfigManifest, 0),
		Network:   make(map[string]*models.SubnetManifest, 0),
		Pods:      make(map[string]*models.PodManifest, 0),
		Volumes:   make(map[string]*models.VolumeManifest, 0),
		Endpoints: make(map[string]*models.EndpointManifest, 0),
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Resolvers = make(map[string]*models.ResolverManifest, 0)
	manifest.Exporter = new(models.ExporterManifest)

	for i, s := range obj.Discovery {
		manifest.Resolvers[i] = s
	}

	manifest.Exporter = obj.Exporter

	for i, s := range obj.Network {
		manifest.Network[i] = s
	}

	for i, s := range obj.Pods {
		manifest.Pods[i] = s
	}

	for i, s := range obj.Volumes {
		manifest.Volumes[i] = s
	}

	for i, s := range obj.Endpoints {
		manifest.Endpoints[i] = s
	}

	for i, s := range obj.Configs {
		manifest.Configs[i] = s
	}

	for i, s := range obj.Secrets {
		manifest.Secrets[i] = s
	}

	return &manifest
}

func (nv *NodeView) NewList(obj *models.NodeList) *NodeList {
	if obj == nil {
		return nil
	}
	nodes := make(NodeList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		nodes[nn.Meta.Name] = nn
	}

	return &nodes
}

func (nv *NodeView) NewManifest(obj *models.NodeManifest) *NodeManifest {

	manifest := NodeManifest{
		Configs:   make(map[string]*models.ConfigManifest, 0),
		Secrets:   make(map[string]*models.SecretManifest, 0),
		Network:   make(map[string]*models.SubnetManifest, 0),
		Pods:      make(map[string]*models.PodManifest, 0),
		Volumes:   make(map[string]*models.VolumeManifest, 0),
		Endpoints: make(map[string]*models.EndpointManifest, 0),
	}

	if obj == nil {
		return nil
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Discovery = obj.Resolvers
	manifest.Exporter = obj.Exporter

	manifest.Configs = obj.Configs
	manifest.Secrets = obj.Secrets
	manifest.Network = obj.Network
	manifest.Pods = obj.Pods
	manifest.Volumes = obj.Volumes
	manifest.Endpoints = obj.Endpoints

	return &manifest
}

func (obj *NodeManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
