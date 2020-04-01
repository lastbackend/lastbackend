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
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type JobList []*Job
type Job views.Job

func (sl *JobList) Print() {

	t := table.New([]string{"NAME", "ENDPOINT", "STATUS", "REPLICAS"})
	t.VisibleHeader = true

	for _, s := range *sl {

		var data = map[string]interface{}{}

		data["NAME"] = s.Meta.Name
		data["STATUS"] = s.Status.State

		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Job) Print() {

	fmt.Printf("Name:\t\t%s/%s\n", s.Meta.Namespace, s.Meta.Name)
	if s.Meta.Description != models.EmptyString {
		fmt.Printf(" Description:\t%s\n", s.Meta.Description)
	}

	fmt.Printf("State:\t\t%s\n", s.Status.State)
	if s.Status.Message != models.EmptyString {
		fmt.Printf("Message:\t\t%s\n", s.Status.Message)
	}

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
	println()

	if len(s.Tasks) > 0 {

		taskTable := table.New([]string{"Name", "State", "Status", "Age", "Message"})
		taskTable.VisibleHeader = true

		for _, t := range s.Tasks {

			var taskRow = map[string]interface{}{}
			got, _ := timeago.TimeAgoWithTime(time.Now(), t.Meta.Created)
			taskRow["Name"] = t.Meta.Name
			taskRow["State"] = t.Status.State
			taskRow["Age"] = got
			taskRow["Message"] = t.Status.Message
			taskTable.AddRow(taskRow)
		}

		taskTable.Print()

	}
	println()
}

func (s *Job) PrintTask(d *views.Task) {

	fmt.Printf(" Name:\t\t%s\n", d.Meta.Name)
	if d.Meta.Description != models.EmptyString {
		fmt.Printf(" Description:\t%s\n", d.Meta.Description)
	}
	fmt.Printf(" State:\t\t%s\n", d.Status.State)
	if d.Status.Message != models.EmptyString {
		fmt.Printf(" Message:\t%s\n", d.Status.Message)
	}
	created, _ := timeago.TimeAgoWithTime(time.Now(), d.Meta.Created)
	updated, _ := timeago.TimeAgoWithTime(time.Now(), d.Meta.Updated)

	fmt.Printf(" Created:\t%s\n", created)
	fmt.Printf(" Updated:\t%s\n", updated)
	println()
	fmt.Printf(" Pods:\n")
	println()

	podTable := table.New([]string{"Name", "Ready", "Status", "Restarts", "Age"})
	podTable.VisibleHeader = true

	var (
		ids = make([]int, 0, len(d.Pods))
	)
	for key := range d.Pods {
		ids = append(ids, key)
	}

	sort.Ints(ids) //sort by key

	for _, id := range ids {
		p := d.Pods[id]

		var ready, restarts int
		for _, c := range p.Status.Runtime.Services {
			if c.Ready {
				ready++
				restarts += c.Restart
			}
		}
		var podRow = map[string]interface{}{}
		got, _ := timeago.TimeAgoWithTime(time.Now(), p.Meta.Created)
		podRow["Name"] = p.Meta.Name
		podRow["Ready"] = string(converter.IntToString(ready) + "/" + converter.IntToString(len(p.Status.Runtime.Services)))
		podRow["Status"] = p.Status.State
		podRow["Restarts"] = restarts
		podRow["Age"] = got
		podTable.AddRow(podRow)
	}

	podTable.Print()
	println()
}

func FromApiJobView(job *views.Job) *Job {

	if job == nil {
		return nil
	}

	item := Job(*job)
	return &item
}

func FromApiJobListView(jobs *views.JobList) *JobList {
	var items = make(JobList, 0)
	for _, job := range *jobs {
		items = append(items, FromApiJobView(job))
	}
	return &items
}
