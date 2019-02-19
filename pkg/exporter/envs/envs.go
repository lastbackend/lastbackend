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

package envs

import (
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/exporter/logger"
	"github.com/lastbackend/lastbackend/pkg/exporter/state"
	"github.com/lastbackend/lastbackend/pkg/network"
)

var _env Env

type Env struct {
	net    *network.Network
	state  *state.State
	client types.ExporterClientV1
	logger *logger.Logger
}

func Get() *Env {
	return &_env
}

func (c *Env) SetNet(n *network.Network) {
	c.net = n
}

func (c *Env) GetNet() *network.Network {
	return c.net
}

func (c *Env) SetState(state *state.State) {
	c.state = state
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetLogger(logger *logger.Logger) {
	c.logger = logger
}

func (c *Env) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Env) SetClient(client types.ExporterClientV1) {
	c.client = client
}

func (c *Env) GetClient() types.ExporterClientV1 {
	return c.client
}
