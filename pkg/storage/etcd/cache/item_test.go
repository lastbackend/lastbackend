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
	"time"
	"github.com/stretchr/testify/assert"
)

func TestExpired(t *testing.T) {

	a := assert.New(t)

	getItem := func(name string, duration time.Duration) *Item {
		item := &Item{data: name}
		expiration := time.Now().Add(time.Second)
		item.expires = &expiration
		return item
	}

	type datasets struct {
		data *Item
	}

	tests := []struct {
		name string
		data datasets
		err  string
		want bool
	}{
		{
			name: "checking expired by default",
			data: datasets{&Item{data: "test"}},
			err:  "item to be expired by default",
			want: true,
		},
		{
			name: "checking that is not expired",
			data: datasets{getItem("test", time.Second)},
			err:  "item to not be expired",
			want: false,
		},
		{
			name: "checking that is not expired",
			data: datasets{getItem("test", 0-time.Second)},
			err:  "item to be expired",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !a.Equal(tc.data.data.expired(), tc.want) {
				t.Errorf(tc.err)
			}
		})
	}

}

func TestSetExpireTime(t *testing.T) {

	a := assert.New(t)

	type datasets struct {
		data *Item
	}

	getItem := func(name string, duration time.Duration) *Item {
		item := &Item{data: name}
		item.setExpireTime(duration)
		return item
	}

	tests := []struct {
		name string
		data datasets
		err  string
		want bool
	}{
		{
			name: "checking that is not expired",
			data: datasets{getItem("test", time.Second)},
			err:  "item to not be expired",
			want: false,
		},
		{
			name: "checking that is not expired",
			data: datasets{getItem("test", 0-time.Second)},
			err:  "item to be expired",
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !a.Equal(tc.data.data.expired(), tc.want) {
				t.Errorf(tc.err)
			}
		})
	}
}
