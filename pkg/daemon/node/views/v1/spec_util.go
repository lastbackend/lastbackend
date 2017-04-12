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
			Meta:  v1.ToPodMeta(pod.Meta),
			Spec:  v1.ToPodSpec(pod.Spec),
			State: v1.ToPodState(pod.State),
		})
	}
	return spec
}

func FromNodeSpec(spec Spec) *types.NodeSpec {

	var s = new(types.NodeSpec)
	for _, item := range spec.Pods {

		pod := new(types.PodNodeSpec)

		pod.Meta = v1.FromPodMeta(item.Meta)
		pod.Spec = v1.FromPodSpec(item.Spec)
		pod.State = v1.FromPodState(item.State)

		s.Pods = append(s.Pods, pod)
	}

	return s
}

func (obj *Spec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
