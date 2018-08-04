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

const (
	logTriggerPrefix = "distribution:trigger"
)

type Trigger struct {
	context context.Context
	storage storage.Storage
}

func (t *Trigger) Get(namespace, service, name string) (*types.Trigger, error) {

	log.V(logLevel).Debugf("%s:get:> get trigger by name %s: %s", logTriggerPrefix, namespace, name)

	trigger := new(types.Trigger)

	err := t.storage.Get(t.context, storage.TriggerKind, t.storage.Key().Trigger(namespace, service, name), &trigger, nil)
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> create trigger err: %v", logTriggerPrefix, err)
		return nil, err
	}

	return trigger, nil
}

func NewTriggerModel(ctx context.Context, stg storage.Storage) *Trigger {
	return &Trigger{ctx, stg}
}
