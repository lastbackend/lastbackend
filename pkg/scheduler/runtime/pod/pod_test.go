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

package pod

import (
	"testing"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/scheduler/envs"
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"encoding/json"
	"reflect"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func TestProvision(t *testing.T) {

	stg, _ := mock.New()
	envs.Get().SetStorage(stg)

	var (
		ns1 = "ns1"
		svc = "svc"
		dp1 = "dp1"
		dp2 = "dp2"
		ctx = context.Background()
		p1  = getPodAsset(ns1, svc, dp1,"test1", "")  // successful
		p2  = getPodAsset(ns1, svc, dp2,"test2", "")  // can not be provisioned by ram
		p3  = getPodAsset(ns1, svc, dp2,"test3", "")  // not found
		n1  = getNodeAsset("node-1", "", true) // limit 512 RAM
		n2  = getNodeAsset("node-2", "", true) // limit 512 RAM
	)

	n1.State.Capacity.Memory 			= 1024
	n1.State.Capacity.Cpu    			= 8
	n1.State.Capacity.Pods   			= 8
	n1.State.Capacity.Containers  = 8

	n1.State.Allocated.Memory 			= 512
	n1.State.Allocated.Cpu    			= 0
	n1.State.Allocated.Pods   			= 7
	n1.State.Allocated.Containers   = 6

	n2.State.Capacity.Memory 			= 1024
	n2.State.Capacity.Cpu    			= 8
	n2.State.Capacity.Pods   			= 8
	n2.State.Capacity.Containers  = 8

	n2.State.Allocated.Memory 			= 0
	n2.State.Allocated.Cpu    			= 0
	n2.State.Allocated.Pods   			= 0
	n2.State.Allocated.Containers   = 0

	var ips = make([]string, 0)
	ips = append(ips, "8.8.8.8")

	// Set pod1 spec
	p1.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	p1spec := types.SpecTemplateContainer{
		Name: "test-template",
		DNS: types.SpecTemplateContainerDNS{
			Server: ips,
			Search: ips,
		},
	}
	p1spec.SetDefault()
	p1spec.Resources = types.SpecTemplateContainerResources{
		Limits: types.SpecTemplateContainerResource{
			CPU: 0,
			RAM: 256,
		},
		Quota: types.SpecTemplateContainerResource{
			CPU: 0,
			RAM: 256,
		},
	}
	p1.Spec.Template.Containers = append(p1.Spec.Template.Containers, p1spec)

	// Set pod2 spec
	p2.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	p2spec := types.SpecTemplateContainer{
		Name: "test-template",
		DNS: types.SpecTemplateContainerDNS{
			Server: ips,
			Search: ips,
		},
	}
	p2spec.SetDefault()
	p2spec.Resources = types.SpecTemplateContainerResources{
		Limits: types.SpecTemplateContainerResource{
			CPU: 0,
			RAM: 1024,
		},
		Quota: types.SpecTemplateContainerResource{
			CPU: 0,
			RAM: 1024,
		},
	}
	p2.Spec.Template.Containers = append(p2.Spec.Template.Containers, p2spec)

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx  context.Context
		pod *types.Pod
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{
			"provision pod failed: pod not found",
			fields{stg},
			args{ctx, &p3, &n1},
			true,
			store.ErrEntityNotFound,
		},
		{
			"provision pod failed: node not found: memory limit",
			fields{stg},
			args{ctx, &p2, &n1},
			true,
			errors.NodeNotFound,
		},
		{
			"provision pod failed: selected another node",
			fields{stg},
			args{ctx, &p1, &n2},
			false,
			"",
		},

		{
			"provision pod success",
			fields{stg},
			args{ctx, &p1, &n2},
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Node().Clear(ctx); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Pod().Insert(ctx, &p1); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Pod().Insert(ctx, &p2); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Node().Insert(ctx, &n1); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Node().Insert(ctx, &n2); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			// Run provision method
			err := Provision(tt.args.pod)
			if err != nil {

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("Provision() = %v, want %v", err, tt.err)
					return
				}

				if !tt.wantErr {
					t.Errorf("Provision() = %v, want no errors", err)
					return
				}

				if tt.err == errors.NodeNotFound {
					got, err := tt.fields.stg.Pod().Get(tt.args.ctx,
						tt.args.pod.Meta.Namespace,
						tt.args.pod.Meta.Service,
						tt.args.pod.Meta.Deployment,
						tt.args.pod.Meta.Name)

					if err != nil {
						t.Errorf("Provision() = %v, want no errors", err)
						return
					}

					if !got.State.Error {
						t.Errorf("Provision() = %v, want error is true", got.State.Error )
						return
					}

					if got.Status.Stage != types.PodStageError {
						t.Errorf("Provision() = %v, want stage %s", got.Status.Stage, types.PodStageError )
						return
					}

					if got.Status.Message != errors.NodeNotFound {
						t.Errorf("Provision() = %v, want error %s", got.Status.Message, errors.NodeNotFound )
						return
					}
				}

				return
			}

			if err == nil && tt.wantErr {
				t.Errorf("Provision() err %v, want %v", err, tt.err)
				return
			}

			got, err := tt.fields.stg.Node().GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("Provision() = %v, want %v", err, tt.err)
				return
			}

			if _, ok := got.Pods[tt.args.pod.SelfLink()]; !ok {
				t.Errorf("Provision() failed: not found: %s", tt.args.pod.SelfLink())
				return
			}

			sp := got.Pods[tt.args.pod.SelfLink()]

			for i, t := range sp.Template.Containers {
				t.Labels = make(map[string]string)
				sp.Template.Containers[i] = t
			}

			g, err := json.Marshal(sp.Template)
			if err != nil {
				t.Errorf("Provision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.args.pod.Spec.Template)
			if err != nil {
				t.Errorf("Provision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("Provision() spec not match = %v, want %v", string(g), string(e))
			}

		})
	}
}

func getPodAsset(namespace, service, deployment, name, desc string) types.Pod {
	p := types.Pod{}

	p.Meta.Name = name
	p.Meta.Description = desc
	p.Meta.Namespace = namespace
	p.Meta.Service = service
	p.Meta.Deployment = deployment
	p.SelfLink()

	return p
}

func getNodeAsset(name, desc string, online bool) types.Node {
	var n = types.Node{
		Meta: types.NodeMeta{
			Region:   "local",
			Token:    "token",
			Provider: "local",
		},
		Info: types.NodeInfo{
			Hostname: name,
		},
		State: types.NodeState{
			Capacity: types.NodeResources{
				Containers: 2,
				Pods:       2,
				Memory:     1024,
				Cpu:        2,
				Storage:    512,
			},
			Allocated: types.NodeResources{
				Containers: 1,
				Pods:       1,
				Memory:     512,
				Cpu:        1,
				Storage:    256,
			},
		},
		Spec: types.NodeSpec{
			Pods:    make(map[string]types.PodSpec),
			Volumes: make(map[string]types.VolumeSpec),
			Routes:  make(map[string]types.RouteSpec),
		},
		Roles: types.NodeRole{},
		Network: types.Subnet{
			Type:   types.NetworkTypeVxLAN,
			Subnet: "10.0.0.1",
			IFace: types.NetworkInterface{
				Index: 1,
				Name:  "lb",
				Addr:  "10.0.0.1",
				HAddr: "dc:a9:04:83:0d:eb",
			},
		},
		Online: online,
	}

	n.Meta.Name = name
	n.Meta.Description = desc

	return n
}
