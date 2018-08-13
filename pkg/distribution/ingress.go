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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logIngressPrefix = "distribution:ingress"
)

type Ingress struct {
	context context.Context
	storage storage.Storage
}

func (n *Ingress) List() (*types.IngressList, error) {
	list := types.NewIngressList()

	if err := n.storage.Map(n.context, n.storage.Collection().Ingress(), "", list, nil); err != nil {
		log.Debugf("%s:list:> get ingress list err: %v", logIngressPrefix, err)
		return nil, err
	}

	return list, nil
}

func (n *Ingress) Create(opts *types.IngressCreateOptions) (*types.Ingress, error) {

	log.Debugf("%s:create:> create ingress in cluster", logIngressPrefix)

	ig := new(types.Ingress)
	ig.Meta.SetDefault()

	ig.Meta.Name = opts.Meta.Name
	ig.Status = opts.Status
	ig.SelfLink()

	if err := n.storage.Put(n.context, n.storage.Collection().Ingress(),
		n.storage.Key().Ingress(ig.Meta.Name), ig, nil); err != nil {
		log.Debugf("%s:create:> insert ingress err: %v", logIngressPrefix, err)
		return nil, err
	}

	return ig, nil
}

func (n *Ingress) Get(name string) (*types.Ingress, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logIngressPrefix, name)

	ingress := new(types.Ingress)

	err := n.storage.Get(n.context, n.storage.Collection().Ingress(), n.storage.Key().Ingress(name), ingress, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: ingress %s not found", logIngressPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get ingress `%s` err: %v", logIngressPrefix, name, err)
		return nil, err
	}

	return ingress, nil
}

func (n *Ingress) SetMeta(ingress *types.Ingress, meta *types.IngressUpdateMetaOptions) error {

	log.V(logLevel).Debugf("%s:setmeta:> update Ingress %#v", logIngressPrefix, meta)
	if meta == nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Ingress err: %v", logIngressPrefix, errors.New(errors.ArgumentIsEmpty))
		return errors.New(errors.ArgumentIsEmpty)
	}

	ingress.Meta.Set(meta)

	if err := n.storage.Set(n.context, n.storage.Collection().Ingress(),
		n.storage.Key().Ingress(ingress.Meta.Name), &ingress, nil); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Ingress meta err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) SetStatus(ingress *types.Ingress, status types.IngressStatus) error {

	ingress.Status = status

	if err := n.storage.Set(n.context, n.storage.Collection().Ingress(),
		n.storage.Key().Ingress(ingress.Meta.Name), &ingress, nil); err != nil {
		log.Errorf("%s:setstatus:> set ingress offline state error: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) Remove(ingress *types.Ingress) error {

	log.V(logLevel).Debugf("%s:remove:> remove ingress %s", logIngressPrefix, ingress.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Ingress(), n.storage.Key().Ingress(ingress.Meta.Name)); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove ingress err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func NewIngressModel(ctx context.Context, stg storage.Storage) *Ingress {
	return &Ingress{ctx, stg}
}
