package cri

import "github.com/lastbackend/lastbackend/libs/model"

type ICRI interface {
	Container()
}

type IContainer interface {
	GetByID(string) (*model.Container)
	List(string) (*model.ContainerList)
	Create() (*model.Container)
	Start() (*model.Container)
	Restart() (*model.Container)
	Stop() (*model.Container)
	Logs() error
	Proxy () error
	Remove() error
}

type IImage interface {
	GetByID(string) (*model.Image)
	List(string) (*model.ImageList)
	Pull(string) (*model.Image)
	Push(string) (*model.Image)
	Remove(string) error
}

type INetwork interface {
	GetByID(string) (*model.Network)
	List(list model.NetworkList) (*model.NetworkList)
	Create(network *model.Network) (*model.ImageList)
	Remove(string) error
}
