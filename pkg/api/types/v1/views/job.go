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
	"time"
)

type JobList []*Job

type Job struct {
	Meta   JobMeta   `json:"meta"`
	Status JobStatus `json:"status"`
	Spec   JobSpec   `json:"spec"`
	Tasks  TaskList  `json:"tasks,omitempty"`
}

type JobMeta struct {
	Meta
	Namespace string `json:"namespace"`
}

type JobStatus struct {
	State   string         `json:"state"`
	Message string         `json:"message"`
	Stats   JobStatusStats `json:"stats"`
	Updated time.Time      `json:"updated"`
}

type JobStatusStats struct {
	Total        int       `json:"total"`
	Active       int       `json:"active"`
	Successful   int       `json:"successful"`
	Failed       int       `json:"failed"`
	LastSchedule time.Time `json:"last_schedule"`
}

type JobSpec struct {
	Enabled     bool               `json:"enabled"`
	Schedule    string             `json:"schedule"`
	Concurrency JobSpecConcurrency `json:"concurrency"`
	Remote      JobSpecRemote      `json:"remote"`
	Task        JobSpecTask        `json:"task"`
}

type JobSpecTask struct {
	Selector ManifestSpecSelector `json:"selector"`
	Runtime  ManifestSpecRuntime  `json:"runtime"`
	Template ManifestSpecTemplate `json:"template"`
}

type JobSpecConcurrency struct {
	Limit    int    `json:"limit"`
	Strategy string `json:"strategy"`
}

type JobSpecRemote struct {
	Timeout  int                  `json:"timeout"`
	Request  JobSpecRemoteRequest `json:"request"`
	Response JobSpecRemoteRequest `json:"response"`
}

type JobSpecRemoteRequest struct {
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
}

func (j *Job) ToJson() ([]byte, error) {
	return json.Marshal(j)
}

func (jl *JobList) ToJson() ([]byte, error) {
	return json.Marshal(jl)
}
