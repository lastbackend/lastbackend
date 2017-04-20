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

package namespace

import (
	"fmt"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func CurrentCmd() {

	namespace, err := Current()

	if err != nil {
		fmt.Print(err)
		return
	}

	if namespace == nil {
		fmt.Print("Namespace didn't select")
		return
	}

	namespace.DrawTable()
}

func Current() (*n.Namespace, error) {

	var (
		err       error
		storage   = c.Get().GetStorage()
		namespace *n.Namespace
	)

	var sName string
	if c.Get().IsMock() {
		sName = "test"
	} else {
		sName = "namespace"
	}
	err = storage.Get(sName, &namespace)
	if err != nil {
		return nil, errors.UnknownMessage
	}

	return namespace, nil
}
