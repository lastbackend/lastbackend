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
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"gopkg.in/yaml.v2"
)

type JobManifest struct {
	Meta JobManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec JobManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type JobManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type JobManifestSpec struct {
	Enabled     bool                        `json:"enabled" yaml:"enabled"`
	Schedule    string                      `json:"schedule" yaml:"schedule"`
	Concurrency JobManifestSpecConcurrency  `json:"concurrency" yaml:"concurrency"`
	Provider    JobManifestSpecProvider     `json:"provider" yaml:"provider"`
	Hook        JobManifestSpecHook         `json:"hook" yaml:"hook"`
	Resources   *JobResourcesOptions        `json:"resources" yaml:"resources"`
	Task        JobManifestSpecTaskTemplate `json:"task" yaml:"task"`
}

type JobManifestSpecTaskTemplate struct {
	Selector *ManifestSpecSelector `json:"selector" yaml:"selector"`
	Runtime  *ManifestSpecRuntime  `json:"runtime" yaml:"runtime"`
	Template *ManifestSpecTemplate `json:"template" yaml:"template"`
}

type JobManifestSpecConcurrency struct {
	Limit    int    `json:"limit" yaml:"limit"`
	Strategy string `json:"strategy" yaml:"strategy"`
}

type JobManifestSpecProvider struct {
	Timeout int                    `json:"timeout" yaml:"timeout"`
	Kind    string                 `json:"kind" yaml:"kind"`
	Config  map[string]interface{} `json:"config" yaml:"config"`
}

type JobManifestSpecHook struct {
	Kind   string                 `json:"kind" yaml:"kind"`
	Config map[string]interface{} `json:"config" yaml:"config"`
}

type JobManifestSpecRemote struct {
	Timeout  int                          `json:"timeout" yaml:"timeout"`
	Request  JobManifestSpecRemoteRequest `json:"request" yaml:"request"`
	Response JobManifestSpecRemoteRequest `json:"response" yaml:"response"`
}

type JobManifestSpecRemoteRequest struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint"`
	Headers  map[string]string `json:"headers" yaml:"headers"`
	Method   string            `json:"method" yaml:"method"`
}

func (j *JobManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, j)
}

func (j *JobManifest) ToJson() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JobManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, j)
}

func (j *JobManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(j)
}

func (j *JobManifest) SetJobMeta(job *types.Job) {

	if job.Meta.Name == types.EmptyString {
		job.Meta.Name = *j.Meta.Name
	}

	if j.Meta.Description != nil {
		job.Meta.Description = *j.Meta.Description
	}

	if j.Meta.Labels != nil {
		job.Meta.Labels = j.Meta.Labels
	}
}

func (j *JobManifest) SetJobSpec(job *types.Job) (err error) {

	job.Spec.Enabled = j.Spec.Enabled

	job.Spec.Concurrency.Limit = j.Spec.Concurrency.Limit
	job.Spec.Concurrency.Strategy = j.Spec.Concurrency.Strategy

	job.Spec.Provider.Kind = j.Spec.Provider.Kind
	job.Spec.Provider.Timeout = j.Spec.Provider.Timeout
	job.Spec.Provider.Config = j.Spec.Provider.Config

	job.Spec.Hook.Kind = j.Spec.Hook.Kind
	job.Spec.Hook.Config = j.Spec.Hook.Config

	if j.Spec.Resources != nil {

		if j.Spec.Resources.Request != nil {

			if j.Spec.Resources.Request.RAM != nil {
				ram, err := resource.DecodeMemoryResource(*j.Spec.Resources.Request.RAM)
				if err != nil {
					return err
				}

				job.Spec.Resources.Request.RAM = ram
			}

			if j.Spec.Resources.Request.CPU != nil {

				cpu, err := resource.DecodeCpuResource(*j.Spec.Resources.Request.CPU)
				if err != nil {
					return err
				}

				job.Spec.Resources.Request.CPU = cpu
			}
		}

		if j.Spec.Resources.Limits != nil {

			if j.Spec.Resources.Limits.RAM != nil {

				ram, err := resource.DecodeMemoryResource(*j.Spec.Resources.Limits.RAM)
				if err != nil {
					return err
				}

				job.Spec.Resources.Limits.RAM = ram
			}

			if j.Spec.Resources.Limits.CPU != nil {

				cpu, err := resource.DecodeCpuResource(*j.Spec.Resources.Limits.CPU)
				if err != nil {
					return err
				}

				job.Spec.Resources.Limits.CPU = cpu
			}
		}

	}

	if j.Spec.Task.Selector != nil {
		j.Spec.Task.Selector.SetSpecSelector(&job.Spec.Task.Selector)
	} else {
		job.Spec.Task.Selector.SetDefault()
	}

	if j.Spec.Task.Runtime != nil {
		j.Spec.Task.Runtime.SetSpecRuntime(&job.Spec.Task.Runtime)
	}

	if j.Spec.Task.Template != nil {

		if err := j.Spec.Task.Template.SetSpecTemplate(&job.Spec.Task.Template); err != nil {
			return err
		}
	}

	return nil
}

type JobResourcesOptions struct {
	Request *JobResourceOptions `json:"request"`
	Limits  *JobResourceOptions `json:"limits"`
}

// swagger:model request_namespace_quotas
type JobResourceOptions struct {
	RAM     *string `json:"ram"`
	CPU     *string `json:"cpu"`
	Storage *string `json:"storage"`
}

type JobRemoveOptions struct {
	Force bool `json:"force"`
}

type JobLogsOptions struct {
	Task      string `json:"task"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Follow    bool   `json:"follow"`
}
