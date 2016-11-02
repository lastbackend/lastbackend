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
}

type IUser interface {
	GetByID(string) (*model.User, *errors.Err)
	Insert(*model.User) (*model.User, *errors.Err)
}

type IBuild interface {
	// Get build and builds
	GetByID(string) (*model.Build, *errors.Err)
	GetByImage(string) (*model.BuildList, *errors.Err)
	// Insert and replace build
	Insert(*model.Build) (*model.Build, *errors.Err)
	Replace(*model.Build) (*model.Build, *errors.Err)
}

type IImage interface {
	GetByID(string) (*model.Image, *errors.Err)
	GetByUser(string) (*model.ImageList, *errors.Err)
	GetByProject(string) (*model.ImageList, *errors.Err)
	GetByService(string) (*model.ImageList, *errors.Err)
	Insert(*model.Image) (*model.Image, *errors.Err)
	Replace(*model.Image) (*model.Image, *errors.Err)
}

type IProject interface {
	GetByID(string) (*model.Project, *errors.Err)
	GetByUser(string) (*model.ProjectList, *errors.Err)
	Insert(*model.Project) (*model.Project, *errors.Err)
	Replace(*model.Project) (*model.Project, *errors.Err)
}
