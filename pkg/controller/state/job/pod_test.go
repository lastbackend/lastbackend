//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package job

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"strings"
)

func getPodAsset(t *types.Task, state, message string) *types.Pod {

	p := new(types.Pod)

	p.Meta.SetDefault()
	p.Meta.Namespace = t.Meta.Namespace
	p.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	p.Meta.Namespace = t.Meta.Namespace

	sl, _ := types.NewPodSelfLink(types.KindTask, t.SelfLink().String(), p.Meta.Name)
	p.Meta.SelfLink = *sl

	p.Status.State = state
	p.Status.Message = message

	if state == types.StateReady {
		p.Status.Running = true
	}

	p.Spec.State = t.Spec.State
	p.Spec.Template = t.Spec.Template

	return p
}

func getPodCopy(pod *types.Pod) *types.Pod {
	p := *pod
	return &p
}
