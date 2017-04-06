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
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const imageStorage string = "images"

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ImageStorage) GetByID(ctx context.Context, id string) (*types.Image, error) {

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, imageStorage, id)
	meta := new(types.ImageMeta)
	if err := client.Get(ctx, key, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	image := new(types.Image)
	image.ID = meta.ID
	image.Name = meta.Name
	image.Description = meta.Description
	image.Labels = meta.Labels
	image.Created = meta.Created
	image.Updated = meta.Updated

	return image, nil
}

// Insert new image into storage
func (s *ImageStorage) Insert(ctx context.Context, source *types.ImageSource) (*types.Image, error) {
	return nil, nil
}

// Update build model
func (s *ImageStorage) Update(ctx context.Context, image *types.Image) (*types.Image, error) {
	return nil, nil
}

func newImageStorage(config store.Config, util IUtil) *ImageStorage {
	s := new(ImageStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
