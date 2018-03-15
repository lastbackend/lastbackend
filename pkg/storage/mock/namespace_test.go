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

func TestNamespaceStorage_GetByName(t *testing.T) {
	type fields struct {
		Namespace storage.Namespace
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Namespace
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NamespaceStorage{
				Namespace: tt.fields.Namespace,
			}
			got, err := s.GetByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceStorage.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceStorage_List(t *testing.T) {
	type fields struct {
		Namespace storage.Namespace
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Namespace
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NamespaceStorage{
				Namespace: tt.fields.Namespace,
			}
			got, err := s.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceStorage_Insert(t *testing.T) {
	type fields struct {
		Namespace storage.Namespace
	}
	type args struct {
		ctx       context.Context
		namespace *types.Namespace
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
			s := &NamespaceStorage{
				Namespace: tt.fields.Namespace,
			}
			if err := s.Insert(tt.args.ctx, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNamespaceStorage_Update(t *testing.T) {
	type fields struct {
		Namespace storage.Namespace
	}
	type args struct {
		ctx       context.Context
		namespace *types.Namespace
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
			s := &NamespaceStorage{
				Namespace: tt.fields.Namespace,
			}
			if err := s.Update(tt.args.ctx, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNamespaceStorage_Remove(t *testing.T) {
	type fields struct {
		Namespace storage.Namespace
	}
	type args struct {
		ctx  context.Context
		name string
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
			s := &NamespaceStorage{
				Namespace: tt.fields.Namespace,
			}
			if err := s.Remove(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("NamespaceStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newNamespaceStorage(t *testing.T) {
	tests := []struct {
		name string
		want *NamespaceStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newNamespaceStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNamespaceStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createNamespace(t *testing.T) {
	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *types.Namespace
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createNamespace(tt.args.name, tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *types.Namespace
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
