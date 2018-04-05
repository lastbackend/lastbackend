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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const (
	ingressStorage = "ingress"
)

// Ingress Service type for interface in interfaces folder
type IngressStorage struct {
	storage.Ingress
}

func (s *IngressStorage) List(ctx context.Context) (map[string]*types.Ingress, error) {

	log.V(logLevel).Debugf("storage:etcd:ingress:> get list ingresss")

	const filter = `\b.+` + ingressStorage + `\/(.+)\/(meta|status|spec)\b`

	ingresss := make(map[string]*types.Ingress, 0)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(ingressStorage)
	if err := client.MapList(ctx, key, filter, ingresss); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> get ingresss list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("storage:etcd:ingress:> get ingresss list result: %d", len(ingresss))

	return ingresss, nil
}

func (s *IngressStorage) Get(ctx context.Context, name string) (*types.Ingress, error) {

	log.V(logLevel).Debugf("storage:etcd:ingress:> get by id: %s", name)

	if len(name) == 0 {
		err := errors.New("ingress can not be empty")
		log.V(logLevel).Errorf("storage:etcd:ingress:> get ingress err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + ingressStorage + `\/.+\/(meta|status|spec)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	ingress := new(types.Ingress)
	key := keyDirCreate(ingressStorage, name)
	if err := client.Map(ctx, key, filter, ingress); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return nil, err
	}

	return ingress, nil
}

func (s *IngressStorage) GetSpec(ctx context.Context, ingress *types.Ingress) (*types.IngressSpec, error) {

	log.V(logLevel).Debugf("storage:etcd:ingress:> get ingress spec: %v", ingress)

	var (
		spec = new(types.IngressSpec)
	)

	spec.Routes = make(map[string]types.RouteSpec)

	if err := s.checkIngressExists(ctx, ingress); err != nil {
		return nil, err
	}

	const filterSpec= `\b.+` + ingressStorage + `\/(.+)\/spec\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	keySpec := keyDirCreate(ingressStorage)
	if err := client.Map(ctx, keySpec, filterSpec, spec); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("storage:etcd:ingress:> get ingress spec: err: %s", err.Error())
		return nil, err
	}

	return spec, nil
}

func (s *IngressStorage) Insert(ctx context.Context, ingress *types.Ingress) error {

	log.V(logLevel).Debugf("storage:etcd:ingress:> insert ingress: %#v", ingress)

	if err := s.checkIngressArgument(ingress); err != nil {
		return err
	}

	ingress.Meta.Created = time.Now()
	ingress.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(ingressStorage, ingress.Meta.Name, "meta")
	if err := tx.Create(keyMeta, &ingress.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> insert ingress err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(ingressStorage, ingress.Meta.Name, "status")
	if err := tx.Create(keyStatus, &ingress.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> insert ingress err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> insert ingress err: %s", err.Error())
		return err
	}

	return nil
}

func (s *IngressStorage) Update(ctx context.Context, ingress *types.Ingress) error {

	log.V(logLevel).Debugf("storage:etcd:ingress:> update ingress: %#v", ingress)

	if err := s.checkIngressExists(ctx, ingress); err != nil {
		return err
	}

	ingress.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(ingressStorage, ingress.Meta.Name, "meta")
	if err := tx.Update(keyMeta, &ingress.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress err: %s", err.Error())
		return err
	}

	return nil
}

func (s *IngressStorage) SetStatus(ctx context.Context, ingress *types.Ingress) error {

	log.V(logLevel).Debugf("storage:etcd:ingress:> update ingress status: %#v", ingress)

	if err := s.checkIngressExists(ctx, ingress); err != nil {
		return err
	}

	ingress.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress status err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(ingressStorage, ingress.Meta.Name, "status")
	if err := tx.Update(key, &ingress.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress status err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> update ingress status err: %s", err.Error())
		return err
	}

	return nil
}

func (s *IngressStorage) Remove(ctx context.Context, ingress *types.Ingress) error {

	log.V(logLevel).Debugf("storage:etcd:ingress:> remove ingress: %#v", ingress)

	if err := s.checkIngressExists(ctx, ingress); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(ingressStorage, ingress.Meta.Name)
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> remove ingress err: %s", err.Error())
		return err
	}

	return nil
}

func (s *IngressStorage) Watch(ctx context.Context, ingress chan *types.Ingress) error {

	log.V(logLevel).Debug("storage:etcd:ingress:> watch ingress")

	const filter = `\b.+` + ingressStorage + `\/(.+)\/(meta|status)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(ingressStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

		n, _ := s.Get(ctx, keys[1])
		if n == nil {
			return
		}

		if action == "PUT" {
			ingress <- n
			return
		}

		if action == "DELETE" {
			ingress <- n
			return
		}

		return
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> watch ingress err: %s", err.Error())
		return err
	}

	return nil
}

func (s *IngressStorage) WatchStatus(ctx context.Context, event chan *types.IngressStatusEvent) error {

	log.V(logLevel).Debug("storage:etcd:ingress:> watch ingress pod spec")

	const filter = `\b.+` + ingressStorage + `\/(.+)\/status\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(ingressStorage)
	cb := func(action, key string, val []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

		e := new(types.IngressStatusEvent)
		e.Event = action
		e.Name = keys[1]

		n, _ := s.Get(ctx, keys[1])
		if n == nil {
			return
		}

		e.Status = n.Status

		event <- e

		return
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> watch ingress err: %s", err.Error())
		return err
	}

	return nil
}


// Clear ingress storage
func (s *IngressStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:ingress:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, ingressStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:ingress:> clear err: %s", err.Error())
		return err
	}

	return nil
}

func newIngressStorage() *IngressStorage {
	s := new(IngressStorage)
	return s
}

func (s *IngressStorage) checkIngressArgument(ingress *types.Ingress) error {
	if ingress == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if ingress.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *IngressStorage) checkIngressExists(ctx context.Context, ingress *types.Ingress) error {

	if err := s.checkIngressArgument(ingress); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:ingress:> check ingress exists")

	if _, err := s.Get(ctx, ingress.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:ingress:> check ingress exists err: %s", err.Error())
		return err
	}

	return nil
}