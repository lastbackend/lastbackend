//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

type DeploymentManifest struct {
	Meta DeploymentManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec DeploymentManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type DeploymentManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type DeploymentManifestSpec struct {
	Replicas *int                  `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Network  *ManifestSpecNetwork  `json:"network,omitempty" yaml:"network,omitempty"`
	Strategy *ManifestSpecStrategy `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

func (s *DeploymentManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *DeploymentManifest) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *DeploymentManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, s)
}

func (s *DeploymentManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *DeploymentManifest) SetDeploymentMeta(dep *types.Deployment) {

	if dep.Meta.Name == types.EmptyString {
		dep.Meta.Name = *s.Meta.Name
	}

	if s.Meta.Description != nil {
		dep.Meta.Description = *s.Meta.Description
	}

	if s.Meta.Labels != nil {
		dep.Meta.Labels = s.Meta.Labels
	}

}

func (s *DeploymentManifest) SetDeploymentSpec(dep *types.Deployment) (err error) {

	defer func() {
		if s.Spec.Replicas != nil {
			dep.Status.State = types.StateProvision
			return
		}
	}()

	if s.Spec.Replicas != nil {
		dep.Spec.Replicas = *s.Spec.Replicas
	}

	if s.Spec.Selector != nil {
		s.Spec.Selector.SetSpecSelector(&dep.Spec.Selector)
	} else {
		dep.Spec.Selector.SetDefault()
	}

	if s.Spec.Template != nil {

		if err := s.Spec.Template.SetSpecTemplate(&dep.Spec.Template); err != nil {
			return err
		}

	}

	return nil
}

func (s *DeploymentManifest) GetManifest() *types.DeploymentManifest {
	sm := new(types.DeploymentManifest)
	return sm
}

// DeploymentUpdateOptions represents options availible to update in deployment
//
// swagger:model request_deployment_update
type DeploymentUpdateOptions struct {
	// Number of replicas
	// required: false
	Replicas *int `json:"replicas"`
	// Deployment status for update
	// required: false
	Status *struct {
		State   string `json:"state"`
		Message string `json:"message"`
	} `json:"status"`
}

type DeploymentRemoveOptions struct {
	Force bool `json:"force"`
}
