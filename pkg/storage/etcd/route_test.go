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

func TestRouteStorage_Get(t *testing.T) {

	initStorage()

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
		ns   string
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
			args{ctx, "test2", ns1},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get route info successful",
			fields{stg},
			args{ctx, "test", ns1},
			&d,
			false,
			"",
		},
		{
			"get route info failed empty namespace",
			fields{stg},
			args{ctx, "test", ""},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get route info failed empty name",
			fields{stg},
			args{ctx, "", ns1},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("RouteStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.name)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("RouteStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !compareRoutes(got, tt.want) {
				t.Errorf("RouteStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestRouteStorage_ListByNamespace(t *testing.T) {

	initStorage()

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
			"get namespace info failed empty namespace",
			fields{stg},
			args{ctx, ""},
			nl,
			true,
			"namespace can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("RouteStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("RouteStorage.ListByNamespace() want error = %v, got none", tt.err)
				return
			}

			if !compareRouteMaps(got, tt.want) {
				t.Errorf("RouteStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteStorage_ListSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test2", "")
	)

	spec1 := types.RouteSpec{
		Domain: "domain1",
	}
	spec2 := types.RouteSpec{
		Domain: "domain2",
	}
	n1.Spec = spec1
	n2.Spec = spec2

	nl0 := map[string]*types.Route{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2

	specs := map[string]*types.RouteSpec{}
	specs[stg.keyGet(&n1)] = &spec1
	specs[stg.keyGet(&n2)] = &spec2

	type fields struct {
		stg storage.Route
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.RouteSpec
		wantErr bool
		err     string
	}{
		{
			"get list spec success",
			fields{stg},
			args{ctx},
			specs,
			false,
			"",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.ListSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("RouteStorage.ListSpec() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListSpec(tt.args.ctx)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.ListSpec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("RouteStorage.ListSpec() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.ListSpec() want error = %v, got none", tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteStorage.ListSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteStorage_SetSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test1", "")
		n3  = getRouteAsset(ns1, "test2", "")
		nl  = make([]*types.Route, 0)
	)

	n2.Spec.Domain = "domain1"

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.SetSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("RouteStorage.SetSpec() storage setup error = %s", err.Error())
					return
				}
			}

			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.SetSpec() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if err != nil {
				t.Errorf("RouteStorage.SetSpec() got Get error = %s", err.Error())
				return
			}
			if !compareRoutes(got, tt.want) {
				t.Errorf("RouteStorage.SetSpec() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newRouteStorage()
		ctx = context.Background()
		n1  = getRouteAsset(ns1, "test1", "")
		n2  = getRouteAsset(ns1, "test1", "")
		n3  = getRouteAsset(ns1, "test2", "")
		nl  = make([]*types.Route, 0)
	)

	n2.Status.State = types.StateReady

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("RouteStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if err != nil {
				t.Errorf("RouteStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !compareRoutes(got, tt.want) {
				t.Errorf("RouteStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_Insert(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Insert() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Insert() want error = %v, got none", tt.err)
				return
			}
		})
	}
}

func TestRouteStorage_Update(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("RouteStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Update() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.route.Meta.Name)
			if err != nil {
				t.Errorf("RouteStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareRoutes(got, tt.want) {
				t.Errorf("RouteStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestRouteStorage_Remove(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("RouteStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RouteStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("RouteStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("RouteStorage.Remove() want error = %v, got none", tt.err)
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

	initStorage()

	var (
		stg    = newRouteStorage()
		ctx    = context.Background()
		err    error
		n      = getRouteAsset("ns1", "test1", "")
		routeC = make(chan *types.Route)
		stopC  = make(chan int)
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
			"check route watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.Watch() storage setup error = %v", err)
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
					t.Errorf("RouteStorage.Watch() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.Watch(ctx, routeC)
					if err != nil {
						t.Errorf("RouteStorage.Watch() storage setup error = %v", err)
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
				t.Errorf("RouteStorage.Watch() update error =%v", err)
				return

			case <-routeC:
				t.Log("RouteStorage.Watch() is working")
				return
			}
		})
	}
}

func TestRouteStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		stg    = newRouteStorage()
		ctx    = context.Background()
		err    error
		n      = getRouteAsset("ns1", "test1", "")
		routeC = make(chan *types.Route)
		stopC  = make(chan int)
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
			"check route watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.WatchSpec() storage setup error = %v", err)
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
					t.Errorf("RouteStorage.WatchSpec() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchSpec(ctx, routeC)
					if err != nil {
						t.Errorf("RouteStorage.WatchSpec() storage setup error = %v", err)
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
				t.Errorf("RouteStorage.WatchSpec() update error =%v", err)
				return

			case <-routeC:
				t.Log("RouteStorage.WatchSpec() is working")
				return
			}
		})
	}
}

func TestRouteStorage_WatchStatus(t *testing.T) {

	initStorage()

	var (
		stg    = newRouteStorage()
		ctx    = context.Background()
		err    error
		n      = getRouteAsset("ns1", "test1", "")
		routeC = make(chan *types.Route)
		stopC  = make(chan int)
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
			"check route watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.WatchStatus() storage setup error = %v", err)
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
					t.Errorf("RouteStorage.WatchStatus() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchStatus(ctx, routeC)
					if err != nil {
						t.Errorf("RouteStorage.WatchStatus() storage setup error = %v", err)
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
				t.Errorf("RouteStorage.WatchStatus() update error =%v", err)
				return

			case <-routeC:
				t.Log("RouteStorage.WatchStatus() is working")
				return
			}
		})
	}
}

func TestRouteStorage_WatchSpecEvents(t *testing.T) {

	initStorage()

	var (
		stg             = newRouteStorage()
		ctx             = context.Background()
		err             error
		n               = getRouteAsset("ns1", "test1", "")
		routeSpecEventC = make(chan *types.RouteSpecEvent)
		stopC           = make(chan int)
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
			"check route watch spec events",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("RouteStorage.WatchSpecEvents() storage setup error = %v", err)
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
					t.Errorf("RouteStorage.WatchSpecEvents() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchSpecEvents(ctx, routeSpecEventC)
					if err != nil {
						t.Errorf("RouteStorage.WatchSpecEvents() storage setup error = %v", err)
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
				t.Errorf("RouteStorage.WatchSpecEvents() update error =%v", err)
				return

			case retSpecEvent := <-routeSpecEventC:
				t.Log("RouteStorage.WatchSpecEvents() is working")
				t.Logf("retSpecEvent=%v\n", retSpecEvent)
				return
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

	n.Meta.Meta.Created = time.Now()

	return n
}

//compare two route structures
func compareRoutes(got, want *types.Route) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		(got.Meta.Security == want.Meta.Security) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Spec, want.Spec) {
		result = true
	}

	return result
}

func compareRouteMaps(got, want map[string]*types.Route) bool {
	for k, v := range got {
		if !compareRoutes(v, want[k]) {
			return false
		}
	}
	return true
}
