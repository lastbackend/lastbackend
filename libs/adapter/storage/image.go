package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ImageTable string = "images"

type IImage interface {
	GetByID(string) (*model.Image, *error)
	GetByUser(string) (*model.ImageList, error)
	GetByProject(string) (*model.ImageList, error)
	GetByService(string) (*model.ImageList, error)
	Insert(*model.Image) (*model.Image, *error)
	Replace(*model.Image) (*model.Image, *error)
}

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
}

func (ImageStorage) GetByID(uuid string) (*model.Image, *e.Err) {

	var err error
	var image = new(model.Image)
	ctx := context.Get()

	res, err := r.Table(ImageTable).Get(uuid).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Image.NotFound(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

func (ImageStorage) GetByUser(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)
	ctx := context.Get()

	res, err := r.Table(ImageTable).Get(id).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (ImageStorage) GetByProject(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)
	ctx := context.Get()

	res, err := r.Table(ImageTable).Get(id).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (ImageStorage) GetByService(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)
	ctx := context.Get()

	res, err := r.Table(ImageTable).Get(id).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

// Insert new image into storage
func (ImageStorage) Insert(image *model.Image) (*model.Image, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(ImageTable).Insert(image, r.InsertOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

// Replace build model
func (ImageStorage) Replace(image *model.Image) (*model.Image, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(ImageTable).Get(image.ID).Replace(image, r.ReplaceOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}
