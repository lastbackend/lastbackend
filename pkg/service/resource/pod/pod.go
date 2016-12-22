package pod

import (
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"github.com/lastbackend/lastbackend/pkg/service/resource/container"
	"k8s.io/client-go/1.5/pkg/api"
)

const kind = "pod"

type PodList struct {
	ListMeta common.ListMeta `json:"meta"`
	Pods     []Pod           `json:"pods"`
}

type PodStatus struct {
	PodPhase        api.PodPhase         `json:"phase"`
	ContainerStates []api.ContainerState `json:"container_states"`
}

type Pod struct {
	ObjectMeta    common.ObjectMeta        `json:"meta"`
	TypeMeta      common.TypeMeta          `json:"spec"`
	PodStatus     PodStatus                `json:"status"`
	RestartCount  int32                    `json:"restart_count"`
	ContainerList *container.ContainerList `json:"containers"`
}

func CreatePodList(pods []api.Pod) *PodList {

	podList := PodList{
		ListMeta: common.ListMeta{Total: len(pods)},
		Pods:     make([]Pod, 0),
	}

	for _, pod := range pods {

		var p = Pod{
			ObjectMeta:    common.NewObjectMeta(pod.ObjectMeta),
			TypeMeta:      common.NewTypeMeta(kind),
			PodStatus:     getPodStatus(pod),
			RestartCount:  getRestartCount(pod),
			ContainerList: container.CreateContainerList(pod.Spec.Containers),
		}

		podList.Pods = append(podList.Pods, p)
	}

	return &podList
}

func getRestartCount(pod api.Pod) (count int32) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		count += containerStatus.RestartCount
	}

	return count
}

func getPodStatus(pod api.Pod) PodStatus {
	var states []api.ContainerState

	for _, containerStatus := range pod.Status.ContainerStatuses {
		states = append(states, containerStatus.State)
	}

	return PodStatus{
		PodPhase:        pod.Status.Phase,
		ContainerStates: states,
	}
}
