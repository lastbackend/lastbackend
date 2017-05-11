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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

const buildStorage string = "builds"

// Service Build type for interface in interfaces folder
type BuildStorage struct {
	IBuild
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get build model by id
func (s *BuildStorage) GetByID(ctx context.Context, imageID, id string) (*types.Build, error) {
	build := &types.Build{}
	return build, nil
}

// Get builds by image
func (s *BuildStorage) ListByImage(ctx context.Context, id string) ([]*types.Build, error) {
	builds := []*types.Build{}
	return builds, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(ctx context.Context, imageName string, build *types.Build) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	iname := strings.Replace(imageName, "/", ":", -1)

	keyImageMeta := s.util.Key(ctx, imageStorage, iname, "meta")
	imeta := new(types.ImageMeta)
	if err := client.Get(ctx, keyImageMeta, imeta); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	imeta.Builds++
	if err := tx.Update(keyImageMeta, imeta, 0); err != nil {
		return err
	}

	keyMeta := s.util.Key(ctx, imageStorage, iname, buildStorage, fmt.Sprintf("%d", imeta.Builds))

	if err := tx.Create(keyMeta, build, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func newBuildStorage(config store.Config, util IUtil) *BuildStorage {
	s := new(BuildStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
