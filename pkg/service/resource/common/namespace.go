package common

import (
	"k8s.io/client-go/pkg/api"
)

type NamespaceQuery struct {
	namespaces []string
}

func NewSameNamespaceQuery(namespace string) *NamespaceQuery {
	return &NamespaceQuery{[]string{namespace}}
}

func (n *NamespaceQuery) ToRequestParam() string {

	if len(n.namespaces) == 1 {
		return n.namespaces[0]
	}

	return api.NamespaceAll
}

func (n *NamespaceQuery) Matches(namespace string) bool {

	if len(n.namespaces) == 0 {
		return true
	}

	for _, queryNamespace := range n.namespaces {
		if namespace == queryNamespace {
			return true
		}
	}

	return false
}
