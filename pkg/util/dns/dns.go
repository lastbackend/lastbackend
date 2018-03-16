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

package dns

import (
	"fmt"
	"github.com/miekg/dns"
)

const (
	UDP = "udp"
	TCP = "tcp"
)

type DNS struct {
	servers []*dns.Server
}

type Tsig struct {
	Name   string
	Secret string
}

func (d *DNS) Start(net string, port int, tsig *Tsig) error {
	var server = &dns.Server{Addr: fmt.Sprintf(":%d", port), Net: net}

	if tsig != nil {
		server.TsigSecret = map[string]string{tsig.Name: tsig.Secret}
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	d.servers = append(d.servers, server)

	return nil
}

func (d *DNS) AddHandler(pattern string, handler dns.HandlerFunc) {
	dns.HandleFunc(pattern, handler)
}

func (d *DNS) Shutdown() {
	for _, server := range d.servers {
		server.Shutdown()
	}
}
