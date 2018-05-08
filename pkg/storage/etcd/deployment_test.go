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

package etcd

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestDeploymentStorage_Get(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		d   = getDeploymentAsset(ns1, svc, "test", "")
	)

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx  context.Context
		name string
		ns   string
		svc  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"get deployment info failed",
			fields{stg},
			args{ctx, "test2", ns1, svc},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get deployment info successful",
			fields{stg},
			args{ctx, "test", ns1, svc},
			&d,
			false,
			"",
		},
		{
			"get deployment info failed empty namespace",
			fields{stg},
			args{ctx, "test", "", svc},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get deployment info failed empty service",
			fields{stg},
			args{ctx, "test", ns1, ""},
			&d,
			true,
			"service can not be empty",
		},
		{
			"get deployment info failed empty name",
			fields{stg},
			args{ctx, "", ns1, svc},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.svc, tt.args.name)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !compareDeployments(got, tt.want) {
				t.Errorf("DeploymentStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestDeploymentStorage_ListByNamespace(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test1", "")
		n2  = getDeploymentAsset(ns1, svc, "test2", "")
		n3  = getDeploymentAsset(ns2, svc, "test1", "")
		nl  = make(map[string]*types.Deployment, 0)
	)

	nl0 := map[string]*types.Deployment{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Deployment{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Deployment{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Deployment
		wantErr bool
		err     string
	}{
		{
			"get namespace list 1 success",
			fields{stg},
			args{ctx, ns1},
			nl1,
			false,
			"",
		},
		{
			"get namespace list 2 success",
			fields{stg},
			args{ctx, ns2},
			nl2,
			false,
			"",
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "empty"},
			nl,
			false,
			"",
		},
		{
			"get namespace list fail empty namespace",
			fields{stg},
			args{ctx, ""},
			nl,
			true,
			"namespace can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("DeploymentStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && (err.Error() != tt.err) {
					t.Errorf("DeploymentStorage.ListByNamespace() error = %v, want err %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("DeploymentStorage.ListByNamespace() want error = %v, got none", tt.err)
				return
			}

			if !compareDeploymentMaps(got, tt.want) {
				t.Errorf("DeploymentStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentStorage_ListByService(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		sv1 = "svc1"
		sv2 = "svc2"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, sv1, "test1", "")
		n2  = getDeploymentAsset(ns1, sv1, "test2", "")
		n3  = getDeploymentAsset(ns1, sv2, "test1", "")
		n4  = getDeploymentAsset(ns2, sv1, "test1", "")
		n5  = getDeploymentAsset(ns2, sv1, "test2", "")
		nl  = make(map[string]*types.Deployment, 0)
	)

	nl0 := map[string]*types.Deployment{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3
	nl0[stg.keyGet(&n4)] = &n4
	nl0[stg.keyGet(&n5)] = &n5

	nl1 := map[string]*types.Deployment{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Deployment{}
	nl2[stg.keyGet(&n3)] = &n3

	nl3 := map[string]*types.Deployment{}
	nl3[stg.keyGet(&n4)] = &n4
	nl3[stg.keyGet(&n5)] = &n5

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx context.Context
		ns  string
		svc string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Deployment
		wantErr bool
		err     string
	}{
		{
			"get namespace 1 service 1 list success",
			fields{stg},
			args{ctx, ns1, sv1},
			nl1,
			false,
			"",
		},
		{
			"get namespace 1 service 2 list success",
			fields{stg},
			args{ctx, ns1, sv2},
			nl2,
			false,
			"",
		},
		{
			"get namespace 2 service 1 list success",
			fields{stg},
			args{ctx, ns2, sv1},
			nl3,
			false,
			"",
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "t", "t"},
			nl,
			false,
			"",
		},
		{
			"get namespace list failed empty namespace",
			fields{stg},
			args{ctx, "", sv1},
			nl,
			true,
			"namespace can not be empty",
		},
		{
			"get namespace list failed empty service",
			fields{stg},
			args{ctx, ns1, ""},
			nl,
			true,
			"service can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.ListByService() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("DeploymentStorage.ListByService() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if err != nil {
				if tt.wantErr && (err.Error() != tt.err) {
					t.Errorf("DeploymentStorage.ListByService() error = %v, want err %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.ListByService() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.ListByService() want error = %v, got none", tt.err)
				return
			}

			if !compareDeploymentMaps(got, tt.want) {
				t.Errorf("DeploymentStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test1", "")
		n2  = getDeploymentAsset(ns1, svc, "test1", "")
		n3  = getDeploymentAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Deployment, 0)
	)

	n2.Status.SetProvision()

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("DeploymentStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.SetStatus() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			if err != nil {
				t.Errorf("DeploymentStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !compareDeployments(got, tt.want) {
				t.Errorf("DeploymentStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestDeploymentStorage_SetSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test1", "")
		n2  = getDeploymentAsset(ns1, svc, "test1", "")
		n3  = getDeploymentAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Deployment, 0)
	)

	n2.Spec.State.Maintenance = true

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.SetSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("DeploymentStorage.SetSpec() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.SetSpec() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			if err != nil {
				t.Errorf("DeploymentStorage.SetSpec() got Get error = %s", err.Error())
				return
			}
			if !compareDeployments(got, tt.want) {
				t.Errorf("DeploymentStorage.SetSpec() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestDeploymentStorage_Insert(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test", "")
		n2  = getDeploymentAsset(ns1, svc, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"test successful insert",
			fields{stg},
			args{ctx, &n1},
			&n1,
			false,
			"",
		},
		{
			"test failed insert: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed insert: invalid structure",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.Insert() want error = %v, got none", tt.err)
				return
			}
		})
	}
}

func TestDeploymentStorage_Update(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test1", "")
		n2  = getDeploymentAsset(ns1, svc, "test1", "test")
		n3  = getDeploymentAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Deployment, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("DeploymentStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			if err != nil {
				t.Errorf("DeploymentStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareDeployments(got, tt.want) {
				t.Errorf("DeploymentStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestDeploymentStorage_Remove(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newDeploymentStorage()
		ctx = context.Background()
		n1  = getDeploymentAsset(ns1, svc, "test1", "")
		n2  = getDeploymentAsset(ns1, svc, "test2", "")
	)

	type fields struct {
		stg storage.Deployment
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"test successful deployment remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil deployment structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: deployment not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("DeploymentStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.Remove() want error = %v, got none", tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("DeploymentStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestDeploymentStorage_Watch(t *testing.T) {

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
	}
	defer destroy()

	var (
		stg         = newDeploymentStorage()
		ctx         = context.Background()
		n           = getDeploymentAsset("ns1", "svc", "test1", "desc")
		deploymentC = make(chan *types.Deployment)
	)

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check deployment watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.Watch(ctxT, deploymentC)
				if err != nil {
					t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
					return
				}
			}()

			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/deployments/ns1:svc:test1/meta"
			value := `{"name":"test1","description":"desc","self_link":"ns1:svc:test1","labels":null,"created":"2018-04-26T11:26:22.999932+03:00","updated":"0001-01-01T00:00:00Z","version":0,"namespace":"ns1","service":"svc","endpoint":"","status":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-deploymentC:
					t.Log("DeploymentStorage.Watch() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("DeploymentStorage.Watch() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestDeploymentStorage_WatchSpec(t *testing.T) {
	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("DeploymentStorage.WatchStatus() storage setup error = %v", err)
	}
	defer destroy()

	var (
		stg         = newDeploymentStorage()
		ctx         = context.Background()
		n           = getDeploymentAsset("ns1", "svc", "test1", "desc")
		deploymentC = make(chan *types.Deployment)
	)

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check deployment watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.WatchSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("DeploymentStorage.WatchSpec() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.WatchSpec(ctxT, deploymentC)
				if err != nil {
					t.Errorf("DeploymentStorage.WatchSpec() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/deployments/ns1:svc:test1/spec"
			value := `{"meta":{"name":"","description":"","self_link":"","labels":null,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"},"replicas":0,"state":{"destroy":false,"maintenance":false},"strategy":{"type":"","rollingOptions":{"period_update":0,"interval":0,"timeout":0,"max_unavailable":0,"max_surge":0},"resources":{},"deadline":0},"triggers":null,"selector":{},"template":{"volumes":null,"container":null,"termination":0}}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-deploymentC:
					t.Log("DeploymentStorage.WatchSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("DeploymentStorage.WatchSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestDeploymentStorage_WatchStatus(t *testing.T) {

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("DeploymentStorage.WatchStatus() storage setup error = %v", err)
	}
	defer destroy()

	var (
		stg         = newDeploymentStorage()
		ctx         = context.Background()
		n           = getDeploymentAsset("ns1", "svc", "test1", "desc")
		deploymentC = make(chan *types.Deployment)
	)
	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check deployment watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.WatchStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("DeploymentStorage.WatchSetStatus() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watcher goroutine
			go func() {
				_ = stg.WatchStatus(ctxT, deploymentC)
			}()

			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/deployments/ns1:svc:test1/status"
			value := `{"state":"","message":""}`
			//t.Log("path =", path)
			err = runEtcdPut(path, key, value)
			//			t.Logf("err = %v", err)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s, path=%s", err.Error(), path)
			}
			//
			for {
				select {
				case <-deploymentC:
					t.Log("DeploymentStorage.WatchStatus() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("DeploymentStorage.WatchSetStatus() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func Test_newDeploymentStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Deployment
	}{
		{"initialize storage",
			newDeploymentStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := newDeploymentStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDeploymentStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getDeploymentAsset(namespace, service, name, desc string) types.Deployment {

	var n = types.Deployment{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Service = service
	n.Meta.Description = desc

	n.SelfLink()
	n.Meta.Created = time.Now()

	return n
}

//compare two deployment structures
func compareDeployments(got, want *types.Deployment) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		(got.Meta.Version == want.Meta.Version) &&
		(got.Meta.Service == want.Meta.Service) &&
		(got.Meta.Endpoint == want.Meta.Endpoint) &&
		(got.Meta.Status == want.Meta.Status) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Replicas, want.Replicas) &&
		reflect.DeepEqual(got.Spec, want.Spec) {
		result = true
	}

	return result
}

func compareDeploymentMaps(got, want map[string]*types.Deployment) bool {
	for k, v := range got {
		if !compareDeployments(v, want[k]) {
			return false
		}
	}
	return true
}
