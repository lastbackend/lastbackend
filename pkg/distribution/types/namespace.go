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

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
)

const (
	DEFAULT_NAMESPACE string = "default"
)

// swagger:ignore
type NamespaceMap struct {
	Runtime
	Items map[string]*Namespace
}

// swagger:ignore
type NamespaceList struct {
	Runtime
	Items []*Namespace
}

// swagger:ignore
type Namespace struct {
	Meta NamespaceMeta `json:"meta"`
	Status NamespaceStatus `json:"status"`
	Spec NamespaceSpec `json:"spec"`
}

// swagger:ignore
type NamespaceEnvs []NamespaceEnv

// swagger:ignore
type NamespaceEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// swagger:ignore
// swagger:model types_namespace_meta
type NamespaceMeta struct {
	Meta     `yaml:",inline"`
	Endpoint string `json:"endpoint"`
	Type     string `json:"type"`
}

// swagger:ignore
type NamespaceSpec struct {
	Resources ResourceRequest `json:"resources"`
	Env       NamespaceEnvs      `json:"env"`
	Domain    NamespaceDomain    `json:"domain"`
}

type NamespaceStatus struct {
	Resources NamespaceStatusResources `json:"resources"`
}

type NamespaceStatusResources struct {
	Allocated ResourceRequestItem `json:"usage"`
}


type ResourceRequest struct {
	Request ResourceRequestItem `json:"request"`
	Limits  ResourceRequestItem `json:"limits"`
}


func (r *ResourceRequest) Equal(rr ResourceRequest) bool {


	if r.Limits.RAM != rr.Limits.RAM {
		return false
	}

	if r.Limits.CPU != rr.Limits.CPU {
		return false
	}

	if r.Limits.Storage != rr.Limits.Storage {
		return false
	}

	if r.Request.RAM != rr.Request.RAM {
		return false
	}

	if r.Request.CPU != rr.Request.CPU {
		return false
	}

	if r.Request.Storage != rr.Request.Storage {
		return false
	}

	return true
}

type NamespaceDomain struct {
	Internal string `json:"internal"`
	External string `json:"external"`
}

// swagger:ignore
type ResourceRequestItem struct {
	RAM     string `json:"ram"`
	CPU     string `json:"cpu"`
	Storage string `json:"storage"`
}



func (n *Namespace) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s", n.Meta.Name)
	}
	return n.Meta.SelfLink
}

func (n *Namespace) ToJson() ([]byte, error) {
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (n *NamespaceList) ToJson() ([]byte, error) {
	if n == nil {
		return []byte("[]"), nil
	}
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// swagger:ignore
type NamespaceCreateOptions struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Domain      *string                   `json:"domain"`
	Resources   *NamespaceResourcesOptions `json:"resources"`
}

// swagger:ignore
type NamespaceUpdateOptions struct {
	Description *string                    `json:"description"`
	Domain      *string                    `json:"domain"`
	Resources   *NamespaceResourcesOptions `json:"resources"`
}

// swagger:ignore
type NamespaceRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
type NamespaceResourcesOptions struct {
	Request *ResourceRequestItemOption `json:"request"`
	Limits  *ResourceRequestItemOption `json:"limits"`
}

// swagger:ignore
type ResourceRequestItemOption struct {
	RAM     *string `json:"ram"`
	CPU     *string `json:"cpu"`
	Storage *string `json:"storage"`
}

func (n *Namespace) AllocateResources(resources ResourceRequest) error {

	var (

		availableRam int64
		availableCpu int64

		allocatedRam int64
		allocatedCpu int64

		requestedRam int64
		requestedCpu int64

	)

	availableRam, _ = resource.DecodeMemoryResource(n.Spec.Resources.Limits.RAM)
	availableCpu, _ = resource.DecodeCpuResource(n.Spec.Resources.Limits.CPU)

	allocatedRam, _ = resource.DecodeMemoryResource(n.Status.Resources.Allocated.RAM)
	allocatedCpu, _ = resource.DecodeCpuResource(n.Status.Resources.Allocated.CPU)

	requestedRam, _ = resource.DecodeMemoryResource(resources.Limits.RAM)
	requestedCpu, _ = resource.DecodeCpuResource(resources.Limits.CPU)

	if (availableRam - allocatedRam - requestedRam) <= 0 {
		return errors.New(errors.ResourcesRamLimitExceeded)
	}

	if (availableCpu - allocatedCpu - requestedCpu) <= 0 {
		return errors.New(errors.ResourcesCpuLimitExceeded)
	}

	allocatedRam += requestedRam
	allocatedCpu += requestedCpu

	n.Status.Resources.Allocated.RAM = resource.EncodeMemoryResource(allocatedRam)
	n.Status.Resources.Allocated.CPU = resource.EncodeCpuResource(allocatedCpu)

	return nil
}

func (n *Namespace) ReleaseResources(resources ResourceRequest) {

	var (
		availableRam int64
		availableCpu int64
		allocatedRam int64
		allocatedCpu int64
		requestedRam int64
		requestedCpu int64
	)

	availableRam, _ = resource.DecodeMemoryResource(n.Spec.Resources.Limits.RAM)
	availableCpu, _ = resource.DecodeCpuResource(n.Spec.Resources.Limits.CPU)

	allocatedRam, _ = resource.DecodeMemoryResource(n.Status.Resources.Allocated.RAM)
	allocatedCpu, _ = resource.DecodeCpuResource(n.Status.Resources.Allocated.CPU)

	requestedRam, _ = resource.DecodeMemoryResource(resources.Limits.RAM)
	requestedCpu, _ = resource.DecodeCpuResource(resources.Limits.CPU)

	if (allocatedRam + requestedRam) > availableRam  {
		allocatedRam = availableRam
	} else {
		allocatedRam+=requestedRam
	}

	if (allocatedCpu + requestedCpu) > availableCpu  {
		allocatedCpu = availableCpu
	} else {
		allocatedCpu+=requestedCpu
	}

	n.Status.Resources.Allocated.RAM = resource.EncodeMemoryResource(allocatedRam)
	n.Status.Resources.Allocated.CPU = resource.EncodeCpuResource(allocatedCpu)
}

func NewNamespaceList() *NamespaceList {
	dm := new(NamespaceList)
	dm.Items = make([]*Namespace, 0)
	return dm
}

func NewNamespaceMap() *NamespaceMap {
	dm := new(NamespaceMap)
	dm.Items = make(map[string]*Namespace)
	return dm
}
