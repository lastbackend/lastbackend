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

package cache

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestGet(t *testing.T) {

	a := assert.New(t)

	cache := NewCache(time.Second)

	toStrPtr := func(data string) *string { return &data }

	type dataset []struct {
		key string
		val string
	}

	type args struct {
		key     string
		waiting time.Duration
	}

	tests := []struct {
		name    string
		dataset dataset
		args    args
		err     string
		want    *string
	}{
		{
			name:    "checking that is exists",
			dataset: dataset{{"test1", "data1"}, {"test2", "data2"}},
			args:    args{"test1", time.Millisecond},
			err:     "item not exists",
			want:    toStrPtr("data1"),
		},
		{
			name:    "checking that is not exists",
			dataset: dataset{{"test1", "data1"}, {"test2", "data2"}},
			args:    args{"test", time.Millisecond},
			err:     "item exists",
			want:    nil,
		},
		{
			name:    "checking that is not exists",
			dataset: dataset{{"test", "data"}},
			args:    args{"test", 2 * time.Second},
			err:     "item not expired",
			want:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer cache.Clear()

			for _, item := range tc.dataset {
				cache.Set(item.key, item.val)
			}

			<-time.After(tc.args.waiting)

			if tc.want == nil {
				if !a.Nil(cache.Get(tc.args.key)) {
					t.Errorf(tc.err)
				}
			}

			if tc.want != nil {
				if !a.Equal(cache.Get(tc.args.key).(string), *tc.want) {
					t.Errorf(tc.err)
				}
			}

		})
	}

}
