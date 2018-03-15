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

func TestNodeStorage_List(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Node
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			got, err := s.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStorage_Get(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			got, err := s.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStorage_Insert(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		node *types.Node
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.Insert(tt.args.ctx, tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_Update(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		node *types.Node
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.Update(tt.args.ctx, tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_InsertPod(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		meta *types.NodeMeta
		pod  *types.Pod
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.InsertPod(tt.args.ctx, tt.args.meta, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.InsertPod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_UpdatePod(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		meta *types.NodeMeta
		pod  *types.Pod
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.UpdatePod(tt.args.ctx, tt.args.meta, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.UpdatePod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_RemovePod(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		meta *types.NodeMeta
		pod  *types.Pod
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.RemovePod(tt.args.ctx, tt.args.meta, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.RemovePod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_Remove(t *testing.T) {
	type fields struct {
		Node storage.Node
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.Remove(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeStorage_Watch(t *testing.T) {
	type fields struct {
		Node storage.Node
	}
	type args struct {
		ctx  context.Context
		node chan *types.Node
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
			s := &NodeStorage{
				Node: tt.fields.Node,
			}
			if err := s.Watch(tt.args.ctx, tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newNodeStorage(t *testing.T) {
	tests := []struct {
		name string
		want *NodeStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newNodeStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNodeStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
