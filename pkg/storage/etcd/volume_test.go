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
		ns   string
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
			args{ctx, "test2", ns1},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get volume info successful",
			fields{stg},
			args{ctx, "test", ns1},
			&d,
			false,
			"",
		},
		{
			"get volume info failed empty namespace",
			fields{stg},
			args{ctx, "test", ""},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get volume info failed empty name",
			fields{stg},
			args{ctx, "", ns1},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("VolumeStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("VolumeStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.Get() want error= %v, got none", tt.err)
				return
			}

			if !compareVolumes(got, tt.want) {
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
			t.Errorf("VolumeStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("VolumeStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("VolumeStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.ListByNamespace() want error = %v, got none", tt.err)
				return
			}

			if !compareVolumeMaps(got, tt.want) {
				t.Errorf("VolumeStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeStorage_ListByService(t *testing.T) {

	initStorage()

	var (
		ns1    = "ns1"
		ns2    = "ns2"
		ns1sv1 = "ns1:svc1"
		ns1sv2 = "ns1:svc2"
		ns2sv1 = "ns2:svc1"
		sv1    = "svc1"
		sv2    = "svc2"
		stg    = newVolumeStorage()
		ctx    = context.Background()
		n1     = getVolumeAsset(ns1sv1, "test1", "")
		n2     = getVolumeAsset(ns1sv2, "test2", "")
		n3     = getVolumeAsset(ns2sv1, "test1", "")
		nl     = make(map[string]*types.Volume, 0)
	)

	nl0 := map[string]*types.Volume{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Volume{}
	nl1[stg.keyGet(&n1)] = &n1
	//nl1[stg.keyGet(&n3)] = &n3

	nl2 := map[string]*types.Volume{}
	nl2[stg.keyGet(&n2)] = &n2

	nl3 := map[string]*types.Volume{}
	nl3[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Volume
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
		want    map[string]*types.Volume
		wantErr bool
		err     string
	}{
		{
			"get namespace 1 service 1 list success",
			fields{stg},
			args{ctx, ns1, sv1},
			nl1,
			false,
			"",
		},
		{
			"get namespace 1 service 2 list success",
			fields{stg},
			args{ctx, ns1, sv2},
			nl2,
			false,
			"",
		},
		{
			"get namespace 2 service 1 list success",
			fields{stg},
			args{ctx, ns2, sv1},
			nl3,
			false,
			"",
		},
		{
			"get namespace empty list success",
			fields{stg},
			args{ctx, "empty", sv1},
			nl,
			false,
			"",
		},
		{
			"get namespace info failed empty namespace",
			fields{stg},
			args{ctx, "", sv1},
			nl,
			true,
			"namespace can not be empty",
		},
		{
			"get namespace info failed empty service",
			fields{stg},
			args{ctx, ns1, ""},
			nl,
			true,
			"service can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.ListByService() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("VolumeStorage.ListByService() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.ListByService() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("VolumeStorage.ListByService() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.ListByService() want error = %v, got none", tt.err)
				return
			}
			if !compareVolumeMaps(got, tt.want) {
				t.Errorf("VolumeStorage.ListByService() = %v, want %v", got, tt.want)
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

	n2.Status.State = types.StateReady

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

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
				t.Errorf("VolumeStorage.SetStatus() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err != nil {
				t.Errorf("VolumeStorage.SetSpec() got Get error = %v", err)
				return
			}
			if !compareVolumes(got, tt.want) {
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

	n2.Spec.State.Destroy = true

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("VolumeStorage.SetSpec() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("VolumeStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("VolumeStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("VolumeStorage.SetSpec() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err != nil {
				t.Errorf("VolumeStorage.SetSpec() got Get error = %v", err)
				return
			}
			if !compareVolumes(got, tt.want) {
				t.Errorf("VolumeStorage.SetSpec() = %v, want %v", got, tt.want)
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("VolumeStorage.Update() storage setup error = %v", err)
					return
				}
			}

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
				t.Errorf("VolumeStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err != nil {
				t.Errorf("VolumeStorage.Update() got Get error = %v", err)
				return
			}
			if !compareVolumes(got, tt.want) {
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("VolumeStorage.Remove() storage setup error = %v", err)
				return
			}

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
				t.Errorf("VolumeStorage.Remove() want error = %v, got none", tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.volume.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("VolumeStorage.Remove() got Get error= %v", err)
				return
			}

		})
	}
}

func TestVolumeStorage_Watch(t *testing.T) {

	var (
		err     error
		stg     = newVolumeStorage()
		ctx     = context.Background()
		n       = getVolumeAsset("ns1", "test1", "")
		volumeC = make(chan *types.Volume)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("VolumeStorage.Watch() storage setup error = %v", err)
	}
	defer destroy()

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
			"check volume watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("VolumeStorage.Watch() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.Watch(ctxT, volumeC)
				if err != nil {
					t.Errorf("VolumeStorage.Watch() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/volumes/ns1:test1/meta"
			value := `{"name":"test1","description":"","self_link":"ns1:test1","labels":null,"created":"2018-04-27T10:02:25.58067+03:00","updated":"0001-01-01T00:00:00Z","namespace":"ns1"}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-volumeC:
					t.Log("VolumeStorage.Watch() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("VolumeStorage.Watch() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestVolumeStorage_WatchSpec(t *testing.T) {

	var (
		err     error
		stg     = newVolumeStorage()
		ctx     = context.Background()
		n       = getVolumeAsset("ns1", "test1", "")
		volumeC = make(chan *types.Volume)
	)
	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("VolumeStorage.WatchSpec() storage setup error = %v", err)
	}
	defer destroy()

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
			"check volume watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.WatchSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("VolumeStorage.WatchSpec() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.WatchSpec(ctxT, volumeC)
				if err != nil {
					t.Errorf("VolumeStorage.WatchSpec() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/volumes/ns1:test1/spec"
			value := `{"state":{"destroy":false}}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-volumeC:
					t.Log("VolumeStorage.WatchSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("VolumeStorage.WatchSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestVolumeStorage_WatchStatus(t *testing.T) {

	var (
		err     error
		stg     = newVolumeStorage()
		ctx     = context.Background()
		n       = getVolumeAsset("ns1", "test1", "")
		volumeC = make(chan *types.Volume)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("VolumeStorage.WatchStatus() storage setup error = %v", err)
	}
	defer destroy()

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
			"check volume watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("VolumeStorage.WatchStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("VolumeStorage.WatchStatus() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.WatchStatus(ctxT, volumeC)
				if err != nil {
					t.Errorf("VolumeStorage.WatchStatus() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/volumes/ns1:test1/status"
			value := `{"state":"","message":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-volumeC:
					t.Log("VolumeStorage.WatchStatus() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("VolumeStorage.WatchStatus() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func getVolumeAsset(namespace, name, desc string) types.Volume {

	var n = types.Volume{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	n.Meta.Created = time.Now()

	return n
}

//compare two volume structures
func compareVolumes(got, want *types.Volume) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		reflect.DeepEqual(got.Spec, want.Spec) &&
		reflect.DeepEqual(got.Status, want.Status) {
		result = true
	}

	return result
}

func compareVolumeMaps(got, want map[string]*types.Volume) bool {
	for k, v := range got {
		if !compareVolumes(v, want[k]) {
			return false
		}
	}
	return true
}
