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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"github.com/satori/go.uuid"
	"time"
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
	return nil, nil
}

// Get builds by image
func (s *BuildStorage) ListByImage(ctx context.Context, id string) (*types.BuildList, error) {
	return nil, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(ctx context.Context, imageID string, source *types.BuildSource) (*types.Build, error) {

	build := new(types.Build)
	build.Meta.ID = uuid.NewV4().String()
	build.Status = types.BuildStatus{
		Step:    types.BuildStepCreate,
		Updated: time.Now(),
	}
	build.Source = *source
	build.Meta.Created = time.Now()
	build.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyImageMeta := s.util.Key(ctx, imageStorage, imageID, "meta")
	imeta := new(types.ImageMeta)
	if err := client.Get(ctx, keyImageMeta, imeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	keyImageMeta = s.util.Key(ctx, imageStorage, imageID, "meta")
	imeta.Builds++
	if err := tx.Update(keyImageMeta, imeta, 0); err != nil {
		if err.Error() == store.ErrKeyExists {
			return nil, nil
		}
		return nil, err
	}

	keyMeta := s.util.Key(ctx, imageStorage, imageID,
		buildStorage, fmt.Sprintf("%s", build.Meta.Created.Unix()))

	if err := tx.Create(keyMeta, build, 0); err != nil {
		if err.Error() == store.ErrKeyExists {
			return nil, nil
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return build, nil
}

func newBuildStorage(config store.Config, util IUtil) *BuildStorage {
	s := new(BuildStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
