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
	Volume() IVolume
}

type IUser interface {
	GetByID(string) (*model.User, *errors.Err)
	GetByUsername(string) (*model.User, *errors.Err)
	GetByEmail(string) (*model.User, *errors.Err)
	Insert(*model.User) (*model.User, *errors.Err)
}

type IBuild interface {
	GetByID(string, string) (*model.Build, *errors.Err)
	ListByImage(string, string) (*model.BuildList, *errors.Err)
	Insert(*model.Build) (*model.Build, *errors.Err)
	Update(*model.Build) (*model.Build, *errors.Err)
}

type IImage interface {
	GetByID(string, string) (*model.Image, *errors.Err)
	GetByUser(string) (*model.ImageList, *errors.Err)
	ListByProject(string, string) (*model.ImageList, *errors.Err)
	ListByService(string, string) (*model.ImageList, *errors.Err)
	Insert(*model.Image) (*model.Image, *errors.Err)
	Update(*model.Image) (*model.Image, *errors.Err)
}

type IProject interface {
	GetByNameOrID(string, string) (*model.Project, *errors.Err)
	GetByName(string, string) (*model.Project, *errors.Err)
	GetByID(string, string) (*model.Project, *errors.Err)
	ListByUser(string) (*model.ProjectList, *errors.Err)
	Insert(*model.Project) (*model.Project, *errors.Err)
	ExistByName(string, string) (bool, error)
	Update(*model.Project) (*model.Project, *errors.Err)
	Remove(string, string) *errors.Err
}

type IService interface {
	CheckExistsByName(string, string, string) (bool, error)
	GetByNameOrID(string, string, string) (*model.Service, *errors.Err)
	GetByName(string, string, string) (*model.Service, *errors.Err)
	GetByID(string, string, string) (*model.Service, *errors.Err)
	ListByUser(string, string) (*model.ServiceList, *errors.Err)
	ListByProject(string, string) (*model.ServiceList, *errors.Err)
	Insert(*model.Service) (*model.Service, *errors.Err)
	Update(*model.Service) (*model.Service, *errors.Err)
	Remove(string, string, string) *errors.Err
	RemoveByProject(string, string) *errors.Err
}

type IHook interface {
	GetByToken(string) (*model.Hook, *errors.Err)
	GetByUser(string) (*model.HookList, *errors.Err)
	ListByImage(string, string) (*model.HookList, *errors.Err)
	ListByService(string, string) (*model.HookList, *errors.Err)
	Insert(*model.Hook) (*model.Hook, *errors.Err)
	Delete(string, string) *errors.Err
}

type IVolume interface {
	GetByToken(string) (*model.Volume, *errors.Err)
	ListByProject(string) (*model.VolumeList, *errors.Err)
	Insert(*model.Volume) (*model.Volume, *errors.Err)
	Remove(string) *errors.Err
}
