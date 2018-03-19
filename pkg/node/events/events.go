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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/api/request/v1"
)

// NewInfoEvent - send node info event after
// node is successful accepted and each hour
func NewInfoEvent(ctx context.Context) error {

	var (
		c = envs.Get().GetClient()
	)

	opts := v1.NodeInfoOpts{}
	return c.SetInfo(ctx, opts)
}

// NewStateEvent - send node state event after
// node is successful accepted and each hour
func NewStateEvent(ctx context.Context) error {
	var (
		c = envs.Get().GetClient()
	)

	opts := v1.NodeStateOpts{}
	return c.SetState(ctx, opts)
}


// NewPodStateEvent - send pod state event after
// node is successful accepted and each hour
func NewPodStateEvent(ctx context.Context, pod *types.Pod) error {

	var (
		c = envs.Get().GetClient()
	)

	if pod == nil {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", pod.Meta.Name)


	opts := v1.NodePodStateOpts{}
	return c.SetPodState(ctx, opts)

	return nil
}

// NewRouteStateEvent - send pod state event after
// node is successful accepted and each hour
func NewRouteStateEvent(ctx context.Context, route *types.Route) error {

	var (
		c = envs.Get().GetClient()
	)

	if route == nil {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", route.Meta.Name)


	opts := v1.NodePodStateOpts{}
	return c.SetPodState(ctx, opts)

	return nil
}

// NewRouteStateEvent - send pod state event after
// node is successful accepted and each hour
func NewVolumeStateEvent(ctx context.Context, volume *types.Volume) error {

	var (
		c = envs.Get().GetClient()
	)

	if volume == nil {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", volume.Meta.Name)


	opts := v1.NodePodStateOpts{}
	return c.SetPodState(ctx, opts)

	return nil
}
