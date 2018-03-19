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

package mock

import (
	"context"
	"reflect"
	"testing"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestPodStorage_Get(t *testing.T) {
	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		d   = getPodAsset(ns1, svc, dp1, "test", "")
	)

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Pod
		wantErr bool
		err     string
	}{
		{
			"get pod info failed",
			fields{stg},
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get pod info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &d); err != nil {
			t.Errorf("PodStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, dp1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestPodStorage_ListByNamespace(t *testing.T) {
	var (
		ns1 = "ns1"
		ns2 = "ns2"
		dp1 = "dp1"
		svc = "svc"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test1", "")
		n2  = getPodAsset(ns1, svc, dp1, "test2", "")
		n3  = getPodAsset(ns2, svc, dp1, "test1", "")
		nl  = make(map[string]*types.Pod, 0)
	)

	nl0 := map[string]*types.Pod{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Pod{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Pod{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Pod
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
			t.Errorf("PodStorage.List() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("PodStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByService(t *testing.T) {
	var (
		ns1 = "ns1"
		ns2 = "ns2"
		sv1 = "svc1"
		sv2 = "svc2"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, sv1, dp1, "test1", "")
		n2  = getPodAsset(ns1, sv1, dp1, "test2", "")
		n3  = getPodAsset(ns1, sv2, dp1, "test1", "")
		n4  = getPodAsset(ns2, sv1, dp1, "test1", "")
		n5  = getPodAsset(ns2, sv1, dp1, "test2", "")
		nl  = make(map[string]*types.Pod, 0)
	)

	nl0 := map[string]*types.Pod{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3
	nl0[stg.keyGet(&n4)] = &n4
	nl0[stg.keyGet(&n5)] = &n5

	nl1 := map[string]*types.Pod{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Pod{}
	nl2[stg.keyGet(&n3)] = &n3

	nl3 := map[string]*types.Pod{}
	nl3[stg.keyGet(&n4)] = &n4
	nl3[stg.keyGet(&n5)] = &n5

	type fields struct {
		stg storage.Pod
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
		want    map[string]*types.Pod
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
			t.Errorf("PodStorage.ListByService() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("PodStorage.ListByService() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.ListByService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByDeployment(t *testing.T) {
	var (
		ns1 = "ns1"
		ns2 = "ns2"
		sv1 = "svc1"
		sv2 = "svc2"
		dp1 = "dp1"
		dp2 = "dp2"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, sv1, dp1, "test1", "")
		n2  = getPodAsset(ns1, sv1, dp1, "test2", "")
		n3  = getPodAsset(ns1, sv2, dp1, "test1", "")
		n4  = getPodAsset(ns2, sv1, dp1, "test1", "")
		n5  = getPodAsset(ns2, sv1, dp1, "test2", "")
		n6  = getPodAsset(ns2, sv2, dp2, "test1", "")
		n7  = getPodAsset(ns2, sv2, dp2, "test2", "")
		nl  = make(map[string]*types.Pod, 0)
	)

	nl0 := map[string]*types.Pod{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3
	nl0[stg.keyGet(&n4)] = &n4
	nl0[stg.keyGet(&n5)] = &n5
	nl0[stg.keyGet(&n6)] = &n6
	nl0[stg.keyGet(&n7)] = &n7

	nl1 := map[string]*types.Pod{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Pod{}
	nl2[stg.keyGet(&n3)] = &n3

	nl3 := map[string]*types.Pod{}
	nl3[stg.keyGet(&n4)] = &n4
	nl3[stg.keyGet(&n5)] = &n5

	nl4 := map[string]*types.Pod{}
	nl4[stg.keyGet(&n6)] = &n6
	nl4[stg.keyGet(&n7)] = &n7

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		ns  string
		svc string
		dp  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Pod
		wantErr bool
	}{
		{
			"get namespace 1 service 1 deployment 1 list success",
			fields{stg},
			args{ctx, ns1, sv1, dp1},
			nl1,
			false,
		},
		{
			"get namespace 1 service 2 deployment 1 list success",
			fields{stg},
			args{ctx, ns1, sv2, dp1},
			nl2,
			false,
		},
		{
			"get namespace 2 service 1 deployment 1 list success",
			fields{stg},
			args{ctx, ns2, sv1, dp1},
			nl3,
			false,
		},
		{
			"get namespace 2 service 2 deployment 2 list success",
			fields{stg},
			args{ctx, ns2, sv2, dp2},
			nl4,
			false,
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "", "", ""},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.ListByDeployment() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("PodStorage.ListByDeployment() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByDeployment(tt.args.ctx, tt.args.ns, tt.args.svc, tt.args.dp)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.ListByDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.ListByDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_SetState(t *testing.T) {
	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test1", "")
		n2  = getPodAsset(ns1, svc, dp1, "test1", "")
		n3  = getPodAsset(ns1, svc, dp1, "test2", "")
		nl  = make([]*types.Pod, 0)
	)

	n2.State.Provision = true
	n2.State.Destroy = true

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		pod *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Pod
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
			t.Errorf("PodStorage.SetState() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("PodStorage.SetState() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetState(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.SetState() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.SetState() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.SetState() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, svc, dp1, tt.args.pod.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.SetState() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_Insert(t *testing.T) {
	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test", "")
		n2  = getPodAsset(ns1, svc, dp1, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		pod *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Pod
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
			t.Errorf("PodStorage.SetState() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestPodStorage_Update(t *testing.T) {
	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test1", "")
		n2  = getPodAsset(ns1, svc, dp1, "test1", "test")
		n3  = getPodAsset(ns1, svc, dp1, "test2", "")
		nl  = make([]*types.Pod, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		pod *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Pod
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
			t.Errorf("PodStorage.Update() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("PodStorage.Update() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, svc, dp1, tt.args.pod.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_Remove(t *testing.T) {
	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test1", "")
		n2  = getPodAsset(ns1, svc, dp1, "test2", "")
	)

	type fields struct {
		stg storage.Pod
	}

	type args struct {
		ctx context.Context
		pod *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Pod
		wantErr bool
		err     string
	}{
		{
			"test successful pod remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil pod structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: pod not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("PodStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Remove() error = %v, want %v", err, tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, svc, dp1, tt.args.pod.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("PodStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_Watch(t *testing.T) {
	var (
		stg = newPodStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Pod
	}
	type args struct {
		ctx context.Context
		pod chan *types.Pod
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
			args{ctx, make(chan *types.Pod)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPodStorage_WatchSpec(t *testing.T) {
	var (
		stg = newPodStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Pod
	}
	type args struct {
		ctx context.Context
		pod chan *types.Pod
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
			args{ctx, make(chan *types.Pod)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.WatchSpec(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newPodStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Pod
	}{
		{"initialize storage",
			newPodStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPodStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPodStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getPodAsset(namespace, service, deployment, name, desc string) types.Pod {
	p := types.Pod{}

	p.Meta.Name = name
	p.Meta.Description = desc
	p.Meta.Namespace = namespace
	p.Meta.Service = service
	p.Meta.Deployment = deployment
	p.SelfLink()

	return p
}
