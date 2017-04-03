package pod

import "github.com/lastbackend/lastbackend/pkg/apis/types"

type Manager interface {
	GetPods() types.PodList
	SetPods(pods types.PodList)

	GetPod(uuid string) *types.Pod
	AddPod(pod *types.Pod)
	SetPod(pod *types.Pod)
	DelPod(pod *types.Pod)
}
