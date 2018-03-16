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

	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

func TestEndpointStorage_Get(t *testing.T) {
	type fields struct {
		Endpoint storage.Endpoint
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EndpointStorage{
				Endpoint: tt.fields.Endpoint,
			}
			got, err := s.Get(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndpointStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpointStorage_Upsert(t *testing.T) {
	type fields struct {
		Endpoint storage.Endpoint
	}
	type args struct {
		ctx  context.Context
		name string
		ips  []string
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
			s := &EndpointStorage{
				Endpoint: tt.fields.Endpoint,
			}
			if err := s.Upsert(tt.args.ctx, tt.args.name, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("EndpointStorage.Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndpointStorage_Remove(t *testing.T) {
	type fields struct {
		Endpoint storage.Endpoint
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
			s := &EndpointStorage{
				Endpoint: tt.fields.Endpoint,
			}
			if err := s.Remove(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("EndpointStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndpointStorage_Watch(t *testing.T) {
	type fields struct {
		Endpoint storage.Endpoint
	}
	type args struct {
		ctx      context.Context
		endpoint chan string
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
			s := &EndpointStorage{
				Endpoint: tt.fields.Endpoint,
			}
			if err := s.Watch(tt.args.ctx, tt.args.endpoint); (err != nil) != tt.wantErr {
				t.Errorf("EndpointStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newEndpointStorage(t *testing.T) {
	tests := []struct {
		name string
		want *EndpointStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEndpointStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEndpointStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
