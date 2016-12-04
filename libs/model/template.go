package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"strings"
)

type TemplateList map[string][]string

type Template struct {
	Secrets                []v1.Secret                `json:"secrets,omitempty"`
	PersistentVolumes      []v1.PersistentVolume      `json:"persistent_volumes,omitempty"`
	PersistentVolumeClaims []v1.PersistentVolumeClaim `json:"persistent_volume_claims,omitempty"`
	ServiceAccounts        []v1.ServiceAccount        `json:"service_accounts,omitempty"`
	Services               []v1.Service               `json:"services,omitempty"`
	ReplicationControllers []v1.ReplicationController `json:"replication_controllers,omitempty"`
	Pods                   []v1.Pod                   `json:"pods,omitempty"`
	Deployments            []v1beta1.Deployment       `json:"deployments,omitempty"`
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

func (tpl *TemplateList) DrawTable() {
	t := table.New([]string{"Name", "Version"})

	for name, versions := range *tpl {
		t.AddRow(map[string]interface{}{
			"Name":    name,
			"Version": strings.Join(versions, "\r\n"),
		})
	}

	t.Print()
}
