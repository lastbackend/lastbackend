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
