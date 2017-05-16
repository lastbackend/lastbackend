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
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func CurrentCmd() {

	ns, err := Current()
	if err != nil {
		fmt.Print(err)
		return
	}

	if ns == nil {
		fmt.Print("Namespace didn't select")
		return
	}

	ns.DrawTable()
}

func Current() (*n.Namespace, error) {

	var (
		err     error
		storage = c.Get().GetStorage()
	)

	ns, err := storage.Namespace().Load()
	if err != nil {
		return nil, errors.UnknownMessage
	}

	return ns, nil
}
