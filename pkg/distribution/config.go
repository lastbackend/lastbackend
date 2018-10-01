//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package distribution

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logConfigPrefix = "distribution:config"
)

type Config struct {
	context context.Context
	storage storage.Storage
}

func (n *Config) Get(namespace, name string) (*types.Config, error) {

	log.V(logLevel).Debugf("%s:get:> get config by id %s/%s", logConfigPrefix, name)

	item := new(types.Config)

	err := n.storage.Get(n.context, n.storage.Collection().Config(), n.storage.Key().Config(namespace, name), &item, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logConfigPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %s", logConfigPrefix, name, err)
		return nil, err
	}

	return item, nil
}

func (n *Config) List(filter string) (*types.ConfigList, error) {

	var f string

	log.V(logLevel).Debugf("%s:list:> get configs list by namespace", logConfigPrefix)

	list := types.NewConfigList()
	if filter != types.EmptyString {
		f = n.storage.Filter().Config().ByNamespace(filter)
	}

	err := n.storage.List(n.context, n.storage.Collection().Config(), f, list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get configs list by namespace err: %s", logConfigPrefix, err)
		return list, err
	}

	log.V(logLevel).Debugf("%s:list:> get configs list by namespace result: %d", logConfigPrefix, len(list.Items))

	return list, nil
}

func (n *Config) Create(namespace *types.Namespace, config *types.Config) (*types.Config, error) {

	log.V(logLevel).Debugf("%s:create:> create config %#v", logConfigPrefix, config.Meta.Name)

	config.Meta.Namespace = namespace.Meta.Name
	if err := n.storage.Put(n.context, n.storage.Collection().Config(),
		n.storage.Key().Config(config.Meta.Namespace, config.Meta.Name), config, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert config err: %v", logConfigPrefix, err)
		return nil, err
	}

	return config, nil
}


func (n *Config) Update(config *types.Config) (*types.Config, error) {

	log.V(logLevel).Debugf("%s:update:> update config %s", logConfigPrefix, config.Meta.Name)


	if err := n.storage.Set(n.context, n.storage.Collection().Config(),
		n.storage.Key().Config(config.Meta.Namespace, config.Meta.Name), config, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> update config err: %s", logConfigPrefix, err)
		return nil, err
	}

	return config, nil
}

func (n *Config) Remove(config *types.Config) error {

	log.V(logLevel).Debugf("%s:remove:> remove config %#v", logConfigPrefix, config)

	if err := n.storage.Del(n.context, n.storage.Collection().Config(),
		n.storage.Key().Config(config.Meta.Namespace, config.Meta.Name)); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove config  err: %s", logConfigPrefix, err)
		return err
	}

	return nil
}

func (n *Config) Watch(ch chan types.ConfigEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch config", logConfigPrefix)

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

				res := types.ConfigEvent{}
				res.Action = e.Action
				res.Name = e.Name

				config := new(types.Config)

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

func NewConfigModel(ctx context.Context, stg storage.Storage) *Config {
	return &Config{ctx, stg}
}
