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
	"strings"
)

const imageStorage string = "images"

// Namespace Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ImageStorage) Get(ctx context.Context, name string) (*types.Image, error) {

	image := new(types.Image)

	client, destroy, err := s.Client()
	if err != nil {
		return image, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, imageStorage, strings.Replace(name, "/", ":", -1), "meta")
	if err := client.Get(ctx, keyMeta, &image.Meta); err != nil {
		return image, err
	}

	keySource := s.util.Key(ctx, imageStorage, image.Name, "source")
	if err := client.Get(ctx, keySource, &image.Source); err != nil {
		return image, err
	}

	return image, nil
}

// Insert new image into storage
func (s *ImageStorage) Insert(ctx context.Context, image *types.Image) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, imageStorage, strings.Replace(image.Meta.Name, "/", ":", -1), "meta")
	if err := tx.Create(keyMeta, &image.Meta, 0); err != nil {
		return err
	}

	keySource := s.util.Key(ctx, imageStorage, strings.Replace(image.Meta.Name, "/", ":", -1), "source")
	if err := tx.Create(keySource, &image.Source, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Update build model
func (s *ImageStorage) Update(ctx context.Context, image *types.Image) error {
	return nil
}

func newImageStorage(config store.Config, util IUtil) *ImageStorage {
	s := new(ImageStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
