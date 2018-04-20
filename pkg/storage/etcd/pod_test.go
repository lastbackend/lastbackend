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

func TestPodStorage_Get(t *testing.T) {

	initStorage()

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
		ns   string
		svc  string
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
			args{ctx, "test2", ns1, svc},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get pod info successful",
			fields{stg},
			args{ctx, "test", ns1, svc},
			&d,
			false,
			"",
		},
		{
			"get pod info failed empty namespace",
			fields{stg},
			args{ctx, "test", "", svc},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get pod info failed empty service",
			fields{stg},
			args{ctx, "test", ns1, ""},
			&d,
			true,
			"service can not be empty",
		},
		{
			"get pod info failed empty name",
			fields{stg},
			args{ctx, "", ns1, svc},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("PodStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.svc, dp1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("PodStorage.Get() = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !comparePods(got, tt.want) {
				t.Errorf("PodStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestPodStorage_ListByNamespace(t *testing.T) {

	initStorage()

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
		err     string
	}{
		{
			"get pod list 1 success",
			fields{stg},
			args{ctx, ns1},
			nl1,
			false,
			"",
		},
		{
			"get pod list 2 success",
			fields{stg},
			args{ctx, ns2},
			nl2,
			false,
			"",
		},
		{
			"get pod empty list success",
			fields{stg},
			args{ctx, "empty"},
			nl,
			false,
			"",
		},
		{
			"get pod info failed empty namespace",
			fields{stg},
			args{ctx, ""},
			nl,
			true,
			"namespace can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.ListByNamespace() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("PodStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("PodStorage.ListByNamespace() want error = %v, got no error", tt.err)
				return
			}

			if !comparePodMaps(got, tt.want) {
				t.Errorf("PodStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByService(t *testing.T) {

	initStorage()

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
		err     string
	}{
		{
			"get pod 1 service 1 list success",
			fields{stg},
			args{ctx, ns1, sv1},
			nl1,
			false,
			"",
		},
		{
			"get pod 1 service 2 list success",
			fields{stg},
			args{ctx, ns1, sv2},
			nl2,
			false,
			"",
		},
		{
			"get pod 2 service 1 list success",
			fields{stg},
			args{ctx, ns2, sv1},
			nl3,
			false,
			"",
		},
		{
			"get pod empty list success",
			fields{stg},
			args{ctx, "t", "t"},
			nl,
			false,
			"",
		},
		{
			"get pod info failed empty namespace",
			fields{stg},
			args{ctx, "", sv1},
			nl,
			true,
			"namespace can not be empty",
		},
		{
			"get pod info failed empty service",
			fields{stg},
			args{ctx, ns1, ""},
			nl,
			true,
			"service can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.ListByService() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.ListByService() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.ListByService() \n%v\nwant\n%v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("PodStorage.ListByService() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("PodStorage.ListByService() want error = %v, got no error", tt.err)
				return
			}

			if !comparePodMaps(got, tt.want) {
				t.Errorf("PodStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByDeployment(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		sv1 = "svc1"
		sv2 = "svc2"
		dp1 = "dp1"
		dp2 = "dp2"
		dp3 = "dp3"
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
		err     string
	}{
		{
			"get pod 1 service 1 deployment 1 list success",
			fields{stg},
			args{ctx, ns1, sv1, dp1},
			nl1,
			false,
			"",
		},
		{
			"get pod 1 service 2 deployment 1 list success",
			fields{stg},
			args{ctx, ns1, sv2, dp1},
			nl2,
			false,
			"",
		},
		{
			"get pod 2 service 1 deployment 1 list success",
			fields{stg},
			args{ctx, ns2, sv1, dp1},
			nl3,
			false,
			"",
		},
		{
			"get pod 2 service 2 deployment 2 list success",
			fields{stg},
			args{ctx, ns2, sv2, dp2},
			nl4,
			false,
			"",
		},
		{
			"get pod 2 service 2 deployment 3 empty list success",
			fields{stg},
			args{ctx, ns2, sv2, dp3},
			nl,
			false,
			"",
		},
		{
			"get pod info failed empty namespace",
			fields{stg},
			args{ctx, "", sv1, dp1},
			nl,
			true,
			"namespace can not be empty",
		},
		{
			"get pod info failed empty service",
			fields{stg},
			args{ctx, ns1, "", dp1},
			nl,
			true,
			"service can not be empty",
		},
		{
			"get pod info failed empty name",
			fields{stg},
			args{ctx, ns1, sv1, ""},
			nl,
			true,
			"service can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.ListByDeployment() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.ListByDeployment() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByDeployment(tt.args.ctx, tt.args.ns, tt.args.svc, tt.args.dp)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.ListByDeployment() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("PodStorage.ListByDeployment() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.ListByDeployment() want error = %v, got no error", tt.err)
				return
			}

			if !comparePodMaps(got, tt.want) {
				t.Errorf("PodStorage.ListByDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_SetStatus(t *testing.T) {

	initStorage()

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

	n2.Status.Stage = types.StateReady

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.SetStatus() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err != nil {
				t.Errorf("PodStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !comparePods(got, tt.want) {
				t.Errorf("PodStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_SetSpec(t *testing.T) {

	initStorage()

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

	n2.Spec.State.Destroy = true

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
			"test successful pod set spec",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed pod set spec: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed pod set spec: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.SetSpec() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.SetSpec() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err != nil {
				t.Errorf("PodStorage.SetSpec() got Get error = %s", err.Error())
				return
			}
			if !comparePods(got, tt.want) {
				t.Errorf("PodStorage.SetSpec() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_SetMeta(t *testing.T) {

	initStorage()

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

	n2.Meta.Endpoint = "127.0.0.1:12345"

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
			"test successful pod set meta",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed pod set meta: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed pod set meta: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.SetMeta() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.SetMeta() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetMeta(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.SetMeta() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.SetMeta() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.SetMeta() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err != nil {
				t.Errorf("PodStorage.SetMeta() got Get error = %s", err.Error())
				return
			}
			if !comparePods(got, tt.want) {
				t.Errorf("PodStorage.SetMeta() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_Insert(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Insert() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Insert() error = %v, want none", tt.err)
				return
			}
		})
	}
}

func TestPodStorage_Update(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("PodStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Update() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Update() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err != nil {
				t.Errorf("PodStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !comparePods(got, tt.want) {
				t.Errorf("PodStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestPodStorage_Remove(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("PodStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Remove() error = %v, want none", tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("PodStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

//TODO attention! test isn't good because of method Destroy
func TestPodStorage_Destroy(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		stg = newPodStorage()
		ctx = context.Background()
		n1  = getPodAsset(ns1, svc, dp1, "test1", "")
		n2  = getPodAsset(ns1, svc, dp1, "test2", "")
	)

	n1.Spec.State.Destroy = true

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
			"test successful pod destroy",
			fields{stg},
			args{ctx, &n1},
			&n1,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed destroy: nil pod structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed destroy: pod not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Destroy() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("PodStorage.Destroy() storage setup error = %v", err)
				return
			}

			t.Log("before destroy=", tt.args.pod)
			err := stg.Destroy(tt.args.ctx, tt.args.pod)
			t.Log("after destroy=", tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PodStorage.Destroy() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("PodStorage.Destroy() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("PodStorage.Destroy() error = %v, want none", tt.err)
				return
			}
			/*TODO fix it when will be correct Destroy
			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.pod.Meta.Namespace, tt.args.pod.Meta.Service, tt.args.pod.Meta.Deployment, tt.args.pod.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("PodStorage.Destroy() = %v, want %v", err, tt.err)
				return
			}
			*/
		})
	}
}

func TestPodStorage_Watch(t *testing.T) {

	initStorage()

	var (
		stg   = newPodStorage()
		ctx   = context.Background()
		n     = getPodAsset("ns1", "svc", "dp1", "test", "")
		err   error
		podC  = make(chan *types.Pod)
		stopC = make(chan int)
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
			"check pod watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err = stg.Insert(ctx, &n)
			startW := make(chan int, 1)
			if err != nil {
				startW <- 2
			} else {
				//start watch after successfull insert
				startW <- 1
			}
			select {
			case res := <-startW:
				if res != 1 {
					t.Errorf("PodStorage.Watch() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.Watch(ctx, podC)
					if err != nil {
						t.Errorf("PodStorage.Watch() storage setup error = %v", err)
						return
					}
				}()
			}
			//run go function to cause watch event
			go func() {
				time.Sleep(1 * time.Second)
				err = stg.Update(ctx, &n)
				time.Sleep(1 * time.Second)
				stopC <- 1
				return
			}()

			//wait for result
			select {
			case <-stopC:
				t.Errorf("PodStorage.Watch() update error =%v", err)
				return

			case <-podC:
				t.Log("PodStorage.Watch() is working")
				return
			}
		})
	}
}

func TestPodStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		stg   = newPodStorage()
		ctx   = context.Background()
		n     = getPodAsset("ns1", "svc", "dp1", "test", "")
		err   error
		podC  = make(chan *types.Pod)
		stopC = make(chan int)
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
			"check pod watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.WatchSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err = stg.Insert(ctx, &n)
			startW := make(chan int, 1)
			if err != nil {
				startW <- 2
			} else {
				//start watch after successfull insert
				startW <- 1
			}
			select {
			case res := <-startW:
				if res != 1 {
					t.Errorf("PodStorage.WatchSpec() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchSpec(ctx, podC)
					if err != nil {
						t.Errorf("PodStorage.WatchSpec() storage setup error = %v", err)
						return
					}
				}()
			}
			//run go function to cause watch event
			go func() {
				time.Sleep(1 * time.Second)
				err = stg.SetSpec(ctx, &n)
				time.Sleep(1 * time.Second)
				stopC <- 1
				return
			}()

			//wait for result
			select {
			case <-stopC:
				t.Errorf("PodStorage.WatchSpec() update error =%v", err)
				return

			case <-podC:
				t.Log("PodStorage.WatchSpec() is working")
				return
			}
		})
	}
}

func TestPodStorage_WatchStatus(t *testing.T) {

	initStorage()

	var (
		stg   = newPodStorage()
		ctx   = context.Background()
		n     = getPodAsset("ns1", "svc", "dp1", "test", "")
		err   error
		podC  = make(chan *types.Pod)
		stopC = make(chan int)
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
			"check pod watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("PodStorage.WatchStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err = stg.Insert(ctx, &n)
			startW := make(chan int, 1)
			if err != nil {
				startW <- 2
			} else {
				//start watch after successfull insert
				startW <- 1
			}
			select {
			case res := <-startW:
				if res != 1 {
					t.Errorf("PodStorage.WatchStatus() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchStatus(ctx, podC)
					if err != nil {
						t.Errorf("PodStorage.WatchStatus() storage setup error = %v", err)
						return
					}
				}()
			}
			//run go function to cause watch event
			go func() {
				time.Sleep(1 * time.Second)
				err = stg.SetStatus(ctx, &n)
				time.Sleep(1 * time.Second)
				stopC <- 1
				return
			}()

			//wait for result
			select {
			case <-stopC:
				t.Errorf("PodStorage.WatchStatus() update error =%v", err)
				return

			case <-podC:
				t.Log("PodStorage.WatchStatus() is working")
				return
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

	p.Meta.Meta.Created = time.Now()

	return p
}

//compare two Pod structures
func comparePods(got, want *types.Pod) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		(got.Meta.Deployment == want.Meta.Deployment) &&
		(got.Meta.Service == want.Meta.Service) &&
		(got.Meta.Node == want.Meta.Node) &&
		(got.Meta.Status == want.Meta.Status) &&
		(got.Meta.Endpoint == want.Meta.Endpoint) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Spec, want.Spec) {
		result = true
	}

	return result
}

func comparePodMaps(got, want map[string]*types.Pod) bool {
	for k, v := range got {
		if !comparePods(v, want[k]) {
			return false
		}
	}
	return true
}
