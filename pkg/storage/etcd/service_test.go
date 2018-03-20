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

func TestServiceStorage_Get(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		d   = getServiceAsset(ns1, "test", "")
	)

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
		wantErr bool
		err     string
	}{
		{
			"get service info failed",
			fields{stg},
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get service info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &d); err != nil {
			t.Errorf("ServiceStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestServiceStorage_ListByNamespace(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test1", "")
		n2  = getServiceAsset(ns1, "test2", "")
		n3  = getServiceAsset(ns2, "test1", "")
		nl  = make(map[string]*types.Service, 0)
	)

	nl0 := map[string]*types.Service{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Service{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Service{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Service
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
			t.Errorf("ServiceStorage.ListByNamespace() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("ServiceStorage.ListByNamespace() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_SetState(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test1", "")
		n2  = getServiceAsset(ns1, "test1", "")
		n3  = getServiceAsset(ns1, "test2", "")
		nl  = make([]*types.Service, 0)
	)

	n2.State.Provision = true
	n2.State.Destroy = true

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
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
			t.Errorf("ServiceStorage.SetState() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("ServiceStorage.SetState() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetState(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.SetState() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.SetState() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.SetState() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.SetState() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestServiceStorage_SetSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test1", "")
		n2  = getServiceAsset(ns1, "test1", "")
		n3  = getServiceAsset(ns1, "test2", "")
		nl  = make([]*types.Service, 0)
	)

	n2.Spec.Template.Termination = 1

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
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
			t.Errorf("ServiceStorage.SetSpec() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("ServiceStorage.SetSpec() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.SetSpec() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestServiceStorage_Insert(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test", "")
		n2  = getServiceAsset(ns1, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
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
			t.Errorf("ServiceStorage.Insert() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestServiceStorage_Update(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test1", "")
		n2  = getServiceAsset(ns1, "test1", "test")
		n3  = getServiceAsset(ns1, "test2", "")
		nl  = make([]*types.Service, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
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
			t.Errorf("ServiceStorage.Update() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("ServiceStorage.Update() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			tt.want.Meta.Updated = got.Meta.Updated
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestServiceStorage_Remove(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newServiceStorage()
		ctx = context.Background()
		n1  = getServiceAsset(ns1, "test1", "")
		n2  = getServiceAsset(ns1, "test2", "")
	)

	type fields struct {
		stg storage.Service
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
		wantErr bool
		err     string
	}{
		{
			"test successful service remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil service structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: service not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("ServiceStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Remove() error = %v, want %v", err, tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("ServiceStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestServiceStorage_Watch(t *testing.T) {

	initStorage()

	var (
		stg = newServiceStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Service
	}
	type args struct {
		ctx     context.Context
		service chan *types.Service
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
			args{ctx, make(chan *types.Service)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func () {
				if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
					t.Errorf("ServiceStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
}

func TestServiceStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		stg = newServiceStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Service
	}
	type args struct {
		ctx     context.Context
		service chan *types.Service
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
			args{ctx, make(chan *types.Service)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func () {
				if err := tt.fields.stg.WatchSpec(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
					t.Errorf("ServiceStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
}

func Test_newServiceStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Service
	}{
		{"initialize storage",
			newServiceStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newServiceStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newServiceStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getServiceAsset(namespace, name, desc string) types.Service {

	var n = types.Service{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}
