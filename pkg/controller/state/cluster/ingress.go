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

package cluster

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"

	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/log"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

const (
	logPrefixIngress = "observer:cluster:endpoint"
)

//TODO: handle node status and update ingress ready state
func ingressHandle(ctx context.Context, cs *ClusterState, node *types.Node) error {

	im := distribution.NewIngressModel(ctx, envs.Get().GetStorage())
	ingress, err := im.Get(node.SelfLink())
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			return nil
		}

		log.Errorf("%s: get ingress error: %s", logPrefixIngress, err.Error())
	}

	if node.Status.Mode.Ingress && ingress == nil {
		return ingressCreate(ctx, cs, node)
	}

	if !node.Status.Mode.Ingress && ingress != nil {
		return ingressRemove(ctx, cs, node)
	}

	return nil
}

// TODO: add ingress to cluster state
func ingressCreate(ctx context.Context, cs *ClusterState, node *types.Node) error {

	log.V(logLevel).Debugf("%s: create ingress for node: %s", logPrefixIngress, node.SelfLink())

	var ingress = new(types.Ingress)
	ingress.Meta.Name = node.SelfLink()
	ingress.Meta.Node = node.SelfLink()
	ingress.Meta.SetDefault()

	im := distribution.NewIngressModel(ctx, envs.Get().GetStorage())

	if err := im.Create(ingress); err != nil {
		log.Errorf("can not create ingress: %s", err.Error())
	}

	return nil
}

// TODO: remove ingress from cluster state
func ingressRemove(ctx context.Context, cs *ClusterState, node *types.Node) error {

	im := distribution.NewIngressModel(ctx, envs.Get().GetStorage())
	ingress, err := im.Get(node.SelfLink())
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			return nil
		}

		log.Errorf("ingress get err: %s", err.Error())
		return err
	}

	if ingress == nil {
		return nil
	}
	if err := im.Remove(ingress); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			return nil
		}

		log.Errorf("ingress remove err: %s", err.Error())
		return err
	}

	return nil
}
