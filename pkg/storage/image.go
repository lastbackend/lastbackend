//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package storage

import (
	"github.com/lastbackend/lastbackend/pkg/api/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const ImageTable string = "images"

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (i *ImageStorage) GetByID(user, id string) (*types.Image, error) {
	return nil, nil
}

func (i *ImageStorage) GetByUser(id string) (*types.ImageList, error) {
	return nil, nil
}

func (i *ImageStorage) ListByProject(user, id string) (*types.ImageList, error) {
	return nil, nil
}

func (i *ImageStorage) ListByService(user, id string) (*types.ImageList, error) {
	return nil, nil
}

// Insert new image into storage
func (i *ImageStorage) Insert(image *types.Image) (*types.Image, error) {
	return nil, nil
}

// Update build model
func (i *ImageStorage) Update(image *types.Image) (*types.Image, error) {
	return nil, nil
}

func newImageStorage(config store.Config) *ImageStorage {
	s := new(ImageStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
