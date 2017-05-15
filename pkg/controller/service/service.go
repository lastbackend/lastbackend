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

package service

import (
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/controller/context"
	"github.com/lastbackend/lastbackend/pkg/controller/pod"
)

func Provision(svc *types.Service) error {

	var (
		log      = context.Get().GetLogger()
		stg      = context.Get().GetStorage()
		replicas int
	)

	log.Debugf("Service Contoller: provision service: %s/%s", svc.Meta.Namespace, svc.Meta.Name)

	for _, pod := range svc.Pods {
		if pod.Spec.State != types.StateDestroy {
			replicas++
		}
	}

	if replicas < svc.Meta.Replicas {
		log.Debug("Service: Replicas: create a new replicas")
		for i := 0; i < (svc.Meta.Replicas - replicas); i++ {
			p := pod.PodCreate(svc)
			svc.Pods[p.Meta.Name] = p
		}
	}

	if replicas > svc.Meta.Replicas {
		log.Debug("Service: Replicas: remove  unneeded replicas")
		names := make([]string, 0, len(svc.Pods))
		for n, p := range svc.Pods {
			if p.Spec.State != types.StateDestroy {
				names = append(names, n)
			}
		}

		for i := 0; i < (replicas - svc.Meta.Replicas); i++ {
			if len(names) > 0 {
				pod.PodRemove(svc.Pods[names[len(names)-1]])
			}
			names = names[0 : len(names)-1]
		}
	}

	for _, p := range svc.Pods {
		log.Debug("Service: provision pods")
		pod.PodSetSpec(p, svc.Spec)
		log.Debug("Service: save new pod spec")
		if err := stg.Pod().Upsert(context.Get().Background(), svc.Meta.Namespace, p); err != nil {
			log.Errorf("Service: save pod spec error: %s", err.Error())
			return err
		}
	}

	return nil
}
