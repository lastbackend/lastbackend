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

package resource

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCpuResource(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int64
		wantErr bool
		err     string
	}{
		{
			name: "parse int",
			args: "1",
			want: 1000000000,
		},
		{
			name: "parse float",
			args: "0.5",
			want: 500000000,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			got, err := DecodeCpuResource(tc.args)

			if tc.wantErr {
				if !assert.Error(t, err, "error should be not nil") {
					return
				}
				return
			}

			if !assert.NoError(t, err, "error should be nil") {
				return
			}

			assert.Equal(t, tc.want, got, "parsed values mismatched")
		})
	}
}

func TestParseMemoryResource(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		want    int64
		wantErr bool
		err     string
	}{
		{
			name: "parse 50B",
			args: "50B",
			want: 50,
		},
		{
			name: "parse 1mb",
			args: "1mb",
			want: 1000 * 1000,
		},
		{
			name: "parse 1mib",
			args: "1mib",
			want: 1024 * 1024,
		},
		{
			name: "parse 10gb",
			args: "10gb",
			want: 10 * 1000 * 1000 * 1000,
		},
		{
			name: "parse 10gib",
			args: "10gib",
			want: 10 * 1024 * 1024 * 1024,
		},
		{
			name: "parse 0.5Gi",
			args: "0.5Gi",
			want: 0.5 * 1024 * 1024 * 1024,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			got, err := ToBytes(tc.args)

			if tc.wantErr {
				if !assert.Error(t, err, "error should be not nil") {
					return
				}
				return
			}

			if !assert.NoError(t, err, "error should be nil") {
				return
			}

			assert.Equal(t, tc.want, got, "parsed values mismatched")
		})
	}

}

func TestEncodeMemoryResource(t *testing.T) {

	tests := []struct {
		name    string
		args    int64
		want    string
		wantErr bool
		err     string
	}{
		{
			name: "encode int",
			args: 1024,
			want: "1KiB",
		},
		{
			name: "encode 1mb",
			args: 1024 * 1024,
			want: "1MiB",
		},
		{
			name: "encode 1gib",
			args: 1024 * 1024 * 1024,
			want: "1GiB",
		},
		{
			name: "encode 1gb",
			args: 1000 * 1000 * 1000,
			want: "953.7MiB",
		},
		{
			name: "encode 12gb",
			args: 12 * 1000 * 1000 * 1000,
			want: "11.18GiB",
		},
		{
			name: "encode 120GiB",
			args: 12 * 1024 * 1024 * 1024,
			want: "12GiB",
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			got := EncodeMemoryResource(tc.args)

			assert.Equal(t, tc.want, got, "parsed values mismatched")
		})
	}

}
