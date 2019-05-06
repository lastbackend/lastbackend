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
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"testing"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPodCreate(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			task *types.Task
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
			name: "successful pod creating",
			args: struct {
				task *types.Task
			}{
				task: task,
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Pod(), "")
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

			pod, err := podCreate(tc.args.task)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "pod create error")
					return
				}
				return
			}

			if pod == nil {
				assert.NotNil(t, pod, "pod can not be nil")
				return
			}

			var stgPod = new(types.Task)
			if err := stg.Get(tc.ctx, stg.Collection().Pod(), pod.SelfLink().String(), stgPod, nil); err != nil {
				t.Log(err)
				return
			}

			assert.Equal(t, pod.Meta.Namespace, stgPod.Meta.Namespace, "pod namespace mismatch")
			assert.Equal(t, pod.Meta.SelfLink.String(), stgPod.Meta.SelfLink.String(), "pod self_link mismatch")
			assert.Equal(t, pod.Status.State, stgPod.Status.State, "pod state mismatch")

		})
	}

}

func TestPodProvision(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			js  *JobState
			pod *types.Pod
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
		jobState := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset("demo", "system", job.Meta.Name)
		task.Status.State = types.StateCreated
		pod := getPodAsset(task, "demo")

		s := suit{
			ctx:  ctx,
			name: "set pod to provision state",
			args: struct {
				js  *JobState
				pod *types.Pod
			}{
				js:  jobState,
				pod: pod,
			},
			wantErr: false,
			preHook: func() error {
				return nil
			},
			postHook: func() error {
				err := stg.Del(ctx, stg.Collection().Pod(), "")
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

			err := podProvision(tc.args.js, tc.args.pod)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "pod create error")
					return
				}
				return
			}

			// TODO

		})
	}

}

func getPodAsset(task *types.Task, name string) *types.Pod {
	p := new(types.Pod)
	p.Meta.SetDefault()
	p.Meta.Name = name
	p.Meta.Namespace = task.SelfLink().Namespace().String()
	sl, _ := types.NewPodSelfLink(types.KindTask, task.SelfLink().String(), name)
	p.Meta.SelfLink = *sl
	return p
}
