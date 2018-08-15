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

package storage

import (
	"context"
	"testing"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/stretchr/testify/assert"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func StorageGetAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		types.Runtime
		Name string `json:"name"`
	}

	type fields struct {
		stg Storage
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
			args{ctx: ctx, key: "n", obj: &obj{Name:"demo"}, out: new(obj)},
			nil,
			true,
			errors.ErrEntityNotFound,
		},
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{Name:"demo"}, out: nil},
			nil,
			true,
			errors.ErrStructOutIsNil,
		},
		{
			"test successful get",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{Name:"demo"}, out: new(obj)},
			&obj{Name:"demo"},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
		if !assert.NoError(t, err) {
			return
		}

		if tt.args.obj != nil {
			err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, tt.args.obj, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Get(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.key, tt.args.out, nil)

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

func StorageListAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		types.Runtime
		Name string `json:"name"`
	}

	type objl struct {
		types.Runtime
		Items []*obj
	}

	type fields struct {
		stg Storage
	}

	type args struct {
		ctx context.Context
		obj objl
		out *objl
		q   string
	}

	outf := func() *objl {
		nl := objl{}
		return &nl
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *objl
		wantErr bool
		err     string
	}{
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, obj: objl{Items:[]*obj{{Name:"demo"}, {Name:"test"}}}, out: nil},
			nil,
			true,
			errors.ErrStructOutIsNil,
		},
		{
			"test successful list with filter",
			fields{stg},
			args{ctx: ctx, obj: objl{Items:[]*obj{{Name:"demo"}, {Name:"test"}}}, out: outf(), q: "demo"},
			&objl{Items:[]*obj{{Name:"demo"}}},
			false,
			"",
		},
		{
			"test successful list",
			fields{stg},
			args{ctx: ctx, obj: objl{Items:[]*obj{{Name:"demo"}, {Name:"test"}}}, out: outf()},
			&objl{Items:[]*obj{{Name:"demo"}, {Name:"test"}}},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
		if !assert.NoError(t, err) {
			return
		}

		for _, o := range tt.args.obj.Items {
			err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), o.Name, o, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.List(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.q, tt.args.out, nil)

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

			if !assert.Equal(t, len(tt.want.Items), len(tt.args.out.Items), "object received invalid length") {
				return
			}

			for _, w := range tt.want.Items {
				var found bool
				for _, a := range tt.args.out.Items {
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

func StorageMapAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		types.Runtime
		Name string `json:"name"`
	}

	type objm struct {
		types.Runtime
		Items map[string]*obj
	}

	type fields struct {
		stg Storage
	}

	type args struct {
		ctx context.Context
		obj []*obj
		out *objm
		q   string
	}

	outf := func() *objm {
		out := objm{}
		out.Items = make(map[string]*obj, 0)
		return &out
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *objm
		wantErr bool
		err     string
	}{
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{Name:"demo"}, {Name:"test"}}, out: nil},
			nil,
			true,
			errors.ErrStructOutIsNil,
		},
		{
			"test successful list with filter",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{Name:"demo"}, {Name:"test"}}, out: outf(), q: "demo"},
			&objm{Items:map[string]*obj{"demo": {Name:"demo"}}},
			false,
			"",
		},
		{
			"test successful map",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{Name:"demo"}, {Name:"test"}}, out: outf()},
			&objm{Items:map[string]*obj{"demo": {Name:"demo"}, "test": {Name:"test"}}},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
		if !assert.NoError(t, err) {
			return
		}

		for _, o := range tt.args.obj {
			err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), o.Name, o, nil)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Map(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.q, tt.args.out, nil)

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

			jse, err := json.Marshal(tt.want.Items)
			if !assert.NoError(t, err) {
				return
			}

			jsa, err := json.Marshal(tt.args.out.Items)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, string(jse), string(jsa), "object received error")
		})
	}

}

func StoragePutAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg Storage
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
			errors.ErrEntityExists,
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

			err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
			if !assert.NoError(t, err) {
				return
			}

			if tt.wantErr && tt.err == errors.ErrEntityExists {
				err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, tt.args.obj, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, tt.args.obj, nil)
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

			err = tt.fields.stg.Get(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.key, tt.args.out, nil)

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

func StorageSetAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	type fields struct {
		stg Storage
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
			errors.ErrEntityNotFound,
		},
		{
			"test successful set when entity not exists",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo", "test"}, out: new(obj)},
			&obj{"demo", "test"},
			false,
			errors.ErrEntityNotFound,
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

			err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
			if !assert.NoError(t, err) {
				return
			}

			log.Info(tt.err)
			if tt.err != errors.ErrEntityNotFound {
				err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, &obj{"demo", "demo"}, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			var opts = GetOpts()

			if !tt.wantErr && tt.err == errors.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Set(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, tt.args.obj, opts)

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

			err = tt.fields.stg.Get(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.key, tt.args.out, nil)

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

func StorageDelAssets(t *testing.T, stg Storage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg Storage
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
			errors.ErrEntityNotFound,
		},
		{
			"test successful del",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{"demo"}, out: new(obj)},
			&obj{"demo"},
			false,
			errors.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), "")
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr && tt.err != errors.ErrEntityNotFound {
				err = tt.fields.stg.Put(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name, tt.args.obj, nil)
				if !assert.NoError(t, err) {
					return
				}
			}

			var opts = GetOpts()

			if !tt.wantErr && tt.err == errors.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Del(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.obj.Name)
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr {

				err := tt.fields.stg.Get(tt.args.ctx, tt.fields.stg.Collection().Test(), tt.args.key, tt.args.out, nil)
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
