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

package service_test

import (
	"testing"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/stretchr/testify/assert"
)

// TestServiceProvision - check that service should create new deployment if no deployments exist
func TestServiceProvision(t *testing.T) {

	var (
		ctx = context.Background()
	)

	type fields struct {

	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{
			"test not found err",
			fields{},
			args{ctx: ctx, },
			true,
			types.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {

		err := tt.fields.stg.Del(tt.args.ctx, TestKind, "")
		if !assert.NoError(t, err) {
			return
		}

		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.stg.Get(tt.args.ctx, TestKind, tt.args.key, tt.args.out)

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

func TestServiceDestroy(t *testing.T) {

}