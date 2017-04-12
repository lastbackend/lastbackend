package v1

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
)

func ToNodeSpec(obj []*types.Pod) Spec {
	spec := Spec{}
	for _, pod := range obj {
		spec.Pods = append(spec.Pods, v1.Pod{
			Meta:  v1.ToPodMeta(pod.Meta),
			Spec:  v1.ToPodSpec(pod.Spec),
			State: v1.ToPodState(pod.State),
		})
	}
	return spec
}

func FromNodeSpec(spec Spec) []*types.Pod {
	var pods []*types.Pod
	for _, s := range spec.Pods {
		pod := types.NewPod()
		pod.Meta = v1.FromPodMeta(s.Meta)
		pod.Spec = v1.FromPodSpec(s.Spec)
		pod.State = v1.FromPodState(s.State)
		pods = append(pods, pod)
	}
	return pods
}

func (obj *Spec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
