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

package views

import (
	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type ExporterView struct{}

func (nv *ExporterView) New(obj *types.Exporter) *Exporter {
	n := Exporter{}
	n.Meta = nv.ToExporterMeta(obj.Meta)
	n.Status = nv.ToExporterStatus(obj.Status)
	return &n
}

func (nv *ExporterView) ToExporterMeta(meta types.ExporterMeta) ExporterMeta {
	m := ExporterMeta{}
	m.Name = meta.Name
	m.Description = meta.Description
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func (nv *ExporterView) ToExporterStatus(status types.ExporterStatus) ExporterStatus {
	return ExporterStatus{
		Ready: status.Ready,
	}
}

func (obj *Exporter) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *ExporterView) NewList(obj *types.ExporterList) *ExporterList {
	if obj == nil {
		return nil
	}
	ingresses := make(ExporterList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		ingresses[nn.Meta.Name] = nn
	}

	return &ingresses
}

func (obj *ExporterList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *ExporterView) NewManifest(obj *types.ExporterManifest) *ExporterManifest {

	manifest := ExporterManifest{}

	if obj == nil {
		return nil
	}

	return &manifest
}

func (obj *ExporterManifest) Decode() *types.ExporterManifest {

	manifest := types.ExporterManifest{}

	return &manifest
}

func (obj *ExporterManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
