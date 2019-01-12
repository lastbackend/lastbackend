//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package distribution

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"time"
)

const (
	logDiscoveryPrefix = "distribution:discovery"
	ttlDiscovery = uint64(5*time.Second)
)

type Discovery struct {
	context context.Context
	storage storage.Storage
}

func (n *Discovery) List() (*types.DiscoveryList, error) {
	list := types.NewDiscoveryList()

	if err := n.storage.List(n.context, n.storage.Collection().Discovery().Info(), "", list, nil); err != nil {
		log.V(logLevel).Errorf("%s:list:> get discovery list err: %v", logDiscoveryPrefix, err)
		return nil, err
	}

	return list, nil
}

func (n *Discovery) Put(discovery *types.Discovery) error {

	log.V(logLevel).Debugf("%s:create:> create discovery in cluster", logDiscoveryPrefix)

	if err := n.storage.Put(n.context, n.storage.Collection().Discovery().Info(),
		n.storage.Key().Discovery(discovery.SelfLink()), discovery, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert discovery err: %v", logDiscoveryPrefix, err)
		return err
	}

	opts := storage.GetOpts()
	opts.Ttl = ttlDiscovery

	if err := n.storage.Put(n.context, n.storage.Collection().Discovery().Status(),
		n.storage.Key().Discovery(discovery.SelfLink()), discovery.Status, opts); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert discovery status err: %v", logDiscoveryPrefix, err)
		return err
	}

	return nil
}

func (n *Discovery) Get(name string) (*types.Discovery, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logDiscoveryPrefix, name)

	discovery := new(types.Discovery)
	err := n.storage.Get(n.context, n.storage.Collection().Discovery().Info(), n.storage.Key().Discovery(name), discovery, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: discovery %s not found", logDiscoveryPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get discovery `%s` err: %v", logDiscoveryPrefix, name, err)
		return nil, err
	}

	return discovery, nil
}

func (n *Discovery) Set(discovery *types.Discovery) error {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logDiscoveryPrefix, discovery.Meta.Name)

	opts := storage.GetOpts()
	opts.Force = true

	err := n.storage.Set(n.context, n.storage.Collection().Discovery().Info(),
		n.storage.Key().Discovery(discovery.Meta.Name), discovery, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:get:> set discovery `%s` err: %v", logDiscoveryPrefix, discovery.Meta.Name, err)
		return err
	}

	if err := n.storage.Set(n.context, n.storage.Collection().Discovery().Status(),
		n.storage.Key().Discovery(discovery.Meta.Name), discovery.Status, nil); err != nil {
		log.V(logLevel).Debugf("%s:get:> set discovery status `%s` err: %v", logDiscoveryPrefix, discovery.Meta.Name, err)
		return err
	}

	return nil
}

func (n *Discovery) SetOnline(discovery *types.Discovery) error {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logDiscoveryPrefix, discovery.Meta.Name)

	opts := storage.GetOpts()
	opts.Force = true

	err := n.storage.Set(n.context, n.storage.Collection().Discovery().Status(),
		n.storage.Key().Discovery(discovery.Meta.Name), discovery.Status, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:get:> set discovery `%s` err: %v", logDiscoveryPrefix, discovery.Meta.Name, err)
		return err
	}

	return nil
}

func (n *Discovery) Remove(discovery *types.Discovery) error {

	log.V(logLevel).Debugf("%s:remove:> remove discovery %s", logDiscoveryPrefix, discovery.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Discovery().Info(), n.storage.Key().Discovery(discovery.SelfLink())); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove discovery err: %v", logDiscoveryPrefix, err)
		return err
	}

	return nil
}


func (n *Discovery) Watch(ch chan types.DiscoveryEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch routes", logDiscoveryPrefix)

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

				res := types.DiscoveryEvent{}
				res.Action = e.Action
				res.Name = e.Name

				discovery := new(types.Discovery)

				if err := json.Unmarshal(e.Data.([]byte), discovery); err != nil {
					log.Errorf("%s:> parse data err: %v", logDiscoveryPrefix, err)
					continue
				}

				res.Data = discovery

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Discovery().Info(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (n *Discovery) WatchOnline(ch chan types.DiscoveryStatusEvent) error {

	log.V(logLevel).Debugf("%s:watch:> watch routes", logDiscoveryPrefix)

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

				res := types.DiscoveryStatusEvent{}
				res.Action = e.Action
				res.Name = e.Name

				discovery := new(types.DiscoveryStatus)

				if err := json.Unmarshal(e.Data.([]byte), discovery); err != nil {
					log.Errorf("%s:> parse data err: %v", logDiscoveryPrefix, err)
					continue
				}

				res.Data = discovery

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Discovery().Status(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewDiscoveryModel(ctx context.Context, stg storage.Storage) *Discovery {
	return &Discovery{ctx, stg}
}

