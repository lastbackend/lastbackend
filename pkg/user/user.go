package user

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
	"github.com/lastbackend/lastbackend/utils"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"time"
)

type IUser interface {
	Create(username, email, password string) *User
	Get(name string) *User
}

type User struct {
	UUID     string    `json:"uuid,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Gravatar string    `json:"gravatar,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}

func Create(username, email string) (*User, error) {

	var ctx = context.Get()
	var tpr = new(v1beta1.ThirdPartyResource)
	tpr.APIVersion = "extensions/v1beta1"
	tpr.Kind = "ThirdPartyResource"
	tpr.ObjectMeta.Name = "users.api.lastbackend.com"
	tpr.Description = "A specification of a User to signup user in the system"
	tpr.Versions = append(tpr.Versions, v1beta1.APIVersion{Name: "v1"})

	ctx.K8S.Extensions().ThirdPartyResources().Create(tpr)

	user := v1.User{
		Spec: v1.UserSpec{
			UUID:     utils.GetUUIDV4(),
			Username: username,
			Email:    email,
			Gravatar: utils.GenerateGravatar(email),
			Updated:  time.Now(),
			Created:  time.Now(),
		},
	}

	user.APIVersion = "api.lastbackend.com/v1"
	user.Kind = "Users"
	tpr.ObjectMeta.Name = "users"

	userK8S, err := ctx.K8S.LB().Users().Create(&user)
	if err != nil {
		return nil, err
	}

	var u = new(User)
	u.UUID = userK8S.Spec.UUID
	u.Username = userK8S.Spec.Username
	u.Email = userK8S.Spec.Email
	u.Gravatar = userK8S.Spec.Gravatar
	u.Updated = userK8S.Spec.Updated
	u.Created = userK8S.Spec.Created

	return u, nil
}

func Get(name string) (*User, error) {

	var ctx = context.Get()

	userK8S, err := ctx.K8S.LB().Users().Get(name)
	if err != nil {
		return nil, err
	}

	var u = new(User)
	u.UUID = userK8S.Spec.UUID
	u.Username = userK8S.Spec.Username
	u.Email = userK8S.Spec.Email
	u.Gravatar = userK8S.Spec.Gravatar
	u.Updated = userK8S.Spec.Updated
	u.Created = userK8S.Spec.Created

	return u, nil
}
