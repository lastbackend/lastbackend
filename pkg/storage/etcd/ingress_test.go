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

func TestIngressStorage_List(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test1", "")
		n2  = getIngressAsset("test2", "")
		nl  = make(map[string]*types.Ingress, 2)
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
		err     string
	}{
		{
			"get ingress list success",
			fields{stg},
			args{ctx},
			nl,
			false,
			"",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		clear()
		defer clear()

		for _, n := range nl {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("IngressStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.stg.List(tt.args.ctx)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.List() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.List() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.List() want error = %v, got none", tt.err)
				return
			}

			if !compareIngressLists(got, tt.want) {
				t.Errorf("IngressStorage.List() = %v\n, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressStorage_Get(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n   = getIngressAsset("test", "desc")
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
			"get ingress info failed",
			fields{stg},
			args{ctx, "test2"},
			&n,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get ingress info successful",
			fields{stg},
			args{ctx, "test"},
			&n,
			false,
			"",
		},
		{
			"get ingress info failed empty name",
			fields{stg},
			args{ctx, ""},
			&n,
			true,
			"ingress can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("IngressStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.name)
			t.Logf("got=%v, err=%v", got, err)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.Get() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("IngressStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !compareIngress(got, tt.want) {
				t.Errorf("IngressStorage.Get() = %v\n, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_Insert(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test", "desc")
		n2  = getIngressAsset("", "")
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.ingress)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Insert() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.Insert() error = %v, wantErr %v", err, tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Insert() want error = %v, got none", tt.err)
				return
			}

			//check
			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if err != nil {
				t.Errorf("IngressStorage.Insert() got Get error = %s", err.Error())
				return
			}
			if !compareIngress(got, tt.want) {
				t.Errorf("IngressStorage.Insert() = %v\n, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_Update(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test", "desc")
		n2  = getIngressAsset("", "")
		n3  = getIngressAsset("test", "new desc")
		nl  = make([]*types.Ingress, 0)
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
			args{ctx, &n1},
			&n3,
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
			"test failed update: invalid structure",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("IngressStorage.Update() storage setup error = %v", err)
					return
				}
			}
			if !tt.wantErr {
				tt.args.ingress.Meta.Meta.Description = "new desc"
			}
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Update() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.Update() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("IngressStorage.Update() want error = %v, got none", tt.err)
				return
			}

			//check
			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if err != nil {
				t.Errorf("IngressStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareIngress(got, tt.want) {
				t.Errorf("IngressStorage.Update() = %v\n, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test", "desc")
		n2  = getIngressAsset("", "")
		n3  = getIngressAsset("test", "desc")
		nl  = make([]*types.Ingress, 0)
	)
	nl0 := append(nl, &n1)
	n3.Status.Ready = true

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
			"test successful set status",
			fields{stg},
			args{ctx, &n1},
			&n3,
			false,
			"",
		},
		{
			"test failed set status: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed set status: invalid structure",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("IngressStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}
			if !tt.wantErr {
				tt.args.ingress.Status.Ready = true
			}
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.SetStatus() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.SetStatus() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("IngressStorage.SetStatus() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ingress.Meta.Name)
			if err != nil {
				t.Errorf("IngressStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !compareIngress(got, tt.want) {
				t.Errorf("IngressStorage.SetStatus() = %v\n, want %v", got, tt.want)
			}

		})
	}
}

func TestIngressStorage_Remove(t *testing.T) {

	initStorage()

	var (
		stg = newIngressStorage()
		ctx = context.Background()
		n1  = getIngressAsset("test", "desc")
		n2  = getIngressAsset("", "")
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
			"test successful remove",
			fields{stg},
			args{ctx, &n1},
			nil,
			false,
			"",
		},
		{
			"test failed remove: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed remove: invalid structure",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("IngressStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.ingress)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.Remove() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.Remove() got error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("IngressStorage.Remove() want error = %v, got none", tt.err)
				return
			}

		})
	}
}

/*
func TestIngressStorage_GetSpec(t *testing.T) {

	initStorage()

	var (
		stg  = newIngressStorage()
		ctx  = context.Background()
		n    = getIngressAsset("test", "desc")
		n1   = getIngressAsset("test1", "")
		spec = types.IngressSpec{}
		rs   = types.RouteSpec{
			Domain: "domain",
		}
	)
	spec.Routes = make(map[string]types.RouteSpec, 1)
	spec.Routes[n.Meta.Name] = rs
	n.Spec = spec

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
			"get ingress spec info failed",
			fields{stg},
			args{ctx, &n1},
			&spec,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get ingress spec info successful",
			fields{stg},
			args{ctx, &n},
			&spec,
			false,
			"",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.GetSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("IngressStorage.GetSpec() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.ingress)
			t.Log("got=", got)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("IngressStorage.GetSpec() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("IngressStorage.GetSpec() error = %v, wantErr %v", err, tt.err)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressStorage.GetSpec() = %v\n, want %v", got, tt.want)
			}

		})
	}
}
*/

func TestIngressStorage_Watch(t *testing.T) {

	var (
		err      error
		stg      = newIngressStorage()
		ctx      = context.Background()
		n        = getIngressAsset("test", "desc")
		ingressC = make(chan *types.Ingress)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("IngressStorage.Watch() storage setup error = %v", err)
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
			"check ingress watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("DeploymentStorage.Watch() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.Watch(ctxT, ingressC)
				if err != nil {
					t.Errorf("IngressStorage.Watch() storage setup error = %v", err)
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
			key := "/lstbknd/ingress/test/meta"
			value := `{"name":"test","description":"desc","self_link":"","labels":null,"created":"2018-04-26T14:50:05.523091+03:00","updated":"2018-04-26T14:50:05.523097+03:00","cluster":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-ingressC:
					t.Log("IngressStorage.Watch() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("IngressStorage.Watch() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestIngressStorage_EventsStatus(t *testing.T) {

	var (
		err                 error
		stg                 = newIngressStorage()
		ctx                 = context.Background()
		n                   = getIngressAsset("test", "desc")
		ingressStatusEventC = make(chan *types.IngressStatusEvent)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("IngressStorage.EventsStatus() storage setup error = %v", err)
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
			"check ingress watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("IngressStorage.EventsStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("DeploymentStorage.EventsStatus() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.WatchStatus(ctxT, ingressStatusEventC)
				if err != nil {
					t.Errorf("IngressStorage.EventsStatus() storage setup error = %v", err)
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
			key := "/lstbknd/ingress/test/status"
			value := `{"ready":false}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-ingressStatusEventC:
					t.Log("IngressStorage.EventsStatus() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("IngressStorage.EventsStatus() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func getIngressAsset(name, desc string) types.Ingress {

	var n = types.Ingress{}

	n.Meta.Name = name
	n.Meta.Description = desc

	n.Meta.Created = time.Now()
	return n
}

//compare two ingress structures
func compareIngress(got, want *types.Ingress) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Cluster == want.Meta.Cluster) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Spec, want.Spec) {
		result = true
	}

	return result
}

func compareIngressLists(got, want map[string]*types.Ingress) bool {
	for k, v := range got {
		if !compareIngress(v, want[k]) {
			return false
		}
	}
	return true
}
