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

func (i *ImageStorage) GetByID(user, id string) (*model.Image, *e.Err) {

	var err error
	var image = new(model.Image)
	var user_filter = r.Row.Field("user").Eq(id)
	res, err := r.Table(ImageTable).Get(id).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, e.Image.NotFound(err)
	}
	res.One(image)

	defer res.Close()
	return image, nil
}

func (i *ImageStorage) GetByUser(id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)

	res, err := r.Table(ImageTable).Get(id).Run(i.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (i *ImageStorage) GetByProject(user, id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)
	var project_filter = r.Row.Field("project").Field("id").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)

	res, err := r.Table(ImageTable).Filter(project_filter).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (i *ImageStorage) GetByService(user, id string) (*model.ImageList, *e.Err) {

	var err error
	var images = new(model.ImageList)

	var project_filter = r.Row.Field("project").Field("id").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(ImageTable).Filter(project_filter).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

// Insert new image into storage
func (i *ImageStorage) Insert(image *model.Image) (*model.Image, *e.Err) {

	res, err := r.Table(ImageTable).Insert(image, r.InsertOpts{ReturnChanges: true}).RunWrite(i.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	image.ID = res.GeneratedKeys[0]

	return image, nil
}

// Update build model
func (i *ImageStorage) Update(image *model.Image) (*model.Image, *e.Err) {
	var user_filter = r.Row.Field("user").Eq(image.User)
	_, err := r.Table(ImageTable).Get(image.ID).Filter(user_filter).Replace(image, r.ReplaceOpts{ReturnChanges: true}).RunWrite(i.Session)
	if err != nil {
		return nil, e.Image.Unknown(err)
	}

	return image, nil
}

func newImageStorage(session *r.Session) *ImageStorage {
	r.TableCreate(ImageTable, r.TableCreateOpts{}).Run(session)
	s := new(ImageStorage)
	s.Session = session
	return s
}
