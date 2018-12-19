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
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"github.com/lastbackend/registry/pkg/distribution/types"
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
	Meta   NamespaceMeta   `json:"meta"`
	Status NamespaceStatus `json:"status"`
	Spec   NamespaceSpec   `json:"spec"`
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
	Env       NamespaceEnvs   `json:"env"`
	Domain    NamespaceDomain `json:"domain"`
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
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Domain      *string                    `json:"domain"`
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

		err error
	)

	var handleErr = func(msg string, e error) error {
		log.Errorf("allocate %s error: %s", msg, e.Error())
		return e
	}

	if n.Spec.Resources.Limits.RAM != types.EmptyString {
		availableRam, err = resource.DecodeMemoryResource(n.Spec.Resources.Limits.RAM)
		if err != nil {
			return handleErr("ns limit ram", err)
		}
	}

	if n.Spec.Resources.Limits.CPU != types.EmptyString {
		availableCpu, err = resource.DecodeCpuResource(n.Spec.Resources.Limits.CPU)
		if err != nil {
			return handleErr("ns limit cpu", err)
		}
	}

	if n.Status.Resources.Allocated.RAM != types.EmptyString {
		allocatedRam, err = resource.DecodeMemoryResource(n.Status.Resources.Allocated.RAM)
		if err != nil {
			return handleErr("ns allocated ram", err)
		}
	}

	if n.Status.Resources.Allocated.CPU != types.EmptyString {
		allocatedCpu, err = resource.DecodeCpuResource(n.Status.Resources.Allocated.CPU)
		if err != nil {
			return handleErr("ns allocated cpu", err)
		}
	}

	if resources.Limits.RAM != types.EmptyString {
		requestedRam, err = resource.DecodeMemoryResource(resources.Limits.RAM)
		if err != nil {
			return handleErr("req limit ram", err)
		}
	}

	if resources.Limits.CPU != types.EmptyString {
		requestedCpu, err = resource.DecodeCpuResource(resources.Limits.CPU)
		if err != nil {
			return handleErr("req limit cpu", err)
		}
	}

	if availableRam > 0 && availableCpu > 0 {

		if requestedRam == 0 {
			return errors.New(errors.ResourcesRamLimitIsRequired)
		}

		if requestedCpu == 0 {
			return errors.New(errors.ResourcesCpuLimitIsRequired)
		}


		if (availableRam - allocatedRam - requestedRam) <= 0 {
			return errors.New(errors.ResourcesRamLimitExceeded)
		}

		if (availableCpu - allocatedCpu - requestedCpu) <= 0 {
			return errors.New(errors.ResourcesCpuLimitExceeded)
		}
	}
	allocatedRam += requestedRam
	allocatedCpu += requestedCpu

	n.Status.Resources.Allocated.RAM = resource.EncodeMemoryResource(allocatedRam)
	n.Status.Resources.Allocated.CPU = resource.EncodeCpuResource(allocatedCpu)

	return nil
}

func (n *Namespace) ReleaseResources(resources ResourceRequest) error {

	var (
		availableRam int64
		availableCpu int64
		allocatedRam int64
		allocatedCpu int64
		requestedRam int64
		requestedCpu int64
		err error
	)

	var handleErr = func(msg string, e error) error {
		log.Errorf("allocate %s error: %s", msg, e.Error())
		return e
	}

	if n.Spec.Resources.Limits.RAM != types.EmptyString {
		availableRam, err = resource.DecodeMemoryResource(n.Spec.Resources.Limits.RAM)
		if err != nil {
			return handleErr("ns limit ram", err)
		}
	}

	if n.Spec.Resources.Limits.CPU != types.EmptyString {
		availableCpu, err = resource.DecodeCpuResource(n.Spec.Resources.Limits.CPU)
		if err != nil {
			return handleErr("ns limit cpu", err)
		}
	}

	if n.Status.Resources.Allocated.RAM != types.EmptyString {
		allocatedRam, err = resource.DecodeMemoryResource(n.Status.Resources.Allocated.RAM)
		if err != nil {
			return handleErr("ns allocated ram", err)
		}
	}

	if n.Status.Resources.Allocated.CPU != types.EmptyString {
		allocatedCpu, err = resource.DecodeCpuResource(n.Status.Resources.Allocated.CPU)
		if err != nil {
			return handleErr("ns allocated cpu", err)
		}
	}

	if resources.Limits.RAM != types.EmptyString {
		requestedRam, err = resource.DecodeMemoryResource(resources.Limits.RAM)
		if err != nil {
			return handleErr("req limit ram", err)
		}
	}

	if resources.Limits.CPU != types.EmptyString {
		requestedCpu, err = resource.DecodeCpuResource(resources.Limits.CPU)
		if err != nil {
			return handleErr("req limit cpu", err)
		}
	}

	if (allocatedRam+requestedRam) > availableRam && (availableRam > 0) {
		allocatedRam = availableRam
	} else {
		allocatedRam -= requestedRam
	}

	if (allocatedCpu+requestedCpu) > availableCpu && (availableRam > 0) {
		allocatedCpu = availableCpu
	} else {
		allocatedCpu -= requestedCpu
	}

	n.Status.Resources.Allocated.RAM = resource.EncodeMemoryResource(allocatedRam)
	n.Status.Resources.Allocated.CPU = resource.EncodeCpuResource(allocatedCpu)

	return nil
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
