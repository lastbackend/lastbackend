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
)

func TestSystemStorage_ProcessSet(t *testing.T) {

	initStorage()

	var (
		stg = newSystemStorage()
		ctx = context.Background()
		p   = types.Process{}
	)

	type fields struct {
		stg storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{"test process set",
			fields{stg},
			args{ctx, &p},
			false,
			"",
		},
		{"test process nil",
			fields{stg},
			args{ctx, nil},
			true,
			"process can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ProcessSet() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := stg.ProcessSet(tt.args.ctx, tt.args.process)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SystemStorage.ProcessSet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("SystemStorage.ProcessSet() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func TestSystemStorage_Elect(t *testing.T) {

	initStorage()

	var (
		stg = newSystemStorage()
		ctx = context.Background()
		p   = types.Process{}
	)

	type fields struct {
		stg storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{"test process elect",
			fields{stg},
			args{ctx, &p},
			false,
			"",
		},
		{"test process nil",
			fields{stg},
			args{ctx, nil},
			true,
			"process can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.Elect() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()
			b, err := stg.Elect(tt.args.ctx, tt.args.process)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SystemStorage.Elect() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("SystemStorage.Elect() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !b {
				t.Errorf("SystemStorage.Elect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemStorage_ElectUpdate(t *testing.T) {

	initStorage()

	var (
		stg = newSystemStorage()
		ctx = context.Background()
		p   = types.Process{}
	)

	type fields struct {
		stg storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{"test process elect update",
			fields{stg},
			args{ctx, &p},
			false,
			"",
		},
		{"test process nil",
			fields{stg},
			args{ctx, nil},
			true,
			"process can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ElectUpdate() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			_, err := stg.Elect(ctx, tt.args.process)
			if err != nil && !tt.wantErr {
				t.Errorf("SystemStorage.ElectUpdate() set storage err = %v", err)
				return
			}

			err = stg.ElectUpdate(tt.args.ctx, tt.args.process)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("SystemStorage.ElectUpdate() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					t.Errorf("SystemStorage.ElectUpdate() got error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestSystemStorage_ElectWait(t *testing.T) {

	var (
		stg = newSystemStorage()
		ctx = context.Background()
		p   = types.Process{}
	)
	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("SystemStorage.ElectWait() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
		stg storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
		ch      chan bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"test process elect wait",
			fields{stg},
			args{ctx, &p, make(chan bool)},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ElectWait() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear()
			defer clear()

			//insert
			_, err := stg.Elect(ctx, tt.args.process)
			if err != nil && !tt.wantErr {
				t.Errorf("SystemStorage.ElectWait() set storage err = %v", err)
				return
			}
			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch gofunction
			go func() {
				err = stg.ElectWait(ctxT, &p, tt.args.ch)
				if err != nil {
					t.Errorf("SystemStorage.ElectWait() storage setup error = %v", err)
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
			key := "/lstbknd/system/lead"
			value := ""
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s, path=%s", err.Error(), path)
			}

			for {
				select {
				case <-tt.args.ch:
					t.Log("SystemStorage.ElectWait() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("SystemStorage.ElectWait() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func Test_newSystemStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.System
	}{
		{"initialize storage",
			newSystemStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSystemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSystemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
