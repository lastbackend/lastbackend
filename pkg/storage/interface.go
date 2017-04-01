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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"golang.org/x/oauth2"
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
	Vendor() IVendor
}

type IUser interface {
	GetByUsername(string) (*types.User, error)
	GetByEmail(string) (*types.User, error)
}

type IBuild interface {
	GetByID(string, string) (*types.Build, error)
	ListByImage(string, string) (*types.BuildList, error)
	Insert(*types.Build) (*types.Build, error)
	Update(*types.Build) (*types.Build, error)
}

type IImage interface {
	GetByID(string, string) (*types.Image, error)
	GetByUser(string) (*types.ImageList, error)
	ListByProject(string, string) (*types.ImageList, error)
	ListByService(string, string) (*types.ImageList, error)
	Insert(*types.Image) (*types.Image, error)
	Update(*types.Image) (*types.Image, error)
}

type IProject interface {
	GetByName(string, string) (*types.Project, error)
	ListByUser(string) (*types.ProjectList, error)
	Insert(string, string, string) (*types.Project, error)
	Update(*types.Project) (*types.Project, error)
	Remove(string, string) error
}

type IService interface {
	CheckExistsByName(string, string) (bool, error)
	GetByNameOrID(string, string) (*types.Service, error)
	GetByName(string, string, string) (*types.Service, error)
	GetByID(string, string) (*types.Service, error)
	ListByUser(string, string) (*types.ServiceList, error)
	ListByProject(string, string) (*types.ServiceList, error)
	Insert(username, name, description string) (*types.Service, error)
	Update(*types.Service) (*types.Service, error)
	Remove(string, string, string) error
	RemoveByProject(string, string) error
}

type IHook interface {
	GetByToken(string) (*types.Hook, error)
	ListByUser(string) (*types.HookList, error)
	ListByImage(string, string) (*types.HookList, error)
	ListByService(string, string) (*types.HookList, error)
	Insert(*types.Hook) (*types.Hook, error)
	Remove(string) error
	RemoveByService(string) error
}

type IVolume interface {
	GetByToken(string) (*types.Volume, error)
	ListByProject(string) (*types.VolumeList, error)
	Insert(*types.Volume) (*types.Volume, error)
	Remove(string) error
}

type IActivity interface {
	Insert(*types.Activity) (*types.Activity, error)
	ListProjectActivity(string, string) (*types.ActivityList, error)
	ListServiceActivity(string, string) (*types.ActivityList, error)
	RemoveByProject(user, project string) error
	RemoveByService(user, service string) error
}

type IVendor interface {
	Insert(string, string, string, string, string, *oauth2.Token) error
	Get(string, string) (*types.Vendor, error)
	List(string) (map[string]*types.Vendor, error)
	Update(string, *types.Vendor) error
	Remove(string, string) error
}
