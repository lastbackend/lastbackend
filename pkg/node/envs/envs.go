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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/node/events/exporter"
	"github.com/lastbackend/lastbackend/pkg/node/ingress"
	"github.com/lastbackend/lastbackend/pkg/node/state"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/runtime/cpi"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/lastbackend/lastbackend/pkg/runtime/iri"
)

var e Env

func Get() *Env {
	return &e
}

type Env struct {
	cri cri.CRI
	iri iri.IRI
	cni cni.CNI
	cpi cpi.CPI
	csi map[string]csi.CSI

	state  *state.State
	client struct {
		node types.NodeClientV1
		rest types.ClientV1
	}

	mode struct {
		provision bool
		ingress   bool
	}

	exporter *exporter.Exporter
	ingress ingress.Ingress
}

func (c *Env) SetCRI(cri cri.CRI) {
	c.cri = cri
}

func (c *Env) GetCRI() cri.CRI {
	return c.cri
}

func (c *Env) SetIRI(iri iri.IRI) {
	c.iri = iri
}

func (c *Env) GetIRI() iri.IRI {
	return c.iri
}

func (c *Env) SetCNI(n cni.CNI) {
	c.cni = n
}

func (c *Env) GetCNI() cni.CNI {
	return c.cni
}

func (c *Env) SetCPI(cpi cpi.CPI) {
	c.cpi = cpi
}

func (c *Env) GetCPI() cpi.CPI {
	return c.cpi
}

func (c *Env) SetIngress(ing ingress.Ingress) {
	c.ingress = ing
}

func (c *Env) GetIngress() ingress.Ingress {
	return c.ingress
}

func (c *Env) SetModeIngress(on bool) {
	c.mode.ingress = on
}

func (c *Env) GetModeIngress() bool {
	return c.mode.ingress
}

func (c *Env) SetPorvision(on bool) {
	c.mode.provision = on
}

func (c *Env) GetProvision() bool {
	return c.mode.provision
}

func (c *Env) SetCSI(kind string, si csi.CSI) {
	c.csi = make(map[string]csi.CSI)
	c.csi[kind] = si
}

func (c *Env) ListCSI () []string {
	var types = []string{}

	for t := range c.csi {
		types = append(types, t)
	}
	return types
}

func (c *Env) GetCSI(kind string) (csi.CSI, error) {
	if _, ok := c.csi[kind]; !ok {
		return nil, errors.New("storage container interface not supported")
	}
	return c.csi[kind], nil
}

func (c *Env) SetState(s *state.State) {
	c.state = s
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetClient(nc types.NodeClientV1, rc types.ClientV1) {
	c.client.node = nc
	c.client.rest = rc
}

func (c *Env) GetNodeClient() types.NodeClientV1 {
	return c.client.node
}

func (c *Env) GetRestClient() types.ClientV1 {
	return c.client.rest
}

func (c *Env) SetExporter(e *exporter.Exporter) {
	c.exporter = e
}

func (c *Env) GetExporter() *exporter.Exporter {
	return c.exporter
}
