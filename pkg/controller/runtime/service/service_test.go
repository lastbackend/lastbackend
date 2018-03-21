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

package service

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"reflect"
	"testing"
)

func TestProvision(t *testing.T) {

	stg, _ := mock.New()
	envs.Get().SetStorage(stg)

	var (
		ns1 = "ns1"
		ctx = context.Background()
		s1  = getServiceAsset(ns1, "test1", "")
		s2  = getServiceAsset(ns1, "test2", "")
	)

	s2.Spec.Replicas = 2
	s2.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)

	spec := types.SpecTemplateContainer{
		Name: "test-template",
	}

	spec.SetDefault()

	s2.Spec.Template.Containers = append(s2.Spec.Template.Containers, spec)

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx     context.Context
		service *types.Service
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Service
		wantErr bool
		err     string
	}{
		{
			"provision service failed: service not found",
			fields{stg},
			args{ctx, &s1},
			&s1,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get service info successful",
			fields{stg},
			args{ctx, &s2},
			&s2,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Service().Clear(ctx); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		if err := stg.Service().Insert(ctx, &s2); err != nil {
			t.Errorf("Provision() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			// Run provision method
			err := Provision(tt.args.service)

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

			got, err := tt.fields.stg.Deployment().ListByService(tt.args.ctx,
				tt.args.service.Meta.Namespace,
				tt.args.service.Meta.Name)

			if err != nil {
				t.Errorf("Provision() = %v, want %v", err, tt.err)
				return
			}

			for _, p := range got {

				if p.Meta.Namespace != tt.args.service.Meta.Namespace {
					t.Errorf("Provision() namespace not match = %v, want %v", p.Meta.Namespace, s2.Meta.Namespace)
					return
				}

				if p.Meta.Service != tt.args.service.Meta.Name {
					t.Errorf("Provision() name not match = %v, want %v", p.Meta.Service, s2.Meta.Name)
					return
				}

				if !reflect.DeepEqual(s2.Spec.State, p.Spec.State) {
					t.Errorf("Provision() state not match = %v, want %v", p.Spec.State, s2.Spec.State)
					return
				}

				if p.Spec.Replicas != tt.args.service.Spec.Replicas {
					t.Errorf("Provision() replicas not match = %v, want %v", p.Spec.State, s2.Spec.State)
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

				e, err := json.Marshal(tt.args.service.Spec.Template)
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

func getServiceAsset(namespace, name, desc string) types.Service {

	var n = types.Service{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}
