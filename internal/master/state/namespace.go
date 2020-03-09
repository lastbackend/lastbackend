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

package state

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/logger"
)

// NamespaceController structure
type NamespaceController struct {
	lock  sync.Mutex
	items map[*types.NamespaceSelfLink]*types.Namespace
}

// List all namespaces in state
func (ns *NamespaceController) List(ctx context.Context) []*types.Namespace {

	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> get namespace list", logPrefix)

	ns.lock.Lock()
	items := make([]*types.Namespace, len(ns.items))
	for _, item := range ns.items {
		items = append(items, item)
	}
	ns.lock.Unlock()

	return items
}

// Map all namespaces in state
func (ns *NamespaceController) Map(ctx context.Context) map[*types.NamespaceSelfLink]*types.Namespace {

	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> get namespace list", logPrefix)

	return ns.items
}

// Set namespace to state
func (ns *NamespaceController) Put(ctx context.Context, item *types.NamespaceManifest) (*types.Namespace, error) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> put namespace", logPrefix)

	/*
		internal, external := handler.Config.DomainInternal, handler.Config.DomainExternal
			ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))


			ns.Spec.Domain.Internal = internal

			if opts.Spec.Domain != nil {
				if len(*opts.Spec.Domain) == 0 {
					ns.Spec.Domain.External = external
				} else {
					ns.Spec.Domain.External = *opts.Spec.Domain
				}
			}
	*/

	return nil, nil
}

// Set namespace to state
func (ns *NamespaceController) Set(ctx context.Context, item *types.NamespaceManifest) (*types.Namespace, error) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> set namespace", logPrefix)

	return nil, nil
}

// Get particular namespace from state
func (ns *NamespaceController) Get(ctx context.Context, selflink *types.NamespaceSelfLink) (*types.Namespace, error) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> get namespace from state", logPrefix)

	ns.lock.Lock()
	item, ok := ns.items[selflink]
	ns.lock.Unlock()

	if !ok {
		return nil, errors.NewResourceNotFound()
	}

	return item, nil
}

// Del namespace in state
func (ns *NamespaceController) Del(ctx context.Context, selflink *types.NamespaceSelfLink) error {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> delete namespace from state", logPrefix)

	ns.lock.Lock()
	_, ok := ns.items[selflink]
	ns.lock.Unlock()

	if !ok {
		return errors.NewResourceNotFound()
	}

	return nil
}

// NewNamespaceController return new instance of namespace controller
func NewNamespaceController(ctx context.Context) *NamespaceController {
	nc := new(NamespaceController)
	return nc
}
