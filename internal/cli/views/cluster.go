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
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/views"
	"github.com/lastbackend/lastbackend/internal/cli/models"
	"github.com/lastbackend/lastbackend/internal/util/table"
	lv "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type ClusterList []*Cluster

type Cluster struct {
	Meta   ClusterMeta      `json:"meta"`
	Status lv.ClusterStatus `json:"status"`
	Spec   ClusterSpec      `json:"spec"`
}

type ClusterMeta struct {
	lv.Meta
	Local bool `json:"local"`
}

type ClusterSpec struct {
	Endpoint string `json:"endpoint"`
}

func (cl *ClusterList) Print() {

	if cl == nil || len(*cl) == 0 {
		fmt.Println("no clusters available")
		println()
		return
	}

	println()

	t := table.New([]string{"NAME", "ENDPOINT"})
	t.VisibleHeader = true

	for _, s := range *cl {
		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		data["ENDPOINT"] = s.Spec.Endpoint
		t.AddRow(data)
	}

	print()
	t.Print()
	println()
}

func (c *Cluster) Print() {
	print()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":     c.Meta.Name,
		"ENDPOINT": c.Spec.Endpoint,
		"CREATED":  c.Meta.Created.Format("15:04:05 _2.01.2006"),
	})
	print()
}

func FromLbApiClusterView(cl *lv.Cluster) *Cluster {
	if cl == nil {
		return nil
	}
	c := new(Cluster)
	c.Meta.Name = cl.Meta.Name
	c.Status = cl.Status
	return c
}

func FromGenesisApiClusterView(cl *views.ClusterView) *Cluster {
	if cl == nil {
		return nil
	}
	c := new(Cluster)
	c.Meta.Name = cl.Meta.SelfLink
	c.Spec.Endpoint = cl.Spec.Endpoint
	c.Status = cl.Status
	return c
}

func FromGenesisApiClusterListView(cl *views.ClusterList) ClusterList {
	if cl == nil {
		list := make(views.ClusterList, 0)
		cl = &list
	}
	list := make(ClusterList, 0)
	for _, item := range *cl {
		list = append(list, FromGenesisApiClusterView(item))
	}
	return list
}

func FromStorageClusterList(cl []*models.Cluster) ClusterList {
	if cl == nil {
		cl = make([]*models.Cluster, 0)
	}
	list := make(ClusterList, 0)
	for _, item := range cl {
		c := new(Cluster)
		c.Meta.Name = item.Name
		c.Spec.Endpoint = item.Endpoint
		list = append(list, c)
	}
	return list
}
