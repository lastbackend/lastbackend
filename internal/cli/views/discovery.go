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
	"sort"
	"time"

	"github.com/ararog/timeago"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type DiscoveryList []*Discovery
type Discovery views.Discovery

func (sl *DiscoveryList) Print() {

	t := table.New([]string{"NAME", "STATUS", "AGE", "VERSION"})
	t.VisibleHeader = true

	for _, s := range *sl {

		var data = map[string]interface{}{}

		data["NAME"] = s.Meta.Name

		if s.Status.Ready {
			data["STATUS"] = models.StateReady
		} else {
			data["STATUS"] = models.StateNotReady
		}

		created, _ := timeago.TimeAgoWithTime(time.Now(), s.Meta.Created)
		data["AGE"] = created
		data["VERSION"] = s.Meta.Version
		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func (s *Discovery) Print() {

	fmt.Printf("Name:\t\t%s\n", s.Meta.Name)
	created, _ := timeago.TimeAgoWithTime(time.Now(), s.Meta.Created)
	updated, _ := timeago.TimeAgoWithTime(time.Now(), s.Meta.Updated)

	fmt.Printf("Created:\t%s\n", created)
	fmt.Printf("Updated:\t%s\n", updated)

	var (
		labels = make([]string, 0, len(s.Meta.Labels))
		out    string
	)

	for key := range s.Meta.Labels {
		labels = append(labels, key)
	}

	sort.Strings(labels) //sort by key
	for _, key := range labels {
		out += key + "=" + s.Meta.Labels[key] + " "
	}

	fmt.Printf("Labels:\t\t%s\n", out)
	println()
}

func FromApiDiscoveryView(discovery *views.Discovery) *Discovery {

	if discovery == nil {
		return nil
	}

	item := Discovery(*discovery)
	return &item
}

func FromApiDiscoveryListView(dl *views.DiscoveryList) *DiscoveryList {
	var items = make(DiscoveryList, 0)
	for _, discovery := range *dl {
		items = append(items, FromApiDiscoveryView(discovery))
	}
	return &items
}
