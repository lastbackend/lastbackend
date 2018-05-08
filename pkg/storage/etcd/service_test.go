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
		ns   string
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
			args{ctx, "test2", ns1},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get service info successful",
			fields{stg},
			args{ctx, "test", ns1},
			&d,
			false,
			"",
		},
		{
			"get service info failed empty namespace",
			fields{stg},
			args{ctx, "test", ""},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get service info failed empty name",
			fields{stg},
			args{ctx, "", ns1},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("ServiceStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.name)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !compareServices(got, tt.want) {
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
			t.Errorf("ServiceStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("ServiceStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("ServiceStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.ListByNamespace() want error = %v, got none", tt.err)
				return
			}

			if !compareServiceMaps(got, tt.want) {
				t.Errorf("ServiceStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_SetStatus(t *testing.T) {

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

	n2.Status.State = types.StateReady

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("ServiceStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.SetStatus() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if err != nil {
				t.Errorf("ServiceStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !compareServices(got, tt.want) {
				t.Errorf("ServiceStorage.SetStatus() = %v, want %v", got, tt.want)
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.SetSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("ServiceStorage.SetSpec() storage setup error = %v", err)
					return
				}
			}

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
				t.Errorf("ServiceStorage.SetSpec() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if err != nil {
				t.Errorf("ServiceStorage.SetSpec() got Get error = %s", err.Error())
				return
			}
			if !compareServices(got, tt.want) {
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Insert() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Insert() want error = %v, got none", tt.err)
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("ServiceStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Update() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, tt.args.service.Meta.Name)
			if err != nil {
				t.Errorf("ServiceStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareServices(got, tt.want) {
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("ServiceStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.service)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ServiceStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ServiceStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ServiceStorage.Remove() want error = %v, got none", tt.err)
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

	var (
		stg    = newServiceStorage()
		ctx    = context.Background()
		err    error
		s      = getServiceAsset("ns1", "test1", "")
		eventC = make(chan *types.Event)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("ServiceStorage.Watch() storage setup error = %v", err)
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
			"check service watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &s); err != nil {
				t.Errorf("ServiceStorage.Watch() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.Watch(ctxT, eventC)
				if err != nil {
					t.Errorf("ServiceStorage.Watch() storage setup error = %v", err)
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
			key := "/lstbknd/services/ns1:test1/meta"
			value := `{"name":"test1","description":"","self_link":"","labels":null,"created":"2018-04-27T09:30:19.278018+03:00","updated":"0001-01-01T00:00:00Z","namespace":"ns1","selflink":"ns1:test1","endpoint":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-eventC:
					t.Log("ServiceStorage.Watch() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("ServiceStorage.Watch() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestServiceStorage_WatchSpec(t *testing.T) {

	var (
		stg      = newServiceStorage()
		ctx      = context.Background()
		err      error
		n        = getServiceAsset("ns1", "test1", "")
		serviceC = make(chan *types.Service)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("ServiceStorage.WatchSpec() storage setup error = %v", err)
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
			"check service watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.WatchSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("ServiceStorage.WatchSpec() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.WatchSpec(ctxT, serviceC)
				if err != nil {
					t.Errorf("ServiceStorage.WatchSpec() storage setup error = %v", err)
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
			key := "/lstbknd/services/ns1:test1/spec"
			value := `{"meta":{"name":"","description":"","self_link":"","labels":null,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"},"replicas":0,"state":{"destroy":false,"maintenance":false},"strategy":{"type":"","rollingOptions":{"period_update":0,"interval":0,"timeout":0,"max_unavailable":0,"max_surge":0},"resources":{},"deadline":0},"triggers":null,"selector":{},"template":{"volumes":null,"container":null,"termination":0}}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-serviceC:
					t.Log("ServiceStorage.WatchSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("ServiceStorage.WatchSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestServiceStorage_WatchStatus(t *testing.T) {

	var (
		stg      = newServiceStorage()
		ctx      = context.Background()
		err      error
		n        = getServiceAsset("ns1", "test1", "")
		serviceC = make(chan *types.Service)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("ServiceStorage.WatchStatus() storage setup error = %v", err)
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
			"check service watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ServiceStorage.WatchStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("ServiceStorage.WatchStatus() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.WatchStatus(ctxT, serviceC)
				if err != nil {
					t.Errorf("ServiceStorage.WatchStatus() storage setup error = %v", err)
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
			key := "/lstbknd/services/ns1:test1/status"
			value := `{"state":"","message":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-serviceC:
					t.Log("ServiceStorage.WatchStatus() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("ServiceStorage.WatchStatus() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
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

	n.Meta.Created = time.Now()

	return n
}

//compare two service structures
func compareServices(got, want *types.Service) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		(got.Meta.Endpoint == want.Meta.Endpoint) &&
		(got.Meta.SelfLink == want.Meta.SelfLink) &&
		reflect.DeepEqual(got.Spec, want.Spec) &&
		reflect.DeepEqual(got.Deployments, want.Deployments) &&
		reflect.DeepEqual(got.Status, want.Status) {
		result = true
	}

	return result
}

func compareServiceMaps(got, want map[string]*types.Service) bool {
	for k, v := range got {
		if !compareServices(v, want[k]) {
			return false
		}
	}
	return true
}
