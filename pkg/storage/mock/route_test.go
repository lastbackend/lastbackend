//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"reflect"
	"testing"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

func Test_newRouteStorage(t *testing.T) {
	tests := []struct {
		name string
		want *RouteStorage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRouteStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRouteStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getRouteAsset(name, desc string) types.Route {
	p := types.Route{}
	p.Meta.Name = name
	p.Meta.Description = desc

	return p
}