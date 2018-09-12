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
)

const (
	logDiscoveryPrefix = "distribution:ingress"
)

type Discovery struct {
	context context.Context
	storage storage.Storage
}

func (n *Discovery) List() (*types.DiscoveryList, error) {
	list := types.NewDiscoveryList()

	if err := n.storage.List(n.context, n.storage.Collection().Discovery(), "", list, nil); err != nil {
		log.V(logLevel).Errorf("%s:list:> get ingress list err: %v", logDiscoveryPrefix, err)
		return nil, err
	}

	return list, nil
}

func (n *Discovery) Create(ingress *types.Discovery) error {

	log.V(logLevel).Debugf("%s:create:> create ingress in cluster", logDiscoveryPrefix)

	if err := n.storage.Put(n.context, n.storage.Collection().Discovery(),
		n.storage.Key().Discovery(ingress.SelfLink()), ingress, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert ingress err: %v", logDiscoveryPrefix, err)
		return err
	}

	return nil
}

func (n *Discovery) Get(name string) (*types.Discovery, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logDiscoveryPrefix, name)

	ingress := new(types.Discovery)
	err := n.storage.Get(n.context, n.storage.Collection().Discovery(), n.storage.Key().Discovery(name), ingress, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: ingress %s not found", logDiscoveryPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get ingress `%s` err: %v", logDiscoveryPrefix, name, err)
		return nil, err
	}

	return ingress, nil
}

func (n *Discovery) Remove(ingress *types.Discovery) error {

	log.V(logLevel).Debugf("%s:remove:> remove ingress %s", logDiscoveryPrefix, ingress.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Discovery(), n.storage.Key().Discovery(ingress.SelfLink())); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove ingress err: %v", logDiscoveryPrefix, err)
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

				ingress := new(types.Discovery)

				if err := json.Unmarshal(e.Data.([]byte), ingress); err != nil {
					log.Errorf("%s:> parse data err: %v", logDiscoveryPrefix, err)
					continue
				}

				res.Data = ingress

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Discovery(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewDiscoveryModel(ctx context.Context, stg storage.Storage) *Discovery {
	return &Discovery{ctx, stg}
}

