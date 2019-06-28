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

package task_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/task"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Testing TaskInfoH handler
func TestTaskInfo(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	j1 := getJobAsset(ns1.Meta.Name, "job1", "")
	j2 := getJobAsset(ns1.Meta.Name, "job2", "")

	t1 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")
	t2 := getTaskAsset(ns1.Meta.Name, j2.Meta.Name, "task2")
	t3 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task3")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		task      *types.Task
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Task
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking get task if namespace not exists",
			handler:      task.TaskInfoH,
			args:         args{ctx, ns2, t1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get task if job not exists",
			handler:      task.TaskInfoH,
			args:         args{ctx, ns1, t2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Job not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get task if not exists",
			handler:      task.TaskInfoH,
			args:         args{ctx, ns1, t3},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Task not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get task successfully",
			handler:      task.TaskInfoH,
			args:         args{ctx, ns1, t1},
			fields:       fields{stg},
			want:         v1.View().Task().New(t1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Job(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Task(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j1.SelfLink().String(), j1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Task(), t1.SelfLink().String(), t1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.args.namespace.Meta.Name, tc.args.task.Meta.Job, tc.args.task.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/job/{job}/task/{task}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				j := new(views.Task)
				err := json.Unmarshal(body, &j)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, j.Meta.Name, "name not equal")
			}

		})
	}

}

// Testing TaskListH handler
func TestTaskList(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("ns1", "")
	ns2 := getNamespaceAsset("ns2", "")

	j1 := getJobAsset(ns1.Meta.Name, "job1", "")
	j2 := getJobAsset(ns1.Meta.Name, "job2", "")

	t1 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")
	t2 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task2")

	tl := types.NewTaskMap()
	tl.Items[t1.SelfLink().String()] = t1
	tl.Items[t2.SelfLink().String()] = t2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		job       *types.Job
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         *types.TaskMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get tasks list if namespace not found",
			args:         args{ctx, ns2, j1},
			fields:       fields{stg},
			handler:      task.TaskListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get tasks list if job not found",
			args:         args{ctx, ns1, j2},
			fields:       fields{stg},
			handler:      task.TaskListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Job not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get tasks list successfully",
			args:         args{ctx, ns1, j1},
			fields:       fields{stg},
			handler:      task.TaskListH,
			want:         tl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Job(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Task(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j1.SelfLink().String(), j1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Task(), t1.SelfLink().String(), t1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Task(), t2.SelfLink().String(), t2, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/job/%s/task", tc.args.namespace.Meta.Name, tc.args.job.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/job/{job}/task", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.RouteList)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				for _, item := range *s {
					if _, ok := tc.want.Items[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
			}

		})
	}

}

// Testing TaskCreateH handler
func TestTaskCreate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("not-found", "")
	ns3 := getNamespaceAsset("limits", "")

	ns3.Spec.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("2GB")
	ns3.Spec.Resources.Limits.CPU, _ = resource.DecodeCpuResource("2")

	j1 := getJobAsset(ns1.Meta.Name, "job1", "")
	j2 := getJobAsset(ns1.Meta.Name, "job2", "")
	j3 := getJobAsset(ns3.Meta.Name, "job3", "")

	j3.Spec.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("1GB")
	j3.Spec.Resources.Limits.CPU, _ = resource.DecodeCpuResource("0.8")

	t1 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "errored")
	t2 := getTaskAsset(ns3.Meta.Name, j3.Meta.Name, "success")
	t3 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "success")

	tm1 := getTaskManifest("errored", "image")
	tm1.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm1.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm1.Spec.Template.Containers[0].Resources.Limits.RAM = "0.5GB"

	tm2 := getTaskManifest("errored", "image")
	tm2.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm2.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm2.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	tm2.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	tm3 := getTaskManifest("errored", "image")
	tm3.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm3.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm3.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	tm3.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	tm4 := getTaskManifest("errored", "image")
	tm4.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm4.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm4.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	tm4.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	tm5 := getTaskManifest("errored", "image")
	tm5.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm5.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm5.Spec.Template.Containers[0].Resources.Limits.RAM = "128MB"
	tm5.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	tm6 := getTaskManifest("success", "image")
	tm6.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm6.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm6.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	tm6.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	tm7 := getTaskManifest("success", "image")
	tm7.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	tm7.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	tm7.Spec.Template.Containers[0].Resources.Limits.RAM = ""
	tm7.Spec.Template.Containers[0].Resources.Limits.CPU = ""

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		job       *types.Job
		tmf       *request.TaskManifest
		task      *types.Task
	}

	tests := []struct {
		name         string
		args         args
		fields       fields
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         *request.TaskManifest
		want         *views.Task
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking create task with not existed namespace",
			args:         args{ctx, ns2, j1, tm1, t1},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         getTaskManifest("task", "redis"),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking create task not existed job",
			args:         args{ctx, ns1, j2, tm1, t1},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         getTaskManifest("task", "redis"),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Job not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking create task with ram limit is higher then job limits",
			args:         args{ctx, ns3, j3, tm2, t1},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         getTaskManifest("task", "redis"),
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create task with cpu limit is higher then job limits",
			args:         args{ctx, ns3, j3, tm3, t1},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         getTaskManifest("task", "redis"),
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create task with ram & cpu limits are higher then job limits",
			args:         args{ctx, ns3, j3, tm4, t1},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         getTaskManifest("task", "redis"),
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check create task success with limits",
			args:         args{ctx, ns3, j3, tm6, t2},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         tm6,
			want:         v1.View().Task().New(t2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "check create task success without limits",
			args:         args{ctx, ns1, j1, tm7, t3},
			fields:       fields{stg},
			handler:      task.TaskCreateH,
			data:         tm6,
			want:         v1.View().Task().New(t3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Job(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Task(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns3.SelfLink().String(), ns3, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j1.SelfLink().String(), j1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j3.SelfLink().String(), j3, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			bd, err := tc.args.tmf.ToJson()
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/job/%s/task", tc.args.namespace.Meta.Name, tc.args.job.Meta.Name), strings.NewReader(string(bd)))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/job/{job}/task", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				job := new(types.Job)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Job(), tc.args.job.SelfLink().String(), job, nil)
				if !assert.NoError(t, err) {
					return
				}

				got := new(types.Task)
				err = tc.fields.stg.Get(tc.args.ctx, stg.Collection().Task(), tc.args.task.SelfLink().String(), got, nil)
				if !assert.NoError(t, err) {
					return
				}

				if got == nil {
					t.Error("can not be not nil")
					return
				}

				assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "name not equal")
				assert.Equal(t, tc.want.Meta.Description, got.Meta.Description, "description not equal")
			}
		})
	}

}

// Testing TaskUpdateH handler
func TestTaskCancel(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("exists", "")
	ns2 := getNamespaceAsset("not-exist", "")

	j1 := getJobAsset(ns1.Meta.Name, "job1", "")
	j2 := getJobAsset(ns1.Meta.Name, "job2", "")

	t1 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")
	t2 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")
	t2.Spec.State.Cancel = true
	t3 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task3")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		job       *types.Job
		task      *types.Task
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Task
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking cancel task in not existing namespace",
			args:         args{ctx, ns2, j1, t1},
			fields:       fields{stg},
			handler:      task.TaskCancelH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking cancel task in not existing job",
			args:         args{ctx, ns1, j2, t1},
			fields:       fields{stg},
			handler:      task.TaskCancelH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Job not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking cancel task in not existing task",
			args:         args{ctx, ns1, j1, t3},
			fields:       fields{stg},
			handler:      task.TaskCancelH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Task not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		// TODO: check another spec parameters
		{
			name:         "check cancel task success",
			args:         args{ctx, ns1, j1, t1},
			fields:       fields{stg},
			handler:      task.TaskCancelH,
			want:         v1.View().Task().New(t1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Job(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Task(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j1.SelfLink().String(), j1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Task(), t1.SelfLink().String(), tc.args.task, nil)
			assert.NoError(t, err)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.args.namespace.Meta.Name, tc.args.job.Meta.Name, tc.args.task.Meta.Name), strings.NewReader(""))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/job/{job}/task/{task}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {
				s := new(views.Task)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "description not equal")
				assert.Equal(t, types.StateCanceled, s.Status.State, "status state is not canceled")
			}
		})
	}

}

// Testing TaskRemoveH handler
func TestTaskRemove(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("exists", "")
	ns2 := getNamespaceAsset("not-exist", "")

	j1 := getJobAsset(ns1.Meta.Name, "job1", "")
	j2 := getJobAsset(ns1.Meta.Name, "job2", "")

	t1 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")
	t2 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task1")

	t2.Status.State = types.StateDestroy
	t2.Spec.State.Destroy = true

	t3 := getTaskAsset(ns1.Meta.Name, j1.Meta.Name, "task3")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		job       *types.Job
		task      *types.Task
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Task
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking cancel task in not existing namespace",
			args:         args{ctx, ns2, j1, t1},
			fields:       fields{stg},
			handler:      task.TaskRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking cancel task in not existing job",
			args:         args{ctx, ns1, j2, t1},
			fields:       fields{stg},
			handler:      task.TaskRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Job not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking remove task in not existing task",
			args:         args{ctx, ns1, j1, t3},
			fields:       fields{stg},
			handler:      task.TaskRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Task not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		// TODO: check another spec parameters
		{
			name:         "check remove task success",
			args:         args{ctx, ns1, j1, t1},
			fields:       fields{stg},
			handler:      task.TaskRemoveH,
			want:         v1.View().Task().New(t2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Job(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Task(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Job(), j1.SelfLink().String(), j1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Task(), t1.SelfLink().String(), tc.args.task, nil)
			assert.NoError(t, err)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.args.namespace.Meta.Name, tc.args.job.Meta.Name, tc.args.task.Meta.Name), strings.NewReader(""))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/job/{job}/task/{task}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {
				s := new(views.Task)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "description not equal")
				assert.Equal(t, types.StateDestroy, s.Status.State, "status state is not canceled")
			}
		})
	}

}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}
	n.Meta.SetDefault()
	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)
	return &n
}

func getJobAsset(namespace, name, desc string) *types.Job {
	var s = types.Job{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	s.Meta.SelfLink = *types.NewJobSelfLink(namespace, name)
	s.Spec.Task.Template.Containers = make(types.SpecTemplateContainers, 0)
	s.Spec.Task.Template.Containers = append(s.Spec.Task.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})
	return &s
}

func getTaskAsset(namespace, job, name string) *types.Task {
	var s = types.Task{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Job = job
	s.Meta.Name = name
	s.Meta.SelfLink = *types.NewTaskSelfLink(namespace, job, name)
	s.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	s.Spec.Template.Containers = append(s.Spec.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})
	return &s
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}

func getTaskManifest(name, image string) *request.TaskManifest {

	var (
		container = request.ManifestSpecTemplateContainer{
			Name: image,
			Image: &request.ManifestSpecTemplateContainerImage{
				Name: image,
			},
			Env: make([]request.ManifestSpecTemplateContainerEnv, 0),
		}
		volume = request.ManifestSpecTemplateVolume{
			Name: "demo",
			Secret: &request.ManifestSpecTemplateSecretVolume{
				Name:  "test",
				Binds: make([]request.ManifestSpecTemplateSecretVolumeBind, 0),
			},
		}
	)

	volume.Secret.Binds = append(volume.Secret.Binds, request.ManifestSpecTemplateSecretVolumeBind{
		Key:  "demo",
		File: "test.txt",
	})

	container.Env = append(container.Env, request.ManifestSpecTemplateContainerEnv{
		Name:  "Demo",
		Value: "test",
	})

	container.Env = append(container.Env, request.ManifestSpecTemplateContainerEnv{
		Name: "Secret",
		Secret: &request.ManifestSpecTemplateContainerEnvSecret{
			Name: "secret-name",
			Key:  "secret-key",
		},
	})

	mf := new(request.TaskManifest)
	mf.Meta.Name = &name
	mf.Spec.Template = new(request.ManifestSpecTemplate)
	mf.Spec.Template.Containers = append(mf.Spec.Template.Containers, container)
	mf.Spec.Template.Volumes = append(mf.Spec.Template.Volumes, volume)
	return mf
}
