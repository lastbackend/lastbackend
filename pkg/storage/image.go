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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"golang.org/x/net/context"
	"time"
)

const ImageTable string = "images"

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ImageStorage) GetByID(user, id string) (*types.Image, error) {
	var (
		image = new(types.Image)
		key   = fmt.Sprintf("/%s/%s", ImageTable, id)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	meta := new(types.ImageMeta)
	if err := client.Get(ctx, key, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	image.ID = meta.ID
	image.Name = meta.Name
	image.Description = meta.Description
	image.Labels = meta.Labels
	image.Created = meta.Created
	image.Updated = meta.Updated

	return image, nil
}

// Insert new image into storage
func (s *ImageStorage) Insert(_ *types.ImageSource) (*types.Image, error) {
	var (
		id    = generator.GetUUIDV4()
		image = new(types.Image)
	)

	image.ID = id

	return image, nil
}

// Update build model
func (s *ImageStorage) Update(image *types.Image) (*types.Image, error) {
	return nil, nil
}

func newImageStorage(config store.Config) *ImageStorage {
	s := new(ImageStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
