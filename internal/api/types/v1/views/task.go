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
	"sort"
)

type TaskList []*Task

func (tl TaskList) Len() int {
	return len(tl)
}

func (tl TaskList) Less(i, j int) bool {
	return tl[j].Meta.Created.Before(tl[i].Meta.Created)
}

func (tl TaskList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

type Task struct {
	Meta   TaskMeta   `json:"meta"`
	Status TaskStatus `json:"status"`
	Spec   TaskSpec   `json:"spec"`
	Pods   PodList    `json:"pods,omitempty"`
}

type TaskMeta struct {
	Meta
	Namespace string `json:"namespace"`
	Job       string `json:"job"`
}

type TaskStatus struct {
	State    string        `json:"state"`
	Message  string        `json:"message"`
	Error    bool          `json:"error"`
	Done     bool          `json:"done"`
	Canceled bool          `json:"canceled"`
	Pod      TaskStatusPod `json:"pod"`
}

type TaskSpec struct {
	Runtime  ManifestSpecRuntime  `json:"runtime"`
	Selector ManifestSpecSelector `json:"selector"`
	Template ManifestSpecTemplate `json:"template"`
}

type TaskStatusPod struct {
	SelfLink string           `json:"self_link"`
	Status   string           `json:"status"`
	State    string           `json:"state"`
	Runtime  PodStatusRuntime `json:"runtime"`
}

func (t *Task) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (tl *TaskList) ToJson() ([]byte, error) {
	sort.Sort(tl)
	return json.Marshal(tl)
}
