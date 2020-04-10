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
// patents in process, and are protected by trade secret or copyright law.
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
	logNamespacePrefix  = "distribution:namespace"
	defaultNamespaceRam = "2GB"
	defaultNamespaceCPU = "200m"
)

type Namespace struct {
	context context.Context
	storage storage.IStorage
}

type NM struct {
	Meta   struct{}
	Entity Namespace
}

func (n *NM) Set(Namespace) error {
	return nil
}

func (n *Namespace) List() (*models.NamespaceList, error) {

	log.V(logLevel).Debugf("%s:list:> get namespaces list", logNamespacePrefix)

	var list = models.NewNamespaceList()

	err := n.storage.List(n.context, n.storage.Collection().Namespace(), "", list, nil)

	if err != nil {
		log.Info(err.Error())
		log.V(logLevel).Error("%s:list:> get namespaces list err: %v", logNamespacePrefix, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> get namespaces list result: %d", logNamespacePrefix, len(list.Items))

	return list, nil
}

func (n *Namespace) Get(name string) (*models.Namespace, error) {

	log.V(logLevel).Infof("%s:get:> get namespace %s", logNamespacePrefix, name)

	if name == "" {
		return nil, errors.New(errors.ArgumentIsEmpty)
	}

	namespace := new(models.Namespace)
	key := models.NewNamespaceSelfLink(name).String()

	err := n.storage.Get(n.context, n.storage.Collection().Namespace(), key, &namespace, nil)
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

func (n *Namespace) Create(ns *models.Namespace) (*models.Namespace, error) {

	log.V(logLevel).Debugf("%s:create:> create Namespace %#v", logNamespacePrefix, ns.Meta.Name)

	ns.Meta.SelfLink = *models.NewNamespaceSelfLink(ns.Meta.Name)

	if err := n.storage.Put(n.context, n.storage.Collection().Namespace(), ns.SelfLink().String(), ns, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert namespace err: %v", logNamespacePrefix, err)
		return nil, err
	}

	return ns, nil
}

func (n *Namespace) Update(namespace *models.Namespace) error {

	log.V(logLevel).Debugf("%s:update:> update Namespace %#v", logNamespacePrefix, namespace)

	if err := n.storage.Set(n.context, n.storage.Collection().Namespace(),
		namespace.SelfLink().String(), namespace, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> namespace update err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

func (n *Namespace) Remove(ns *models.Namespace) error {

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logNamespacePrefix, ns.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Namespace(), ns.SelfLink().String()); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove namespace err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

// Watch namespace changes
func (n *Namespace) Watch(ch chan models.NamespaceEvent) error {

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

				res := models.NamespaceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(models.Namespace)

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

func NewNamespaceModel(ctx context.Context, stg storage.IStorage) *Namespace {
	return &Namespace{ctx, stg}
}
