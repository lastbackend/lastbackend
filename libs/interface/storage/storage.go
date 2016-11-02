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
	Hook() IHook
}

type IUser interface {
	GetByID(string) (*model.User, *errors.Err)
	Insert(*model.User) (*model.User, *errors.Err)
}

type IBuild interface {
	// Get build and builds
	GetByID(string, string) (*model.Build, *errors.Err)
	GetByImage(string, string) (*model.BuildList, *errors.Err)
	// Insert and replace build
	Insert(*model.Build) (*model.Build, *errors.Err)
	Replace(*model.Build) (*model.Build, *errors.Err)
}

type IImage interface {
	GetByID(string, string) (*model.Image, *errors.Err)
	GetByUser(string) (*model.ImageList, *errors.Err)
	GetByProject(string, string) (*model.ImageList, *errors.Err)
	GetByService(string, string) (*model.ImageList, *errors.Err)
	Insert(*model.Image) (*model.Image, *errors.Err)
	Replace(*model.Image) (*model.Image, *errors.Err)
}

type IProject interface {
	GetByID(string, string) (*model.Project, *errors.Err)
	GetByUser(string) (*model.ProjectList, *errors.Err)
	Insert(*model.Project) (*model.Project, *errors.Err)
	Replace(*model.Project) (*model.Project, *errors.Err)
}

type IHook interface {
	GetByToken(string) (*model.Hook, *errors.Err)
	GetByUser(string) (*model.HookList, *errors.Err)
	GetByImage(string, string) (*model.HookList, *errors.Err)
	GetByService(string, string) (*model.HookList, *errors.Err)
	Insert(*model.Hook) (*model.Hook, *errors.Err)
	Delete(string, string) *errors.Err
}
