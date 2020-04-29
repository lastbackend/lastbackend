//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type ResourceController struct {
}

// List method returns slice of namespace resources
func (ns *ResourceController) List(ctx context.Context, filter *ResourceFilter) (models.NamespaceResource, error) {
	return nil, nil
}

// Get method returns namespace resource
func (ns *ResourceController) Get(ctx context.Context, selflink models.SelfLink) (models.NamespaceResource, error) {
	return nil, nil
}

// Put method insert namespace resource
func (ns *ResourceController) Put(ctx context.Context, manifest models.NamespaceResourceManifest) (models.NamespaceResource, error) {
	return nil, nil
}

// Set method updates namespace resource
func (ns *ResourceController) Set(ctx context.Context, manifest models.NamespaceResourceManifest) (models.NamespaceResource, error) {
	return nil, nil
}

// Del method removes namespace resource
func (ns *ResourceController) Del(ctx context.Context, selflink models.SelfLink) (models.NamespaceResource, error) {
	return nil, nil
}

type ResourceFilter struct {
	namespace []string
	kind      []string
}

func (rf *ResourceFilter) WithNamespace(namespace ...string) *ResourceFilter {
	for _, ns := range namespace {
		if len(ns) <= 0 {
			continue
		}
		rf.namespace = append(rf.namespace, ns)
	}

	return rf
}

func (rf *ResourceFilter) WithKind(kind ...string) *ResourceFilter {
	for _, k := range kind {

		if len(k) <= 0 {
			continue
		}

		rf.kind = append(rf.kind, k)
	}

	return rf
}

func NewResourceFilter() *ResourceFilter {
	return new(ResourceFilter)
}
