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

package discovery

import (
	"github.com/lastbackend/lastbackend/internal/discovery/resources"
	"github.com/lastbackend/lastbackend/internal/util/dns"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logLevel = 3

func Listen(host string, port int) (*dns.DNS, error) {

	var d = dns.DNS{}

	log.V(logLevel).Debug(`Init discovery resources`)

	for pattern, resource := range resources.Map {
		d.AddHandler(pattern, resource)
	}

	go func() {
		log.V(logLevel).Debugf(`Start discovery %s service on %d port`, dns.TCP, port)
		if err := d.Start(dns.TCP, host, port, nil); err != nil {
			log.Errorf(`Start discovery %s service on %d port error: %s`, dns.TCP, port, err)
			return
		}
	}()

	go func() {
		log.V(logLevel).Debugf(`Start discovery %s service on %d port`, dns.UDP, port)
		if err := d.Start(dns.UDP, host, port, nil); err != nil {
			log.Errorf(`Start discovery %s service on %d port error: %s`, dns.TCP, port, err)
			return
		}
	}()

	return &d, nil
}
