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

package pod

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
)

func HandleStatus(p *types.Pod) error {

	var (
		stg       = envs.Get().GetStorage()
		logPrefix = "controller:pod:controller:status>"
		status    = make(map[string]int)
		message   string
	)

	pm := distribution.NewPodModel(context.Background(), stg)
	dm := distribution.NewDeploymentModel(context.Background(), stg)

	d, err := dm.Get(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
	if err != nil {
		log.Errorf("%s get deployment err: %v", logPrefix, err)
		return err
	}
	if d == nil {
		log.Warnf("%s: deployment already removed: [%s:%s:%s]", logPrefix, p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
		if err := pm.Remove(context.Background(), p); err != nil {
			log.Errorf("%s: remove pod [%s] err: %v", logPrefix, p.SelfLink())
			return err
		}
		return nil
	}

	pl, err := pm.ListByDeployment(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
	if err != nil {
		log.Errorf("%s> get pod list err: %v", logPrefix, err)
		return err
	}

	for _, ps := range pl {

		switch ps.Status.Stage {

		case types.StateError:
			status[types.StateError] += 1
			// TODO: check if many pods contains different errors: create an error map
			message = ps.Status.Message
			break
		case types.StateProvision:
			status[types.StateProvision] += 1
			break
		case types.StatePull:
			status[types.StateProvision] += 1
			break
		case types.StateCreated:
			status[types.StateProvision] += 1
			break
		case types.StateStarting:
			status[types.StateProvision] += 1
			break
		case types.StateStarted:
			status[types.StateRunning] += 1
			break
		case types.StateRunning:
			status[types.StateRunning] += 1
			break
		case types.StateStopped:
			status[types.StateStopped] += 1
			break
		case types.StateExited:
			status[types.StateStopped] += 1
			break
		case types.StateDestroy:
			status[types.StateDestroy] += 1
			break
		case types.StateDestroyed:
			status[types.StateDestroyed] += 1
			break
		case types.StateWarning:
			status[types.StateWarning] += 1
			break
		}
	}

	opts := new(request.DeploymentUpdateOptions)

	switch true {
	case status[types.StateError] > 0:
		opts.Status.State = types.StateError
		opts.Status.Message = message
		break
	case status[types.StateProvision] > 0:
		opts.Status.State = types.StateProvision
		opts.Status.Message = ""
		break
	case status[types.StateDestroy] > 0:
		opts.Status.State = types.StateDestroy
		opts.Status.Message = ""
		break
	case status[types.StateWarning] > 0 && d.Status.State != types.StateDestroy:
		opts.Status.State = types.StateWarning
		opts.Status.Message = p.Status.Message
		break
	case status[types.StateRunning] == d.Spec.Replicas:
		opts.Status.State = types.StateRunning
		break
	case status[types.StateStopped] == d.Spec.Replicas:
		opts.Status.State = types.StateStopped
		break
	case status[types.StateDestroyed] == 1 && p.Status.Stage == types.StateDestroyed:
		opts.Status.State = types.StateDestroyed
		break
	}

	if p.Status.Stage == types.StateDestroyed {
		if err := pm.Remove(context.Background(), p); err != nil {
			log.Errorf("%s: remove pod [%s] err: %v", logPrefix, p.SelfLink(), err)
			return err
		}
	}

	if err := dm.Update(d, opts); err != nil {
		log.Errorf("%s> update deployment err: %v", logPrefix, err)
		return err
	}

	return nil
}
