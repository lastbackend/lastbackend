package pod

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"github.com/lastbackend/lastbackend/pkg/service/resource/container"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"time"
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
	Status        PodStatus                `json:"status"`
	RestartCount  int32                    `json:"restart_count"`
	IP            string                   `json:"ip"`
	StartTime     time.Time                `json:"startTime"`
	RestartPolicy api.RestartPolicy        `json:"restartPolicy,omitempty"`
	ContainerList *container.ContainerList `json:"containers"`
}

func (p *Pod) Remove(client k8s.IK8S) error {
	var opts = new(v1.DeleteOptions)
	return client.CoreV1().Pods(p.ObjectMeta.Namespace).Delete(p.ObjectMeta.Name, opts)
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
			Status:        getPodStatus(pod),
			IP:            pod.Status.PodIP,
			StartTime:     pod.Status.StartTime.Time,
			RestartPolicy: pod.Spec.RestartPolicy,
			RestartCount:  getRestartCount(pod),
			ContainerList: container.CreateContainerList(pod.Spec.Containers, pod.Status.ContainerStatuses),
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
