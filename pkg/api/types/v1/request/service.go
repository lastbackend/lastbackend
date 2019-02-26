//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
	"time"
)

type ServiceManifest struct {
	Meta ServiceManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec ServiceManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ServiceManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type ServiceManifestSpec struct {
	Replicas *int                  `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Network  *ManifestSpecNetwork  `json:"network,omitempty" yaml:"network,omitempty"`
	Strategy *ManifestSpecStrategy `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

func (s *ServiceManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *ServiceManifest) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ServiceManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, s)
}

func (s *ServiceManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *ServiceManifest) SetServiceMeta(svc *types.Service) {

	if svc.Meta.Name == types.EmptyString {
		svc.Meta.Name = *s.Meta.Name
	}

	if s.Meta.Description != nil {
		svc.Meta.Description = *s.Meta.Description
	}

	if s.Meta.Labels != nil {
		svc.Meta.Labels = s.Meta.Labels
	}

}

func (s *ServiceManifest) SetServiceSpec(svc *types.Service) (err error) {

	tn := svc.Spec.Network.Updated
	tc := svc.Spec.Template.Updated

	defer func() {
		if s.Spec.Replicas != nil {
			svc.Status.State = types.StateProvision
			return
		}

		if tn.Before(svc.Spec.Network.Updated) || tc.Before(svc.Spec.Template.Updated) {
			svc.Status.State = types.StateProvision
			return
		}
	}()

	if s.Spec.Replicas != nil {
		svc.Spec.Replicas = *s.Spec.Replicas
	}

	if s.Spec.Network != nil {

		if s.Spec.Network.IP != nil {
			svc.Spec.Network.IP = *s.Spec.Network.IP
		}

		if len(s.Spec.Network.Ports) > 0 {

			svc.Spec.Network.Ports = make(map[uint16]string, 0)

			for _, p := range s.Spec.Network.Ports {
				mp := strings.Split(p, ":")
				var base = 10
				var size = 16
				port, err := strconv.ParseUint(mp[0], base, size)
				if err != nil {
					continue
				}
				if len(mp) == 1 {
					svc.Spec.Network.Ports[uint16(port)] = mp[0]
				}

				if len(mp) == 2 {
					svc.Spec.Network.Ports[uint16(port)] = mp[1]
				}

			}
		}

		svc.Spec.Network.Updated = time.Now()
	}

	if s.Spec.Selector != nil {
		s.Spec.Selector.SetSpecSelector(&svc.Spec.Selector)
	} else {
		svc.Spec.Selector.SetDefault()
	}

	if s.Spec.Strategy != nil {
		if s.Spec.Strategy.Type != nil {
			svc.Spec.Strategy.Type = *s.Spec.Strategy.Type
		}
	}

	if s.Spec.Template != nil {

		if err := s.Spec.Template.SetSpecTemplate(&svc.Spec.Template); err != nil {
			return err
		}

	}

	return nil
}

func (s *ServiceManifest) GetManifest() *types.ServiceManifest {
	sm := new(types.ServiceManifest)
	return sm
}

// swagger:ignore
// swagger:model request_service_remove
type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
// swagger:model request_service_logs
type ServiceLogsOptions struct {
	Tail       int    `json:"tail"`
	Deployment string `json:"deployment"`
	Pod        string `json:"pod"`
	Container  string `json:"container"`
	Follow     bool   `json:"follow"`
}
