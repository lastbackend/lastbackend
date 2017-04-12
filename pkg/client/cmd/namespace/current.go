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
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func CurrentCmd() {

	var (
		log = c.Get().GetLogger()
	)

	namespace, err := Current()

	if err != nil {
		log.Error(err)
		return
	}

	if namespace == nil {
		log.Info("Namespace didn't select")
		return
	}

	namespace.DrawTable()
}

func Current() (*n.Namespace, error) {

	var (
		err       error
		storage   = c.Get().GetStorage()
		namespace = new(n.Namespace)
	)

	err = storage.Get("namespace", namespace)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return namespace, nil
}
