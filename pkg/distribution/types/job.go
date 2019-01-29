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

const (
	DEFAULT_JOB_MEMORY      int64 = 128
	DEFAULT_JOB_PARALLELISM int   = 1
)

type Job struct {
	Runtime
	Meta   JobMeta   `json:"meta"`
	Status JobStatus `json:"status"`
	Spec   JobSpec   `json:"spec"`
}

type JobMap struct {
	Runtime
	Items map[string]*Job
}

type JobList struct {
	Runtime
	Items []*Job
}

type JobMeta struct {
	Meta
	Namespace string `json:"namespace"`
	SelfLink  string `json:"self_link"`
}

type JobStatus struct {
}

type JobSpec struct {
}

type JobSpecTemplate struct {
}

type JobRunner struct {
	Runtime
	Meta   JobRunnerMeta   `json:"meta"`
	Status JobRunnerStatus `json:"status"`
	Spec   JobRunnerSpec   `json:"spec"`
}

type JobRunnerMap struct {
	Runtime
	Items map[string]*JobRunner
}

type JobRunnerList struct {
	Runtime
	Items []*JobRunner
}

type JobRunnerMeta struct {
	Meta
	Namespace string `json:"namespace"`
	SelfLink  string `json:"self_link"`
}

type JobRunnerStatus struct {
}

type JobRunnerSpec struct {
}
