package v1

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
)

type Build struct {
	unversioned.TypeMeta
	api.ObjectMeta
}

func (obj *Build) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
