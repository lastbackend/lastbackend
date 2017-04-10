package docker

import (
	"context"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"strings"
)

func (r *Runtime) PodList() ([]*types.Pod, error) {

	var (
		err  error
		list []*types.Pod
	)

	pods := make(map[string]*types.Pod)

	items, err := r.client.ContainerList(context.Background(), docker.ContainerListOptions{
		All: true,
	})

	if err != nil {
		return list, err
	}

	for _, c := range items {

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

		pod, ok := pods[meta.ID]
		if !ok {
			pod = types.NewPod()
			pods[meta.ID] = pod
		}
		pod.Meta = meta
		pod.Spec.ID = pod.Meta.Spec

		inspected, _ := r.client.ContainerInspect(context.Background(), c.ID)
		if container := GetContainer(c, inspected); container != nil {
			pod.AddContainer(container)
		}

	}

	for _, p := range pods {
		p.UpdateState()
		list = append(list, p)
	}

	return list, err
}
