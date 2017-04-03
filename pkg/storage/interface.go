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
	Activity() IActivity
	Build() IBuild
	Hook() IHook
	Image() IImage
	Project() IProject
	Service() IService
	User() IUser
	Vendor() IVendor
	Volume() IVolume
}

type IActivity interface {
	Insert(*types.Activity) (*types.Activity, error)
	ListProjectActivity(string, string) (*types.ActivityList, error)
	ListServiceActivity(string, string) (*types.ActivityList, error)
	RemoveByProject(user, project string) error
	RemoveByService(user, service string) error
}

type IBuild interface {
	GetByID(string, string) (*types.Build, error)
	ListByImage(string, string) (*types.BuildList, error)
	Insert(*types.Build) (*types.Build, error)
	Update(*types.Build) (*types.Build, error)
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

type IProject interface {
	GetByID(username, id string) (*types.Project, error)
	GetByName(username, name string) (*types.Project, error)
	ListByUser(username string) (*types.ProjectList, error)
	Insert(username, name, description string) (*types.Project, error)
	Update(username string, project *types.Project) (*types.Project, error)
	Remove(username, name string) error
}

type IService interface {
	GetByID(username, project, id string) (*types.Service, error)
	GetByName(username, project, name string) (*types.Service, error)
	ListByProject(username, project string) (*types.ServiceList, error)
	Insert(username, project, name, description string, source *types.ServiceSource, config *types.ServiceConfig) (*types.Service, error)
	Update(username, project string, service *types.Service) (*types.Service, error)
	Remove(username, project, name string) error
	RemoveByProject(username, project string) error
}

type IImage interface {
	GetByID(string, string) (*types.Image, error)
	GetByUser(string) (*types.ImageList, error)
	ListByProject(string, string) (*types.ImageList, error)
	ListByService(string, string) (*types.ImageList, error)
	Insert(*types.Image) (*types.Image, error)
	Update(*types.Image) (*types.Image, error)
}

type IUser interface {
	GetByUsername(username string) (*types.User, error)
	GetByEmail(email string) (*types.User, error)
}

type IVendor interface {
	Insert(string, string, string, string, string, *oauth2.Token) error
	Get(string, string) (*types.Vendor, error)
	List(string) (map[string]*types.Vendor, error)
	Update(string, *types.Vendor) error
	Remove(string, string) error
}

type IVolume interface {
	GetByToken(string) (*types.Volume, error)
	ListByProject(string) (*types.VolumeList, error)
	Insert(*types.Volume) (*types.Volume, error)
	Remove(string) error
}
