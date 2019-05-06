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
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const taskManifest = `
{
  "Meta": {
    "Name": "demo"
  },
  "Spec": {
    "Runtime": {
      "services": [
        "s_etcd",
        "s_dind"
      ],
      "tasks": [
        {
          "name": "clone:github.com/lastbackend/lastbackend",
          "container": "r_builder",
          "commands": [
            "lb clone -v github -o lastbackend -n lastbackend -b master /data/"
          ]
        },
        {
          "name": "step:test",
          "container": "p_test",
          "commands": [
            "apt-get -y install openssl",
            "mkdir -p ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "cp -r /data/. ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "cd ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "make deps",
            "make test"
          ]
        },
        {
          "name": "build:hub.lstbknd.net/lastbackend/lastbackend:master",
          "container": "r_builder",
          "commands": [
            "lb build -i hub.lstbknd.net/lastbackend/lastbackend:master -f ./images/lastbackend/Dockerfile .",
            "lb push hub.lstbknd.net/lastbackend/lastbackend:master"
          ]
        }
      ]
    },
    "Template": {
      "containers": [
        {
          "name": "s_etcd",
          "command": "/usr/local/bin/etcd --data-dir=/etcd-data --name node --initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 --advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://127.0.0.1:2379 --initial-cluster node=http://127.0.0.1:2380",
          "image": {
            "name": "quay.io/coreos/etcd:latest"
          }
        },
        {
          "name": "p_test",
          "workdir": "/data/",
          "env": [
            {
              "name": "ENV_GIT_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "github"
              }
            },
            {
              "name": "ENV_DOCKER_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "docker"
              }
            },
            {
              "name": "DOCKER_HOST",
              "value": "tcp://127.0.0.1:2375"
            }
          ],
          "volumes": [
            {
              "name": "data",
              "path": "/data/"
            }
          ],
          "image": {
            "name": "golang:stretch"
          },
          "security": {
            "privileged": false
          }
        },
        {
          "name": "s_dind",
          "image": {
            "name": "docker:dind"
          },
          "security": {
            "privileged": true
          }
        },
        {
          "name": "r_builder",
          "workdir": "/data/",
          "env": [
            {
              "name": "ENV_GIT_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "github"
              }
            },
            {
              "name": "ENV_DOCKER_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "docker"
              }
            },
            {
              "name": "DOCKER_HOST",
              "value": "tcp://127.0.0.1:2375"
            }
          ],
          "volumes": [
            {
              "name": "data",
              "path": "/data/"
            }
          ],
          "image": {
            "name": "index.lstbknd.net/lastbackend/builder"
          }
        }
      ],
      "volumes": [
        {
          "name": "data"
        }
      ]
    }
  }
}`

func TestTaskCreate(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			j  *types.Job
			mf *types.TaskManifest
		}
		want     *types.Task
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		job := getJobAsset("demo", "system")
		manifest := getTaskManifestAsset("demo")

		ctx := context.Background()

		s := suit{
			ctx:  ctx,
			name: "successful task creation",
			args: struct {
				j  *types.Job
				mf *types.TaskManifest
			}{
				j:  job,
				mf: manifest,
			},
			want: &types.Task{
				Meta: types.TaskMeta{
					Namespace: job.Meta.SelfLink.Namespace().String(),
					Job:       job.Meta.SelfLink.String(),
					SelfLink:  *types.NewTaskSelfLink(job.Meta.Namespace, job.Meta.Name, *manifest.Meta.Name),
				},
				Status: types.TaskStatus{
					State: types.StateCreated,
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Task(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			task, err := taskCreate(tc.args.j, tc.args.mf)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task create error")
					return
				}
				return
			}

			if tc.want != nil && task == nil {
				assert.NotNil(t, task, "task can not be nil")
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, task.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, task.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), task.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, task.Status.State, "task state mismatch")

			var stgTask = new(types.Task)

			if err := stg.Get(tc.ctx, stg.Collection().Task(), task.SelfLink().String(), stgTask, nil); err != nil {
				t.Log(err)
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, stgTask.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, stgTask.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), stgTask.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, stgTask.Status.State, "task state mismatch")
		})
	}

}

func TestTaskQueue(t *testing.T) {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			js *JobState
			t  *types.Task
		}
		want struct {
			js  types.JobStatus
			jts JobTaskState
			jps JobPodState
			t   *types.Task
		}
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateCreated

		s := suit{
			ctx:  ctx,
			name: "set task as queue state",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: NewJobState(cluster.NewClusterState(), job),
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateWaiting,
				},
				jts: JobTaskState{Active: 0, Queue: 0, List: 0, Finished: 0},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateQueued,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateQueued

		s := suit{
			ctx:  ctx,
			name: "add task in queue and set in provision state",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: NewJobState(cluster.NewClusterState(), job),
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateRunning,
				},
				jts: JobTaskState{Active: 0, Queue: 1, List: 0, Finished: 0},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateProvision,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskQueue(tc.args.js, tc.args.t)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task queue error")
					return
				}
				return
			}

			assert.Equal(t, tc.want.js.State, tc.args.js.job.Status.State, "job state mismatch")

			assert.Equal(t, tc.want.jts.Active, len(tc.args.js.task.active), "job state task active counts mismatch")
			assert.Equal(t, tc.want.jts.Queue, len(tc.args.js.task.queue), "job state task queue counts mismatch")
			assert.Equal(t, tc.want.jts.List, len(tc.args.js.task.list), "job state task list counts mismatch")
			assert.Equal(t, tc.want.jts.Finished, len(tc.args.js.task.finished), "job state task finish counts mismatch")

			assert.Equal(t, tc.want.t.Status.State, tc.args.t.Status.State, "task state mismatch")

			if len(tc.args.js.task.queue) > 0 {
				var stgJob = new(types.Job)

				if err := stg.Get(tc.ctx, stg.Collection().Job(), tc.args.t.Meta.Job, stgJob, nil); err != nil {
					t.Log(err)
					return
				}

				assert.Equal(t, tc.want.js.State, stgJob.Status.State, "job state mismatch")
			}

		})
	}
}

func TestTaskProvision(t *testing.T) {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			js *JobState
			t  *types.Task
		}
		want struct {
			js  types.JobStatus
			jts JobTaskState
			jps JobPodState
			t   *types.Task
		}
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateQueued

		s := suit{
			ctx:  ctx,
			name: "set task as provision state",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: NewJobState(cluster.NewClusterState(), job),
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateWaiting,
				},
				jts: JobTaskState{Active: 0, Queue: 0, List: 0, Finished: 0},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateProvision,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Task(), "")
				if !assert.NoError(t, err) {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		jobState := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		jobState.task.queue[task.SelfLink().String()] = task
		task.Status.State = types.StateProvision

		s := suit{
			ctx:  ctx,
			name: "add task in active and create pod from task",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: jobState,
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateRunning,
				},
				jts: JobTaskState{Active: 1, Queue: 0, List: 0, Finished: 0},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateProvision,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Task(), "")
				if !assert.NoError(t, err) {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		jobState := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision
		jobState.task.active[task.SelfLink().String()] = task

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		jobState.pod.list[task.SelfLink().String()] = pod

		s := suit{
			ctx:  context.Background(),
			name: "check pod exists for task",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: jobState,
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateRunning,
				},
				jts: JobTaskState{Active: 1, Queue: 0, List: 0, Finished: 0},
				jps: JobPodState{List: 1},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateProvision,
					},
				},
			},
			wantErr: true,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Task(), "")
				if !assert.NoError(t, err) {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskProvision(tc.args.js, tc.args.t)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task provision error")
					return
				}
				return
			}

			assert.Equal(t, tc.want.js.State, tc.args.js.job.Status.State, "job state mismatch")

			assert.Equal(t, tc.want.jts.Active, len(tc.args.js.task.active), "job state task active counts mismatch")
			assert.Equal(t, tc.want.jts.Queue, len(tc.args.js.task.queue), "job state task queue counts mismatch")
			assert.Equal(t, tc.want.jts.List, len(tc.args.js.task.list), "job state task list counts mismatch")
			assert.Equal(t, tc.want.jts.Finished, len(tc.args.js.task.finished), "job state task finish counts mismatch")

			assert.Equal(t, tc.want.t.Status.State, tc.args.t.Status.State, "task state mismatch")

		})
	}
}

func TestTaskDestroy(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			j    *types.Job
			js   *JobState
			task *types.Task
		}
		want     *types.Task
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("demo", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision
		jobState := NewJobState(cluster.NewClusterState(), job)
		jobState.pod.list = make(map[string]*types.Pod, 0)

		s := suit{
			ctx:  ctx,
			name: "set task to destroy state",
			args: struct {
				j    *types.Job
				js   *JobState
				task *types.Task
			}{
				j:    job,
				js:   jobState,
				task: task,
			},
			want: &types.Task{
				Meta: types.TaskMeta{
					Namespace: task.Meta.Namespace,
					Job:       task.Meta.Job,
					SelfLink:  task.Meta.SelfLink,
				},
				Status: types.TaskStatus{
					State: types.StateDestroy,
				},
			},
			wantErr: false,
			preHook: func() error {
				if err := stg.Set(ctx, stg.Collection().Task(), task.SelfLink().String(), task, nil); err != nil {
					return err
				}
				return nil
			},
			postHook: func() error {
				if err := stg.Del(context.Background(), stg.Collection().Task(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("demo", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateDestroy
		jobState := NewJobState(cluster.NewClusterState(), job)
		jobState.pod.list = make(map[string]*types.Pod, 0)

		s := suit{
			ctx:  ctx,
			name: "destroy task to destroyed",
			args: struct {
				j    *types.Job
				js   *JobState
				task *types.Task
			}{
				j:    job,
				js:   jobState,
				task: task,
			},
			want: &types.Task{
				Meta: types.TaskMeta{
					Namespace: task.Meta.Namespace,
					Job:       task.Meta.Job,
					SelfLink:  task.Meta.SelfLink,
				},
				Status: types.TaskStatus{
					State: types.StateDestroyed,
				},
			},
			wantErr: false,
			preHook: func() error {
				if err := stg.Set(ctx, stg.Collection().Task(), task.SelfLink().String(), task, nil); err != nil {
					return err
				}
				return nil
			},
			postHook: func() error {
				if err := stg.Del(context.Background(), stg.Collection().Task(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("demo", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateDestroy
		jobState := NewJobState(cluster.NewClusterState(), job)
		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		jobState.pod.list[task.SelfLink().String()] = pod

		s := suit{
			ctx:  ctx,
			name: "destroy task",
			args: struct {
				j    *types.Job
				js   *JobState
				task *types.Task
			}{
				j:    job,
				js:   jobState,
				task: task,
			},
			want: &types.Task{
				Meta: types.TaskMeta{
					Namespace: task.Meta.Namespace,
					Job:       task.Meta.Job,
					SelfLink:  task.Meta.SelfLink,
				},
				Status: types.TaskStatus{
					State: types.StateDestroy,
				},
			},
			wantErr: false,
			preHook: func() error {
				if err := stg.Set(ctx, stg.Collection().Task(), task.SelfLink().String(), task, nil); err != nil {
					return err
				}
				return nil
			},
			postHook: func() error {
				if err := stg.Del(context.Background(), stg.Collection().Task(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskDestroy(tc.args.js, tc.args.task)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}
				return
			}

			var stgTask = new(types.Task)

			if err := stg.Get(tc.ctx, stg.Collection().Task(), tc.args.task.SelfLink().String(), stgTask, nil); err != nil {
				if tc.want == nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
					return
				}
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, stgTask.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, stgTask.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), stgTask.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, stgTask.Status.State, "task state mismatch")

		})
	}

}

func TestTaskUpdate(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			j    *types.Job
			task *types.Task
		}
		want     *types.Task
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("demo", "system")
		taskOld := getTaskAsset("demo", "system", job.Meta.Name)

		taskNew := types.Task(*taskOld)
		taskNew.Status.State = types.StateDestroy

		s := suit{
			ctx:  ctx,
			name: "successful task remove",
			args: struct {
				j    *types.Job
				task *types.Task
			}{
				j:    job,
				task: &taskNew,
			},
			want:    &taskNew,
			wantErr: false,
			preHook: func() error {
				if err := stg.Set(ctx, stg.Collection().Task(), taskOld.SelfLink().String(), taskOld, nil); err != nil {
					return err
				}
				return nil
			},
			postHook: func() error {
				if err := stg.Del(context.Background(), stg.Collection().Task(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			ts := tc.args.task.Meta.Updated
			tc.args.task.Meta.Updated = time.Now()

			err := taskUpdate(tc.args.task, ts)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}
				return
			}

			var stgTask = new(types.Task)

			if err := stg.Get(tc.ctx, stg.Collection().Task(), tc.args.task.SelfLink().String(), stgTask, nil); err != nil {
				if tc.want == nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
					return
				}
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, stgTask.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, stgTask.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), stgTask.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, stgTask.Status.State, "task state mismatch")

		})
	}

}

func TestTaskRemove(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			j    *types.Job
			task *types.Task
		}
		want     *types.Task
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("demo", "system")
		task := getTaskAsset("demo", "system", job.Meta.Name)

		s := suit{
			ctx:  ctx,
			name: "successful task remove",
			args: struct {
				j    *types.Job
				task *types.Task
			}{
				j:    job,
				task: task,
			},
			want:    nil,
			wantErr: false,
			preHook: func() error {
				if err := stg.Set(ctx, stg.Collection().Task(), task.SelfLink().String(), task, nil); err != nil {
					return err
				}
				return nil
			},
			postHook: func() error {
				if err := stg.Del(context.Background(), stg.Collection().Task(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskRemove(tc.args.task)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}
				return
			}

			var task = new(types.Task)
			if err := stg.Get(tc.ctx, stg.Collection().Task(), tc.args.task.SelfLink().String(), task, nil); err != nil {
				if tc.want == nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
					return
				}
				return
			}

		})
	}

}

func TestTaskFinish(t *testing.T) {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			js *JobState
			t  *types.Task
		}
		want struct {
			js  types.JobStatus
			jts JobTaskState
			jps JobPodState
			t   *types.Task
		}
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision

		s := suit{
			ctx:  ctx,
			name: "set task as queue state",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: NewJobState(cluster.NewClusterState(), job),
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateRunning,
				},
				jts: JobTaskState{Active: 0, Queue: 0, List: 0, Finished: 1},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateExited,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}
				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)

		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		js.pod.list[task.SelfLink().String()] = pod

		s := suit{
			ctx:  ctx,
			name: "set task as queue state",
			args: struct {
				js *JobState
				t  *types.Task
			}{
				js: js,
				t:  task,
			},
			want: struct {
				js  types.JobStatus
				jts JobTaskState
				jps JobPodState
				t   *types.Task
			}{
				js: types.JobStatus{
					State: types.StateRunning,
				},
				jts: JobTaskState{Active: 0, Queue: 0, List: 0, Finished: 5},
				jps: JobPodState{List: 0},
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateExited,
					},
				},
			},
			wantErr: false,
			preHook: func() error {

				task1 := getTaskAsset("demo1", "system", job.Meta.Name)
				task1.Status.State = types.StateExited

				task2 := getTaskAsset("demo2", "system", job.Meta.Name)
				task2.Status.State = types.StateExited

				task3 := getTaskAsset("demo3", "system", job.Meta.Name)
				task3.Status.State = types.StateExited

				task4 := getTaskAsset("demo4", "system", job.Meta.Name)
				task4.Status.State = types.StateExited

				task5 := getTaskAsset("demo5", "system", job.Meta.Name)
				task5.Status.State = types.StateExited

				js.task.finished = append(js.task.finished, task1)
				js.task.finished = append(js.task.finished, task2)
				js.task.finished = append(js.task.finished, task3)
				js.task.finished = append(js.task.finished, task4)
				js.task.finished = append(js.task.finished, task5)

				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskFinish(tc.args.js, tc.args.t)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task finish error")
					return
				}
				return
			}

			assert.Equal(t, tc.want.js.State, tc.args.js.job.Status.State, "job state mismatch")

			assert.Equal(t, tc.want.jts.Active, len(tc.args.js.task.active), "job state task active counts mismatch")
			assert.Equal(t, tc.want.jts.Queue, len(tc.args.js.task.queue), "job state task queue counts mismatch")
			assert.Equal(t, tc.want.jts.List, len(tc.args.js.task.list), "job state task list counts mismatch")
			assert.Equal(t, tc.want.jts.Finished, len(tc.args.js.task.finished), "job state task finish counts mismatch")

			assert.Equal(t, tc.want.t.Status.State, tc.args.t.Status.State, "task state mismatch")

			_, ok := tc.args.js.pod.list[tc.args.t.SelfLink().String()]
			assert.False(t, ok)

		})
	}
}

func TestTaskStatusState(t *testing.T) {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			js *JobState
			t  *types.Task
			p  *types.Pod
		}
		want struct {
			t *types.Task
		}
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateCreated

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateProvision

		s := suit{
			ctx:  ctx,
			name: "set task created -> provision state (pod: provision)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateProvision,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateDestroy

		s := suit{
			ctx:  ctx,
			name: "set task provision -> exited state (pod: destroy)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateExited,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateProvision

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateError

		s := suit{
			ctx:  ctx,
			name: "set task provision -> error state (pod: error)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateError,
						Error: true,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateCanceled
		task.Status.Canceled = true

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateExited

		s := suit{
			ctx:  ctx,
			name: "set task canceled -> exited state (pod: exited)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State:    types.StateExited,
						Canceled: true,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateError
		task.Status.Error = true

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateExited

		s := suit{
			ctx:  ctx,
			name: "set task error -> exited state (pod: exited)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State: types.StateExited,
						Error: true,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateCanceled
		task.Status.Canceled = true

		pod, err := podCreate(task)
		if err != nil {
			t.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
		}
		assert.NotNil(t, pod)
		pod.Status.State = types.StateExited

		s := suit{
			ctx:  ctx,
			name: "set task canceled -> exited state (pod: exited)",
			args: struct {
				js *JobState
				t  *types.Task
				p  *types.Pod
			}{
				js: js,
				t:  task,
				p:  pod,
			},
			want: struct {
				t *types.Task
			}{
				t: &types.Task{
					Meta: types.TaskMeta{
						Namespace: task.Meta.Namespace,
						Job:       task.Meta.Job,
						SelfLink:  task.Meta.SelfLink,
					},
					Status: types.TaskStatus{
						State:    types.StateExited,
						Canceled: true,
					},
				},
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {

				if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
					return err
				}

				err = stg.Del(ctx, stg.Collection().Pod(), "")
				if !assert.NoError(t, err) {
					return err
				}

				return nil
			},
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskStatusState(tc.args.js, tc.args.t, tc.args.p)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task finish error")
					return
				}
				return
			}

			assert.Equal(t, tc.want.t.Status.State, tc.args.t.Status.State, "task state mismatch")
			assert.Equal(t, tc.want.t.Status.Canceled, tc.args.t.Status.Canceled, "task canceled mismatch")
			assert.Equal(t, tc.want.t.Status.Error, tc.args.t.Status.Error, "task error mismatch")

		})
	}
}

func getJobAsset(name, namespace string) *types.Job {
	j := new(types.Job)
	j.Meta.SetDefault()
	j.Meta.Name = name
	j.Meta.Namespace = types.NewNamespaceSelfLink(namespace).String()
	j.Meta.SelfLink = *types.NewJobSelfLink(namespace, name)
	j.Meta.Labels = map[string]string{"app": "lb", "type": "job"}
	j.Status.State = types.StateWaiting
	j.Spec.Enabled = true

	return j
}

func getTaskManifestAsset(name string) *types.TaskManifest {
	t := new(types.TaskManifest)
	json.Unmarshal([]byte(taskManifest), t)

	if name != types.EmptyString {
		t.Meta.Name = &name
	}

	return t
}

func getTaskAsset(name, namespace, job string) *types.Task {
	t := new(types.Task)
	t.Meta.SetDefault()
	t.Meta.Name = name
	t.Meta.Namespace = types.NewNamespaceSelfLink(namespace).String()
	t.Meta.Job = types.NewJobSelfLink(namespace, job).String()
	t.Meta.SelfLink = *types.NewTaskSelfLink(namespace, job, name)
	return t
}
