package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
)

type TemplateList map[string][]string

type Template struct {
	Namespaces             []api.Namespace             `json:"namespaces,omitempty"`
	PersistentVolumes      []api.PersistentVolume      `json:"persistent_volumes,omitempty"`
	PersistentVolumeClaims []api.PersistentVolumeClaim `json:"persistent_volume_claims,omitempty"`
	ServiceAccounts        []api.ServiceAccount        `json:"service_accounts,omitempty"`
	Services               []api.Service               `json:"services,omitempty"`
	ReplicationControllers []api.ReplicationController `json:"replication_controllers,omitempty"`
	Pods                   []api.Pod                   `json:"pods,omitempty"`
	Deployments            []extensions.Deployment     `json:"deployments,omitempty"`
}

func (t *Template) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, e.Template.Unknown(err)
	}

	return buf, nil
}

func (t *TemplateList) ToJson() ([]byte, *e.Err) {

	if t == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(t)
	if err != nil {
		return nil, e.Template.Unknown(err)
	}

	return buf, nil
}
