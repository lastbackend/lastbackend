package mock

import (
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

func (i *ImageMock) GetByID(user, id string) (*model.Image, error) {
	return nil, nil
}

func (i *ImageMock) ListByUser(id string) (*model.ImageList, error) {
	return nil, nil
}

func (i *ImageMock) ListByProject(user, id string) (*model.ImageList, error) {
	return nil, nil
}

func (i *ImageMock) ListByService(user, id string) (*model.ImageList, error) {
	return nil, nil
}

// Insert new image into storage
func (i *ImageMock) Insert(image *model.Image) (*model.Image, error) {
	return nil, nil
}

// Update image model
func (i *ImageMock) Update(image *model.Image) (*model.Image, error) {
	return nil, nil
}

func newImageMock(mock *r.Mock) *ImageMock {
	s := new(ImageMock)
	s.Mock = mock
	return s
}
