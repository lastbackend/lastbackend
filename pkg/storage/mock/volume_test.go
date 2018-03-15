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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestVolumeStorage_GetByID(t *testing.T) {
	type fields struct {
		Volume storage.Volume
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Volume
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &VolumeStorage{
				Volume: tt.fields.Volume,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumeStorage.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeStorage_ListByProject(t *testing.T) {
	type fields struct {
		Volume storage.Volume
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Volume
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &VolumeStorage{
				Volume: tt.fields.Volume,
			}
			got, err := s.ListByProject(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumeStorage.ListByProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeStorage.ListByProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeStorage_Insert(t *testing.T) {
	type fields struct {
		Volume storage.Volume
	}
	type args struct {
		ctx    context.Context
		volume *types.Volume
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
			s := &VolumeStorage{
				Volume: tt.fields.Volume,
			}
			if err := s.Insert(tt.args.ctx, tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("VolumeStorage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVolumeStorage_Remove(t *testing.T) {
	type fields struct {
		Volume storage.Volume
	}
	type args struct {
		ctx context.Context
		id  string
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
			s := &VolumeStorage{
				Volume: tt.fields.Volume,
			}
			if err := s.Remove(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("VolumeStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newVolumeStorage(t *testing.T) {
	type args struct {
		config store.Config
	}
	tests := []struct {
		name string
		args args
		want *VolumeStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newVolumeStorage(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newVolumeStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
