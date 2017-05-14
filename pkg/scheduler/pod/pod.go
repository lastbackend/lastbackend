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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/scheduler/context"
)

func Provision (p *types.Pod) {

	var (
		log      = context.Get().GetLogger()
	)

	log.Debugf("Generate pod spec: %s", p.Meta.Name)



	log.Debugf("Allocate node for pod: %s", p.Meta.Name)
}