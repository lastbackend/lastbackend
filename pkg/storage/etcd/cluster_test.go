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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

// Test cluster storage set status method
func TestClusterStorage_SetStatus(t *testing.T) {

	initStorage()

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c = getClusterAsset()
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := stg.Clear(ctx)
			if err != nil {
				t.Errorf("ClusterStorage.Insert() storage setup error = %v", err)
				return
			}

			err = tt.fields.stg.SetStatus(tt.args.ctx, &tt.args.cluster.Status)
			if err != nil {
				if tt.wantErr && (err.Error() != tt.err) {
					t.Errorf("ClusterStorage.Insert() error = %v, want errorr %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("ClusterStorage.Insert() error = %v, want no errorr", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ClusterStorage.Insert() want error = %v, got none", tt.err)
				return
			}
			//check that really inserted
			got, err := tt.fields.stg.Get(tt.args.ctx)
			if err != nil {
				t.Errorf("ClusterStorage.SetSpec() got Get error = %s", err.Error())
				return
			}
			if !compareClusters(got, tt.want) {
				t.Errorf("ClusterStorage.Insert() = %v, want %v", got, tt.want)
			}

		})
	}
}

// Test cluster storage get return method
func TestClusterStorage_Get(t *testing.T) {

	initStorage()

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c   = getClusterAsset()
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
		err     string
	}{
		{
			"get cluster info successful",
			fields{stg},
			args{ctx},
			&c,
			false,
			"",
		},
		{
			"test failed no entity",
			fields{stg},
			args{ctx},
			nil,
			true,
			store.ErrEntityNotFound,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ClusterStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if !tt.wantErr {
				if err := stg.SetStatus(ctx, &c.Status); err != nil {
					t.Errorf("ClusterStorage.Get() storage setup error = %v", err)
					return
				}
			}

			got, err := tt.fields.stg.Get(tt.args.ctx)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ClusterStorage.Get() error = %v, wantErr %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("ClusterStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("ClusterStorage.Get() want error = %v, got none", tt.err)
				return
			}

			if !compareClusters(got, tt.want) {
				t.Errorf("ClusterStorage.Get() = %v, want %v", got, tt.want)
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

//compare two cluster structures
func compareClusters(got, want *types.Cluster) bool {
	result := false
	if reflect.DeepEqual(got.Status, want.Status) {
		result = true
	}

	return result
}

func getClusterAsset() types.Cluster {
	var c = types.Cluster{
		Status: types.ClusterStatus{
			Nodes: types.ClusterStatusNodes{
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

	return c
}
