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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

type ITrigger interface {
	Get(namespace, service, name string) (*types.Trigger, error)
}

type Trigger struct {
	context context.Context
	storage storage.Storage
}

func (h *Trigger) Get(namespace, service, name string) (*types.Trigger, error) {

	log.V(logLevel).Debugf("Trigger: Get: get Trigger by name %s: %s", namespace, name)

	hook, err := h.storage.Trigger().Get(h.context, namespace, service, name)
	if err != nil {
		log.V(logLevel).Errorf("Trigger: Get: create Trigger err: %s", err)
		return nil, err
	}

	return hook, nil
}

func NewTriggerModel(ctx context.Context, stg storage.Storage) ITrigger {
	return &Trigger{ctx, stg}
}
