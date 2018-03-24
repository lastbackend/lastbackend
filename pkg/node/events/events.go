//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package events

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/pkg/errors"
)

// NewConnectEventt - send node info event after
// node is successful accepted and each hour
func NewConnectEventt(ctx context.Context) error {

	var (
		c = envs.Get().GetClient()
	)

	opts := v1.Request().Node().NodeConnectOptions()
	opts.Info = envs.Get().GetState().Node().Info
	opts.Status = envs.Get().GetState().Node().Status

	return c.Connect(ctx, opts)
}

// NewStatusEvent - send node state event after
// node is successful accepted and each hour
func NewStatusEvent(ctx context.Context) error {
	var (
		c = envs.Get().GetClient()
	)

	opts := v1.Request().Node().NodeStatusOptions()
	opts.Capacity = envs.Get().GetState().Node().Status.Capacity
	opts.Allocated = envs.Get().GetState().Node().Status.Allocated

	return c.SetStatus(ctx, opts)
}

// NewPodStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewPodStatusEvent(ctx context.Context, pod string) error {

	var (
		c = envs.Get().GetClient()
		p = envs.Get().GetState().Pods().GetPod(pod)
	)

	if pod == "" {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", pod)

	opts := v1.Request().Node().NodePodStatusOptions()
	opts.Stage = p.State
	opts.Message = p.Message
	opts.Containers = p.Containers
	opts.Network = p.Network
	opts.Steps = p.Steps

	return c.SetPodStatus(ctx, pod, opts)
}

// NewRouteStatusEvent - send route state event after
// node is successful accepted and each hour
func NewRouteStatusEvent(ctx context.Context, route string) error {

	var (
		c = envs.Get().GetClient()
	)

	if route == "" {
		log.Errorf("Event: route state event: route is empty")
		return errors.New("Event: route state event: route is empty")
	}

	log.Debugf("Event: route state event state: %s", route)

	opts := v1.Request().Node().NodeRouteStatusOptions()
	return c.SetRouteStatus(ctx, route, opts)
}

// NewRouteStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewVolumeStatusEvent(ctx context.Context, volume string) error {

	var (
		c = envs.Get().GetClient()
	)

	if volume == "" {
		log.Errorf("Event: volume state event: volume is empty")
		return errors.New("Event: volume state event: volume is empty")
	}

	log.Debugf("Event: volume state event state: %s", volume)

	opts := v1.Request().Node().NodeVolumeStatusOptions()
	return c.SetVolumeStatus(ctx, volume, opts)
}
