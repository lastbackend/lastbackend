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
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

type PodListChannel struct {
	List  chan *api.PodList
	Error chan error
}

func GetPodListChannelWithOptions(client k8s.IK8S, nsQuery *NamespaceQuery, options v1.ListOptions, limit int) PodListChannel {

	channel := PodListChannel{
		List:  make(chan *api.PodList, limit),
		Error: make(chan error, limit),
	}

	go func() {
		list, err := client.CoreV1().Pods(nsQuery.ToRequestParam()).List(options)

		var items []api.Pod
		var apiPodList = new(api.PodList)

		err = converter.Convert_PodList_v1_to_api(list, apiPodList)
		if err != nil {
			channel.List <- nil
			channel.Error <- err
			return
		}

		for _, item := range apiPodList.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				items = append(items, item)
			}
		}

		apiPodList.Items = items

		for i := 0; i < limit; i++ {
			channel.List <- apiPodList
			channel.Error <- err
		}

	}()

	return channel
}
