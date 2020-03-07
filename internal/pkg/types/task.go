//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	Namespace string       `json:"namespace"`
	Job       string       `json:"job"`
	SelfLink  TaskSelfLink `json:"self_link"`
}

type TaskStatus struct {
	State        string             `json:"state"`
	Message      string             `json:"message"`
	Error        bool               `json:"error"`
	Canceled     bool               `json:"canceled"`
	Done         bool               `json:"done"`
	Dependencies StatusDependencies `json:"dependencies"`
	Pod          TaskStatusPod      `json:"pod"`
}

type TaskStatusPod struct {
	SelfLink string           `json:"self_link"`
	State    string           `json:"state"`
	Status   string           `json:"status"`
	Runtime  PodStatusRuntime `json:"runtime"`
}

type TaskSpec struct {
	State    SpecState    `json:"state" yaml:"state"`
	Runtime  SpecRuntime  `json:"runtime" yaml:"runtime"`
	Selector SpecSelector `json:"selector" yaml:"selector"`
	Template SpecTemplate `json:"template" yaml:"template"`
}

type TaskManifest struct {
	Meta TaskManifestMeta
	Spec TaskManifestSpec
}

type TaskManifestMeta struct {
	Name        *string
	Description *string
	Labels      map[string]string
}

type TaskManifestSpec struct {
	Runtime  *ManifestSpecRuntime
	Selector *ManifestSpecSelector
	Template *ManifestSpecTemplate
}

func (t *TaskManifest) SetTaskMeta(task *Task) {
	if task.Meta.Name == EmptyString {
		task.Meta.Name = *t.Meta.Name
	}

	if t.Meta.Labels != nil {
		task.Meta.Labels = t.Meta.Labels
	}
}

func (t *TaskManifest) SetTaskSpec(task *Task) error {

	if t.Spec.Runtime != nil {
		t.Spec.Runtime.SetSpecRuntime(&task.Spec.Runtime)
	}

	if t.Spec.Selector != nil {
		t.Spec.Selector.SetSpecSelector(&task.Spec.Selector)
	}

	if t.Spec.Template != nil {
		if err := t.Spec.Template.SetSpecTemplate(&task.Spec.Template); err != nil {
			return err
		}
	}

	return nil
}

func (ts *TaskStatus) CheckDeps() bool {

	for _, d := range ts.Dependencies.Volumes {
		if d.Status != StateReady {
			return false
		}
	}

	for _, d := range ts.Dependencies.Secrets {
		if d.Status != StateReady {
			return false
		}
	}

	for _, d := range ts.Dependencies.Configs {
		if d.Status != StateReady {
			return false
		}
	}

	return true
}

func (t *Task) SelfLink() *TaskSelfLink {
	return &t.Meta.SelfLink
}

func (t *Task) JobLink() *JobSelfLink {
	return t.Meta.SelfLink.parent.SelfLink.(*JobSelfLink)
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
