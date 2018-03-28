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
	"testing"

	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"reflect"
)

func TestSecretStorage_Get(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newSecretStorage()
		ctx = context.Background()
		d   = getSecretAsset(ns1, "test")
	)

	type fields struct {
		stg storage.Secret
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Secret
		wantErr bool
		err     string
	}{
		{
			"get secret info failed",
			fields{stg},
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get secret info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SecretStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("SecretStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SecretStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("SecretStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSecretStorage_ListByNamespace(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		stg = newSecretStorage()
		ctx = context.Background()
		n1  = getSecretAsset(ns1, "test1")
		n2  = getSecretAsset(ns1, "test2")
		n3  = getSecretAsset(ns2, "test1")
		nl  = make(map[string]*types.Secret, 0)
	)

	nl0 := map[string]*types.Secret{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Secret{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Secret{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Secret
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Secret
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SecretStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("SecretStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretStorage_Insert(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newSecretStorage()
		ctx = context.Background()
		n1  = getSecretAsset(ns1, "test")
		n2  = getSecretAsset(ns1, "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Secret
	}

	type args struct {
		ctx    context.Context
		secret *types.Secret
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Secret
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
			t.Errorf("SecretStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.secret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SecretStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SecretStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("SecretStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestSecretStorage_Update(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newSecretStorage()
		ctx = context.Background()
		n1  = getSecretAsset(ns1, "test1")
		n2  = getSecretAsset(ns1, "test1")
		n3  = getSecretAsset(ns1, "test2")
		nl  = make([]*types.Secret, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Secret
	}

	type args struct {
		ctx    context.Context
		secret *types.Secret
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Secret
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
			t.Errorf("SecretStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("SecretStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.secret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SecretStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SecretStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("SecretStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.secret.Meta.Name)
			tt.want.Meta.Updated = got.Meta.Updated
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestSecretStorage_Remove(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newSecretStorage()
		ctx = context.Background()
		n1  = getSecretAsset(ns1, "test1")
		n2  = getSecretAsset(ns1, "test2")
	)

	type fields struct {
		stg storage.Secret
	}

	type args struct {
		ctx    context.Context
		secret *types.Secret
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Secret
		wantErr bool
		err     string
	}{
		{
			"test successful secret remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil secret structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: secret not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SecretStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("SecretStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.secret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SecretStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SecretStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("SecretStorage.Remove() error = %v, want %v", err, tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.secret.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("SecretStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func Test_newSecretStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Secret
	}{
		{"initialize storage",
			newSecretStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSecretStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSecretStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getSecretAsset(namespace, name string) types.Secret {

	var n = types.Secret{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace

	return n
}
