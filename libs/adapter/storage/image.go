package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ImageTable string = "images"

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	Session *r.Session
	storage.IImage
}

func (s *ImageStorage) GetByID(uuid string) (*model.Image, *e.Err) {

	var err error
	var image = new(model.Image)

	res, err := r.Table(ImageTable).Get(uuid).Run(s.Session)
	if err != nil {
		return nil, e.Image.NotFound(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

func (s *ImageStorage) GetByUser(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)

	res, err := r.Table(ImageTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (s *ImageStorage) GetByProject(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)

	res, err := r.Table(ImageTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (s *ImageStorage) GetByService(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)

	res, err := r.Table(ImageTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

// Insert new image into storage
func (s *ImageStorage) Insert(image *model.Image) (*model.Image, *e.Err) {

	res, err := r.Table(ImageTable).Insert(image, r.InsertOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

// Replace build model
func (s *ImageStorage) Replace(image *model.Image) (*model.Image, *e.Err) {

	res, err := r.Table(ImageTable).Get(image.ID).Replace(image, r.ReplaceOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

func newImageStorage(session *r.Session) *ImageStorage {
	s := new(ImageStorage)
	s.Session = session
	return s
}
