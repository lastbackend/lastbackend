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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type namespace struct {
	Context context.Context
}

func New(ctx context.Context) *namespace {
	return &namespace{ctx}
}

func (ns *namespace) List() (*types.NamespaceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Namespace().List(ns.Context)
}

func (ns *namespace) Get(id string) (*types.Namespace, error) {
	var (
		err     error
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	if validator.IsUUID(id) {
		n, err = storage.Namespace().GetByID(ns.Context, id)
	} else {
		n, err = storage.Namespace().GetByName(ns.Context, id)
	}

	return n, err
}

func (ns *namespace) Create(rq *request.RequestNamespaceCreateS) (*types.Namespace, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	n, err = storage.Namespace().Insert(ns.Context, rq.Name, rq.Description)
	if err != nil {
		log.Errorf("Error: insert namespace to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *namespace) Update(n *types.Namespace) (*types.Namespace, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	n, err = storage.Namespace().Update(ns.Context, n)
	if err != nil {
		log.Errorf("Error: update namespace to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *namespace) Remove(id string) error {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)
	err = storage.Namespace().Remove(ns.Context, id)
	if err != nil {
		log.Error("Error: remove namespace from db", err)
		return err
	}

	return nil
}
