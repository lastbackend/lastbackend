package pod

import (
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"k8s.io/client-go/1.5/pkg/api"
)

// ReplicationSetList contains a list of Pods in the cluster.
type PodList struct {
	ListMeta common.ListMeta `json:"listMeta"`

	// Unordered list of Pods.
	Pods []Pod `json:"pods"`
	//CumulativeMetrics []metric.Metric `json:"cumulativeMetrics"`
}

type PodStatus struct {
	// Status of the Pod. See Kubernetes API for reference.
	PodPhase api.PodPhase `json:"podPhase"`

	ContainerStates []api.ContainerState `json:"containerStates"`
}

type PodCell api.Pod

// Pod is a presentation layer view of Kubernetes Pod resource. This means
// it is Pod plus additional augumented data we can get from other sources
// (like services that target it).
type Pod struct {
	ObjectMeta common.ObjectMeta `json:"objectMeta"`
	TypeMeta   common.TypeMeta   `json:"typeMeta"`
	// More info on pod status
	PodStatus PodStatus `json:"podStatus"`
	// Count of containers restarts.
	RestartCount int32 `json:"restartCount"`
}

func CreatePodList(pods []api.Pod) PodList {

	podList := PodList{
		Pods:     make([]Pod, 0),
		ListMeta: common.ListMeta{TotalItems: len(pods)},
	}

	for _, pod := range pods {
		podDetail := Pod{
			ObjectMeta:   common.NewObjectMeta(pod.ObjectMeta),
			TypeMeta:     common.NewTypeMeta(common.ResourceKindPod),
			PodStatus:    getPodStatus(pod),
			RestartCount: getRestartCount(pod),
		}
		podList.Pods = append(podList.Pods, podDetail)

	}

	return podList
}

// Gets restart count of given pod (total number of its containers restarts).
func getRestartCount(pod api.Pod) int32 {
	var restartCount int32 = 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount += containerStatus.RestartCount
	}
	return restartCount
}

// getPodStatus returns a PodStatus object containing a summary of the pod's status.
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
