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

package mock

import (
	"context"
	"reflect"
	"testing"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"encoding/json"
)

func TestNodeStorage_List(t *testing.T) {
	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", false)
		nl  = make(map[string]*types.Node, 0)
	)

	nl[n1.Meta.Name] = &n1
	nl[n2.Meta.Name] = &n2

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.Node
		wantErr bool
	}{
		{
			"get node list success",
			fields{stg},
			args{ctx},
			nl,
			false,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.List() storage setup error = %v", err)
			return
		}

		for _, n := range nl {
			if err := stg.Insert(ctx, n); err != nil {
				t.Errorf("NodeStorage.List() storage setup error = %v", err)
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := stg.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStorage_Get(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n   = getNodeAsset("test", "", true)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"get node info failed",
			fields{stg},
			args{ctx, "test2"},
			&n,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get node info successful",
			fields{stg},
			args{ctx, "test"},
			&n,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n); err != nil {
			t.Errorf("NodeStorage.Get() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.name)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNodeStorage_GetSpec(t *testing.T) {

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"
		stg = newNodeStorage()
		ctx = context.Background()
		n1   = getNodeAsset("test1", "", true)
		n2   = getNodeAsset("test1", "", true)
		n3   = getNodeAsset("test2", "", true)
		p1  = getPodAsset(ns, svc, dp, "test1", "")
		p2  = getPodAsset(ns, svc, dp, "test2", "")
	)

	n2.Spec.Pods    = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes  = make(map[string]types.RouteSpec)

	n2.Spec.Pods[p1.SelfLink()] = p1.Spec
	n2.Spec.Pods[p2.SelfLink()] = p2.Spec

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.NodeSpec
		wantErr bool
		err     string
	}{
		{
			"get node spec info failed",
			fields{stg},
			args{ctx, &n3},
			&n3.Spec,
			true,
			store.ErrEntityNotFound,
		},
		{
			"get node spec info successful",
			fields{stg},
			args{ctx, &n1},
			&n2.Spec,
			false,
			"",
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Get() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.GetSpec() storage setup error = %v", err)
			return
		}

		if err := stg.InsertPod(ctx, &n1, &p1); err != nil {
			t.Errorf("NodeStorage.GetSpec() storage setup error = %v", err)
			return
		}

		if err := stg.InsertPod(ctx, &n1, &p2); err != nil {
			t.Errorf("NodeStorage.GetSpec() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)

			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Get() error = %v, wantErr %v", err, tt.err)
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.GetSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("NodeStorage.GetSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.GetSpec() = %v, want %v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_Insert(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test", "", true)
		n2  = getNodeAsset("", "", true)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful insert",
			fields{stg},
			args{ctx, &n1},
			&n1,
			false,
			"",
		},
		{
			"test failed insert",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed insert",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Insert() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Insert() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Insert() error = %v, want %v", err.Error(), tt.err)
				return
			}
		})
	}
}

func TestNodeStorage_Update(t *testing.T) {
	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "desc", true)
		n3  = getNodeAsset("test2", "", false)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Update() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.Update() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Update(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Update() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Update() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.Update() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetStatus(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
	)

	n2.Status.Capacity.Containers++

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetStatus() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.SetStatus() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetStatus() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.SetStatus() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetInfo(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
	)

	n2.Info.Hostname = "demo"

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetInfo() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.SetInfo() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetInfo(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetInfo() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetInfo() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetInfo() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.SetInfo() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetNetwork(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
	)

	n2.Network.IFace.Index++

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetNetwork() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.SetNetwork() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetNetwork(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetNetwork() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetNetwork() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetNetwork() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.SetNetwork() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetOnline(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", false)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetOnline() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.SetOnline() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetOnline(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetOnline() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetOnline() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetOnline() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.SetOnline() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetOffline(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", false)
		n3  = getNodeAsset("test2", "", false)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful update",
			fields{stg},
			args{ctx, &n2},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil structure",
			fields{stg},
			args{ctx, nil},
			&n1,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetOffline() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.SetOffline() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.SetOffline(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetOffline() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetOffline() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetOffline() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStorage.SetOffline() = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_InsertPod(t *testing.T) {

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		p1  = getPodAsset(ns, svc, dp, "test1", "")
		p2  = getPodAsset(ns, svc, dp, "test1", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods[p1.SelfLink()] = p1.Spec
	p2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
		pod  *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful pod insert",
			fields{stg},
			args{ctx, &n1, &p1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &p1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil pod structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: invalid pod structure",
			fields{stg},
			args{ctx, &n1, &p2},
			&n2,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.InsertPod() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.InsertPod() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.InsertPod(tt.args.ctx, tt.args.node, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.InsertPod() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.InsertPod() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.InsertPod() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if !reflect.DeepEqual(*got, tt.want.Spec) {
				t.Errorf("NodeStorage.InsertPod() = %v, want %v", got, tt.want.Spec)
				return
			}

		})
	}
}

func TestNodeStorage_RemovePod(t *testing.T) {

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		p1  = getPodAsset(ns, svc, dp, "test1", "")
		p2  = getPodAsset(ns, svc, dp, "test2", "")
		p3  = getPodAsset(ns, svc, dp, "test2", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	p2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
		pod  *types.Pod
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful pod remove",
			fields{stg},
			args{ctx, &n1, &p1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &p1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil pod structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: node not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: pod arg is invalid",
			fields{stg},
			args{ctx, &n1, &p2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: pod not found",
			fields{stg},
			args{ctx, &n1, &p3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}

		if err := stg.InsertPod(ctx, &n1, &p1); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.RemovePod(tt.args.ctx, tt.args.node, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.RemovePod() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.RemovePod() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.RemovePod() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, _ := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if !reflect.DeepEqual(*got, tt.want.Spec) {
				t.Errorf("NodeStorage.RemovePod() = %v, want %v", got, tt.want.Spec)
				return
			}

		})
	}
}

func TestNodeStorage_InsertVolume(t *testing.T) {

	var (
		ns  = "ns"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		v1  = getVolumeAsset(ns, "test1", "")
		v2  = getVolumeAsset(ns, "test1", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Volumes[v1.SelfLink()] = v1.Spec
	v2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx    context.Context
		node   *types.Node
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful volume insert",
			fields{stg},
			args{ctx, &n1, &v1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &v1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil volume structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: invalid volume structure",
			fields{stg},
			args{ctx, &n1, &v2},
			&n2,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.InsertVolume() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.InsertVolume() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.InsertVolume(tt.args.ctx, tt.args.node, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.InsertVolume() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.InsertVolume() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.InsertVolume() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.InsertVolume() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.InsertVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.InsertVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.InsertVolume() = %v, want %v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_RemoveVolume(t *testing.T) {

	var (
		ns  = "ns"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		v1  = getVolumeAsset(ns, "test1", "")
		v2  = getVolumeAsset(ns, "test2", "")
		v3  = getVolumeAsset(ns, "test2", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	v2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx    context.Context
		node   *types.Node
		volume *types.Volume
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful volume remove",
			fields{stg},
			args{ctx, &n1, &v1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &v1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil volume structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: node not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: volume arg is invalid",
			fields{stg},
			args{ctx, &n1, &v2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: volume not found",
			fields{stg},
			args{ctx, &n1, &v3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.RemoveVolume() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.RemoveVolume() storage setup error = %v", err)
			return
		}

		if err := stg.InsertVolume(ctx, &n1, &v1); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.RemoveVolume(tt.args.ctx, tt.args.node, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.RemoveVolume() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.RemoveVolume() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.RemoveVolume() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.RemoveVolume() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.RemoveVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.RemoveVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.RemoveVolume() = %v, want %v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_InsertRoute(t *testing.T) {

	var (
		ns  = "ns"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		r1  = getRouteAsset(ns, "test1", "")
		r2  = getRouteAsset(ns, "test1", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Routes[r1.SelfLink()] = r1.Spec
	r2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx   context.Context
		node  *types.Node
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful route insert",
			fields{stg},
			args{ctx, &n1, &r1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &r1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil route structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: invalid route structure",
			fields{stg},
			args{ctx, &n1, &r2},
			&n2,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: entity not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.InsertRoute() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.InsertRoute() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.InsertRoute(tt.args.ctx, tt.args.node, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.InsertRoute() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.InsertRoute() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.InsertRoute() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.InsertRoute() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.InsertRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.InsertRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.InsertRoute() = %v, want %v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_RemoveRoute(t *testing.T) {

	var (
		ns  = "ns"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", false)
		r1  = getRouteAsset(ns, "test1", "")
		r2  = getRouteAsset(ns, "test2", "")
		r3  = getRouteAsset(ns, "test2", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)

	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)
	n2.Spec.Routes = make(map[string]types.RouteSpec)

	n3.Spec.Pods = make(map[string]types.PodSpec)

	r2.Meta.Name = ""

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx   context.Context
		node  *types.Node
		route *types.Route
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful route remove",
			fields{stg},
			args{ctx, &n1, &r1},
			&n2,
			false,
			"",
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil, &r1},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: nil route structure",
			fields{stg},
			args{ctx, &n1, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: node not found",
			fields{stg},
			args{ctx, &n3, nil},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: route arg is invalid",
			fields{stg},
			args{ctx, &n1, &r2},
			&n1,
			true,
			store.ErrStructArgIsInvalid,
		},
		{
			"test failed update: route not found",
			fields{stg},
			args{ctx, &n1, &r3},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.RemoveRoute() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.RemoveRoute() storage setup error = %v", err)
			return
		}

		if err := stg.InsertRoute(ctx, &n1, &r1); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.RemoveRoute(tt.args.ctx, tt.args.node, tt.args.route)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.RemoveRoute() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.RemoveRoute() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.RemoveRoute() error = %v, want %v", err.Error(), tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.RemoveRoute() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.RemoveRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.RemoveRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.RemoveRoute() = %v, want %v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_Remove(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", true)
	)

	type fields struct {
		stg storage.Node
	}

	type args struct {
		ctx  context.Context
		node *types.Node
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Node
		wantErr bool
		err     string
	}{
		{
			"test successful node remove",
			fields{stg},
			args{ctx, &n1},
			&n2,
			false,
			store.ErrEntityNotFound,
		},
		{
			"test failed update: nil node structure",
			fields{stg},
			args{ctx, nil},
			&n2,
			true,
			store.ErrStructArgIsNil,
		},
		{
			"test failed update: node not found",
			fields{stg},
			args{ctx, &n2},
			&n1,
			true,
			store.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Remove() storage setup error = %v", err)
			return
		}

		if err := stg.Insert(ctx, &n1); err != nil {
			t.Errorf("NodeStorage.Remove() storage setup error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Remove() error = %v, want %v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Remove() error = %v, want %v", err.Error(), tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("NodeStorage.Remove() = %v, want %v", err, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_Watch(t *testing.T) {

	var (
		stg = newNodeStorage()
		ctx = context.Background()
	)

	type fields struct {
		stg storage.Node
	}
	type args struct {
		ctx  context.Context
		node chan *types.Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check watch",
			fields{stg},
			args{ctx, make(chan *types.Node)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.stg.Watch(tt.args.ctx, tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newNodeStorage(t *testing.T) {
	tests := []struct {
		name string
		want storage.Node
	}{
		{"initialize storage",
			newNodeStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newNodeStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNodeStorage() = %v, want %v", got, tt.want)
			}
		})
	}
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
		Status: types.NodeStatus{
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
