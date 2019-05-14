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

package job

import (
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
)

func init() {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)
}

func getJobAsset(state, message string) *types.Job {
	j := new(types.Job)

	j.Meta.Namespace = "test"
	j.Meta.Name = "job"
	j.Meta.SelfLink = *types.NewJobSelfLink(j.Meta.Namespace, j.Meta.Name)

	j.Status.State = state
	j.Status.Message = message

	j.Spec.Enabled = true
	j.Spec.Concurrency.Limit = 1

	return j
}

func getJobStateAsset(job *types.Job) *JobState {

	n := new(types.Node)

	n.Meta.Name = "node"
	n.Meta.Hostname = "node.local"
	n.Status.Capacity = types.NodeResources{
		Containers: 10,
		Pods:       10,
		RAM:        1000,
		CPU:        1,
		Storage:    1000,
	}
	n.Meta.SelfLink = *types.NewNodeSelfLink(n.Meta.Hostname)

	cs := cluster.NewClusterState()
	cs.SetNode(n)
	s := NewJobState(cs, job)

	return s
}

func getJobStateCopy(js *JobState) *JobState {

	j := *js.job

	njs := NewJobState(js.cluster, &j)

	njs.task.list = make(map[string]*types.Task, 0)
	for k, t := range js.task.list {
		njs.task.list[k] = &(*t)
	}

	njs.task.active = make(map[string]*types.Task, 0)
	for k, t := range js.task.active {
		task := *t
		njs.task.active[k] = &task
	}

	njs.task.queue = make(map[string]*types.Task, 0)
	for k, t := range js.task.queue {
		task := *t
		njs.task.queue[k] = &task
	}

	njs.task.finished = make([]*types.Task, 0)
	for _, t := range js.task.finished {
		task := *t
		njs.task.finished = append(njs.task.finished, &task)
	}

	njs.pod.list = make(map[string]*types.Pod, 0)
	for k, p := range js.pod.list {
		pod := *p
		njs.pod.list[k] = &pod
	}

	return njs
}
