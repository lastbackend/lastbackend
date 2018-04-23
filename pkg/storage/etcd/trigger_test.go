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

func TestTriggerStorage_Get(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		d   = getTriggerAsset(ns1, svc, "test", "")
	)

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx  context.Context
		name string
		ns   string
		svc  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
		wantErr bool
		err     string
	}{
		{
			"get trigger info failed",
			fields{stg},
			args{ctx, "test2", ns1, svc},
			&d,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get trigger info successful",
			fields{stg},
			args{ctx, "test", ns1, svc},
			&d,
			false,
			"",
		},
		{
			"get trigger info failed empty namespace",
			fields{stg},
			args{ctx, "test", "", svc},
			&d,
			true,
			"namespace can not be empty",
		},
		{
			"get trigger info failed empty service",
			fields{stg},
			args{ctx, "test", ns1, ""},
			&d,
			true,
			"service can not be empty",
		},
		{
			"get trigger info failed empty name",
			fields{stg},
			args{ctx, "", ns1, svc},
			&d,
			true,
			"name can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("TriggerStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &d); err != nil {
				t.Errorf("TriggerStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.ns, tt.args.svc, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("TriggerStorage.Get() error = %v, want no error", err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.Get() wantErr %v, got none", tt.err)
				return
			}

			if !compareTriggers(got, tt.want) {
				t.Errorf("TriggerStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestTriggerStorage_ListByNamespace(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test1", "")
		n2  = getTriggerAsset(ns1, svc, "test2", "")
		n3  = getTriggerAsset(ns2, svc, "test1", "")
		nl  = make(map[string]*types.Trigger, 0)
	)

	nl0 := map[string]*types.Trigger{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3

	nl1 := map[string]*types.Trigger{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Trigger{}
	nl2[stg.keyGet(&n3)] = &n3

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx context.Context
		ns  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Trigger
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
			t.Errorf("TriggerStorage.ListByNamespace() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("TriggerStorage.ListByNamespace() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByNamespace(tt.args.ctx, tt.args.ns)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("TriggerStorage.ListByNamespace() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("TriggerStorage.ListByNamespace() want error = %v, got none", tt.err)
				return
			}

			if !compareTriggerMaps(got, tt.want) {
				t.Errorf("TriggerStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTriggerStorage_ListByService(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		ns2 = "ns2"
		sv1 = "svc1"
		sv2 = "svc2"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, sv1, "test1", "")
		n2  = getTriggerAsset(ns1, sv1, "test2", "")
		n3  = getTriggerAsset(ns1, sv2, "test1", "")
		n4  = getTriggerAsset(ns2, sv1, "test1", "")
		n5  = getTriggerAsset(ns2, sv1, "test2", "")
		nl  = make(map[string]*types.Trigger, 0)
	)

	nl0 := map[string]*types.Trigger{}
	nl0[stg.keyGet(&n1)] = &n1
	nl0[stg.keyGet(&n2)] = &n2
	nl0[stg.keyGet(&n3)] = &n3
	nl0[stg.keyGet(&n4)] = &n4
	nl0[stg.keyGet(&n5)] = &n5

	nl1 := map[string]*types.Trigger{}
	nl1[stg.keyGet(&n1)] = &n1
	nl1[stg.keyGet(&n2)] = &n2

	nl2 := map[string]*types.Trigger{}
	nl2[stg.keyGet(&n3)] = &n3

	nl3 := map[string]*types.Trigger{}
	nl3[stg.keyGet(&n4)] = &n4
	nl3[stg.keyGet(&n5)] = &n5

	type fields struct {
		stg storage.Trigger
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
		want    map[string]*types.Trigger
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
			args{ctx, "t", "t"},
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
			t.Errorf("TriggerStorage.ListByService() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("TriggerStorage.ListByService() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.ListByService(tt.args.ctx, tt.args.ns, tt.args.svc)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.ListByService() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr {
					t.Errorf("TriggerStorage.ListByService() error = %v, want no error", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("TriggerStorage.ListByService() want error = %v, got none", tt.err)
				return
			}

			if !compareTriggerMaps(got, tt.want) {
				t.Errorf("TriggerStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTriggerStorage_SetSpec(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test1", "")
		n2  = getTriggerAsset(ns1, svc, "test1", "")
		n3  = getTriggerAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Trigger, 0)
	)

	//TODO change spec for Trigger
	n2.Spec = types.TriggerSpec{}

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx     context.Context
		trigger *types.Trigger
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
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
			t.Errorf("TriggerStorage.SetSpec() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("TriggerStorage.SetSpec() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetSpec(tt.args.ctx, tt.args.trigger)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("TriggerStorage.SetSpec() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.SetSpec() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.SetSpec() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.trigger.Meta.Namespace, tt.args.trigger.Meta.Service, tt.args.trigger.Meta.Name)
			if err != nil {
				t.Errorf("TriggerStorage.SetSpec() got Get error %v", err)
				return
			}
			if !compareTriggers(got, tt.want) {
				t.Errorf("TriggerStorage.SetSpec() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestTriggerStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test1", "")
		n2  = getTriggerAsset(ns1, svc, "test1", "")
		n3  = getTriggerAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Trigger, 0)
	)

	n2.Status.State = types.StateReady

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx     context.Context
		trigger *types.Trigger
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
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
			t.Errorf("TriggerStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("TriggerStorage.SetStatus() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.trigger)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("TriggerStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.SetStatus() want error = %v, got none %v", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.trigger.Meta.Namespace, tt.args.trigger.Meta.Service, tt.args.trigger.Meta.Name)
			if err != nil {
				t.Errorf("TriggerStorage.SetStatus() got Get error %v", err)
				return
			}
			if !compareTriggers(got, tt.want) {
				t.Errorf("TriggerStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestTriggerStorage_Insert(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test", "")
		n2  = getTriggerAsset(ns1, svc, "", "")
	)

	n2.Meta.Name = ""

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx     context.Context
		trigger *types.Trigger
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
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
			t.Errorf("TriggerStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.trigger)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("TriggerStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.Insert() want error = %v, got none", tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.Insert() error = %v, want %v", err, tt.err)
				return
			}
		})
	}
}

func TestTriggerStorage_Update(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test1", "")
		n2  = getTriggerAsset(ns1, svc, "test1", "test")
		n3  = getTriggerAsset(ns1, svc, "test2", "")
		nl  = make([]*types.Trigger, 0)
	)

	nl0 := append(nl, &n1)

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx     context.Context
		trigger *types.Trigger
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
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
			t.Errorf("TriggerStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl0 {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("TriggerStorage.Update() storage setup error = %v", err)
					return
				}
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.trigger)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("TriggerStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.Update() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.trigger.Meta.Name)
			if err != nil {
				t.Errorf("TriggerStorage.Update() got Get error %v", err)
				return
			}
			if !compareTriggers(got, tt.want) {
				t.Errorf("TriggerStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestTriggerStorage_Remove(t *testing.T) {

	initStorage()

	var (
		ns1 = "ns1"
		svc = "svc"
		stg = newTriggerStorage()
		ctx = context.Background()
		n1  = getTriggerAsset(ns1, svc, "test1", "")
		n2  = getTriggerAsset(ns1, svc, "test2", "")
	)

	type fields struct {
		stg storage.Trigger
	}

	type args struct {
		ctx     context.Context
		trigger *types.Trigger
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Trigger
		wantErr bool
		err     string
	}{
		{
			"test successful trigger remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil trigger structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: trigger not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("TriggerStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("TriggerStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.trigger)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("TriggerStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("TriggerStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("TriggerStorage.Remove() want error = %v, got none", tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, ns1, svc, tt.args.trigger.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("TriggerStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestTriggerStorage_Watch(t *testing.T) {

	initStorage()

	var (
		err      error
		stg      = newTriggerStorage()
		ctx      = context.Background()
		n        = getTriggerAsset("ns1", "svc", "test1", "")
		triggerC = make(chan *types.Trigger)
		stopC    = make(chan int)
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
			"check trigger watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("TriggerStorage.Watch() storage setup error = %v", err)
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
					t.Errorf("TriggerStorage.Watch() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.Watch(ctx, triggerC)
					if err != nil {
						t.Errorf("TriggerStorage.Watch() storage setup error = %v", err)
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
				t.Errorf("TriggerStorage.Watch() update error =%v", err)
				return

			case <-triggerC:
				t.Log("TriggerStorage.Watch() is working")
				return
			}
		})
	}
}

func TestTriggerStorage_WatchSpec(t *testing.T) {

	initStorage()

	var (
		err      error
		stg      = newTriggerStorage()
		ctx      = context.Background()
		n        = getTriggerAsset("ns1", "svc", "test1", "")
		triggerC = make(chan *types.Trigger)
		stopC    = make(chan int)
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
			"check trigger watch spec",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("TriggerStorage.WatchSpec() storage setup error = %v", err)
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
					t.Errorf("TriggerStorage.WatchSpec() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchSpec(ctx, triggerC)
					if err != nil {
						t.Errorf("TriggerStorage.WatchSpec() storage setup error = %v", err)
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
				t.Errorf("TriggerStorage.WatchSpec() update error =%v", err)
				return

			case <-triggerC:
				t.Log("TriggerStorage.WatchSpec() is working")
				return
			}
		})
	}
}

func TestTriggerStorage_WatchStatus(t *testing.T) {

	initStorage()

	var (
		err      error
		stg      = newTriggerStorage()
		ctx      = context.Background()
		n        = getTriggerAsset("ns1", "svc", "test1", "")
		triggerC = make(chan *types.Trigger)
		stopC    = make(chan int)
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
			"check trigger watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("TriggerStorage.WatchStatus() storage setup error = %v", err)
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
					t.Errorf("TriggerStorage.WatchStatus() insert error = %v", err)
					return
				}
				//run watch go function
				go func() {
					err = stg.WatchStatus(ctx, triggerC)
					if err != nil {
						t.Errorf("TriggerStorage.WatchStatus() storage setup error = %v", err)
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
				t.Errorf("TriggerStorage.WatchStatus() update error =%v", err)
				return

			case <-triggerC:
				t.Log("TriggerStorage.WatchStatus() is working")
				return
			}
		})
	}
}

func Test_newTriggerStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Trigger
	}{
		{"initialize storage",
			newTriggerStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTriggerStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTriggerStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTriggerAsset(namespace, service, name, desc string) types.Trigger {

	var n = types.Trigger{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Service = service
	n.Meta.Description = desc

	n.Meta.Created = time.Now()

	return n
}

//compare two secret structures
func compareTriggers(got, want *types.Trigger) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Meta.Service == want.Meta.Service) &&
		(got.Meta.Namespace == want.Meta.Namespace) &&
		reflect.DeepEqual(got.Spec, want.Spec) &&
		reflect.DeepEqual(got.Status, want.Status) {
		result = true
	}

	return result
}

func compareTriggerMaps(got, want map[string]*types.Trigger) bool {
	for k, v := range got {
		if !compareTriggers(v, want[k]) {
			return false
		}
	}
	return true
}
