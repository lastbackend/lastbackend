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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
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
	Namespace string      `json:"namespace"`
	SelfLink  JobSelfLink `json:"self_link"`
}

type JobStatus struct {
	State     string             `json:"state"`
	Message   string             `json:"message"`
	Stats     JobStatusStats     `json:"stats"`
	Resources JobStatusResources `json:"resources"`
	Updated   time.Time          `json:"updated"`
}

type JobStatusStats struct {
	Total        int       `json:"total"`
	Active       int       `json:"active"`
	Successful   int       `json:"successful"`
	Failed       int       `json:"failed"`
	LastSchedule time.Time `json:"last_schedule"`
}

type JobStatusResources struct {
	Allocated ResourceItem `json:"allocated"`
	Total     ResourceItem `json:"total"`
}

type JobSpec struct {
	State       SpecState          `json:"state"`
	Enabled     bool               `json:"enabled"`
	Provider    JobSpecProvider    `json:"provider"`
	Hook        JobSpecHook        `json:"hook"`
	Concurrency JobSpecConcurrency `json:"concurrency"`
	Resources   ResourceRequest    `json:"resources"`
	Task        JobSpecTask        `json:"task"`
}

type JobSpecTask struct {
	Selector SpecSelector `json:"selector"`
	Runtime  SpecRuntime  `json:"runtime"`
	Template SpecTemplate `json:"template"`
}

type JobSpecConcurrency struct {
	Limit    int    `json:"limit"`
	Strategy string `json:"strategy"`
}

type JobSpecProvider struct {
	Kind    string                 `json:"kind"`
	Timeout int                    `json:"timeout"`
	Config  map[string]interface{} `json:"config"`
}
type JobSpecHook struct {
	Kind   string                 `json:"kind"`
	Config map[string]interface{} `json:"config"`
}

type JobSpecKindHttpConfig struct {
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

func (js *JobStatus) GetResourceAvailable() ResourceItem {

	var (
		RAM = js.Resources.Total.RAM - js.Resources.Allocated.RAM
		CPU = js.Resources.Total.CPU - js.Resources.Allocated.CPU
		STG = js.Resources.Total.Storage - js.Resources.Allocated.Storage
	)

	if RAM < 0 {
		RAM = 0
	}

	if CPU < 0 {
		CPU = 0
	}

	return ResourceItem{
		RAM: RAM, CPU: CPU, Storage: STG,
	}
}

func (js *JobSpec) GetResourceRequest() ResourceRequest {
	rr := ResourceRequest{}

	var (
		limitsRAM int64
		limitsCPU int64

		requestRAM int64
		requestCPU int64
	)

	for _, c := range js.Task.Template.Containers {

		limitsCPU += c.Resources.Limits.CPU
		limitsRAM += c.Resources.Limits.RAM

		requestCPU += c.Resources.Request.CPU
		requestRAM += c.Resources.Request.RAM
	}

	if requestRAM > 0 {
		requestRAM = int64(js.Concurrency.Limit) * requestRAM
		rr.Request.RAM = requestRAM
	}

	if requestCPU > 0 {
		requestCPU = int64(js.Concurrency.Limit) * requestCPU
		rr.Request.CPU = requestCPU
	}

	if limitsRAM > 0 {
		limitsRAM = int64(js.Concurrency.Limit) * limitsRAM
		rr.Limits.RAM = limitsRAM
	}

	if limitsCPU > 0 {
		limitsCPU = int64(js.Concurrency.Limit) * limitsCPU
		rr.Limits.CPU = limitsCPU
	}

	return rr
}

func (j *Job) SelfLink() *JobSelfLink {
	return &j.Meta.SelfLink
}

func (j *Job) AllocateResources(resources ResourceRequest) error {

	var (
		availableRam int64
		availableCpu int64

		allocatedRam int64
		allocatedCpu int64

		requestedRam int64
		requestedCpu int64
	)

	availableRam = j.Spec.Resources.Limits.RAM
	availableCpu = j.Spec.Resources.Limits.CPU

	allocatedRam = j.Status.Resources.Allocated.RAM
	allocatedCpu = j.Status.Resources.Allocated.CPU

	requestedRam = resources.Limits.RAM
	requestedCpu = resources.Limits.CPU

	if availableRam > 0 && availableCpu > 0 {

		if requestedRam == 0 {
			return errors.New(errors.ResourcesRamLimitIsRequired)
		}

		if requestedCpu == 0 {
			return errors.New(errors.ResourcesCpuLimitIsRequired)
		}

		if (availableRam - allocatedRam - requestedRam) <= 0 {
			return errors.New(errors.ResourcesRamLimitExceeded)
		}

		if (availableCpu - allocatedCpu - requestedCpu) <= 0 {
			return errors.New(errors.ResourcesCpuLimitExceeded)
		}
	}

	allocatedRam += requestedRam
	allocatedCpu += requestedCpu

	j.Status.Resources.Allocated.RAM = allocatedRam
	j.Status.Resources.Allocated.CPU = allocatedCpu

	return nil
}

func (j *Job) ReleaseResources(resources ResourceRequest) {

	var (
		availableRam int64
		availableCpu int64
		allocatedRam int64
		allocatedCpu int64
		requestedRam int64
		requestedCpu int64
	)

	availableRam = j.Spec.Resources.Limits.RAM
	availableCpu = j.Spec.Resources.Limits.CPU

	allocatedRam = j.Status.Resources.Allocated.RAM
	allocatedCpu = j.Status.Resources.Allocated.CPU

	requestedRam = resources.Limits.RAM
	requestedCpu = resources.Limits.CPU

	if (allocatedRam+requestedRam) > availableRam && (availableRam > 0) {
		allocatedRam = availableRam
	} else {
		allocatedRam -= requestedRam
	}

	if (allocatedCpu+requestedCpu) > availableCpu && (availableRam > 0) {
		allocatedCpu = availableCpu
	} else {
		allocatedCpu -= requestedCpu
	}

	j.Status.Resources.Allocated.RAM = allocatedRam
	j.Status.Resources.Allocated.CPU = allocatedCpu

}

func NewJobList() *JobList {
	jrl := new(JobList)
	jrl.Items = make([]*Job, 0)
	return jrl
}

func NewJobMap() *JobMap {
	jrl := new(JobMap)
	jrl.Items = make(map[string]*Job, 0)
	return jrl
}
