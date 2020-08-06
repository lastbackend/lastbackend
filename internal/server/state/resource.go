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

// ResourceController defines base resource controller methods
type ResourceController interface {
	// Run controller state
	Run(ctx context.Context) error
	// Restore controller state
	Restore(ctx context.Context) error
	// List method returns slice of namespace resources
	List(ctx context.Context, filter *ResourceFilter) (models.NamespaceResourceList, error)
	// Get method returns namespace resource
	Get(ctx context.Context, selflink models.SelfLink) (models.NamespaceResource, error)
	// Put method insert namespace resource
	Put(ctx context.Context, manifest models.NamespaceResourceManifest) (models.NamespaceResource, error)
	// Set method updates namespace resource
	Set(ctx context.Context, manifest models.NamespaceResourceManifest) (models.NamespaceResource, error)
	// Del method removes namespace resource
	Del(ctx context.Context, selflink models.SelfLink) (models.NamespaceResource, error)
}

type ResourceFilter struct {
	namespace []string
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

func NewResourceFilter() *ResourceFilter {
	return new(ResourceFilter)
}
