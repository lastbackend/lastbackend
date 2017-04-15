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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
)

func ToNodeSpec(obj types.NodeSpec) *Spec {
	spec := &Spec{}
	for _, pod := range obj.Pods {
		spec.Pods = append(spec.Pods, v1.Pod{
			Meta: v1.ToPodMeta(pod.Meta),
			Spec: v1.ToPodSpec(pod.Spec),
		})
	}
	return spec
}

func FromNodeSpec(spec Spec) *types.NodeSpec {

	var s = new(types.NodeSpec)
	for _, item := range spec.Pods {

		pod := types.PodNodeSpec{}

		pod.Meta = v1.FromPodMeta(item.Meta)
		pod.Spec = v1.FromPodSpec(item.Spec)

		s.Pods = append(s.Pods, pod)
	}

	return s
}

func (obj *Spec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
