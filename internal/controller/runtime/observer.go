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

package runtime

import (
	"github.com/lastbackend/lastbackend/internal/controller/envs"
	"github.com/lastbackend/lastbackend/internal/controller/state"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
)

const (
	logPrefix = "controller:runtime:observer"
)

type Observer struct {
	rev   *int64
	stg   storage.Storage
	state *state.State
}

func (o *Observer) Loop() {
	o.state.Loop()
}

func (o *Observer) Stop() {
	o.state = nil
}

func NewObserver() *Observer {

	o := new(Observer)

	o.stg = envs.Get().GetStorage()
	o.state = state.NewState()

	return o
}
