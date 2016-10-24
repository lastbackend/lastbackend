package v1

import (
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/api"
)

type Account struct {
	unversioned.TypeMeta
	api.ObjectMeta
}

func (obj *Account) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
