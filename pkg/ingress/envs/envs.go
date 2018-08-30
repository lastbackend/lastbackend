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

package envs

import (
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/ingress/events/exporter"
	"github.com/lastbackend/lastbackend/pkg/ingress/state"
)

var e Env

func Get() *Env {
	return &e
}

type Env struct {
	state    *state.State
	client   types.IngressClientV1
	exporter *exporter.Exporter
}

func (c *Env) SetState(s *state.State) {
	c.state = s
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetClient(cl types.IngressClientV1) {
	c.client = cl
}

func (c *Env) GetClient() types.IngressClientV1 {
	return c.client
}

func (c *Env) SetExporter(e *exporter.Exporter) {
	c.exporter = e
}

func (c *Env) GetExporter() *exporter.Exporter {
	return c.exporter
}
