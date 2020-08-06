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
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logResolversPrefix = "state:resolvers:>"

type ResolverState struct {
	lock      sync.RWMutex
	resolvers map[string]*models.ResolverManifest
}

func (n *ResolverState) GetResolvers() map[string]*models.ResolverManifest {
	return n.resolvers
}

func (n *ResolverState) AddResolver(cidr string, sn *models.ResolverManifest) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s add resolver: %s", logResolversPrefix, cidr)
	n.SetResolver(cidr, sn)
}

func (n *ResolverState) SetResolver(cidr string, sn *models.ResolverManifest) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s set resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}

	n.resolvers[cidr] = sn
}

func (n *ResolverState) GetResolver(cidr string) *models.ResolverManifest {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.resolvers[cidr]
	if !ok {
		return nil
	}
	return s
}

func (n *ResolverState) DelResolver(cidr string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s del resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}
}
