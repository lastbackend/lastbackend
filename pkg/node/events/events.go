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
)

// NewInfoEvent - send node info event after
// node is successful accepted and each hour
func NewInfoEvent(ctx context.Context) error {

	return nil
}

// NewStateEvent - send node state event after
// node is successful accepted and each hour
func NewStateEvent(ctx context.Context) error {

	return nil
}

// NewPodStateEvent - send pod state event after
// node is successful accepted and each hour
func NewPodStateEvent(ctx context.Context, pod *types.Pod) error {

	if pod == nil {
		log.Errorf("Event: Pod state event: pod is empty")
		return errors.New("Event: Pod state event: pod is empty")
	}

	log.Debugf("Event: Pod state event state: %s", pod.Meta.Name)

	return nil
}

// NewRouteStateEvent - send pod state event after
// node is successful accepted and each hour
func NewRouteStateEvent(ctx context.Context, route, status string) error {

	if route == "" {
		log.Errorf("Event: Route state event: route is empty")
		return errors.New("Event: Route state event: route is empty")
	}

	log.Debugf("Event: Route state event state: %s", route)

	return nil
}
