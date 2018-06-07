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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	logIngressPrefix = "distribution:ingress"
)

type IIngress interface {
	List() (map[string]*types.Ingress, error)
	Create(opts *types.IngressCreateOptions) (*types.Ingress, error)
	Get(name string) (*types.Ingress, error)
	GetSpec(ingress *types.Ingress) (*types.IngressSpec, error)
	SetMeta(ingress *types.Ingress, meta *types.IngressUpdateMetaOptions) error
	SetStatus(ingress *types.Ingress, state types.IngressStatus) error
	Remove(ingress *types.Ingress) error
}

type Ingress struct {
	context context.Context
	storage storage.Storage
}

func (n *Ingress) List() (map[string]*types.Ingress, error) {
	return n.storage.Ingress().List(n.context)
}

func (n *Ingress) Create(opts *types.IngressCreateOptions) (*types.Ingress, error) {

	log.Debugf("%s:create:> create ingress in cluster", logIngressPrefix)

	ig := new(types.Ingress)
	ig.Meta.SetDefault()

	ig.Meta.Name = opts.Meta.Name
	ig.Status = opts.Status

	if err := n.storage.Ingress().Insert(n.context, ig); err != nil {
		log.Debugf("%s:create:> insert ingress err: %s", logIngressPrefix, err.Error())
		return nil, err
	}

	return ig, nil
}

func (n *Ingress) Get(name string) (*types.Ingress, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logIngressPrefix, name)

	ingress, err := n.storage.Ingress().Get(n.context, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> get: ingress %s not found", logIngressPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get ingress `%s` err: %s", logIngressPrefix, name, err.Error())
		return nil, err
	}

	return ingress, nil
}

func (n *Ingress) GetSpec(ingress *types.Ingress) (*types.IngressSpec, error) {

	log.V(logLevel).Debugf("%s:getspec:> get ingress spec: %s", logIngressPrefix, ingress.Meta.Name)

	spec, err := n.storage.Ingress().GetSpec(n.context, ingress)
	if err != nil {
		log.V(logLevel).Debugf("%s:getspec:> get Ingress `%s` err: %s", logIngressPrefix, ingress.Meta.Name, err.Error())
		return nil, err
	}

	return spec, nil
}

func (n *Ingress) SetMeta(ingress *types.Ingress, meta *types.IngressUpdateMetaOptions) error {

	log.V(logLevel).Debugf("%s:setmeta:> update Ingress %#v", logIngressPrefix, meta)
	if meta == nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Ingress err: %s", logIngressPrefix, errors.New(errors.ArgumentIsEmpty))
		return errors.New(errors.ArgumentIsEmpty)
	}

	ingress.Meta.Set(meta)

	if err := n.storage.Ingress().Update(n.context, ingress); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Ingress meta err: %s", logIngressPrefix, err.Error())
		return err
	}

	return nil
}

func (n *Ingress) SetStatus(ingress *types.Ingress, status types.IngressStatus) error {

	ingress.Status = status

	if err := n.storage.Ingress().SetStatus(n.context, ingress); err != nil {
		log.Errorf("%s:setstatus:> set ingress offline state error: %s", logIngressPrefix, err.Error())
		return err
	}

	return nil
}

func (n *Ingress) Remove(ingress *types.Ingress) error {

	log.V(logLevel).Debugf("%s:remove:> remove ingress %s", logIngressPrefix, ingress.Meta.Name)

	if err := n.storage.Ingress().Remove(n.context, ingress); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove ingress err: %s", logIngressPrefix, err.Error())
		return err
	}

	return nil
}

func NewIngressModel(ctx context.Context, stg storage.Storage) IIngress {
	return &Ingress{ctx, stg}
}
