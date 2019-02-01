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
	Runtime  *ManifestSpecRuntime  `json:"runtime" yaml:"runtime"`
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
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

type TaskLogsOptions struct {
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Follow    bool   `json:"follow"`
}
