package common

import (
	"k8s.io/client-go/pkg/api"
)

func FilterNamespacedPodsBySelector(pods []api.Pod, namespace string, selector map[string]string) []api.Pod {

	var result []api.Pod

	for _, pod := range pods {
		if pod.ObjectMeta.Namespace == namespace && IsSelectorMatching(selector, pod.Labels) {
			result = append(result, pod)
		}
	}

	return result
}
