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

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func PodCreate(d *types.Deployment) (*types.Pod, error) {
	dm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	return dm.Create(d)
}

func PodDestroy(p *types.Pod) error {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if p.Meta.Node == types.EmptyString {
		p.Status.State = types.StateDestroyed
	} else {
		p.Status.State = types.StateDestroy
	}

	if err := pm.Update(p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// Remove pod from storage, remove manifest and release node if leased
func PodRemove(p *types.Pod) error {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if p.Meta.Node == types.EmptyString {
		if err := pm.Remove(p); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}
