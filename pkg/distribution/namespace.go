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
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	logNamespacePrefix     = "distribution:namespace"
	defaultNamespaceRam    = 4096
	defaultNamespaceRoutes = 1
)

type INamespace interface {
	List() (map[string]*types.Namespace, error)
	Get(name string) (*types.Namespace, error)
	Create(opts *types.NamespaceCreateOptions) (*types.Namespace, error)
	Update(namespace *types.Namespace, opts *types.NamespaceUpdateOptions) error
	Remove(namespace *types.Namespace) error
}

type Namespace struct {
	context context.Context
	storage storage.Storage
}

func (n *Namespace) List() (map[string]*types.Namespace, error) {

	log.V(logLevel).Debug("%s:list:> get namespaces list", logNamespacePrefix)

	items, err := n.storage.Namespace().List(n.context)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get namespaces list err: %s", logNamespacePrefix, err.Error())
		return items, err
	}

	log.V(logLevel).Debugf("%s:list:> get namespaces list result: %d", logNamespacePrefix, len(items))

	return items, nil
}

func (n *Namespace) Get(name string) (*types.Namespace, error) {

	log.V(logLevel).Infof("%s:get:> get namespace %s", logNamespacePrefix, name)

	if name == "" {
		return nil, errors.New(errors.ArgumentIsEmpty)
	}

	namespace, err := n.storage.Namespace().Get(n.context, name)
	if err != nil {
		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> namespace by name `%s` not found", logNamespacePrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get namespace by name `%s` err: %s", logNamespacePrefix, name, err.Error())
		return nil, err
	}

	return namespace, nil
}

func (n *Namespace) Create(opts *types.NamespaceCreateOptions) (*types.Namespace, error) {

	log.V(logLevel).Debugf("%s:create:> create Namespace %#v", logNamespacePrefix, opts)

	var ns = new(types.Namespace)
	ns.Meta.SetDefault()
	ns.Meta.Name = opts.Name
	ns.Meta.Description = opts.Description
	ns.Meta.Endpoint = strings.ToLower(opts.Name)
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

	if err := n.storage.Namespace().Insert(n.context, ns); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert namespace err: %s", logNamespacePrefix, err.Error())
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

	if err := n.storage.Namespace().Update(n.context, namespace); err != nil {
		log.V(logLevel).Errorf("%s:update:> namespace update err: %s", logNamespacePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Namespace) Remove(namespace *types.Namespace) error {

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logNamespacePrefix, namespace.Meta.Name)

	if err := n.storage.Namespace().Remove(n.context, namespace); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove namespace err: %s", logNamespacePrefix, err.Error())
		return err
	}

	return nil
}

func NewNamespaceModel(ctx context.Context, stg storage.Storage) INamespace {
	return &Namespace{ctx, stg}
}
