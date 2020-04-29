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

package exporter

import (
	"fmt"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/proxy"
)

type Exporter struct {
	ready  bool
	srv    *proxy.Server
	client *proxy.Client
}

func (c *Exporter) Proxy(msg models.ProxyMessage) error {

	if !c.ready {
		// TODO: cache messages due to reconnect
		return nil
	}

	return c.client.Send(msg.Line)
}

func (c *Exporter) Listen() {
	for {
		if err := c.srv.Listen(c.Proxy); err != nil {
			fmt.Println(err.Error())
		}
		<-time.NewTimer(3 * time.Second).C
	}
}

func (c *Exporter) Reconnect(addr string) {
	c.client.Reconnect(addr)
}

func NewExporter(name, addr string) (*Exporter, error) {

	var err error

	c := new(Exporter)
	c.srv, err = proxy.NewServer(proxy.DefaultServer)
	if err != nil {
		return nil, err
	}

	c.client = proxy.NewClient(name, addr, nil)
	c.ready = true

	return c, nil
}
