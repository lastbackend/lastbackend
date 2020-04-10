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

package views

import (
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type VolumeList []*Volume
type Volume views.Volume

func (rl *VolumeList) Print() {

	t := table.New([]string{"NAMESPACE", "NAME", "STATUS", "CAPACITY"})
	t.VisibleHeader = true

	for _, r := range *rl {
		var data = map[string]interface{}{}
		data["NAMESPACE"] = r.Meta.Namespace
		data["NAME"] = r.Meta.Name
		data["STATUS"] = r.Status
		data["CAPACITY"] = r.Spec.Capacity.Storage
		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func (r *Volume) Print() {
	var data = map[string]interface{}{}
	data["NAME"] = r.Meta.Name
	data["NAMESPACE"] = r.Meta.Namespace
	data["STATUS"] = r.Status
	data["CAPACITY"] = r.Spec.Capacity.Storage
	println()
	table.PrintHorizontal(data)
	println()
}

func FromApiVolumeView(volume *views.Volume) *Volume {
	item := Volume(*volume)
	return &item
}

func FromApiVolumeListView(volumes *views.VolumeList) *VolumeList {
	var items = make(VolumeList, 0)
	for _, volume := range *volumes {
		items = append(items, FromApiVolumeView(volume))
	}
	return &items
}
