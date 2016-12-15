package common

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
)

type PodListChannel struct {
	List  chan *api.PodList
	Error chan error
}

func GetPodListChannelWithOptions(client k8s.IK8S, nsQuery *NamespaceQuery, options api.ListOptions, numReads int) PodListChannel {

	channel := PodListChannel{
		List:  make(chan *api.PodList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.Core().Pods(nsQuery.ToRequestParam()).List(options)

		var filteredItems []api.Pod
		var apiPodList = new(api.PodList)

		err = v1.Convert_v1_PodList_To_api_PodList(list, apiPodList, nil)
		if err != nil {
			channel.List <- nil
			channel.Error <- err
			return
		}

		for _, item := range apiPodList.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}

		apiPodList.Items = filteredItems

		for i := 0; i < numReads; i++ {
			channel.List <- apiPodList
			channel.Error <- err
		}

	}()

	return channel
}
