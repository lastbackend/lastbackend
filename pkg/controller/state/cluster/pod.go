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

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func PodProvision(p *types.Pod, cs *ClusterState) error {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	// Check node is allocated
	if p.Meta.Node == types.EmptyString {
		// Allocate Node

		var RAM int64

		for _, s := range p.Spec.Template.Containers {
			RAM += s.Resources.Request.RAM
		}

		opts := types.NodeLeaseOptions{
			Memory: &RAM,
		}

		node, err := cs.Lease(opts)
		if err != nil {
			return err
		}

		if node == nil {
			p.Status.State = types.StateError
			p.Status.Message = "node not found"

			if err := pm.Update(p); err != nil {
				log.Errorf("%s", err.Error())
				return err
			}

			return nil
		}
		p.Meta.Node = node.SelfLink()

	}

	// Check manifest
	mm := distribution.NewManifestModel(context.Background(), envs.Get().GetStorage())
	m, err := mm.PodManifestGet(p.Meta.Node, p.Meta.SelfLink)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	if m == nil {
		pm := types.PodManifest(p.Spec)

		if err := mm.PodManifestAdd(p.Meta.Node, p.Meta.SelfLink, &pm); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	p.Status.State = types.StateProvision
	if err := pm.Update(p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func PodDestroy(p *types.Pod) error {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if p.Meta.Node == types.EmptyString {
		p.Status.State = types.StateDestroyed
		if err := pm.Update(p); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		return nil
	}

	mm := distribution.NewManifestModel(context.Background(), envs.Get().GetStorage())
	m, err := mm.PodManifestGet(p.Meta.Node, p.Meta.SelfLink)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	if m == nil {
		p.Status.State = types.StateDestroyed
		if err := pm.Update(p); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		return nil
	}

	// Update manifest
	*m = types.PodManifest(p.Spec)
	if err := mm.PodManifestSet(p.Meta.Node, p.Meta.SelfLink, m); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// Remove pod from storage, remove manifest and release node if leased
func PodRemove(p *types.Pod, cs *ClusterState) error {

	var RAM int64

	if p.Meta.Node == types.EmptyString {
		return nil
	}

	// Remove manifest
	mm := distribution.NewManifestModel(context.Background(), envs.Get().GetStorage())
	err := mm.PodManifestDel(p.Meta.Node, p.Meta.SelfLink)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}
	// Release node

	for _, s := range p.Spec.Template.Containers {
		RAM += s.Resources.Request.RAM
	}

	opts := types.NodeLeaseOptions{
		Node:   &p.Meta.SelfLink,
		Memory: &RAM,
	}

	if _, err := cs.Release(opts); err != nil {
		return err
	}

	p.Meta.Node = types.EmptyString
	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if err := pm.Update(p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}
