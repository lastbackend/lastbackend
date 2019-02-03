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

package types

import "fmt"

type Task struct {
	System
	Meta   TaskMeta   `json:"meta"`
	Status TaskStatus `json:"status"`
	Spec   TaskSpec   `json:"spec"`
}

type TaskMap struct {
	System
	Items map[string]*Task
}

type TaskList struct {
	System
	Items []*Task
}

type TaskMeta struct {
	Meta
	Namespace string `json:"namespace"`
	Job       string `json:"job"`
	SelfLink  string `json:"self_link"`
}

type TaskStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

type TaskSpec struct {
	State    SpecState    `json:"state" yaml:"state"`
	Runtime  SpecRuntime  `json:"runtime" yaml:"runtime"`
	Selector SpecSelector `json:"selector" yaml:"selector"`
	Template SpecTemplate `json:"template" yaml:"template"`
}

func (j *TaskStatus) SetCreated() {
	j.State = StateCreated
	j.Message = ""
}

func (j *TaskStatus) SetProvision() {
	j.State = StateProvision
	j.Message = ""
}

func (j *TaskStatus) SetStarted() {
	j.State = StateStarted
	j.Message = ""
}

func (j *TaskStatus) SetFinished() {
	j.State = StateExited
	j.Message = ""
}

func (j *TaskStatus) SetCancel() {
	j.State = StateCancel
	j.Message = ""
}

func (j *TaskStatus) SetDestroy() {
	j.State = StateDestroy
	j.Message = ""
}

func (j *TaskStatus) SetError(message string) {
	j.State = StateError
	j.Message = message
}

func (j *Task) SelfLink() string {
	if j.Meta.SelfLink == EmptyString {
		j.Meta.SelfLink = j.CreateSelfLink(j.Meta.Namespace, j.Meta.Job, j.Meta.Name)
	}
	return j.Meta.SelfLink
}

func (j *Task) CreateSelfLink(namespace, job, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, job, name)
}

// GetResourceRequest - request resources for task creation
// Use replica later when multi-pod tasks will be implemented
func (ts *TaskSpec) GetResourceRequest() ResourceRequest {

	rr := ResourceRequest{}

	var (
		limitsRAM int64
		limitsCPU int64

		requestRAM int64
		requestCPU int64
	)

	for _, c := range ts.Template.Containers {

		limitsCPU += c.Resources.Limits.CPU
		limitsRAM += c.Resources.Limits.RAM

		requestCPU += c.Resources.Request.CPU
		requestRAM += c.Resources.Request.RAM
	}

	if requestRAM > 0 {
		rr.Request.RAM = requestRAM
	}

	if requestCPU > 0 {
		rr.Request.CPU = requestCPU
	}

	if limitsRAM > 0 {
		rr.Limits.RAM = limitsRAM
	}

	if limitsCPU > 0 {
		rr.Limits.CPU = limitsCPU
	}

	return rr
}

func NewTaskList() *TaskList {
	jl := new(TaskList)
	jl.Items = make([]*Task, 0)
	return jl
}

func NewTaskMap() *TaskMap {
	jm := new(TaskMap)
	jm.Items = make(map[string]*Task, 0)
	return jm
}