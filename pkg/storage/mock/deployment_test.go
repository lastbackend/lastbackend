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

func TestDeploymentStorage_Get(t *testing.T) {
	type fields struct {
		Deployment storage.Deployment
	}
	type args struct {
		ctx       context.Context
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DeploymentStorage{
				Deployment: tt.fields.Deployment,
			}
			got, err := s.Get(tt.args.ctx, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentStorage_updateState(t *testing.T) {
	type fields struct {
		Deployment storage.Deployment
	}
	type args struct {
		ctx        context.Context
		deployment *types.Deployment
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
			s := &DeploymentStorage{
				Deployment: tt.fields.Deployment,
			}
			if err := s.updateState(tt.args.ctx, tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.updateState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeploymentStorage_SpecWatch(t *testing.T) {
	type fields struct {
		Deployment storage.Deployment
	}
	type args struct {
		ctx        context.Context
		deployment chan *types.Deployment
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
			s := &DeploymentStorage{
				Deployment: tt.fields.Deployment,
			}
			if err := s.SpecWatch(tt.args.ctx, tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentStorage.SpecWatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newDeploymentStorage(t *testing.T) {
	tests := []struct {
		name string
		want *DeploymentStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDeploymentStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDeploymentStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
