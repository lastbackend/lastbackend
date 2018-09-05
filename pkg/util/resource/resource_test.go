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

package resource


import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseResource(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		want    int64
		wantErr bool
		err     string
	}{
		{
			name: "parse int",
			args: "1024",
			want: 1024*1024,
		},
		{
			name: "parse 1mb",
			args: "1mb",
			want: 1000*1000,
		},
		{
			name: "parse 1mib",
			args: "1mib",
			want: 1024*1024,
		},
		{
			name: "parse 1gb",
			args: "1gb",
			want: 1000*1000*1000,
		},
		{
			name: "parse 1gib",
			args: "1gib",
			want: 1024*1024*1024,
		},
		{
			name: "parse 12gb",
			args: "12gb",
			want: 12*1000*1000*1000,
		},
		{
			name: "parse 12000mib",
			args: "12000mib",
			want: 12*1000*1024*1024,
		},
		{
			name: "parse mib",
			args: "mib",
			want: 1024*1024,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			got, err := DecodeResource(tc.args)

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