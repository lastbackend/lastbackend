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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const activityStorage string = "activities"

// Activity Service type for interface in interfaces folder
type ActivityStorage struct {
	IActivity
	log    logger.ILogger
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ActivityStorage) Insert(ctx context.Context, activity *types.Activity) error {
	return nil
}

func (s *ActivityStorage) ListProjectActivity(ctx context.Context, projectID string) ([]*types.Activity, error) {
	var activities = []*types.Activity{}
	return activities, nil
}

func (s *ActivityStorage) ListServiceActivity(ctx context.Context, serviceID string) ([]*types.Activity, error) {
	var activities = []*types.Activity{}
	return activities, nil
}

func (s *ActivityStorage) RemoveByProject(ctx context.Context, projectID string) error {
	return nil
}

func (s *ActivityStorage) RemoveByService(ctx context.Context, serviceID string) error {
	return nil
}

func newActivityStorage(config store.Config, log logger.ILogger) *ActivityStorage {
	s := new(ActivityStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
