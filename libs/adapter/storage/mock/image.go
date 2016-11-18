package mock

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ImageTable string = "images"

// Project Service type for interface in interfaces folder
type ImageMock struct {
	Mock *r.Mock
	storage.IImage
}

func (i *ImageMock) GetByID(user, id string) (*model.Image, *e.Err) {
	return nil, nil
}

func (i *ImageMock) GetByUser(id string) (*model.ImageList, *e.Err) {
	return nil, nil
}

func (i *ImageMock) GetByProject(user, id string) (*model.ImageList, *e.Err) {
	return nil, nil
}

func (i *ImageMock) GetByService(user, id string) (*model.ImageList, *e.Err) {
	return nil, nil
}

// Insert new image into storage
func (i *ImageMock) Insert(image *model.Image) (*model.Image, *e.Err) {
	return nil, nil
}

// Update build model
func (i *ImageMock) Update(image *model.Image) (*model.Image, *e.Err) {
	return nil, nil
}

func newImageMock(mock *r.Mock) *ImageMock {
	s := new(ImageMock)
	s.Mock = mock
	return s
}
