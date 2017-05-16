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

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/api/pod/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

func ToNodeSpec(obj types.NodeSpec) *Spec {
	spec := &Spec{}
	for _, pod := range obj.Pods {
		spec.Pods = append(spec.Pods, v1.Pod{
			Meta:  v1.ToPodMeta(pod.Meta),
			State: v1.ToPodState(pod.State),
			Spec:  v1.ToPodSpec(pod.Spec),
		})
	}
	return spec
}

func FromNodeSpec(spec Spec) *types.NodeSpec {

	var s = new(types.NodeSpec)
	s.Pods = make(map[string]types.PodNodeSpec, len(spec.Pods))
	for _, item := range spec.Pods {

		pod := types.PodNodeSpec{}

		pod.Meta = v1.FromPodMeta(item.Meta)
		pod.State = v1.FromPodState(item.State)
		pod.Spec = v1.FromPodSpec(item.Spec)

		s.Pods[pod.Meta.Name] = pod
	}

	return s
}

func (obj *Spec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
