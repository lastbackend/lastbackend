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

type JobManifest struct {
	Meta JobManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec JobManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type JobManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type JobManifestSpec struct {
	Enabled     bool                       `json:"enabled" yaml:"enabled"`
	Schedule    string                     `json:"schedule" yaml:"schedule"`
	Runtime     ManifestSpecRuntime        `json:"runtime" yaml:"runtime"`
	Template    ManifestSpecTemplate       `json:"template" yaml:"template"`
	Concurrency JobManifestSpecConcurrency `json:"concurrency" yaml:"concurrency"`
	Remote      JobManifestSpecRemote      `json:"remote" yaml:"remote"`
}

type JobManifestSpecConcurrency struct {
	Limit    int    `json:"limit" yaml:"limit"`
	Strategy string `json:"strategy" yaml:"limit"`
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

	if j.Meta.Labels != nil {
		job.Meta.Labels = j.Meta.Labels
	}
}

func (j *JobManifest) SetJobSpec(job *types.Job) (err error) {

	return nil
}

type JobRemoveOptions struct {
	Force bool `json:"force"`
}
