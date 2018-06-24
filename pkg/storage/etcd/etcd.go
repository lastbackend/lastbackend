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

package etcd

import (
	"context"
	"regexp"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/store"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/spf13/viper"
)

const (
	logLevel  = 6
	logPrefix = "storage:etcd"
)

type Storage struct {
	client *client
}

type client struct {
	store store.Store
	dfunc store.DestroyFunc
}

func New() (*Storage, error) {

	log.Debug("Etcd: define storage")

	var (
		err    error
		s      = new(Storage)
		config *v3.Config
	)

	if err := viper.UnmarshalKey("etcd", &config); err != nil {
		log.Errorf("%s: error parsing etcd config: %v", logPrefix, err)
		return nil, err
	}

	s.client = new(client)

	if s.client.store, s.client.dfunc, err = v3.GetClient(config); err != nil {
		log.Errorf("%s: store initialize err: %v", err)
		return nil, err
	}

	return s, nil
}

func (s Storage) Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error {
	return s.client.store.Get(ctx, keyCreate(kind.String(), name), obj)
}

func (s Storage) List(ctx context.Context, kind types.Kind, query string, obj interface{}) error {
	return s.client.store.List(ctx, keyCreate(kind.String(), query), "", obj)
}

func (s Storage) Map(ctx context.Context, kind types.Kind, query string, obj interface{}) error {
	return s.client.store.Map(ctx, keyCreate(kind.String(), query), "", obj)
}

func (s Storage) Create(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {

	if opts == nil {
		opts = new(types.Opts)
	}

	return s.client.store.Create(ctx, keyCreate(kind.String(), name), obj, nil, opts.Ttl)
}

func (s Storage) Update(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {

	if opts == nil {
		opts = new(types.Opts)
	}

	return s.client.store.Update(ctx, keyCreate(kind.String(), name), obj, nil, opts.Ttl)
}

func (s Storage) Upsert(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {

	if opts == nil {
		opts = new(types.Opts)
	}

	return s.client.store.Upsert(ctx, keyCreate(kind.String(), name), obj, nil, opts.Ttl)
}

func (s Storage) Remove(ctx context.Context, kind types.Kind, name string) error {
	return s.client.store.Delete(ctx, keyCreate(kind.String(), name))
}

func (s Storage) Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error {

	log.V(logLevel).Debug("%s:> watch %s", logPrefix, kind.String())

	const filter = `\b.+\/(.+)\b`

	client, destroy, err := s.getClient()
	if err != nil {
		log.V(logLevel).Errorf("%s:> watch err: %v", logPrefix, err)
		return err
	}
	defer destroy()

	watcher, err := client.Watch(ctx, keyCreate(kind.String()), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:> watch err: %v", logPrefix, err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Debugf("%s:> the user interrupted watch", logPrefix)
			watcher.Stop()
			return nil
		case res := <-watcher.ResultChan():

			if res == nil {
				continue
			}

			if res.Type == store.STORAGEERROREVENT {
				err := res.Object.(error)
				log.Errorf("%s:> watch err: %v", logPrefix, err)
				return err
			}

			r, _ := regexp.Compile(filter)
			keys := r.FindStringSubmatch(res.Key)
			if len(keys) == 0 {
				continue
			}

			e := new(types.WatcherEvent)
			e.Action = res.Type

			match := strings.Split(res.Key, ":")

			if len(match) > 0 {
				e.Name = match[len(match)-1]
			} else {
				e.Name = keys[0]
			}

			if res.Type == store.STORAGEDELETEEVENT {
				e.Data = nil
				event <- e
				continue
			}

			e.Data = res.Object

			event <- e
		}
	}

	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s Storage) getClient() (store.Store, store.DestroyFunc, error) {
	return s.client.store, s.client.dfunc, nil
}
