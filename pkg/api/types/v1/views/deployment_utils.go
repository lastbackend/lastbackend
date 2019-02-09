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

package views

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type DeploymentView struct{}

func (dv *DeploymentView) New(obj *types.Deployment, pl *types.PodList) *Deployment {
	d := Deployment{}
	d.SetMeta(obj.Meta)
	d.SetStatus(obj.Status)
	d.SetSpec(obj.Spec)

	d.Pods = make(map[string]Pod, 0)
	if pl != nil {
		d.JoinPods(pl)
	}

	return &d
}

func (d *Deployment) SetMeta(obj types.DeploymentMeta) {
	meta := DeploymentMeta{}
	meta.Name = obj.Name
	meta.Description = obj.Description
	meta.Version = obj.Version
	meta.SelfLink = obj.SelfLink.String()
	meta.Namespace = obj.Namespace
	meta.Service = obj.Service
	meta.Endpoint = obj.Endpoint
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	d.Meta = meta
}

func (d *Deployment) SetStatus(obj types.DeploymentStatus) {
	d.Status = DeploymentStatusInfo{
		State:   obj.State,
		Message: obj.Message,
	}
}

func (d *Deployment) SetSpec(obj types.DeploymentSpec) {
	mv := new(ManifestView)
	var spec = DeploymentSpec{
		Replicas: obj.Replicas,
		Template: mv.NewManifestSpecTemplate(obj.Template),
		Selector: mv.NewManifestSpecSelector(obj.Selector),
	}

	d.Spec = spec
}

func (d *Deployment) JoinPods(obj *types.PodList) {
	for _, p := range obj.Items {

		if p.Meta.Namespace != d.Meta.Namespace {
			continue
		}

		k, sl := p.SelfLink().Parent()
		if k != types.KindDeployment {
			continue
		}

		if sl.String() != d.Meta.SelfLink {
			continue
		}
		d.Pods[p.SelfLink().String()] = new(PodView).New(p)
	}
}

func (d *Deployment) ToJson() ([]byte, error) {
	return json.Marshal(d)
}

func (dv *DeploymentView) NewList(obj *types.DeploymentList, pods *types.PodList) *DeploymentList {
	dl := make(DeploymentList, 0)
	for _, d := range obj.Items {
		dv := new(DeploymentView)
		dp := dv.New(d, pods)
		dl = append(dl, dp)
	}
	return &dl
}

func (di *DeploymentList) ToJson() ([]byte, error) {
	return json.Marshal(di)
}
