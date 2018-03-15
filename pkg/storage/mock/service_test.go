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

func TestServiceStorage_GetByName(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx  context.Context
		app  string
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			got, err := s.GetByName(tt.args.ctx, tt.args.app, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_GetByPodName(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			got, err := s.GetByPodName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.GetByPodName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.GetByPodName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_ListByNamespace(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx context.Context
		app string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Service
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			got, err := s.ListByNamespace(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_CountByNamespace(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx context.Context
		app string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			got, err := s.CountByNamespace(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.CountByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceStorage.CountByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStorage_Insert(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.Insert(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_Update(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.Update(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_UpdateSpec(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.UpdateSpec(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.UpdateSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_Remove(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.Remove(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_RemoveByNamespace(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx context.Context
		app string
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.RemoveByNamespace(tt.args.ctx, tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.RemoveByNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_Watch(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service chan *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.Watch(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_SpecWatch(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service chan *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.SpecWatch(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.SpecWatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_PodsWatch(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service chan *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.PodsWatch(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.PodsWatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStorage_updateState(t *testing.T) {
	type fields struct {
		Service storage.Service
	}
	type args struct {
		ctx     context.Context
		service *types.Service
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
			s := &ServiceStorage{
				Service: tt.fields.Service,
			}
			if err := s.updateState(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("ServiceStorage.updateState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newServiceStorage(t *testing.T) {
	tests := []struct {
		name string
		want *ServiceStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newServiceStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newServiceStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
