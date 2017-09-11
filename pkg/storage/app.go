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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/pkg/errors"
	"time"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const appStorage = "apps"

// App Service type for interface in interfaces folder
type AppStorage struct {
	IApp
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get app by name
func (s *AppStorage) GetByName(ctx context.Context, name string) (*types.App, error) {

	log.V(logLevel).Debugf("Storage: App: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: App: get app err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: App: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	app := new(types.App)
	keyMeta := keyCreate(appStorage, name, "meta")
	if err := client.Get(ctx, keyMeta, &app.Meta); err != nil {
		log.V(logLevel).Errorf("Storage: App: get app `%s` meta err: %s", name, err.Error())
		return nil, err
	}

	return app, nil
}

// List projects
func (s *AppStorage) List(ctx context.Context) ([]*types.App, error) {

	log.V(logLevel).Debug("Storage: App: get app list")

	const filter = `\b(.+)` + appStorage + `\/.+\/(meta)\b`

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: App: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	apps := []*types.App{}
	keyApps := keyCreate(appStorage)
	if err := client.List(ctx, keyApps, filter, &apps); err != nil {
		log.V(logLevel).Errorf("Storage: App: get apps list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Storage: App: get app list result: %d", len(apps))

	return apps, nil
}

// Insert new app into storage
func (s *AppStorage) Insert(ctx context.Context, app *types.App) error {

	log.V(logLevel).Debug("Storage: App: insert app: %#v", app)

	if app == nil {
		err := errors.New("app can not be nil")
		log.V(logLevel).Errorf("Storage: App: insert app err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: App: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(appStorage, app.Meta.Name, "meta")
	if err := client.Create(ctx, keyMeta, app.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: App: insert app err: %s", err.Error())
		return err
	}

	return nil
}

// Update app model
func (s *AppStorage) Update(ctx context.Context, app *types.App) error {

	log.V(logLevel).Debugf("Storage: App: update app: %#v", app)

	if app == nil {
		err := errors.New("app can not be nil")
		log.V(logLevel).Errorf("Storage: App: update app err: %s", err.Error())
		return err
	}

	app.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: App: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	meta := types.Meta{}
	meta = app.Meta
	meta.Updated = time.Now()

	keyMeta := keyCreate(appStorage, app.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: App: update app meta err: %s", err.Error())
		return err
	}

	return nil
}

// Remove app model
func (s *AppStorage) Remove(ctx context.Context, name string) error {

	log.V(logLevel).Debugf("Storage: App: remove app: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: App: remove app err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: App: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyApp := keyCreate(appStorage, name)
	if err := client.DeleteDir(ctx, keyApp); err != nil {
		log.V(logLevel).Errorf("Storage: App: remove app `%s` err: %s", name, err.Error())
		return err
	}

	return nil
}

func newAppStorage(config store.Config) *AppStorage {
	s := new(AppStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
