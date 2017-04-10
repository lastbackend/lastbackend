package docker

import (
	"encoding/json"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"strings"
)

func (r *Runtime) PodList() (map[string]*types.Pod, error) {
	log := context.Get().GetLogger()
	log.Debug("Docker: retrieve pod list")

	var err error
	var pods types.PodMap

	pods = types.PodMap{
		Items: make(map[string]*types.Pod),
	}

	items, err := r.client.ContainerList(context.Background(), docker.ContainerListOptions{
		All: true,
	})

	for _, c := range items {

		log.Debug("Check container:", c.ID)

		// Check container is managed by LB
		// Meta: owner/project/service/pod/spec
		label, ok := c.Labels["LB_META"]
		if !ok {
			continue
		}

		info := strings.Split(label, "/")

		meta := types.PodMeta{
			Owner:   info[1],
			Project: info[2],
			Service: info[3],
			Spec:    info[5],
		}
		meta.ID = info[4]

		pod, ok := pods.Items[meta.ID]
		if !ok {
			pod = types.NewPod()
			pods.Items[meta.ID] = pod
		}
		pod.Meta = meta
		pod.Spec.ID = pod.Meta.Spec

		inspected, _ := r.client.ContainerInspect(context.Background(), c.ID)
		if container := GetContainer(c, inspected); container != nil {
			log.Debugf("Add container %s to pod %s", container.ID, pod.Meta.ID)
			pod.AddContainer(container)
		}

	}

	pds, err := json.Marshal(pods.Items)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(string(pds))

	if err != nil {
		return pods.Items, err
	}

	return pods.Items, err
}
