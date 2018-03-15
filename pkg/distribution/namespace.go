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
	"fmt"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

type INamespace interface {
	List() ([]*types.Namespace, error)
	Get(name string) (*types.Namespace, error)
	Create(opts *types.NamespaceCreateOptions) (*types.Namespace, error)
	Update(namespace *types.Namespace, opts *types.NamespaceUpdateOptions) error
	Remove(id string) error
}

type Namespace struct {
	context context.Context
	storage storage.Storage
}

func (n *Namespace) List() ([]*types.Namespace, error) {

	log.V(logLevel).Debug("Namespace: List: list namespace")

	items, err := n.storage.Namespace().List(n.context)
	if err != nil {
		log.V(logLevel).Error("Namespace: List: list namespace err: %s", err.Error())
		return items, err
	}

	log.V(logLevel).Debugf("Namespace: List: list namespace result: %d", len(items))

	return items, nil
}

func (n *Namespace) Get(name string) (*types.Namespace, error) {

	log.V(logLevel).Debugf("Namespace: Get: get namespace %s", name)

	namespace, err := n.storage.Namespace().GetByName(n.context, name)
	if err != nil {
		log.V(logLevel).Errorf("Namespace: Get: get namespace by name `%s` err: %s", name, err.Error())
		return nil, err
	}
	if namespace == nil {
		log.V(logLevel).Warnf("Namespace: Get: namespace by name `%s` not found", name)
		return nil, nil
	}

	return namespace, nil
}

func (n *Namespace) Create(opts *types.NamespaceCreateOptions) (*types.Namespace, error) {

	log.V(logLevel).Debugf("Namespace: Create: create Namespace %#v", opts)

	var ns = new(types.Namespace)
	ns.Meta.SetDefault()
	ns.Meta.Name = opts.Name
	ns.Meta.Description = opts.Description
	ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s", opts.Name))
	ns.Meta.SelfLink = fmt.Sprintf("/namespace/%s", opts.Name)

	if opts.Quotas != nil {
		ns.Quotas.RAM = opts.Quotas.RAM
		ns.Quotas.Routes = opts.Quotas.Routes
		ns.Quotas.Disabled = opts.Quotas.Disabled
	} else {
		ns.Quotas.Disabled = true
		ns.Quotas.RAM = 4096
		ns.Quotas.Routes = 1
	}

	if err := n.storage.Namespace().Insert(n.context, ns); err != nil {
		log.V(logLevel).Errorf("Namespace: Create: insert namespace err: %s", err.Error())
		return nil, err
	}

	return ns, nil
}

func (n *Namespace) Update(namespace *types.Namespace, opts *types.NamespaceUpdateOptions) error {

	log.V(logLevel).Debugf("Namespace: Update: update Namespace %#v", namespace)

	if opts.Description != nil {
		namespace.Meta.Description = *opts.Description
	}

	if err := n.storage.Namespace().Update(n.context, namespace); err != nil {
		log.V(logLevel).Errorf("Namespace: update Namespace err: %s", err.Error())
		return err
	}

	return nil
}

func (n *Namespace) Remove(id string) error {

	log.V(logLevel).Debugf("Namespace: Remove: remove namespace %s", id)

	if err := n.storage.Namespace().Remove(n.context, id); err != nil {
		log.V(logLevel).Errorf("Namespace: remove namespace err: %s", err.Error())
		return err
	}

	return nil
}

func NewNamespaceModel(ctx context.Context, stg storage.Storage) INamespace {
	return &Namespace{ctx, stg}
}
