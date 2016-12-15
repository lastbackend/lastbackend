package common

import (
	"k8s.io/client-go/1.5/pkg/api"
)

func FilterNamespacedPodsBySelector(pods []api.Pod, namespace string, resourceSelector map[string]string) []api.Pod {

	var result []api.Pod

	for _, pod := range pods {
		if pod.ObjectMeta.Namespace == namespace && IsSelectorMatching(resourceSelector, pod.Labels) {
			result = append(result, pod)
		}
	}

	return result
}
