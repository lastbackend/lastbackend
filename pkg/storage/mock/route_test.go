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

package mock

import (
	"testing"

	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"reflect"
)

func TestRouteStorage_Get(t *testing.T) {
	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		d   = getRouteAsset(ns1, "test", "")
	)

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Route
		wantErr bool
		err     string
	}{
		{
			"get route info failed",
			fields{stg},
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get route info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &d); err != nil {
			t.Errorf("RouteStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestRouteStorage_ListByNamespace(t *testing.T) {
	var (
		ns1 = "ns1"
		ns2 = "ns2"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test2", "")
		n3  = getRouteAsset(ns2, "test1", "")
		nl  = make(map[string]*types.Route, 0)
	)

	nl0 := map[string]*types.Route{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Route{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Route{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Route
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
			t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteStorage_SetStatus(t *testing.T) {
	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test1", "")
		n3  = getRouteAsset(ns1, "test2", "")
		nl  = make([]*types.Route, 0)
	)

	n2.Status.Stage = types.StateReady

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx   context.Context
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Route
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
			t.Errorf("RouteStorage.SetStatus() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("RouteStorage.SetStatus() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_Insert(t *testing.T) {
	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test", "")
		n2  = getRouteAsset(ns1, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx   context.Context
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Route
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
			t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestRouteStorage_Update(t *testing.T) {
	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test1", "test")
		n3  = getRouteAsset(ns1, "test2", "")
		nl  = make([]*types.Route, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx   context.Context
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Route
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
			t.Errorf("RouteStorage.Update() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("RouteStorage.Update() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_Remove(t *testing.T) {
	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test2", "")
	)

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx   context.Context
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Route
		wantErr bool
		err     string
	}{
		{
			"test successful route remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil route structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: route not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("RouteStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Remove() error = %v, want %v", err, tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("RouteStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_Watch(t *testing.T) {
	var (
		stg = newRouteStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Route
	}
	type args struct {
		ctx   context.Context
		route chan *types.Route
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
			args{ctx, make(chan *types.Route)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.route); (err != nil) != tt.wantErr {
				t.Errorf("RouteStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouteStorage_WatchSpec(t *testing.T) {
	var (
		stg = newRouteStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Route
	}
	type args struct {
		ctx   context.Context
		route chan *types.Route
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
			args{ctx, make(chan *types.Route)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.WatchSpec(tt.args.ctx, tt.args.route); (err != nil) != tt.wantErr {
				t.Errorf("RouteStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newRouteStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Route
	}{
		{"initialize storage",
			newRouteStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRouteStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRouteStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getRouteAsset(namespace, name, desc string) types.Route {

	var n = types.Route{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}
