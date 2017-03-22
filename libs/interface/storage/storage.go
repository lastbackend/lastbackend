//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package storage

import (
	"github.com/lastbackend/lastbackend/libs/model"
)

type IStorage interface {
	User() IUser
	Build() IBuild
	Image() IImage
	Project() IProject
	Service() IService
	Hook() IHook
	Volume() IVolume
	Activity() IActivity
}

type IUser interface {
	GetByID(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	GetByUsernameOrEmail(string) (*model.User, error)
	Insert(*model.User) (*model.User, error)
}

type IBuild interface {
	GetByID(string, string) (*model.Build, error)
	ListByImage(string, string) (*model.BuildList, error)
	Insert(*model.Build) (*model.Build, error)
	Update(*model.Build) (*model.Build, error)
}

type IImage interface {
	GetByID(string, string) (*model.Image, error)
	GetByUser(string) (*model.ImageList, error)
	ListByProject(string, string) (*model.ImageList, error)
	ListByService(string, string) (*model.ImageList, error)
	Insert(*model.Image) (*model.Image, error)
	Update(*model.Image) (*model.Image, error)
}

type IProject interface {
	GetByNameOrID(string, string) (*model.Project, error)
	GetByName(string, string) (*model.Project, error)
	GetByID(string, string) (*model.Project, error)
	ListByUser(string) (*model.ProjectList, error)
	Insert(*model.Project) (*model.Project, error)
	ExistByName(string, string) (bool, error)
	Update(*model.Project) (*model.Project, error)
	Remove(string, string) error
}

type IService interface {
	CheckExistsByName(string, string) (bool, error)
	GetByNameOrID(string, string) (*model.Service, error)
	GetByName(string, string) (*model.Service, error)
	GetByID(string, string) (*model.Service, error)
	ListByUser(string, string) (*model.ServiceList, error)
	ListByProject(string, string) (*model.ServiceList, error)
	Insert(*model.Service) (*model.Service, error)
	Update(*model.Service) (*model.Service, error)
	Remove(string, string) error
	RemoveByProject(string, string) error
}

type IHook interface {
	GetByToken(string) (*model.Hook, error)
	ListByUser(string) (*model.HookList, error)
	ListByImage(string, string) (*model.HookList, error)
	ListByService(string, string) (*model.HookList, error)
	Insert(*model.Hook) (*model.Hook, error)
	Remove(string) error
	RemoveByService(string) error
}

type IVolume interface {
	GetByToken(string) (*model.Volume, error)
	ListByProject(string) (*model.VolumeList, error)
	Insert(*model.Volume) (*model.Volume, error)
	Remove(string) error
}

type IActivity interface {
	Insert(*model.Activity) (*model.Activity, error)
	ListProjectActivity(string, string) (*model.ActivityList, error)
	ListServiceActivity(string, string) (*model.ActivityList, error)
	RemoveByProject(user, project string) error
	RemoveByService(user, service string) error
}
