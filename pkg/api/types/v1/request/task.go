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

package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"gopkg.in/yaml.v2"
)

type TaskManifest struct {
	Meta TaskManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec TaskManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type TaskManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type TaskManifestSpec struct {
	Runtime  *ManifestSpecRuntime  `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

func (t *TaskManifest) SetTaskManifestMeta(task *types.TaskManifest) {

	task.Meta.Name = t.Meta.Name
	task.Meta.Labels = t.Meta.Labels
	task.Meta.Description = t.Meta.Description
}

func (t *TaskManifest) SetTaskManifestSpec(task *types.TaskManifest) error {

	if t.Spec.Runtime != nil {
		if task.Spec.Runtime == nil {
			task.Spec.Runtime = new(types.ManifestSpecRuntime)
		}
		t.Spec.Runtime.SetManifestSpecRuntime(task.Spec.Runtime)
	}

	if t.Spec.Selector != nil {
		if task.Spec.Selector == nil {
			task.Spec.Selector = new(types.ManifestSpecSelector)
		}
		t.Spec.Selector.SetManifestSpecSelector(task.Spec.Selector)
	}

	if t.Spec.Template != nil {
		if task.Spec.Template == nil {
			task.Spec.Template = new(types.ManifestSpecTemplate)
		}
		if err := t.Spec.Template.SetManifestSpecTemplate(task.Spec.Template); err != nil {
			return err
		}
	}

	return nil
}

func (t *TaskManifest) SetTaskMeta(task *types.Task) {
	if task.Meta.Name == types.EmptyString {
		task.Meta.Name = *t.Meta.Name
	}

	if t.Meta.Labels != nil {
		task.Meta.Labels = t.Meta.Labels
	}
}

func (t *TaskManifest) SetTaskSpec(task *types.Task) error {

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

func (t *TaskManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *TaskManifest) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TaskManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, t)
}

func (t *TaskManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(t)
}

type TaskCancelOptions struct {
	Force bool `json:"force"`
}

func (t *TaskCancelOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

type TaskLogsOptions struct {
	Tail      int    `json:"tail"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Follow    bool   `json:"follow"`
}

func (t *TaskLogsOptions) ToJson() ([]byte, error) {
	return json.Marshal(t)
}
