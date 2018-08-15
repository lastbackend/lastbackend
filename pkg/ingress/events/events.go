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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

// NewConnectEventt - send ingress info event after
// ingress is successful accepted and each hour
func NewConnectEvent(ctx context.Context) error {

	var (
		c = envs.Get().GetClient()
	)

	opts := v1.Request().Ingress().IngressConnectOptions()
	opts.Status.Ready = true

	return c.Connect(ctx, opts)

}

// NewRouteStatusEvent - send route state event
func NewRouteStatusEvent(ctx context.Context, route string) error {

	var (
		c = envs.Get().GetClient()
	)

	if route == "" {
		log.Errorf("Event: route state event: route is empty")
		return errors.New("Event: route state event: route is empty")
	}

	log.V(logLevel).Debugf("Event: route state event state: %s", route)

	opts := v1.Request().Ingress().IngressRouteStatusOptions()
	return c.SetRouteStatus(ctx, route, opts)
}
