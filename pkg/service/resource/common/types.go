package common

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
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
