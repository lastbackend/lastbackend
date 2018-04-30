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
	"context"
	"reflect"
	"testing"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestIngressStorage_List(t *testing.T) {
	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test1", "", true)
		n2  = getIngressAsset("test2", "", false)
		nl  = make(map[string]*types.Ingress, 0)
	)

	nl[n1.Meta.Name] = &n1
	nl[n2.Meta.Name] = &n2

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Ingress
		wantErr bool
	}{
		{
			"get ingress list success",
			fields{stg},
			args{ctx},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.List() storage setup error = %v", err)
			return
		}

		for _, n := range nl {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("IngressStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.stg.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngressStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressStorage_Get(t *testing.T) {

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n   = getIngressAsset("test", "", true)
	)

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Ingress
		wantErr bool
		err     string
	}{
		{
			"get Ingress info failed",
			fields{stg},
			args{ctx, "test2"},
			&n,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get Ingress info successful",
			fields{stg},
			args{ctx, "test"},
			&n,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n); err != nil {
			t.Errorf("IngressStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_GetSpec(t *testing.T) {

	var (
		stg  = newIngressStorage()
		ctx  = context.Background()
		n    = getIngressAsset("test", "", true)
		n1   = getIngressAsset("", "", true)
		n2   = getIngressAsset("test2", "", true)
		spec = types.IngressSpec{}
		rs   = types.RouteSpec{
			Domain: "domain",
		}
	)
	spec.Routes = make(map[string]types.RouteSpec, 1)
	spec.Routes[n.Meta.Name] = rs

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx     context.Context
		ingress *types.Ingress
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.IngressSpec
		wantErr bool
		err     string
	}{
		{
			"get Ingress spec failed args invalid",
			fields{stg},
			args{ctx, &n1},
			&spec,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"get Ingress spec failed not found",
			fields{stg},
			args{ctx, &n2},
			&spec,
			true,
			store.ErrEntityNotFound,
		},

		{
			"get Ingress spec successful",
			fields{stg},
			args{ctx, &n},
			&spec,
			false,
			"",
		},
	}

	for _, tt := range tests {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.GetSpec() storage setup (clear) error = %v", err)
			return
		}
		if err := stg.Insert(ctx, &n); err != nil {
			t.Errorf("IngressStorage.GetSpec() storage setup (insert) error = %v", err)
			return
		}
		//add spec info to ingressAsset
		if err := insertRouteSpec(&n, spec); err != nil {
			t.Errorf("IngressStorage.GetSpec() storage setup (sub insert) error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.ingress)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.GetSpec() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.GetSpec() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.GetSpec() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_Insert(t *testing.T) {
	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test", "", true)
		n2  = getIngressAsset("", "", true)
	)

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx     context.Context
		ingress *types.Ingress
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Ingress
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
			t.Errorf("IngressStorage.Insert() storage setup (clear) error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IngressStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}

			//check inserted item
			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.Insert() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestIngressStorage_Update(t *testing.T) {
	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test1", "", true)
		n2  = getIngressAsset("test1", "desc", true)
		n3  = getIngressAsset("test3", "", true)

		nl = make([]*types.Ingress, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx     context.Context
		ingress *types.Ingress
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Ingress
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
			t.Errorf("IngressStorage.Update() storage setup (clear) error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("IngressStorage.Update() storage setup (insert) error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IngressStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			//check updated item
			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestIngressStorage_SetStatus(t *testing.T) {
	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test1", "", true)
		n2  = getIngressAsset("test1", "", true)
		n3  = getIngressAsset("test3", "", true)

		nl = make([]*types.Ingress, 0)
	)

	n2.Status.Ready = false

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx     context.Context
		ingress *types.Ingress
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Ingress
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
			t.Errorf("IngressStorage.SetStatus() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("IngressStorage.SetStatus() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IngressStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.SetStatus()got = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestIngressStorage_Remove(t *testing.T) {

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test1", "", true)
		n2  = getIngressAsset("test2", "", true)
	)

	type fields struct {
		stg storage.Ingress
	}

	type args struct {
		ctx     context.Context
		ingress *types.Ingress
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Ingress
		wantErr bool
		err     string
	}{
		{
			"test successful ingress remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed remove: nil ingress structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed remove: ingress not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Remove() storage setup (clear) error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("IngressStorage.Remove() storage setup (insert) error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IngressStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("IngressStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestIngressStorage_Watch(t *testing.T) {

	var (
		stg = newIngressStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Ingress
	}
	type args struct {
		ctx     context.Context
		ingress chan *types.Ingress
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
			args{ctx, make(chan *types.Ingress)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.ingress); (err != nil) != tt.wantErr {
				t.Errorf("IngressStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//add test routespec into ingress asset
func insertRouteSpec(ingress *types.Ingress, spec types.IngressSpec) error {
	ingress.Spec = spec
	return nil
}

func getIngressAsset(name, desc string, ready bool) types.Ingress {
	var n = types.Ingress{
		Meta: types.IngressMeta{
			Meta: types.Meta{
				Name:        name,
				Description: desc,
			},
		},
		Status: types.IngressStatus{
			Ready: ready,
		},
		Spec: types.IngressSpec{
			Routes: make(map[string]types.RouteSpec),
		},
	}

	/*


		var rs = types.RouteSpec{
			Domain: "domain1",
		}

			n.Meta.Name = name
			n.Meta.Description = desc
			n.Status.Ready = ready
			n.Spec = make(map[string]*types.RouteSpec)
	*/
	//n.Spec.Routes[] = rs

	return n
}
