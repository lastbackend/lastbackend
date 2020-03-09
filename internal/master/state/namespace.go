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

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/logger"
)

// NamespaceController structure
type NamespaceController struct {
	items []*types.Namespace
}

// List all namespaces in state
func (ns *NamespaceController) List(ctx context.Context) []*types.Namespace {

	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> get namespace list", logPrefix)

	return ns.items
}

// Set namespace to state
func (ns *NamespaceController) Set(ctx context.Context, mf *types.NamespaceManifest) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> set namespace", logPrefix)
}

// Get particular namespace from state
func (ns *NamespaceController) Get(ctx context.Context) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> get namespace from state", logPrefix)
}

// Del namespace in state
func (ns *NamespaceController) Del(ctx context.Context) {
	log := logger.WithContext(ctx)
	log.Debugf("%s:list:> delete namespace from state", logPrefix)
}

// NewNamespaceController return new instance of namespace controller
func NewNamespaceController(ctx context.Context) *NamespaceController {
	nc := new(NamespaceController)
	return nc
}
