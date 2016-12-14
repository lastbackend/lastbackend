package common

import (
	"k8s.io/client-go/1.5/pkg/api"
)

// FilterNamespacedPodsBySelector returns pods targeted by given resource label selector in given
// namespace.
func FilterNamespacedPodsBySelector(pods []api.Pod, namespace string, resourceSelector map[string]string) []api.Pod {

	var matchingPods []api.Pod
	for _, pod := range pods {
		if pod.ObjectMeta.Namespace == namespace && IsSelectorMatching(resourceSelector, pod.Labels) {
			matchingPods = append(matchingPods, pod)
		}
	}

	return matchingPods
}
