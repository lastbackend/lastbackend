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

func TestHookStorage_Get(t *testing.T) {
	type fields struct {
		Hook storage.Hook
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Hook
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HookStorage{
				Hook: tt.fields.Hook,
			}
			got, err := s.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("HookStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HookStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHookStorage_Insert(t *testing.T) {
	type fields struct {
		Hook storage.Hook
	}
	type args struct {
		ctx  context.Context
		hook *types.Hook
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
			s := &HookStorage{
				Hook: tt.fields.Hook,
			}
			if err := s.Insert(tt.args.ctx, tt.args.hook); (err != nil) != tt.wantErr {
				t.Errorf("HookStorage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHookStorage_Remove(t *testing.T) {
	type fields struct {
		Hook storage.Hook
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
			s := &HookStorage{
				Hook: tt.fields.Hook,
			}
			if err := s.Remove(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("HookStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newHookStorage(t *testing.T) {
	tests := []struct {
		name string
		want *HookStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHookStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHookStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
