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

func TestVolumeStorage_Get(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		d   = getVolumeAsset(ns1, "test", "")
	)

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
		wantErr bool
		err     string
	}{
		{
			"get volume info failed",
			fields{stg},
			args{ctx, "test2"},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get volume info successful",
			fields{stg},
			args{ctx, "test"},
			&d,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &d); err != nil {
			t.Errorf("VolumeStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestVolumeStorage_ListByNamespace(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test1", "")
		n2  = getVolumeAsset(ns1, "test2", "")
		n3  = getVolumeAsset(ns2, "test1", "")
		nl  = make(map[string]*types.Volume, 0)
	)

	nl0 := map[string]*types.Volume{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Volume{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Volume{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Volume
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
			t.Errorf("VolumeStorage.ListByNamespace() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("VolumeStorage.ListByNamespace() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumeStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test1", "")
		n2  = getVolumeAsset(ns1, "test1", "")
		n3  = getVolumeAsset(ns1, "test2", "")
		nl  = make([]*types.Volume, 0)
	)

	n2.Status.Stage = types.StageReady

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx    context.Context
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
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
			t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.SetStatus() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestVolumeStorage_SetSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test1", "")
		n2  = getVolumeAsset(ns1, "test1", "")
		n3  = getVolumeAsset(ns1, "test2", "")
		nl  = make([]*types.Volume, 0)
	)

	n2.Status.Stage = types.StageReady

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx    context.Context
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
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
			t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.SetStatus() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestVolumeStorage_Insert(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test", "")
		n2  = getVolumeAsset(ns1, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx    context.Context
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
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
			t.Errorf("VolumeStorage.Insert() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestVolumeStorage_Update(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test1", "")
		n2  = getVolumeAsset(ns1, "test1", "test")
		n3  = getVolumeAsset(ns1, "test2", "")
		nl  = make([]*types.Volume, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx    context.Context
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
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
			t.Errorf("VolumeStorage.Update() storage setup error = %v", err)
			return
		}

		for _, n := range nl0 {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("VolumeStorage.Update() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.Update() error = %v, want %v", err, tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err != nil {
				t.Errorf("VolumeStorage.Update() error = %v, want no error", err.Error())
				return
			}
			tt.want.Meta.Updated = got.Meta.Updated
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestVolumeStorage_Remove(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		stg = newVolumeStorage()
		ctx = context.Background()
		n1  = getVolumeAsset(ns1, "test1", "")
		n2  = getVolumeAsset(ns1, "test2", "")
	)

	type fields struct {
		stg storage.Volume
	}

	type args struct {
		ctx    context.Context
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
		wantErr bool
		err     string
	}{
		{
			"test successful volume remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil volume structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: volume not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("VolumeStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.Remove() error = %v, want %v", err, tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("VolumeStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestVolumeStorage_Watch(t *testing.T) {

	initStorage()

	var (
		stg = newVolumeStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Volume
	}
	type args struct {
		ctx    context.Context
		volume chan *types.Volume
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
			args{ctx, make(chan *types.Volume)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestVolumeStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		stg = newVolumeStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Volume
	}
	type args struct {
		ctx    context.Context
		volume chan *types.Volume
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
			args{ctx, make(chan *types.Volume)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func Test_newVolumeStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Volume
	}{
		{"initialize storage",
			newVolumeStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newVolumeStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newVolumeStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getVolumeAsset(namespace, name, desc string) types.Volume {

	var n = types.Volume{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}
