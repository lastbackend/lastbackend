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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/network"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/router"
)

type Runtime struct{}


func Restore(ctx context.Context) {
	network.Restore(ctx)
	router.Restore(ctx)
	pod.Restore(ctx)
}

func Provision(ctx context.Context, spec types.NodeNamespace) {
	log.Debugf("node spec namespace: %s", spec.Meta.Name)

	for _, r := range spec.Spec.Routes {
		log.Debugf("route: %v", r)
	}

	for _, p := range spec.Spec.Pods {
		log.Debugf("pod: %v", p)
	}

	for _, v := range spec.Spec.Volumes {
		log.Debugf("volume: %v", v)
	}

	for _, s := range spec.Spec.Secrets {
		log.Debugf("secret: %v", s)
	}

}


func Subscribe(ctx context.Context) {

	log.Debug("Runtime subscribe state")
	pc := make(chan *types.Pod)

	go func() {

		for {
			select {
			case p := <-pc:
				log.Debugf("Send new pod state event: %#v", p)
				events.NewPodStateEvent(ctx, p)
			}
		}
	}()

	envs.Get().GetCri().Subscribe(ctx, envs.Get().GetState().Pods(), pc)
}
