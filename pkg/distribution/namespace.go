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
	"fmt"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
)

const (
	logNamespacePrefix     = "distribution:namespace"
	defaultNamespaceRam    = 4096
	defaultNamespaceRoutes = 1
)

type Namespace struct {
	context context.Context
	storage storage.Storage
}

type NM struct {
	Meta   struct{}
	Entity Namespace
}

func (n *NM) Set(Namespace) error {
	return nil
}

func (n *Namespace) List() (*types.NamespaceList, error) {

	log.V(logLevel).Debugf("%s:list:> get namespaces list", logNamespacePrefix)

	var list = types.NewNamespaceList()

	err := n.storage.List(n.context, n.storage.Collection().Namespace(), "", list, nil)

	if err != nil {
		log.Info(err.Error())
		log.V(logLevel).Error("%s:list:> get namespaces list err: %v", logNamespacePrefix, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> get namespaces list result: %d", logNamespacePrefix, len(list.Items))

	return list, nil
}

func (n *Namespace) Get(name string) (*types.Namespace, error) {

	log.V(logLevel).Infof("%s:get:> get namespace %s", logNamespacePrefix, name)

	if name == "" {
		return nil, errors.New(errors.ArgumentIsEmpty)
	}

	namespace := new(types.Namespace)

	err := n.storage.Get(n.context, n.storage.Collection().Namespace(), n.storage.Key().Namespace(name), &namespace, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> namespace by name `%s` not found", logNamespacePrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get namespace by name `%s` err: %v", logNamespacePrefix, name, err)
		return nil, err
	}

	return namespace, nil
}

func (n *Namespace) Create(opts *types.NamespaceCreateOptions) (*types.Namespace, error) {

	log.V(logLevel).Debugf("%s:create:> create Namespace %#v", logNamespacePrefix, opts)

	var ns = new(types.Namespace)
	ns.Meta.SetDefault()
	ns.Meta.Name = strings.ToLower(opts.Name)
	ns.Meta.Description = opts.Description
	ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", opts.Name, viper.GetString("domain.internal")))
	ns.SelfLink()

	if opts.Quotas != nil {
		ns.Spec.Quotas.RAM = opts.Quotas.RAM
		ns.Spec.Quotas.Routes = opts.Quotas.Routes
		ns.Spec.Quotas.Disabled = opts.Quotas.Disabled
	} else {
		ns.Spec.Quotas.Disabled = true
		ns.Spec.Quotas.RAM = defaultNamespaceRam
		ns.Spec.Quotas.Routes = defaultNamespaceRoutes
	}

	if err := n.storage.Put(n.context, n.storage.Collection().Namespace(), n.storage.Key().Namespace(ns.Meta.Name), ns, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert namespace err: %v", logNamespacePrefix, err)
		return nil, err
	}

	return ns, nil
}

func (n *Namespace) Update(namespace *types.Namespace, opts *types.NamespaceUpdateOptions) error {

	log.V(logLevel).Debugf("%s:update:> update Namespace %#v", logNamespacePrefix, namespace)

	if opts.Description != nil {
		namespace.Meta.Description = *opts.Description
	}

	if opts.Quotas != nil {
		namespace.Spec.Quotas.RAM = opts.Quotas.RAM
		namespace.Spec.Quotas.Routes = opts.Quotas.Routes
		namespace.Spec.Quotas.Disabled = opts.Quotas.Disabled
	}

	if err := n.storage.Set(n.context, n.storage.Collection().Namespace(),
		n.storage.Key().Namespace(namespace.Meta.Name), namespace, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> namespace update err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

func (n *Namespace) Remove(namespace *types.Namespace) error {

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logNamespacePrefix, namespace.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Namespace(), n.storage.Key().Namespace(namespace.Meta.Name)); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove namespace err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

// Watch namespace changes
func (n *Namespace) Watch(ch chan types.NamespaceEvent) error {

	log.V(logLevel).Debugf("%s:watch:> watch namespace", logNamespacePrefix)

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

				res := types.NamespaceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(types.Namespace)

				if err := json.Unmarshal(e.Data.([]byte), &obj); err != nil {
					log.Errorf("%s:watch:> parse json", logNamespacePrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Namespace(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewNamespaceModel(ctx context.Context, stg storage.Storage) *Namespace {
	return &Namespace{ctx, stg}
}
