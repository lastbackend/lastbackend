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

package docker

import (
	docker "github.com/docker/docker/api/types"
	ctx "github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"strings"
)

func (r *Runtime) PodList(c context.Context) ([]*types.Pod, error) {

	var (
		err  error
		list []*types.Pod
	)

	pods := make(map[string]*types.Pod)

	items, err := r.client.ContainerList(c.Background(), docker.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return list, err
	}

	for _, container := range items {

		// Check container is managed by LB
		// Meta: owner/namespace/service/pod/spec
		label, ok := container.Labels["LB_META"]
		if !ok {
			continue
		}

		info := strings.Split(label, "/")

		pod, ok := pods[info[0]]
		if !ok {
			pod = types.NewPod()
			pods[info[0]] = pod
		}
		pod.Node.Hostname, _ = system.GetHostname()
		pod.Node.ID = *ctx.Get().GetID()
		pod.Meta.Name = info[0]
		pod.Spec.ID = info[1]
		pod.Spec.Containers = make(map[string]*types.ContainerSpec)

		container, err := r.ContainerInspect(c, container.ID)
		if err != nil || container == nil {
			continue
		}
		pod.Spec.Containers[container.Spec] = new(types.ContainerSpec)

		pod.State.Provision = false
		pod.State.Ready = true

		pod.AddContainer(container)
	}

	for _, p := range pods {
		p.UpdateState()
		list = append(list, p)
	}

	return list, err
}
