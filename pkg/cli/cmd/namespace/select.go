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

package namespace

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
)

func SelectCmd(name string) {

	ns, err := Select(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The namespace `%s` was selected as the current\n", ns.Meta.Name)
}

func Select(name string) (*view.Namespace, error) {

	stg := context.Get().GetStorage()

	ns, err := Fetch(name)
	if err != nil {
		return nil, err
	}

	if err := stg.Namespace().Save(ns); err != nil {
		return nil, err
	}

	return ns, nil
}
