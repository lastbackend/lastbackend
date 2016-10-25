package v1

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
)

type User struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 UserSpec `json:"spec,omitempty"`
}

// UserSpec has the information to represent a user and also additional
// information about a user
type UserSpec struct {
	UUID     string            `json:"id,omitempty"`
	Username string            `json:"username,omitempty"`
	Email    string            `json:"email,omitempty"`
	Gravatar string            `json:"gravatar,omitempty"`
	Updated  *unversioned.Time `json:"updated,omitempty"`
	Created  *unversioned.Time `json:"created,omitempty"`
}

func (obj *User) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
