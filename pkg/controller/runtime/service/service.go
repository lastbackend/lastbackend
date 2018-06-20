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

package service

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
	"github.com/spf13/viper"
)

// Provision service
// Remove deployment or cancel if service is market for destroy
// Remove service if no active deployments present and service is marked for destroy
func Provision(svc *types.Service) error {

	var (
		stg = envs.Get().GetStorage()
	)

	sm := distribution.NewServiceModel(context.Background(), stg)

	if d, err := sm.Get(svc.Meta.Namespace, svc.Meta.Name); d == nil || err != nil {
		if d == nil {
			return errors.New(store.ErrEntityNotFound)
		}
		log.Errorf("%s:provision:> get deployment error: %v", logPrefix, err)
		return err
	}

	log.Debugf("%s:provision:> provision service: %v", logPrefix, svc.SelfLink())

	em := distribution.NewEndpointModel(context.Background(), stg)
	endpoint, err := em.Get(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s:provision:> get endpoint error: %v", logPrefix, err)
		return err
	}

	// Check service ports
	switch true {
	case len(svc.Spec.Template.Network.Ports) == 0 && endpoint != nil:
		if err := em.Remove(endpoint); err != nil {
			log.Errorf("%s:provision:> get endpoint error: %v", logPrefix, err)
			return err
		}
		svc.Status.Network.IP = ""
		break

	case len(svc.Spec.Template.Network.Ports) != 0 && endpoint == nil:
		opts := types.EndpointCreateOptions{
			IP:            svc.Spec.Template.Network.IP,
			Ports:         svc.Spec.Template.Network.Ports,
			Policy:        svc.Spec.Template.Network.Policy,
			BindStrategy:  svc.Spec.Template.Network.Strategy.Bind,
			RouteStrategy: svc.Spec.Template.Network.Strategy.Route,
			Domain:        fmt.Sprintf("%s-%s.%s", svc.Meta.Name, svc.Meta.Namespace, viper.GetString("domain.internal")),
		}

		ept, err := em.Create(svc.Meta.Namespace, svc.Meta.Name, &opts)
		if err != nil {
			log.Errorf("%s:provision:> get endpoint error: %v", logPrefix, err)
			return err
		}

		svc.Status.Network.IP = ept.Spec.IP

		break

	case len(svc.Spec.Template.Network.Ports) != 0 && endpoint != nil:

		var equal = true

		if (svc.Spec.Template.Network.IP != types.EmptyString && svc.Spec.Template.Network.IP != endpoint.Spec.IP) ||
			(svc.Spec.Template.Network.Policy != endpoint.Spec.Policy) ||
			(svc.Spec.Template.Network.Strategy.Bind != endpoint.Spec.Strategy.Bind) ||
			(svc.Spec.Template.Network.Strategy.Route != endpoint.Spec.Strategy.Route) {
			equal = false
		}

		// check ports equal
		for ext, pm := range svc.Spec.Template.Network.Ports {
			if cpm, ok := endpoint.Spec.PortMap[ext]; !ok || pm != cpm {
				equal = false
				break
			}
		}

		// check if some ports are deleted from spec but presents in endpoint spec
		if !equal {
			for ext, pm := range endpoint.Spec.PortMap {
				if cpm, ok := svc.Spec.Template.Network.Ports[ext]; !ok || pm != cpm {
					equal = false
					break
				}
			}
		}

		if equal {
			break
		}

		opts := types.EndpointUpdateOptions{
			Ports:         svc.Spec.Template.Network.Ports,
			Policy:        svc.Spec.Template.Network.Policy,
			BindStrategy:  svc.Spec.Template.Network.Strategy.Bind,
			RouteStrategy: svc.Spec.Template.Network.Strategy.Route,
		}

		if svc.Spec.Template.Network.IP != types.EmptyString && svc.Spec.Template.Network.IP != endpoint.Spec.IP {
			opts.IP = &svc.Spec.Template.Network.IP
		}

		if endpoint, err = em.Update(endpoint, &opts); err != nil {
			log.Errorf("%s:provision:> get endpoint error: %v", logPrefix, err)
			return err
		}

		svc.Status.Network.IP = endpoint.Spec.IP

		break
	}

	// Get all deployments per service
	dm := distribution.NewDeploymentModel(context.Background(), stg)
	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s:provision:> get deployments list error: %v", logPrefix, err)
		return err
	}

	if len(dl) == 0 && svc.Status.State == types.StateDestroy {
		if err := sm.Remove(svc); err != nil {
			log.Errorf("%s:provision:> remove service err: %v", logPrefix, err)
			return nil
		}
	}

	activeDeploymentExists := false

	for _, d := range dl {
		if d.Spec.State.Destroy {
			continue
		}

		if !svc.Spec.State.Destroy && d.Spec.Meta.Name == svc.Spec.Meta.Name {
			activeDeploymentExists = true
			continue
		}

		if svc.Spec.State.Destroy && !d.Spec.State.Destroy {
			if err := dm.Destroy(d); err != nil {
				log.Errorf("%s> remove deployment err: %s", err.Error())
				continue
			}
		}
	}

	// Check service is marked for destroy
	if activeDeploymentExists || svc.Spec.State.Destroy {
		return nil
	}

	// Create new deployment
	if _, err := dm.Create(svc); err != nil {
		log.Errorf("%s:provision:> create deployment err: %v", logPrefix, err)
		svc.Status.State = types.StateError
		svc.Status.Message = err.Error()
	}

	// Update service state
	svc.Status.State = types.StateProvision
	if err := distribution.NewServiceModel(context.Background(), stg).SetStatus(svc); err != nil {
		log.Errorf("%s:provision:> service set state err: %v", logPrefix, err)
		return err
	}

	return nil
}

// HandleStatus handles status of service
func HandleStatus(svc *types.Service) error {

	var (
		stg    = envs.Get().GetStorage()
		status = make(map[string]int)
	)

	if svc == nil {
		log.Errorf("%s:shandle_status:> service is nil", logPrefix)
		return nil
	}

	log.Debugf("%s:shandle_status:> handle service [%s] status", logPrefix, svc.SelfLink())

	dm := distribution.NewDeploymentModel(context.Background(), stg)
	sm := distribution.NewServiceModel(context.Background(), stg)

	// Skip state handle
	if svc.Status.State == types.StateDestroy {
		log.Debugf("%s:shandle_status:> skip service status [%s] handle: %s", logPrefix, svc.SelfLink())
		return nil
	}

	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s:shandle_status:> get pod list err: %v", logPrefix, err)
		return err
	}

	for _, di := range dl {
		switch di.Status.State {
		case types.StateDestroyed:
			status[types.StateDestroyed] += 1
			break
		}
	}

	if svc.Spec.State.Destroy && len(dl) == status[types.StateDestroyed] {
		log.Debugf("%s:shandle_status:> remove destroyed service: %s", logPrefix, svc.SelfLink())
		if err := sm.Remove(svc); err != nil {
			log.Errorf("%s:handle_status:> remove destroyed service [%s] err: %v", logPrefix, err)
			return err
		}

		return nil
	}

	return nil
}



func Update() error {
	return nil
}
