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

type IHook interface {
	Get(id string) (*types.Hook, error)
}

type Hook struct {
	context context.Context
	storage storage.Storage
}

func (h *Hook) Get(id string) (*types.Hook, error) {

	log.V(logLevel).Debugf("Hook: Get: get Hook by id %s", id)

	hook, err := h.storage.Hook().Get(h.context, id)
	if err != nil {
		log.V(logLevel).Errorf("Hook: Get: create Hook err: %s", err)
		return nil, err
	}

	return hook, nil
}

func NewHookModel(ctx context.Context, stg storage.Storage) IHook {
	return &Hook{ctx, stg}
}
