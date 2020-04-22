//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func StorageGetAssets(t *testing.T, stg IStorage) {

	var ctx = context.Background()

	type obj struct {
		models.System
		Name string `json:"name"`
	}

	type fields struct {
		stg IStorage
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
			args{ctx: ctx, key: "n", obj: &obj{Name: "demo"}, out: new(obj)},
			nil,
			true,
			errors.ErrEntityNotFound,
		},
		{
			"out struct is nil",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{Name: "demo"}, out: nil},
			nil,
			true,
			errors.ErrStructOutIsNil,
		},
		{
			"test successful get",
			fields{stg},
			args{ctx: ctx, key: "demo", obj: &obj{Name: "demo"}, out: new(obj)},
			&obj{Name: "demo"},
			false,
			"",
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del("demo", tt.args.obj.Name)
		if !assert.NoError(t, err) {
			return
		}

		if tt.args.obj != nil {
			err := tt.fields.stg.Set("demo", tt.args.obj.Name, tt.args.obj)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Get("demo", tt.args.key, tt.args.out)
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

func StorageListAssets(t *testing.T, stg IStorage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg IStorage
	}

	type args struct {
		ctx context.Context
		obj []*obj
		out []*obj
		q   string
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
			"test successful list",
			fields{stg},
			args{ctx: ctx, obj: []*obj{{Name: "demo"}, {Name: "test"}}, out: []*obj{}},
			[]*obj{{Name: "demo"}, {Name: "test"}},
			false,
			"",
		},
	}

	for _, tt := range tests {

		for _, o := range tt.args.obj {
			err := tt.fields.stg.Set("demo", o.Name, o)
			if !assert.NoError(t, err) {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.List("demo", &tt.args.out)

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

			if !assert.Equal(t, len(tt.want), len(tt.args.out), "object received invalid length") {
				return
			}

			for _, w := range tt.want {
				var found bool
				for _, a := range tt.args.out {
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

func StoragePutAssets(t *testing.T, stg IStorage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg IStorage
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

			err := tt.fields.stg.Del("demo", tt.args.obj.Name)
			if !assert.NoError(t, err) {
				return
			}

			if tt.wantErr && tt.err == errors.ErrEntityExists {
				err = tt.fields.stg.Set("demo", tt.args.obj.Name, tt.args.obj)
				if !assert.NoError(t, err) {
					return
				}
			}

			err = tt.fields.stg.Put("demo", tt.args.obj.Name, tt.args.obj)
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

			err = tt.fields.stg.Get("demo", tt.args.key, tt.args.out)

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

func StorageSetAssets(t *testing.T, stg IStorage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	type fields struct {
		stg IStorage
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

			err := tt.fields.stg.Del("demo", "test")
			if !assert.NoError(t, err) {
				return
			}

			var opts = GetOpts()

			if !tt.wantErr && tt.err == errors.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Set("demo", tt.args.obj.Name, tt.args.obj)

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

			err = tt.fields.stg.Get("demo", tt.args.key, tt.args.out)
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

func StorageDelAssets(t *testing.T, stg IStorage) {

	var ctx = context.Background()

	type obj struct {
		Name string `json:"name"`
	}

	type fields struct {
		stg IStorage
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

			err := tt.fields.stg.Del("demo", "test")
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr && tt.err != errors.ErrEntityNotFound {
				err = tt.fields.stg.Put("demo", tt.args.obj.Name, tt.args.obj)
				if !assert.NoError(t, err) {
					return
				}
			}

			var opts = GetOpts()

			if !tt.wantErr && tt.err == errors.ErrEntityNotFound {
				opts.Force = true
			}

			err = tt.fields.stg.Del("demo", tt.args.obj.Name)
			if !assert.NoError(t, err) {
				return
			}

			if !tt.wantErr {

				err := tt.fields.stg.Get("demo", tt.args.key, tt.args.out)
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
