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

func TestNamespaceStorage_Get(t *testing.T) {
	var (
		stg = newNamespaceStorage()
		ctx = context.Background()
		n   = getNamespaceAsset("test", "")
	)

	type fields struct {
		stg storage.Namespace
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Namespace
		wantErr bool
		err     string
	}{
		{
			"get namespace info failed",
			fields{stg},
			args{ctx, "test2"},
			&n,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get namespace info successful",
			fields{stg},
			args{ctx, "test"},
			&n,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NamespaceStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n); err != nil {
			t.Errorf("NamespaceStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NamespaceStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NamespaceStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNamespaceStorage_List(t *testing.T) {
	var (
		stg = newNamespaceStorage()
		ctx = context.Background()
		n1  = getNamespaceAsset("test1", "")
		n2  = getNamespaceAsset("test2", "")
		nl  = make(map[string]*types.Namespace, 0)
	)

	nl[n1.Meta.Name] = &n1
	nl[n2.Meta.Name] = &n2

	type fields struct {
		stg storage.Namespace
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Namespace
		wantErr bool
	}{
		{
			"get namespace list success",
			fields{stg},
			args{ctx},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NamespaceStorage.List() storage setup error = %v", err)
			return
		}

		for _, n := range nl {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("NamespaceStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceStorage_Insert(t *testing.T) {
	var (
		stg = newNamespaceStorage()
		ctx = context.Background()
		n1  = getNamespaceAsset("test", "")
		n2  = getNamespaceAsset("", "")
	)

	type fields struct {
		stg storage.Namespace
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Namespace
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
			"test failed insert",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed insert",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NamespaceStorage.Insert() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.namespace)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NamespaceStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NamespaceStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NamespaceStorage.Insert() error = %v, want %v", err.Error(), tt.err)
				return
			}
		})
	}
}

func TestNamespaceStorage_Update(t *testing.T) {

	var (
		stg = newNamespaceStorage()
		ctx = context.Background()
		n1  = getNamespaceAsset("test1", "")
		n2  = getNamespaceAsset("test1", "desc")
		n3  = getNamespaceAsset("test2", "")
	)

	type fields struct {
		stg storage.Namespace
	}

	type args struct {
		ctx     context.Context
		naspace *types.Namespace
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Namespace
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
			t.Errorf("NamespaceStorage.Update() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NamespaceStorage.Update() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.naspace)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NamespaceStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NamespaceStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NamespaceStorage.Update() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.naspace.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNamespaceStorage_Remove(t *testing.T) {

	var (
		stg = newNamespaceStorage()
		ctx = context.Background()
		n1  = getNamespaceAsset("test1", "")
		n2  = getNamespaceAsset("test2", "")
	)

	type fields struct {
		stg storage.Namespace
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Namespace
		wantErr bool
		err     string
	}{
		{
			"test successful namespace remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil namespace structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: namespace not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NamespaceStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NamespaceStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.namespace)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NamespaceStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NamespaceStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NamespaceStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.namespace.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("NamespaceStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func Test_newNamespaceStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Namespace
	}{
		{"initialize storage",
			newNamespaceStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newNamespaceStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNamespaceStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getNamespaceAsset(name, desc string) types.Namespace {
	var n = types.Namespace{}

	n.Meta.Name = name
	n.Meta.Description = desc

	return n
}
