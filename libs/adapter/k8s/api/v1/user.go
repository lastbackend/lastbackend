package v1

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"time"
)

type User struct {
	unversioned.TypeMeta `json:",inline"`
	// metadata for Build
	api.ObjectMeta `json:"metadata,omitempty"`
	// user Spec
	Spec UserSpec `json:"spec,omitempty"`
}

// UserSpec has the information to represent a user and also additional
// information about a user
type UserSpec struct {
	// User id, generate automatically
	UUID string `json:"id,omitempty"`
	// User username, need to set
	Username string `json:"username,omitempty"`
	// User email, need to set
	Email string `json:"email,omitempty"`
	// User gravatar hash, generate automatically
	Gravatar string `json:"gravatar,omitempty"`
	// User status updated time
	Updated time.Time `json:"updated,omitempty"`
	// User status create time
	Created time.Time `json:"created,omitempty"`
}

func (obj *User) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
