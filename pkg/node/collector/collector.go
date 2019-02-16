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

package collector

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/proxy"
	"time"
)

type Collector struct {
	srv    *proxy.Server
	client *proxy.Client
}

func (c *Collector) Handler(msg types.ProxyMessage) error {
	log.Debugf("collector: receive message: %s", msg.Line)
	return nil
}

func (c *Collector) Listen() {
	for {
		if err := c.srv.Listen(c.Handler); err != nil {
			log.Errorf(err.Error())
		}
		<-time.NewTimer(3 * time.Second).C
	}
}

func NewCollector() (*Collector, error) {

	var err error

	c := new(Collector)
	c.srv, err = proxy.NewServer(proxy.DefaultServer)
	if err != nil {
		return nil, err
	}

	return c, nil
}
