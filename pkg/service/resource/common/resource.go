package common

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"k8s.io/client-go/1.5/pkg/api"
)

type PodListChannel struct {
	List  chan *api.PodList
	Error chan error
}

func GetPodListChannelWithOptions(client k8s.IK8S, nsQuery *NamespaceQuery, options api.ListOptions, limit int) PodListChannel {

	channel := PodListChannel{
		List:  make(chan *api.PodList, limit),
		Error: make(chan error, limit),
	}

	go func() {
		list, err := client.Core().Pods(nsQuery.ToRequestParam()).List(options)

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
