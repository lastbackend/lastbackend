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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/agent/exporter"
	"github.com/lastbackend/lastbackend/pkg/agent/state"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/spf13/viper"
)

var e Env

func Get() *Env {
	return &e
}

type config struct {
	Verbose  string `mapstructure:"verbose" json:"verbose" yaml:"verbose"`
	Token    string `mapstructure:"token" json:"token" yaml:"token"`
	Workdir  string `mapstructure:"workdir" json:"workdir" yaml:"workdir"`
	Manifest struct {
		Dir string `mapstructure:"dir" json:"dir" yaml:"dir"`
	} `mapstructure:"manifest" json:"manifest" yaml:"manifest"`
	Network struct {
		Interface string `mapstructure:"interface" json:"interface" yaml:"interface"`
		Cpi       struct {
			Type      string `mapstructure:"type" json:"type" yaml:"type"`
			Interface struct {
				External interface{} `mapstructure:"external" json:"external" yaml:"external"`
				Internal interface{} `mapstructure:"internal" json:"internal" yaml:"internal"`
			} `mapstructure:"interface" json:"interface" yaml:"interface"`
		} `mapstructure:"cpi" json:"cpi" yaml:"cpi"`
		Cni struct {
			Type string `mapstructure:"type" json:"type" yaml:"type"`
		} `mapstructure:"cni" json:"cni" yaml:"cni"`
	} `mapstructure:"network" json:"network" yaml:"network"`
	Container struct {
		Cri struct {
			Type   string `mapstructure:"type" json:"type" yaml:"type"`
			Docker struct {
				Version string `mapstructure:"version" json:"version" yaml:"version"`
			} `mapstructure:"docker" json:"docker" yaml:"docker"`
		} `mapstructure:"cri" json:"cri" yaml:"cri"`
		Csi struct {
			Dir struct {
				Root string `mapstructure:"root" json:"root" yaml:"root"`
			} `mapstructure:"dir" json:"dir" yaml:"dir"`
		} `mapstructure:"csi" json:"csi" yaml:"csi"`
		Iri struct {
			Type   string `mapstructure:"type" json:"type" yaml:"type"`
			Docker struct {
				Version string `mapstructure:"version" json:"version" yaml:"version"`
			} `mapstructure:"docker" json:"docker" yaml:"docker"`
		} `mapstructure:"iri" json:"iri" yaml:"iri"`
		ExtraHosts []string `mapstructure:"extra_hosts" json:"extra_hosts" yaml:"extra_hosts"`
	} `mapstructure:"container" json:"container" yaml:"container"`
	Server struct {
		Host string `mapstructure:"host" json:"host" yaml:"host"`
		Port string `mapstructure:"port" json:"port" yaml:"port"`
		TLS  struct {
			Insecure string `mapstructure:"insecure" json:"insecure" yaml:"insecure"`
		} `mapstructure:"tls" json:"tls" yaml:"tls"`
		Cert string `mapstructure:"cert" json:"cert" yaml:"cert"`
		Key  string `mapstructure:"key" json:"key" yaml:"key"`
		Ca   string `mapstructure:"ca" json:"ca" yaml:"ca"`
	} `mapstructure:"server" json:"server" yaml:"server"`
	API struct {
		URI interface{} `mapstructure:"uri" json:"uri" yaml:"uri"`
		TLS struct {
			Insecure string `mapstructure:"insecure" json:"insecure" yaml:"insecure"`
			Cert     string `mapstructure:"cert" json:"cert" yaml:"cert"`
			Key      string `mapstructure:"key" json:"key" yaml:"key"`
			Ca       string `mapstructure:"ca" json:"ca" yaml:"ca"`
		} `mapstructure:"tls" json:"tls" yaml:"tls"`
	} `mapstructure:"api" json:"api" yaml:"api"`
}

type Env struct {
	config config

	cri cri.CRI
	cii cii.CII
	csi map[string]csi.CSI

	state *state.State
	net   *network.Network

	client struct {
		node types.NodeClientV1
		rest types.ClientV1
	}

	exporter *exporter.Exporter
	mode     struct {
		provision bool
	}
}

func (c *Env) SetConfig(v *viper.Viper) error {
	cfg := new(config)

	if err := v.Unmarshal(cfg); err != nil {
		return err
	}
	c.config = *cfg
	return nil
}

func (c *Env) GetConfig() config {
	return c.config
}

func (c *Env) SetCRI(cri cri.CRI) {
	c.cri = cri
}

func (c *Env) GetCRI() cri.CRI {
	return c.cri
}

func (c *Env) SetCII(iri cii.CII) {
	c.cii = iri
}

func (c *Env) GetCII() cii.CII {
	return c.cii
}

func (c *Env) SetNet(net *network.Network) {
	c.net = net
}

func (c *Env) GetNet() *network.Network {
	return c.net
}

func (c *Env) SetProvision(on bool) {
	c.mode.provision = on
}

func (c *Env) GetProvision() bool {
	return c.mode.provision
}

func (c *Env) SetCSI(kind string, si csi.CSI) {
	c.csi = make(map[string]csi.CSI)
	c.csi[kind] = si
}

func (c *Env) ListCSI() []string {
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

func (c *Env) SetExporter(s *exporter.Exporter) {
	c.exporter = s
}

func (c *Env) GetExporter() *exporter.Exporter {
	return c.exporter
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
