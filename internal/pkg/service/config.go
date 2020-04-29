//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade config or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package service

import (
	"context"
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logConfigPrefix = "distribution:config"
)

type Config struct {
	context context.Context
	storage storage.IStorage
}

func (n *Config) Runtime() (*models.System, error) {

	log.Debugf("%s:get:> get config runtime info", logConfigPrefix)
	runtime, err := n.storage.Info(n.context, n.storage.Collection().Config(), "")
	if err != nil {
		log.Errorf("%s:get:> get runtime info error: %s", logConfigPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

func (n *Config) Get(namespace, name string) (*models.Config, error) {

	log.Debugf("%s:get:> get config by id %s/%s", logConfigPrefix, name)

	item := new(models.Config)

	err := n.storage.Get(n.context, n.storage.Collection().Config(), models.NewConfigSelfLink(namespace, name).String(), &item, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.Warnf("%s:get:> in namespace %s by name %s not found", logConfigPrefix, name)
			return nil, nil
		}

		log.Errorf("%s:get:> in namespace %s by name %s error: %s", logConfigPrefix, name, err)
		return nil, err
	}

	return item, nil
}

func (n *Config) List(filter string) (*models.ConfigList, error) {

	var f string

	log.Debugf("%s:list:> get configs list by namespace", logConfigPrefix)

	list := models.NewConfigList()
	if filter != models.EmptyString {
		f = n.storage.Filter().Config().ByNamespace(filter)
	}

	err := n.storage.List(n.context, n.storage.Collection().Config(), f, list, nil)
	if err != nil {
		log.Error("%s:list:> get configs list by namespace err: %s", logConfigPrefix, err)
		return list, err
	}

	log.Debugf("%s:list:> get configs list by namespace result: %d", logConfigPrefix, len(list.Items))

	return list, nil
}

func (n *Config) Create(namespace *models.Namespace, config *models.Config) (*models.Config, error) {

	log.Debugf("%s:create:> create config %#v", logConfigPrefix, config.Meta.Name)

	config.Meta.SetDefault()
	config.Meta.Namespace = namespace.Meta.Name
	config.SelfLink()

	if err := n.storage.Put(n.context, n.storage.Collection().Config(),
		config.SelfLink().String(), config, nil); err != nil {
		log.Errorf("%s:create:> insert config err: %v", logConfigPrefix, err)
		return nil, err
	}

	return config, nil
}

func (n *Config) Update(config *models.Config) (*models.Config, error) {

	log.Debugf("%s:update:> update config %s", logConfigPrefix, config.Meta.Name)

	if err := n.storage.Set(n.context, n.storage.Collection().Config(),
		config.SelfLink().String(), config, nil); err != nil {
		log.Errorf("%s:update:> update config err: %s", logConfigPrefix, err)
		return nil, err
	}

	return config, nil
}

func (n *Config) Remove(config *models.Config) error {

	log.Debugf("%s:remove:> remove config %#v", logConfigPrefix, config)

	if err := n.storage.Del(n.context, n.storage.Collection().Config(),
		config.SelfLink().String()); err != nil {
		log.Errorf("%s:remove:> remove config  err: %s", logConfigPrefix, err)
		return err
	}

	return nil
}

func (n *Config) Watch(ch chan models.ConfigEvent, rev *int64) error {

	log.Debugf("%s:watch:> watch config", logConfigPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-n.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := models.ConfigEvent{}
				res.Action = e.Action
				res.Name = e.Name

				config := new(models.Config)

				if err := json.Unmarshal(e.Data.([]byte), config); err != nil {
					log.Errorf("%s:> parse data err: %v", logConfigPrefix, err)
					continue
				}

				res.Data = config

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Config(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewConfigModel(ctx context.Context, stg storage.IStorage) *Config {
	return &Config{ctx, stg}
}
