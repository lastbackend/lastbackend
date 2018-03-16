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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestStorage_Cluster(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Cluster
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Cluster(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Cluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Deployment(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Deployment
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Deployment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Deployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Hook(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Hook
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Hook(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Hook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Node(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Node
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Node(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Namespace(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Namespace
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Namespace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Namespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Route(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Route
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Route(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Pod(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Pod
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Pod(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Pod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Service(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Service
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Service(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Volume(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Volume
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Volume(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Volume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Endpoint(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Endpoint
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.Endpoint(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_System(t *testing.T) {
	type fields struct {
		Context           context.Context
		CancelFunc        context.CancelFunc
		ClusterStorage    *ClusterStorage
		DeploymentStorage *DeploymentStorage
		EndpointStorage   *EndpointStorage
		HookStorage       *HookStorage
		NodeStorage       *NodeStorage
		NamespaceStorage  *NamespaceStorage
		PodStorage        *PodStorage
		ServiceStorage    *ServiceStorage
		RouteStorage      *RouteStorage
		VolumeStorage     *VolumeStorage
		SystemStorage     *SystemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.System
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Context:           tt.fields.Context,
				CancelFunc:        tt.fields.CancelFunc,
				ClusterStorage:    tt.fields.ClusterStorage,
				DeploymentStorage: tt.fields.DeploymentStorage,
				EndpointStorage:   tt.fields.EndpointStorage,
				HookStorage:       tt.fields.HookStorage,
				NodeStorage:       tt.fields.NodeStorage,
				NamespaceStorage:  tt.fields.NamespaceStorage,
				PodStorage:        tt.fields.PodStorage,
				ServiceStorage:    tt.fields.ServiceStorage,
				RouteStorage:      tt.fields.RouteStorage,
				VolumeStorage:     tt.fields.VolumeStorage,
				SystemStorage:     tt.fields.SystemStorage,
			}
			if got := s.System(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.System() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyCreate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyCreate(tt.args.args...); got != tt.want {
				t.Errorf("keyCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Storage
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClient(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    store.Store
		want1   store.DestroyFunc
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getClient(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClient() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
