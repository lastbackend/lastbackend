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
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
)

const (
	logIngressPrefix = "distribution:ingress"
)

type IIngress interface {
	List() (map[string]*types.Ingress, error)
	Create(opts *types.IngressCreateOptions) (*types.Ingress, error)
	Get(name string) (*types.Ingress, error)
	SetMeta(ingress *types.Ingress, meta *types.IngressUpdateMetaOptions) error
	SetStatus(ingress *types.Ingress, state types.IngressStatus) error
	Remove(ingress *types.Ingress) error
}

type Ingress struct {
	context context.Context
	storage storage.Storage
}

func (n *Ingress) List() (map[string]*types.Ingress, error) {
	list := make(map[string]*types.Ingress, 0)

	if err := n.storage.Map(n.context, storage.IngressKind, "", &list); err != nil {
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

	if err := n.storage.Create(n.context, storage.IngressKind, ig.Meta.SelfLink, ig, nil); err != nil {
		log.Debugf("%s:create:> insert ingress err: %v", logIngressPrefix, err)
		return nil, err
	}

	return ig, nil
}

func (n *Ingress) Get(name string) (*types.Ingress, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logIngressPrefix, name)

	ingress := new(types.Ingress)

	err := n.storage.Get(n.context, storage.IngressKind, name, &ingress)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
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

	if err := n.storage.Update(n.context, storage.IngressKind, ingress.Meta.SelfLink, &ingress, nil); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Ingress meta err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) SetStatus(ingress *types.Ingress, status types.IngressStatus) error {

	ingress.Status = status

	if err := n.storage.Update(n.context, storage.IngressKind, ingress.Meta.SelfLink, &ingress, nil); err != nil {
		log.Errorf("%s:setstatus:> set ingress offline state error: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) Remove(ingress *types.Ingress) error {

	log.V(logLevel).Debugf("%s:remove:> remove ingress %s", logIngressPrefix, ingress.Meta.Name)

	if err := n.storage.Remove(n.context, storage.IngressKind, ingress.Meta.SelfLink); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove ingress err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func NewIngressModel(ctx context.Context, stg storage.Storage) IIngress {
	return &Ingress{ctx, stg}
}
