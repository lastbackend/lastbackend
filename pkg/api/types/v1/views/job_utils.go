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

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

type JobView struct{}

func (jw *JobView) New(obj *types.Job) *Job {
	j := Job{}
	return &j
}

func (jw *JobView) NewWithTasks(obj *types.Job, tasks *types.TaskList) *Job {
	j := Job{}
	return &j
}

func (jw *JobView) NewList(obj *types.JobList) *JobList {
	jl := make(JobList, 0)
	return &jl
}

func (jw *JobView) NewListWithTasks(obj *types.JobList, tl *types.TaskList) *JobList {
	jl := make(JobList, 0)
	return &jl
}
