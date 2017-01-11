package common

import (
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/apis/extensions"
)

type ObjectMeta struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Created   unversioned.Time  `json:"created,omitempty"`
}

type TypeMeta struct {
	Kind Kind `json:"kind,omitempty"`
}

type Spec struct {
	Replicas             int32                         `json:"replicas,omitempty"`
	Selector             *unversioned.LabelSelector    `json:"selector,omitempty"`
	Template             api.PodTemplateSpec           `json:"template"`
	Strategy             extensions.DeploymentStrategy `json:"strategy,omitempty"`
	MinReadySeconds      int32                         `json:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int32                        `json:"revisionHistoryLimit,omitempty"`
	Paused               bool                          `json:"paused,omitempty"`
	RollbackTo           *extensions.RollbackConfig    `json:"rollbackTo,omitempty"`
}

type Kind string

func NewObjectMeta(k8SObjectMeta api.ObjectMeta) ObjectMeta {
	return ObjectMeta{
		Name:      k8SObjectMeta.Name,
		Namespace: k8SObjectMeta.Namespace,
		Labels:    k8SObjectMeta.Labels,
		Created:   k8SObjectMeta.CreationTimestamp,
	}
}

type ListMeta struct {
	Total int `json:"total"`
}

func NewTypeMeta(kind Kind) TypeMeta {
	return TypeMeta{
		Kind: kind,
	}
}

func NewSpec(k8sSpec extensions.DeploymentSpec) Spec {
	return Spec{
		Replicas:             k8sSpec.Replicas,
		Selector:             k8sSpec.Selector,
		Template:             k8sSpec.Template,
		Strategy:             k8sSpec.Strategy,
		MinReadySeconds:      k8sSpec.MinReadySeconds,
		RevisionHistoryLimit: k8sSpec.RevisionHistoryLimit,
		Paused:               k8sSpec.Paused,
		RollbackTo:           k8sSpec.RollbackTo,
	}
}

func IsSelectorMatching(labelSelector map[string]string, labels map[string]string) bool {

	if len(labelSelector) == 0 {
		return false
	}

	for key, val := range labelSelector {
		if item, ok := labels[key]; !ok || item != val {
			return false
		}
	}

	return true
}
