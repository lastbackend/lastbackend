package docker

import (
	"encoding/json"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func (r *Runtime) PodList() (map[types.PodID]*types.Pod, error) {
	log := context.Get().GetLogger()
	log.Debug("Docker: retrieve pod list")

	var err error
	var pods types.PodMap

	pods = types.PodMap{
		Items: make(map[types.PodID]*types.Pod),
	}

	items, err := r.client.ContainerList(context.Background(), docker.ContainerListOptions{
		All: true,
	})

	for _, c := range items {

		log.Debug("Check container:", c.ID)

		// Check container is managed by LB
		_, ok := c.Labels["LB_MANAGED"]
		if !ok {
			continue
		}

		meta := GetPodMetaFromContainer(c)

		pod, ok := pods.Items[meta.ID]
		if !ok {
			pod = new(types.Pod)
			pods.Items[meta.ID] = pod
			pod.Meta = meta
			pod.Containers = make(map[types.ContainerID]types.Container)
		}

		info, err := r.client.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			continue
		}

		pod.Spec.Containers = append(pod.Spec.Containers, GetContainerSpecFromContainer(info))
		pod.AddContainer(GetContainer(c))
	}

	pds, err := json.Marshal(pods.Items)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debugf(string(pds))

	if err != nil {
		return pods.Items, err
	}

	return pods.Items, err
}
