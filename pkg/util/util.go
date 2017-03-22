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

package util

import (
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"reflect"
)

// SetZeroValue would set the object of objPtr to zero value of its type.
func SetZeroValue(objPtr interface{}) error {

	v, err := converter.EnforcePtr(objPtr)
	if err != nil {
		return err
	}

	v.Set(reflect.Zero(v.Type()))

	return nil
}
