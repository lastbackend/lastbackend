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

// Test cluster storage insert method
func TestClusterStorage_Insert(t *testing.T) {

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c   = types.Cluster{}
	)

	type fields struct {
		stg storage.Cluster
	}

	type args struct {
		ctx     context.Context
		cluster *types.Cluster
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Cluster
		wantErr bool
		err     string
	}{
		{
			"test successful insert",
			fields{stg},
			args{ctx, &c},
			&c,
			false,
			"",
		},
		{
			"test failed insert",
			fields{stg},
			args{ctx, nil},
			&c,
			true,
			store.ErrStructArgIsNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Insert(tt.args.ctx, tt.args.cluster); (err != nil) != tt.wantErr || (tt.wantErr && (err.Error() != tt.err)) {
				t.Errorf("ClusterStorage.Insert() error = %v, want errorr %v", err, tt.err)
			}
		})
	}
}

// Test cluster storage info return method
func TestClusterStorage_Info(t *testing.T) {

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c   = getClusterAsset("test")
	)

	type fields struct {
		stg storage.Cluster
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Cluster
		wantErr bool
	}{
		{
			"get cluster info successful",
			fields{stg},
			args{ctx},
			&c,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Insert(ctx, &c); err != nil {
			t.Errorf("ClusterStorage.Info() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Info(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterStorage.Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterStorage.Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test cluster storage update method
func TestClusterStorage_Update(t *testing.T) {
	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c1   = getClusterAsset("test1")
		c2   = getClusterAsset("test2")
	)

	type fields struct {
		stg storage.Cluster
	}

	type args struct {
		ctx context.Context
		cluster *types.Cluster
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Cluster
		wantErr bool
		err     string
	}{
		{
			"update cluster info failed",
			fields{stg},
			args{ctx, nil},
			&c2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"update cluster info successful",
			fields{stg},
			args{ctx, &c2},
			&c2,
			false,
			"",
		},
	}
	for _, tt := range tests {

		if err := stg.Insert(ctx, &c1); err != nil {
			t.Errorf("ClusterStorage.Update() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.cluster)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ClusterStorage.Update() = %v, want %v", err, tt.err)
					return
				}
				return
			} else {
				if tt.wantErr {
					t.Errorf("ClusterStorage.Update() error = %v, wantErr %v", err, tt.err)
					return
				}
			}

			got, _ := tt.fields.stg.Info(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterStorage.Update() = %v, want %v", got, tt.want)
			}

		})
	}
}

// Test storage initialization
func Test_newClusterStorage(t *testing.T) {


	tests := []struct {
		name string
		want storage.Cluster
	}{
		{"initialize storage",
			newClusterStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newClusterStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClusterStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getClusterAsset(name string) types.Cluster {
	var c = types.Cluster{
		Meta: types.ClusterMeta{
			Region:   "local",
			Token:    "token",
			Provider: "local",
			Shared:   false,
		},
		State: types.ClusterState{
			Nodes: types.ClusterStateNodes{
				Total:   2,
				Online:  1,
				Offline: 1,
			},
			Capacity: types.ClusterResources{
				Containers: 1,
				Pods:       1,
				Memory:     1024,
				Cpu:        1,
				Storage:    512,
			},
			Allocated: types.ClusterResources{
				Containers: 1,
				Pods:       1,
				Memory:     1024,
				Cpu:        1,
				Storage:    512,
			},
			Deleted: false,
		},
	}

	c.Meta.Name = name

	return c
}