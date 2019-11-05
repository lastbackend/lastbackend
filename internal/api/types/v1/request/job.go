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
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/resource"
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
	Timeout  string                           `json:"timeout" yaml:"timeout"`
	Http     *JobManifestSpecProviderHTTP     `json:"http" yaml:"http"`
	Cron     *JobManifestSpecProviderCron     `json:"cron" yaml:"cron"`
	RabbitMQ *JobManifestSpecProviderRabbitMQ `json:"rabbitmq" yaml:"rabbitmq"`
}

type JobManifestSpecProviderHTTP struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint"`
	Method   string            `json:"method" yaml:"method"`
	Headers  map[string]string `json:"headers" yaml:"headers"`
}

type JobManifestSpecProviderCron struct {
}

type JobManifestSpecProviderRabbitMQ struct {
}

type JobManifestSpecHook struct {
	Http *JobManifestSpecProviderHTTP `json:"http" yaml:"http"`
}

type JobManifestSpecHookHTTP struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint"`
	Method   string            `json:"method" yaml:"method"`
	Headers  map[string]string `json:"headers" yaml:"headers"`
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

	job.Spec.Provider.Timeout = j.Spec.Provider.Timeout

	if j.Spec.Provider.Http != nil {
		job.Spec.Provider.Http = new(types.JobSpecProviderHTTP)
		job.Spec.Provider.Http.Endpoint = j.Spec.Provider.Http.Endpoint
		job.Spec.Provider.Http.Method = j.Spec.Provider.Http.Method
		job.Spec.Provider.Http.Headers = j.Spec.Provider.Http.Headers
	} else {
		job.Spec.Provider.Http = nil
	}

	if j.Spec.Provider.Cron != nil {

	}

	if j.Spec.Provider.RabbitMQ != nil {

	}

	if j.Spec.Hook.Http != nil {
		job.Spec.Hook.Http = new(types.JobSpecHookHTTP)
		job.Spec.Hook.Http.Endpoint = j.Spec.Hook.Http.Endpoint
		job.Spec.Hook.Http.Method = j.Spec.Hook.Http.Method
		job.Spec.Hook.Http.Headers = j.Spec.Hook.Http.Headers
	} else {
		job.Spec.Hook.Http = nil
	}

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
	Tail      int    `json:"tail"`
	Task      string `json:"task"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Follow    bool   `json:"follow"`
}
