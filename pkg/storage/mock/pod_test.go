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

func TestPodStorage_GetByName(t *testing.T) {
	type fields struct {
		Pod storage.Pod
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
		want    *types.Pod
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			got, err := s.Get(tt.args.ctx, tt.args.app, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByNamespace(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx context.Context
		app string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Pod
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			got, err := s.ListByNamespace(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.ListByNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.ListByNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_ListByService(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx       context.Context
		namespace string
		service   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Pod
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			got, err := s.ListByService(tt.args.ctx, tt.args.namespace, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.ListByService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodStorage.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodStorage_Upsert(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx context.Context
		pod *types.Pod
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
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			if err := s.Upsert(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPodStorage_Update(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx context.Context
		pod *types.Pod
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
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			if err := s.Update(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPodStorage_Remove(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx context.Context
		pod *types.Pod
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
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			if err := s.Remove(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPodStorage_Watch(t *testing.T) {
	type fields struct {
		Pod storage.Pod
	}
	type args struct {
		ctx context.Context
		pod chan *types.Pod
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
			s := &PodStorage{
				Pod: tt.fields.Pod,
			}
			if err := s.Watch(tt.args.ctx, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("PodStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newPodStorage(t *testing.T) {
	tests := []struct {
		name string
		want *PodStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPodStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPodStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getPodAsset(name, desc string) types.Pod {
	p := types.Pod{}
	p.Meta.Name = name
	p.Meta.Description = desc

	return p
}