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
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

func TestNodeStorage_List(t *testing.T) {

	initStorage()

	var (
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", true)
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.List() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			for _, n := range nl {
				if err := stg.Insert(ctx, n); err != nil {
					t.Errorf("NodeStorage.List() storage setup error = %v", err)
					return
				}
			}

			got, err := stg.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("NodeStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.List() \n%v\nwant\n%v", string(g), string(e))
			}
		})
	}
}

func TestNodeStorage_Get(t *testing.T) {

	initStorage()

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
		{
			"get node info failed empty name",
			fields{stg},
			args{ctx, ""},
			&n,
			true,
			"node can not be empty",
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("NodeStorage.Get() storage setup error = %v", err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.name)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Get() = %v, want %v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("NodeStorage.Get() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Get() error = %v, want none", tt.err)
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("NodeStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.Get() \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_GetSpec(t *testing.T) {

	initStorage()

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"
		stg = newNodeStorage()
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test1", "", true)
		n3  = getNodeAsset("test2", "", true)
		n4  = getNodeAsset("test3", "", true)
		p1  = getPodAsset(ns, svc, dp, "test1", "")
		p2  = getPodAsset(ns, svc, dp, "test2", "")
	)

	n1.Network.Range = "10.0.1.0"
	n2.Network.Range = "10.0.1.0"
	n3.Network.Range = "10.0.2.0"

	n2.Spec.Network = make(map[string]types.NetworkSpec)
	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)

	n2.Spec.Pods[p1.SelfLink()] = p1.Spec
	n2.Spec.Pods[p2.SelfLink()] = p2.Spec

	n2.Spec.Network[n3.Meta.Name] = n3.Network

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
			args{ctx, &n4},
			&n4.Spec,
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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Get() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.GetSpec() storage setup error = %v", err)
				return
			}
			if err := stg.Insert(ctx, &n3); err != nil {
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

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			//t.Logf("got get spec=%v\n from node = %v", got, tt.args.node)
			if err != nil {
				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.GetSpec() \n%v\nwant\n%v", err, tt.err)
					return
				}
				if !tt.wantErr {
					t.Errorf("NodeStorage.GetSpec() error = %v, want no error", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.GetSpec() error = %v, want none", tt.err)
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
				t.Errorf("NodeStorage.GetSpec() = \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_Insert(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Insert() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tt.fields.stg.Insert(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Insert() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Insert() error \n%v\nwant\n%v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Insert() error = %v, want none", tt.err)
				return
			}
		})
	}
}

func TestNodeStorage_Update(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Update() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.Update() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Update(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Update() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Update() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Update() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.Update() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.Update() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetStatus(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetStatus() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.SetStatus() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.SetStatus(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetStatus() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetStatus() error \n%v\nwant\n%v", err.Error(), tt.err)
					return
				}

				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetStatus() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.SetStatus() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.SetStatus() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetInfo(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetInfo() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.SetInfo() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.SetInfo(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetInfo() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetInfo() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetInfo() error = %v,  want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.SetInfo() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.SetInfo() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetNetwork(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetNetwork() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.SetNetwork() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.SetNetwork(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetNetwork() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetNetwork() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return

			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetNetwork() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.SetNetwork() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.SetNetwork() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetOnline(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetOnline() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.SetOnline() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.SetOnline(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetOnline() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetOnline() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetOnline() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.SetOnline() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.SetOnline() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_SetOffline(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.SetOffline() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.SetOffline() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.SetOffline(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.SetOffline() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.SetOffline() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.SetOffline() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err != nil {
				t.Errorf("NodeStorage.SetOffline() got Get error = %s", err.Error())
				return
			}
			if !compareNodes(got, tt.want) {
				t.Errorf("NodeStorage.SetOffline() \n%v\nwant\n%v", got, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_InsertPod(t *testing.T) {

	initStorage()

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

	n2.Spec.Network = make(map[string]types.NetworkSpec)
	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)

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

	clear := func() {

		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.InsertPod() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.InsertPod() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.InsertPod(tt.args.ctx, tt.args.node, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.InsertPod() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.InsertPod() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.InsertPod() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.InsertPod() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.InsertPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.InsertPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.InsertPod() \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_RemovePod(t *testing.T) {

	initStorage()

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

	n2.Spec.Network = make(map[string]types.NetworkSpec)
	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
				return
			}

			if err := stg.InsertPod(ctx, &n1, &p1); err != nil {
				t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.RemovePod(tt.args.ctx, tt.args.node, tt.args.pod)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.RemovePod() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.RemovePod() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.RemovePod() error = %v, want none", tt.err)
				return
			}

			got, err := tt.fields.stg.GetSpec(tt.args.ctx, tt.args.node)
			if err != nil {
				t.Errorf("NodeStorage.RemovePod() error = %v, want no error", err.Error())
				return
			}

			g, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NodeStorage.RemovePod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			e, err := json.Marshal(tt.want.Spec)
			if err != nil {
				t.Errorf("NodeStorage.RemovePod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(g, e) {
				t.Errorf("NodeStorage.RemovePod() \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_InsertVolume(t *testing.T) {

	initStorage()

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

	n2.Spec.Network = make(map[string]types.NetworkSpec)
	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.InsertVolume() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.InsertVolume() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.InsertVolume(tt.args.ctx, tt.args.node, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.InsertVolume() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.InsertVolume() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.InsertVolume() error = %v, want none", tt.err)
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
				t.Errorf("NodeStorage.InsertVolume() \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_RemoveVolume(t *testing.T) {

	initStorage()

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

	n2.Spec.Network = make(map[string]types.NetworkSpec)
	n2.Spec.Pods = make(map[string]types.PodSpec)
	n2.Spec.Volumes = make(map[string]types.VolumeSpec)

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.RemoveVolume() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.RemoveVolume() storage setup error = %v", err)
				return
			}

			if err := stg.InsertVolume(ctx, &n1, &v1); err != nil {
				t.Errorf("NodeStorage.RemovePod() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.RemoveVolume(tt.args.ctx, tt.args.node, tt.args.volume)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.RemoveVolume() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.RemoveVolume() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.RemoveVolume() error = %v, want none", tt.err)
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
				t.Errorf("NodeStorage.RemoveVolume() \n%v\nwant\n%v", string(g), string(e))
			}

		})
	}
}

func TestNodeStorage_Remove(t *testing.T) {

	initStorage()

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

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Remove() storage setup error = %v", err)
			return
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n1); err != nil {
				t.Errorf("NodeStorage.Remove() storage setup error = %v", err)
				return
			}

			err := tt.fields.stg.Remove(tt.args.ctx, tt.args.node)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NodeStorage.Remove() error = %v, want no error", err.Error())
					return
				}

				if tt.wantErr && tt.err != err.Error() {
					t.Errorf("NodeStorage.Remove() error \n%v\nwant\n%v", err.Error(), tt.err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("NodeStorage.Remove() error = %v, want none", tt.err)
				return
			}

			_, err = tt.fields.stg.Get(tt.args.ctx, tt.args.node.Meta.Name)
			if err == nil || tt.err != err.Error() {
				t.Errorf("NodeStorage.Remove() \n%v\nwant\n%v", err, tt.want)
				return
			}

		})
	}
}

func TestNodeStorage_Watch(t *testing.T) {

	var (
		stg   = newNodeStorage()
		ctx   = context.Background()
		err   error
		n     = getNodeAsset("test1", "desc1", true)
		nodeC = make(chan *types.Node)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("NodeStorage.Watch() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check node watch",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.Watch() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("NodeStorage.Watch() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()

			//run watch go function
			go func() {
				err = stg.Watch(ctxT, nodeC)
				if err != nil {
					t.Errorf("NodeStorage.Watch() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/node/test1/online"
			value := "true"
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-nodeC:
					t.Log("NodeStorage.Watch() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("NodeStorage.Watch() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestNodeStorage_WatchStatus(t *testing.T) {

	var (
		stg              = newNodeStorage()
		ctx              = context.Background()
		err              error
		n                = getNodeAsset("test1", "desc1", true)
		nodeStatusEventC = make(chan *types.NodeStatusEvent)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("NodeStorage.EventsStatus() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check node watch status",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.EventsStatus() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("NodeStorage.EventsStatus() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.EventStatus(ctxT, nodeStatusEventC)
				if err != nil {
					t.Errorf("NodeStorage.EventsStatus() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/node/test1/online"
			value := "true"
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-nodeStatusEventC:
					t.Log("NodeStorage.EventsStatus() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("NodeStorage.EventsStatus() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestNodeStorage_EventNetworkSpec(t *testing.T) {

	var (
		stg                   = newNodeStorage()
		ctx                   = context.Background()
		err                   error
		n                     = getNodeAsset("test1", "desc1", true)
		nodeNetworkSpecEventC = make(chan *types.NetworkSpecEvent)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("NodeStorage.EventNetworkSpec() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check node watch network spec event",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.EventNetworkSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			if err := stg.Insert(ctx, &n); err != nil {
				t.Errorf("NodeStorage.EventNetworkSpec() storage setup error = %v", err)
				return
			}

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.EventNetworkSpec(ctxT, nodeNetworkSpecEventC)
				if err != nil {
					t.Errorf("NodeStorage.EventNetworkSpec() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/node/test1/network"
			value := `{"type":"vxlan","range":"10.0.0.1","iface":{"index":1,"name":"lb","addr":"10.0.0.1","HAddr":"dc:a9:04:83:0d:eb"},"addr":""}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-nodeNetworkSpecEventC:
					t.Log("NodeStorage.EventNetworkSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("NodeStorage.EventNetworkSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestNodeStorage_EventPodSpec(t *testing.T) {

	var (
		stg               = newNodeStorage()
		ctx               = context.Background()
		err               error
		n                 = getNodeAsset("test1", "desc1", true)
		p                 = getPodAsset("ns1", "svc", "dp1", "test1", "desc")
		nodePodSpecEventC = make(chan *types.PodSpecEvent)
	)
	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("NodeStorage.EventPodSpec() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check node watch pod spec event",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.EventPodSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err = stg.Insert(ctx, &n)
			err = stg.InsertPod(ctx, &n, &p)

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose()
			//run watch go function
			go func() {
				err = stg.EventPodSpec(ctxT, nodePodSpecEventC)
				if err != nil {
					t.Errorf("NodeStorage.EventPodSpec() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/node/test1/spec/pods/ns1:svc:dp1:test1"
			value := `{"state":{"destroy":false,"maintenance":false},"template":{"volumes":null,"container":null,"termination":0}}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-nodePodSpecEventC:
					t.Log("NodeStorage.EventPodSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("NodeStorage.EventPodSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
		})
	}
}

func TestNodeStorage_EventVolumeSpec(t *testing.T) {

	var (
		stg                  = newNodeStorage()
		ctx                  = context.Background()
		err                  error
		n                    = getNodeAsset("test1", "desc1", true)
		v                    = getVolumeAsset("ns1", "test1", "desc")
		nodeVolumeSpecEventC = make(chan *types.VolumeSpecEvent)
	)

	etcdCtl, destroy, err := initStorageWatch()
	if err != nil {
		t.Errorf("NodeStorage.EventVolumeSpec() storage setup error = %v", err)
	}
	defer destroy()

	type fields struct {
	}
	type args struct {
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"check node watch volume spec event",
			fields{},
			args{},
			false,
		},
	}

	clear := func() {
		if err := stg.Clear(ctx); err != nil {
			t.Errorf("NodeStorage.EventVolumeSpec() storage setup error = %v", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clear()
			defer clear()

			err = stg.Insert(ctx, &n)
			err = stg.InsertVolume(ctx, &n, &v)

			//create timeout context
			ctxT, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			defer etcdCtl.WatchClose() //run watch go function

			go func() {
				err = stg.EventVolumeSpec(ctxT, nodeVolumeSpecEventC)
				if err != nil {
					t.Errorf("NodeStorage.EventVolumeSpec() storage setup error = %v", err)
					return
				}
			}()
			//wait for result
			time.Sleep(1 * time.Second)

			//make etcd key put through etcdctrl
			path := getEtcdctrl()
			if path == "" {
				t.Skipf("skip watch test: not found etcdctl path=%s", path)
			}
			key := "/lstbknd/node/test1/spec/volumes/ns1:test1"
			value := `{"state":{"destroy":false}}`
			err = runEtcdPut(path, key, value)
			if err != nil {
				t.Skipf("skip watch test: exec etcdctl err=%s", err.Error())
			}

			for {
				select {
				case <-nodeVolumeSpecEventC:
					t.Log("NodeStorage.EventVolumeSpec() is working")
					return
				case <-ctxT.Done():
					t.Log("ctxT done=", ctxT.Err(), "time=", time.Now())
					t.Error("NodeStorage.EventVolumeSpec() NO watch event happen")
					return
				case <-time.After(500 * time.Millisecond):
					//wait for 500 ms
				}
			}
			t.Log("successfull!")
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
				t.Errorf("newNodeStorage() \n%v\nwant\n%v", got, tt.want)
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
			Hostname:   name,
			InternalIP: "0.0.0.0",
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
		Spec:  types.NodeSpec{},
		Roles: types.NodeRole{},
		Network: types.NetworkSpec{
			Type:  types.NetworkTypeVxLAN,
			Range: "10.0.0.1",
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
	n.Meta.Created = time.Now()

	return n
}

//compare two node structures
func compareNodes(got, want *types.Node) bool {
	result := false
	if compareMeta(got.Meta.Meta, want.Meta.Meta) &&
		(got.Online == want.Online) &&
		(got.Meta.Cluster == want.Meta.Cluster) &&
		(got.Meta.Token == want.Meta.Token) &&
		(got.Meta.Region == want.Meta.Region) &&
		(got.Meta.Provider == want.Meta.Provider) &&
		reflect.DeepEqual(got.Status, want.Status) &&
		reflect.DeepEqual(got.Info, want.Info) &&
		reflect.DeepEqual(got.Roles, want.Roles) &&
		reflect.DeepEqual(got.Network, want.Network) &&
		reflect.DeepEqual(got.Spec, want.Spec) {
		result = true
	}

	return result
}

func compareNodeMaps(got, want map[string]*types.Node) bool {
	for k, v := range got {
		if !compareNodes(v, want[k]) {
			return false
		}
	}
	return true
}
