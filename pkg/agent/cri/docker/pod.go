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

		pod, ok := pods[info[0]]
		if !ok {
			pod = types.NewPod()
			pods[info[0]] = pod
		}
		pod.Meta.ID = info[0]
		pod.Spec.ID = info[1]

		container, _, err := r.ContainerInspect(c.ID)
		if err != nil || container == nil {
			continue
		}

		pod.AddContainer(container)
	}

	for _, p := range pods {
		p.UpdateState()
		list = append(list, p)
	}

	return list, err
}
