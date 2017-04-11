package v1

import "github.com/lastbackend/lastbackend/pkg/apis/types"

func NewNamespace(obj *types.Namespace) *Namespace {
	return New(obj)
}

func NewNamespaceList(obj *types.NamespaceList) *NamespaceList {
	return NewList(obj)
}
