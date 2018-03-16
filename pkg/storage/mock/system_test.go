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
)

func TestSystemStorage_ProcessSet(t *testing.T) {
	type fields struct {
		System storage.System
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SystemStorage{
				System: tt.fields.System,
			}
			if err := s.ProcessSet(tt.args.ctx, tt.args.process); (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.ProcessSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemStorage_Elect(t *testing.T) {
	type fields struct {
		System storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SystemStorage{
				System: tt.fields.System,
			}
			got, err := s.Elect(tt.args.ctx, tt.args.process)
			if (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.Elect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SystemStorage.Elect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemStorage_ElectUpdate(t *testing.T) {
	type fields struct {
		System storage.System
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SystemStorage{
				System: tt.fields.System,
			}
			if err := s.ElectUpdate(tt.args.ctx, tt.args.process); (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.ElectUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemStorage_ElectWait(t *testing.T) {
	type fields struct {
		System storage.System
	}
	type args struct {
		ctx     context.Context
		process *types.Process
		lead    chan bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SystemStorage{
				System: tt.fields.System,
			}
			if err := s.ElectWait(tt.args.ctx, tt.args.process, tt.args.lead); (err != nil) != tt.wantErr {
				t.Errorf("SystemStorage.ElectWait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newSystemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *SystemStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSystemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSystemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
