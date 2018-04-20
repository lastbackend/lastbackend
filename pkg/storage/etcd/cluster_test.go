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
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

// Test cluster storage insert method
func TestClusterStorage_Insert(t *testing.T) {

	initStorage()

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		//c   = types.Cluster{}
		c = getClusterAsset("test")
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

			err := stg.Clear(ctx)
			if err != nil {
				t.Errorf("ClusterStorage.Insert() storage setup error = %v", err)
				return
			}

			err = tt.fields.stg.Insert(tt.args.ctx, tt.args.cluster)
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
				if err := stg.Insert(ctx, &c); err != nil {
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

// Test cluster storage update method
func TestClusterStorage_Update(t *testing.T) {

	initStorage()

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c1  = getClusterAsset("test1")
		c2  = getClusterAsset("test2")
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("ClusterStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &c1); err != nil {
				t.Errorf("ClusterStorage.Update() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.cluster)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ClusterStorage.Update() error = %v, want no error", err.Error())
					return
				}
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("ClusterStorage.Update() error = %v, want %v", err, tt.err)
					return
				}
				return
			}
			if tt.wantErr {
				t.Errorf("ClusterStorage.Update() want error = %v, got none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx)
			if err != nil {
				t.Errorf("ClusterStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareClusters(got, tt.want) {
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

//compare two cluster structures
func compareClusters(got, want *types.Cluster) bool {
	result := false
	if (got.Meta.Name == want.Meta.Name) &&
		(got.Meta.Description == want.Meta.Description) &&
		(got.Meta.SelfLink == want.Meta.SelfLink) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Quotas, want.Quotas) &&
		reflect.DeepEqual(got.Meta.Labels, want.Meta.Labels) {
		result = true
	}

	return result
}

func getClusterAsset(name string) types.Cluster {
	var c = types.Cluster{
		Meta: types.ClusterMeta{
			Region:   "local",
			Token:    "token",
			Provider: "local",
			Shared:   false,
		},
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

	c.Meta.Name = name
	c.Meta.Created = time.Now()

	return c
}
