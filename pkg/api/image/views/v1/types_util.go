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

package v1

import "github.com/lastbackend/lastbackend/pkg/common/types"

func ToImageSpec(spec types.ImageSpec) ImageSpec {
	return ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
	}
}

func FromImageSpec(spec ImageSpec) types.ImageSpec {
	return types.ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
	}
}
