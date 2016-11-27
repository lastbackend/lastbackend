package storage

import (
	"github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type IStorage interface {
	User() IUser
	Build() IBuild
	Image() IImage
	Project() IProject
	Service() IService
	Hook() IHook
}

type IUser interface {
	GetByID(string) (*model.User, *errors.Err)
	GetByUsername(string) (*model.User, *errors.Err)
	GetByEmail(string) (*model.User, *errors.Err)
	Insert(*model.User) (*model.User, *errors.Err)
}

type IBuild interface {
	// Get build and builds
	GetByID(string, string) (*model.Build, *errors.Err)
	GetByImage(string, string) (*model.BuildList, *errors.Err)
	// Insert and replace build
	Insert(*model.Build) (*model.Build, *errors.Err)
	Update(*model.Build) (*model.Build, *errors.Err)
}

type IImage interface {
	GetByID(string, string) (*model.Image, *errors.Err)
	GetByUser(string) (*model.ImageList, *errors.Err)
	GetByProject(string, string) (*model.ImageList, *errors.Err)
	GetByService(string, string) (*model.ImageList, *errors.Err)
	Insert(*model.Image) (*model.Image, *errors.Err)
	Update(*model.Image) (*model.Image, *errors.Err)
}

type IProject interface {
	GetByName(string, string) (*model.Project, *errors.Err)
	GetByID(string, string) (*model.Project, *errors.Err)
	GetByUser(string) (*model.ProjectList, *errors.Err)
	Insert(*model.Project) (*model.Project, *errors.Err)
	ExistByName(string, string) (bool, error)
	Update(*model.Project) (*model.Project, *errors.Err)
	Remove(string) *errors.Err
}

type IService interface {
	CheckExistsByName(string, string, string) (bool, error)
	GetByName(string, string) (*model.Service, *errors.Err)
	GetByID(string, string) (*model.Service, *errors.Err)
	GetByUser(string) (*model.ServiceList, *errors.Err)
	Insert(*model.Service) (*model.Service, *errors.Err)
	Update(*model.Service) (*model.Service, *errors.Err)
	Remove(string) *errors.Err
}

type IHook interface {
	GetByToken(string) (*model.Hook, *errors.Err)
	GetByUser(string) (*model.HookList, *errors.Err)
	GetByImage(string, string) (*model.HookList, *errors.Err)
	GetByService(string, string) (*model.HookList, *errors.Err)
	Insert(*model.Hook) (*model.Hook, *errors.Err)
	Delete(string, string) *errors.Err
}
