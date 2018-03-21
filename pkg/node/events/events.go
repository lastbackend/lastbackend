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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/pkg/errors"
)

// NewInfoEvent - send node info event after
// node is successful accepted and each hour
func NewInfoEvent(ctx context.Context) error {

	var (
		c = envs.Get().GetClient()
	)

	opts := v1.Request().Node().NodeInfoOptions()

	return c.SetInfo(ctx, opts)
}

// NewStatusEvent - send node state event after
// node is successful accepted and each hour
func NewStatusEvent(ctx context.Context) error {
	var (
		c = envs.Get().GetClient()
	)

	opts := v1.Request().Node().NodeStatusOptions()
	return c.SetStatus(ctx, opts)
}


// NewPodStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewPodStatusEvent(ctx context.Context, pod *types.Pod) error {

	var (
		c = envs.Get().GetClient()
	)

	if pod == nil {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", pod.Meta.Name)

	opts := v1.Request().Node().NodePodStatusOptions()
	return c.SetPodStatus(ctx, opts)
}

// NewRouteStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewRouteStatusEvent(ctx context.Context, route *types.Route) error {

	var (
		c = envs.Get().GetClient()
	)

	if route == nil {
		log.Errorf("Event: route state event: pod is empty")
		return errors.New("Event: route state event: pod is empty")
	}

	log.Debugf("Event: route state event state: %s", route.Meta.Name)

	opts := v1.Request().Node().NodeRouteStatusOptions()
	return c.SetRouteStatus(ctx, opts)
}

// NewRouteStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewVolumeStatusEvent(ctx context.Context, volume *types.Volume) error {

	var (
		c = envs.Get().GetClient()
	)

	if volume == nil {
		log.Errorf("Event: volume state event: pod is empty")
		return errors.New("Event: volume state event: pod is empty")
	}

	log.Debugf("Event: volume state event state: %s", volume.Meta.Name)

	opts := v1.Request().Node().NodeVolumeStatusOptions()
	return c.SetVolumeStatus(ctx, opts)
}
