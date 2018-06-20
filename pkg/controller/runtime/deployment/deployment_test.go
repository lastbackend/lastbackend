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

package deployment

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
	"reflect"
	"testing"
)

func TestProvision(t *testing.T) {

	stg, _ := mock.New()
	envs.Get().SetStorage(stg)

	var (
		ns1 = "ns1"
		svc = "svc"
		ctx = context.Background()
		d1  = getDeploymentAsset(ns1, svc, "test1", "")
		d2  = getDeploymentAsset(ns1, svc, "test2", "")
	)

	var ips = make([]string, 0)
	ips = append(ips, "8.8.8.8")

	d2.Spec.Replicas = 2
	d2.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)

	spec := types.SpecTemplateContainer{
		Name: "test-template",
		DNS: types.SpecTemplateContainerDNS{
			Server: ips,
			Search: ips,
		},
	}

	spec.SetDefault()

	d2.Spec.Template.Containers = append(d2.Spec.Template.Containers, spec)

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx        context.Context
		deployment *types.Deployment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Deployment
		wantErr bool
		err     string
	}{
		{
			"provision deployment failed: deployment not found",
			fields{stg},
			args{ctx, &d1},
			&d1,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get deployment info successful",
			fields{stg},
			args{ctx, &d2},
			&d2,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Deployment().Clear(ctx); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Deployment().Insert(ctx, &d2); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			// Run provision method
			err := Provision(tt.args.deployment)

			if err != nil {

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("Provision() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				return
			}

			got, err := tt.fields.stg.Pod().ListByDeployment(tt.args.ctx,
				tt.args.deployment.Meta.Namespace,
				tt.args.deployment.Meta.Service,
				tt.args.deployment.Meta.Name)

			if err != nil {
				t.Errorf("Provision() = %v, want %v", err, tt.err)
				return
			}

			if len(got) != d2.Spec.Replicas {
				t.Errorf("Provision() replicas mismatch: %v, want %v", len(got), d2.Spec.Replicas)
				return
			}

			for _, p := range got {

				if p.Meta.Namespace != tt.args.deployment.Meta.Namespace {
					t.Errorf("Provision() namespace not match = %v, want %v", p.Meta.Namespace, d2.Meta.Namespace)
					return
				}

				if p.Meta.Service != tt.args.deployment.Meta.Service {
					t.Errorf("Provision() service not match = %v, want %v", p.Meta.Service, d2.Meta.Service)
					return
				}

				if p.Meta.Deployment != tt.args.deployment.Meta.Name {
					t.Errorf("Provision() name not match = %v, want %v", p.Meta.Service, d2.Meta.Service)
					return
				}

				if !reflect.DeepEqual(d2.Spec.State, p.Spec.State) {
					t.Errorf("Provision() state not match = %v, want %v", p.Spec.State, d2.Spec.State)
					return
				}

				for i, t := range p.Spec.Template.Containers {
					t.Labels = make(map[string]string)
					p.Spec.Template.Containers[i] = t
				}

				g, err := json.Marshal(p.Spec.Template)
				if err != nil {
					t.Errorf("Provision() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				e, err := json.Marshal(d2.Spec.Template)
				if err != nil {
					t.Errorf("Provision() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(g, e) {
					t.Errorf("Provision() spec not match = %v, want %v", string(g), string(e))
				}
			}
		})
	}
}

func getDeploymentAsset(namespace, service, name, desc string) types.Deployment {

	var n = types.Deployment{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Service = service
	n.Meta.Description = desc

	n.SelfLink()

	return n
}
