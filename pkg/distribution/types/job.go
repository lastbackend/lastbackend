//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

import (
	"fmt"
	"time"
)

const (
	DEFAULT_JOB_MEMORY      int64 = 128
	DEFAULT_JOB_PARALLELISM int   = 1
)

type Job struct {
	System
	Meta   JobMeta   `json:"meta"`
	Status JobStatus `json:"status"`
	Spec   JobSpec   `json:"spec"`
}

type JobMap struct {
	System
	Items map[string]*Job
}

type JobList struct {
	System
	Items []*Job
}

type JobMeta struct {
	Meta
	Namespace string `json:"namespace"`
	SelfLink  string `json:"self_link"`
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
	Template    SpecTemplate       `json:"template"`
	Concurrency JobSpecConcurrency `json:"concurrency"`
	Remote      JobSpecRemote      `json:"remote"`
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

func (js *JobStatus) SetCreated() {
	js.State = StateCreated
	js.Message = ""
}

func (js *JobStatus) SetProvision() {
	js.State = StateProvision
	js.Message = ""
}

func (js *JobStatus) SetRunning() {
	js.State = StateRunning
	js.Message = ""
}

func (js *JobStatus) SetPaused() {
	js.State = StatePaused
	js.Message = ""
}

func (js *JobStatus) SetDestroy() {
	js.State = StateDestroy
	js.Message = ""
}

func (js *JobStatus) SetError(message string) {
	js.State = StateError
	js.Message = message
}

func (j *Job) SelfLink() string {
	if j.Meta.SelfLink == EmptyString {
		j.Meta.SelfLink = j.CreateSelfLink(j.Meta.Namespace, j.Meta.Name)
	}
	return j.Meta.SelfLink
}

func (j *Job) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func NewJobList() *JobList {
	jrl := new(JobList)
	jrl.Items = make([]*Job, 0)
	return jrl
}
