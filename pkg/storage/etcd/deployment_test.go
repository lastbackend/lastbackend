//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"testing"

	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"reflect"
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
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get deployment info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &d); err != nil {
			t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
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
	}{
		{
			"get namespace list 1 success",
			fields{stg},
			args{ctx, ns1},
			nl1,
			false,
		},
		{
			"get namespace list 2 success",
			fields{stg},
			args{ctx, ns2},
			nl2,
			false,
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "empty"},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Get() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("DeploymentStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentStorage.List() = %v, want %v", got, tt.want)
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
	}{
		{
			"get namespace 1 service 1 list success",
			fields{stg},
			args{ctx, ns1, sv1},
			nl1,
			false,
		},
		{
			"get namespace 1 service 2 list success",
			fields{stg},
			args{ctx, ns1, sv2},
			nl2,
			false,
		},
		{
			"get namespace 2 service 1 list success",
			fields{stg},
			args{ctx, ns2, sv1},
			nl3,
			false,
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "t", "t"},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.ListByService() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("DeploymentStorage.ListByService() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.ListByService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentStorage_SetState(t *testing.T) {

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

	n2.State.Provision = true
	n2.State.Ready = true

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

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.SetState() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("DeploymentStorage.SetState() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetState(tt.args.ctx, tt.args.deployment)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeploymentStorage.SetState() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("DeploymentStorage.SetState() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("DeploymentStorage.SetState() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			tt.want.Meta.Updated = got.Meta.Updated
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentStorage.SetState() = %v, want %v", got, tt.want)
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

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Insert() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("DeploymentStorage.Insert() error = %v, want %v", err, tt.err)
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

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Update() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("DeploymentStorage.Update() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("DeploymentStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.deployment.Meta.Name)
			tt.want.Meta.Updated = got.Meta.Updated
			if !reflect.DeepEqual(got, tt.want) {
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

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("DeploymentStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("DeploymentStorage.Remove() error = %v, want %v", err, tt.err)
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

	initStorage()

	var (
		stg = newDeploymentStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Deployment
	}
	type args struct {
		ctx        context.Context
		deployment chan *types.Deployment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check watch",
			fields{stg},
			args{ctx, make(chan *types.Deployment)},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
			return
		}

		go t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeploymentStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		stg = newDeploymentStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Deployment
	}
	type args struct {
		ctx        context.Context
		deployment chan *types.Deployment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check watch",
			fields{stg},
			args{ctx, make(chan *types.Deployment)},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("DeploymentStorage.WatchSpec() storage setup error = %v", err)
			return
		}

		go t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.WatchSpec(tt.args.ctx, tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.WatchSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
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

	return n
}
