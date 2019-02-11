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

func (p TaskList) Len() int {
	return len(p)
}

func (p TaskList) Less(i, j int) bool {
	return p[j].Meta.Created.Before(p[i].Meta.Created)
}

func (p TaskList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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
	State   string `json:"state"`
	Message string `json:"message"`
}

type TaskSpec struct {
	Runtime  ManifestSpecRuntime  `json:"runtime"`
	Selector ManifestSpecSelector `json:"selector"`
	Template ManifestSpecTemplate `json:"template"`
}

func (t *Task) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (tl *TaskList) ToJson() ([]byte, error) {
	sort.Sort(tl)
	return json.Marshal(tl)
}
