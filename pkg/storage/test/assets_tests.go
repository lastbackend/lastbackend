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

package test_test

import (
	"context"
	"testing"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/stretchr/testify/assert"
)

func StorageGetAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		key string
		obj *obj
		out *obj
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *obj
		wantErr bool
		err     string
	}{
		{
			"test not found err",
			fields{stg},
			args{ctx: ctx, key: "n", obj: &obj{"demo"}, out: new(obj)},
			nil,
			true,
			types.ErrEntityNotFound,
		},
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: nil},
			nil,
			true,
			types.ErrStructOutIsNil,
		},
		{
			"test successful get",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
		if !assert.NoError(t, err) {
			return
		}

		if tt.args.obj != nil {
			err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, tt.args.obj.Name, tt.args.obj, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Get(tt.args.ctx, storage.TestKind, tt.args.key, tt.args.out)

			if tt.wantErr {
				if !assert.Error(t, err, "expected err") {
					return
				}
				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			if !assert.NotNil(t, tt.args.out, "expected pointer") {
				return
			}

			assert.Equal(t, tt.want.Name, tt.args.out.Name, "object received error")
		})
	}

}

func StorageListAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		obj []*obj
		out *[]*obj
		q   string
	}

	outf := func() *[]*obj {
		out := make([]*obj, 0)
		return &out
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*obj
		wantErr bool
		err     string
	}{
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: nil},
			nil,
			true,
			types.ErrStructOutIsNil,
		},
		{
			"test successful list with filter",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: outf(), q: "demo"},
			[]*obj{{"demo"}},
			false,
			"",
		},
		{
			"test successful list",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: outf()},
			[]*obj{{"demo"}, {"test"}},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
		if !assert.NoError(t, err) {
			return
		}

		for _, o := range tt.args.obj {
			err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, o.Name, o, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.List(tt.args.ctx, storage.TestKind, tt.args.q, tt.args.out)

			if tt.wantErr {
				if !assert.Error(t, err, "expected err") {
					return
				}
				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			if !assert.NotNil(t, tt.args.out, "expected pointer") {
				return
			}

			if !assert.Equal(t, len(tt.want), len(*tt.args.out), "object received invalid length") {
				return
			}

			for _, w := range tt.want {
				var found bool
				for _, a := range *tt.args.out {
					if w.Name == a.Name {
						found = true
					}
				}

				if !assert.True(t, found, "can not found expected value in actual response") {
					return
				}
			}

		})
	}

}

func StorageMapAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		obj []*obj
		out *map[string]*obj
		q   string
	}

	outf := func() *map[string]*obj {
		out := make(map[string]*obj, 0)
		return &out
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*obj
		wantErr bool
		err     string
	}{
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: nil},
			nil,
			true,
			types.ErrStructOutIsNil,
		},
		{
			"test successful list with filter",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: outf(), q: "demo"},
			map[string]*obj{"demo": {"demo"}},
			false,
			"",
		},
		{
			"test successful map",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{"demo"}, {"test"}}, out: outf()},
			map[string]*obj{"demo": {"demo"}, "test": {"test"}},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
		if !assert.NoError(t, err) {
			return
		}

		for _, o := range tt.args.obj {
			err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, o.Name, o, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Map(tt.args.ctx, storage.TestKind, tt.args.q, tt.args.out)

			if tt.wantErr {
				if !assert.Error(t, err, "expected err") {
					return
				}
				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			if !assert.NotNil(t, tt.args.out, "expected pointer") {
				return
			}

			jse, err := json.Marshal(tt.want)
			if !assert.NoError(t, err) {
				return
			}

			jsa, err := json.Marshal(tt.args.out)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, string(jse), string(jsa), "object received error")
		})
	}

}

func StoragePutAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		key string
		obj *obj
		out *obj
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *obj
		wantErr bool
		err     string
	}{
		{
			"test put err entity exists",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			true,
			types.ErrEntityExists,
		},
		{
			"test successful put",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			false,
			"",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
			if !assert.NoError(t, err) {
				return
			}

			if tt.wantErr && tt.err == types.ErrEntityExists {
				err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, tt.args.obj.Name, tt.args.obj, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, tt.args.obj.Name, tt.args.obj, nil)
			if tt.wantErr {
				if !assert.Error(t, err, "expected err") {
					return
				}
				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			err = tt.fields.stg.Get(tt.args.ctx, storage.TestKind, tt.args.key, tt.args.out)

			if !assert.NoError(t, err) {
				return
			}

			if !assert.NotNil(t, tt.args.out, "expected pointer") {
				return
			}

			assert.Equal(t, tt.args.obj.Name, tt.args.out.Name, "object received error")
		})
	}

}

func StorageSetAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		key string
		obj *obj
		out *obj
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *obj
		wantErr bool
		err     string
	}{
		{
			"test set err entity not found",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo", "test"}, out: new(obj)},
			&obj{"demo", "test"},
			true,
			types.ErrEntityNotFound,
		},
		{
			"test successful set when entity not exists",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo", "test"}, out: new(obj)},
			&obj{"demo", "test"},
			false,
			types.ErrEntityNotFound,
		},
		{
			"test successful set",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo", "test"}, out: new(obj)},
			&obj{"demo", "test"},
			false,
			"",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
			if !assert.NoError(t, err) {
				return
			}

			log.Info(tt.err)
			if tt.err != types.ErrEntityNotFound {
				err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, tt.args.obj.Name, &obj{"demo", "demo"}, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			var opts = storage.GetOpts()

			if !tt.wantErr && tt.err == types.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Set(tt.args.ctx, storage.TestKind, tt.args.obj.Name, tt.args.obj, opts)

			if tt.wantErr {
				if !assert.Error(t, err, "expected err") {
					return
				}
				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			err = tt.fields.stg.Get(tt.args.ctx, storage.TestKind, tt.args.key, tt.args.out)

			if !assert.NoError(t, err) {
				return
			}

			if !assert.NotNil(t, tt.args.out, "expected pointer") {
				return
			}

			assert.Equal(t, tt.want.Desc, tt.args.out.Desc, "object received error")
		})
	}

}

func StorageDelAssets(t *testing.T, stg storage.Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
		key string
		obj *obj
		out *obj
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *obj
		wantErr bool
		err     string
	}{
		{
			"test del err entity not found",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			true,
			types.ErrEntityNotFound,
		},
		{
			"test successful del",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			false,
			types.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Del(tt.args.ctx, storage.TestKind, "")
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr && tt.err != types.ErrEntityNotFound {
				err = tt.fields.stg.Put(tt.args.ctx, storage.TestKind, tt.args.obj.Name, tt.args.obj, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			var opts = storage.GetOpts()

			if !tt.wantErr && tt.err == types.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Del(tt.args.ctx, storage.TestKind, tt.args.obj.Name)
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr {

				err := tt.fields.stg.Get(tt.args.ctx, storage.TestKind, tt.args.key, tt.args.out)
				if !assert.Error(t, err, "expected err") {
					return
				}

				assert.Equal(t, tt.err, err.Error(), "err message different")
				return
			}

			if assert.NoError(t, err) {
				return
			}

		})
	}

}
