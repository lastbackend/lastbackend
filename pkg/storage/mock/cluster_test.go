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

func TestClusterStorage_Info(t *testing.T) {

	var (
		stg = newClusterStorage()
		ctx = context.Background()
		c   = types.Cluster{}
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
			"1",
			{stg},
			{ctx},
			&c,
			false,
		},
	}

	for _, tt := range tests {
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

func TestClusterStorage_Update(t *testing.T) {
	type fields struct {
		Cluster storage.Cluster
	}
	type args struct {
		ctx     context.Context
		cluster *types.Cluster
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
			s := &ClusterStorage{
				Cluster: tt.fields.Cluster,
			}
			if err := s.Update(tt.args.ctx, tt.args.cluster); (err != nil) != tt.wantErr {
				t.Errorf("ClusterStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newClusterStorage(t *testing.T) {
	tests := []struct {
		name string
		want *ClusterStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newClusterStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClusterStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
