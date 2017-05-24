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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

const imageStorage string = "images"

// Image Service type for interface in interfaces folder
type ImageStorage struct {
	IImage
	log    logger.ILogger
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ImageStorage) Get(ctx context.Context, name string) (*types.Image, error) {

	s.log.V(debugLevel).Debugf("Storage: Image: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Image: get namespace err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Image: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	image := new(types.Image)
	keyMeta := s.util.Key(ctx, imageStorage, strings.Replace(name, "/", ":", -1), "meta")
	if err := client.Get(ctx, keyMeta, &image.Meta); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Pod: get image meta err: %s", err.Error())
		return nil, err
	}

	keySource := s.util.Key(ctx, imageStorage, image.Name, "source")
	if err := client.Get(ctx, keySource, &image.Source); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Pod: get image source err: %s", err.Error())
		return nil, err
	}

	return image, nil
}

// Insert new image into storage
func (s *ImageStorage) Insert(ctx context.Context, image *types.Image) error {

	s.log.V(debugLevel).Debugf("Storage: Image: insert image: %#v", image)

	if image == nil {
		err := errors.New("image can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Image: insert image err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Image: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, imageStorage, strings.Replace(image.Meta.Name, "/", ":", -1), "meta")
	if err := tx.Create(keyMeta, &image.Meta, 0); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Image: create image meta err: %s", err.Error())
		return err
	}

	keySource := s.util.Key(ctx, imageStorage, strings.Replace(image.Meta.Name, "/", ":", -1), "source")
	if err := tx.Create(keySource, &image.Source, 0); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Image: create image source err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Image: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

// Update build model
func (s *ImageStorage) Update(ctx context.Context, image *types.Image) error {

	s.log.V(debugLevel).Debugf("Storage: Image: update image: %#v", image)

	if image == nil {
		err := errors.New("image can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Image: update image err: %s", err.Error())
		return err
	}

	return nil
}

func newImageStorage(config store.Config, log logger.ILogger, util IUtil) *ImageStorage {
	s := new(ImageStorage)
	s.log = log
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
