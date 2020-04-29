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

const logEndpointPrefix = "state:endpoints:>"

type EndpointState struct {
	lock      sync.RWMutex
	hash      string
	endpoints map[string]*models.EndpointState
}

func (es *EndpointState) GetHash() string {
	return es.hash
}

func (es *EndpointState) SetHash(hash string) {
	es.hash = hash
}

func (es *EndpointState) GetEndpoints() map[string]*models.EndpointState {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get endpoints", logEndpointPrefix)
	return es.endpoints
}

func (es *EndpointState) SetEndpoints(endpoints map[string]*models.EndpointState) {
	es.lock.Lock()
	defer es.lock.Unlock()

	for key, endpoint := range endpoints {
		es.endpoints[key] = endpoint
	}
}

func (es *EndpointState) GetEndpoint(key string) *models.EndpointState {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get endpoint: %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()

	ep, ok := es.endpoints[key]
	if !ok {
		return nil
	}

	return ep
}

func (es *EndpointState) AddEndpoint(key string, endpoint *models.EndpointState) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: add endpoint %s: %s", logEndpointPrefix, key, endpoint.IP)
	es.lock.Lock()
	defer es.lock.Unlock()
	es.endpoints[key] = endpoint
}

func (es *EndpointState) SetEndpoint(key string, endpoint *models.EndpointState) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: set endpoint %s: %s", logEndpointPrefix, key, endpoint.IP)
	es.lock.Lock()
	defer es.lock.Unlock()
	es.endpoints[key] = endpoint
}

func (es *EndpointState) DelEndpoint(key string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: del endpoint %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()
	delete(es.endpoints, key)
}
