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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"sort"
	"time"

	"github.com/ararog/timeago"
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type Task views.Task

func (s *Task) Print() {

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

	if len(s.Status.Pod.Runtime.Services) > 0 {
		fmt.Println("Services:")
		fmt.Println()
		taskTable := table.New([]string{"Name", "Status", "Age", "Message"})
		taskTable.VisibleHeader = true

		for _, svc := range s.Status.Pod.Runtime.Services {

			var taskRow = map[string]interface{}{}
			got, _ := timeago.TimeAgoWithTime(time.Now(), svc.State.Created.Timestamp)
			taskRow["Name"] = svc.Name
			taskRow["Status"] = svc.Ready
			taskRow["Age"] = got
			taskTable.AddRow(taskRow)
		}

		taskTable.Print()

	}

	fmt.Println()
	fmt.Println("Pipeline:")
	fmt.Println()

	if len(s.Status.Pod.Runtime.Pipeline) > 0 {

		for _, step := range s.Status.Pod.Runtime.Pipeline {

			fmt.Printf("Step: %s Status: %s\n", step.Name, step.Status)
			fmt.Println()

			if step.Error && step.Message != models.EmptyString {
				fmt.Printf("Error: %s\n", step.Message)
			}

			taskTable := table.New([]string{"Command", "State", "Age", "Message"})
			taskTable.VisibleHeader = true

			for _, cmd := range step.Commands {
				var taskRow = map[string]interface{}{}
				got, _ := timeago.TimeAgoWithTime(time.Now(), cmd.State.Created.Timestamp)
				taskRow["Name"] = cmd.Name
				taskRow["Status"] = cmd.Ready
				taskRow["Age"] = got
				taskTable.AddRow(taskRow)
			}

			taskTable.Print()
			fmt.Println()

		}

	} else {
		fmt.Println("no commands executed")
	}

	println()
}

func FromApiTaskView(task *views.Task) *Task {

	if task == nil {
		return nil
	}

	item := Task(*task)
	return &item
}
