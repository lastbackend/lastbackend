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
	}{
		{"dummy test",
			fields{stg},
			args{ctx, &p},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ProcessSet() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := stg.ProcessSet(tt.args.ctx, tt.args.process); (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.ProcessSet() error = %v, wantErr %v", err, tt.wantErr)
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
	}{
		{"dummy test",
			fields{stg},
			args{ctx, &p},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.Elect() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			if b, err := stg.Elect(tt.args.ctx, tt.args.process); !b || (err != nil) != tt.wantErr {
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

	p.ID = "test"

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
	}{
		{"dummy test",
			fields{stg},
			args{ctx, &p},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ElectUpdate() storage setup error = %v", err)
			return
		}

		if _, err := stg.Elect(ctx, tt.args.process); err != nil {
			t.Errorf("SystemStorage.ElectUpdate() set storage err = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := stg.ElectUpdate(tt.args.ctx, tt.args.process); (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.ElectUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemStorage_ElectWait(t *testing.T) {

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
		cn      chan bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"dummy test",
			fields{stg},
			args{ctx, &p, make(chan bool)},
			false,
		},
	}
	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("SystemStorage.ElectUpdate() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

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
