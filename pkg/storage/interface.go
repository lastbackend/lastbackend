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
	"golang.org/x/net/context"
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
	Insert(ctx context.Context, activity *types.Activity) (*types.Activity, error)
	ListProjectActivity(ctx context.Context, username, id string) (*types.ActivityList, error)
	ListServiceActivity(ctx context.Context, username, id string) (*types.ActivityList, error)
	RemoveByProject(ctx context.Context, username, id string) error
	RemoveByService(ctx context.Context, username, id string) error
}

type IBuild interface {
	GetByID(ctx context.Context, username, id string) (*types.Build, error)
	ListByImage(ctx context.Context, username, id string) (*types.BuildList, error)
	Insert(ctx context.Context, build *types.Build) (*types.Build, error)
	Update(ctx context.Context, build *types.Build) (*types.Build, error)
}

type IHook interface {
	GetByToken(ctx context.Context, token string) (*types.Hook, error)
	ListByUser(ctx context.Context, username string) (*types.HookList, error)
	ListByImage(ctx context.Context, username, id string) (*types.HookList, error)
	ListByService(ctx context.Context, username, id string) (*types.HookList, error)
	Insert(ctx context.Context, hook *types.Hook) (*types.Hook, error)
	Remove(ctx context.Context, id string) error
	RemoveByService(ctx context.Context, id string) error
}

type IProject interface {
	GetByID(ctx context.Context, username, id string) (*types.Project, error)
	GetByName(ctx context.Context, username, name string) (*types.Project, error)
	ListByUser(ctx context.Context, username string) (*types.ProjectList, error)
	Insert(ctx context.Context, username, name, description string) (*types.Project, error)
	Update(ctx context.Context, username string, project *types.Project) (*types.Project, error)
	Remove(ctx context.Context, username, name string) error
}

type IService interface {
	GetByID(ctx context.Context, username, project, id string) (*types.Service, error)
	GetByName(ctx context.Context, username, project, name string) (*types.Service, error)
	ListByProject(ctx context.Context, username, project string) (*types.ServiceList, error)
	Insert(ctx context.Context, username, project, name, description, image string, config *types.ServiceConfig) (*types.Service, error)
	Update(ctx context.Context, username, project string, service *types.Service) (*types.Service, error)
	Remove(ctx context.Context, username, project, name string) error
	RemoveByProject(ctx context.Context, username, project string) error
}

type IImage interface {
	GetByID(ctx context.Context, username, id string) (*types.Image, error)
	Insert(ctx context.Context, source *types.ImageSource) (*types.Image, error)
	Update(ctx context.Context, image *types.Image) (*types.Image, error)
}

type IUser interface {
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
}

type IVendor interface {
	Insert(ctx context.Context, username, vendorUsername, vendorName, vendorHost, serviceID string, token *oauth2.Token) error
	Get(ctx context.Context, username, vendorName string) (*types.Vendor, error)
	List(ctx context.Context, username string) (map[string]*types.Vendor, error)
	Update(ctx context.Context, username string, vendor *types.Vendor) error
	Remove(ctx context.Context, username string, vendorName string) error
}

type IVolume interface {
	GetByToken(ctx context.Context, token string) (*types.Volume, error)
	ListByProject(ctx context.Context, project string) (*types.VolumeList, error)
	Insert(ctx context.Context, volume *types.Volume) (*types.Volume, error)
	Remove(ctx context.Context, id string) error
}
