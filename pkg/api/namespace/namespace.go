//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package namespace

import (
	"context"
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const logLevel = 3

type namespace struct {
	Context context.Context
}

func New(ctx context.Context) *namespace {
	return &namespace{ctx}
}

func (ns *namespace) List() (types.NamespaceList, error) {
	var (
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
		list    = types.NamespaceList{}
	)

	log.V(logLevel).Debug("Namespace: list namespace")

	items, err := storage.Namespace().List(ns.Context)
	if err != nil {
		log.V(logLevel).Error("Namespace: list namespace err: %s", err.Error())
		return list, err
	}

	log.V(logLevel).Debugf("Namespace: list namespace result: %d", len(items))

	for _, item := range items {
		var ns = item
		list = append(list, ns)
	}

	return list, nil
}

func (ns *namespace) Get(name string) (*types.Namespace, error) {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Namespace: get namespace %s", name)

	n, err := storage.Namespace().GetByName(ns.Context, name)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Warnf("Namespace: namespace by name `%s` not found", name)
			return nil, nil
		}
		log.V(logLevel).Errorf("Namespace: get namespace by name `%s` err: %s", name, err.Error())
		return nil, err
	}

	return n, nil
}

func (ns *namespace) Create(rq *request.RequestNamespaceCreateS) (*types.Namespace, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Namespace: create namespace %#v", rq)

	var nsp = types.Namespace{}
	nsp.Meta.SetDefault()
	nsp.Meta.Name = rq.Name
	nsp.Meta.Description = rq.Description

	if err = storage.Namespace().Insert(ns.Context, &nsp); err != nil {
		log.V(logLevel).Errorf("Namespace: insert namespace err: %s", err.Error())
		return &nsp, err
	}

	return &nsp, nil
}

func (ns *namespace) Update(n *types.Namespace) (*types.Namespace, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Namespace: update namespace %#v", n)

	if err = storage.Namespace().Update(ns.Context, n); err != nil {
		log.V(logLevel).Errorf("Namespace: update namespace err: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *namespace) Remove(name string) error {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Namespace: remove namespace %s", name)

	err = storage.Namespace().Remove(ns.Context, name)
	if err != nil {
		log.V(logLevel).Errorf("Namespace: remove namespace err: %s", err.Error())
		return err
	}

	return nil
}

func (ns *namespace) WatchService(service chan *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage().Service()
	)

	log.V(logLevel).Debugf("Namespace: watch services in namespace")

	if err := storage.PodsWatch(ns.Context, service); err != nil {
		log.V(logLevel).Errorf("Namespace: watch services in namespace err: %s", err.Error())
		return err
	}

	return nil
}
