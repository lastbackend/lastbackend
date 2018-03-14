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
	"github.com/lastbackend/lastbackend/pkg/discovery/dns/resources"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/dns"
)

func Listen(port int) (*dns.DNS, error) {

	var d = dns.DNS{}

	log.Debug(`Init discovery resources`)

	d.AddHandler(`lb.local.`, resources.LbLocal)
	d.AddHandler(`lstbknd.io.`, resources.LstbkndIo)

	go func() {
		log.Debugf(`Start discovery %s service on %d port`, dns.TCP, port)
		if err := d.Start(dns.TCP, port, nil); err != nil {
			log.Errorf(`Start discovery %s service on %d port error: %s`, dns.TCP, port, err)
			return
		}
	}()

	go func() {
		log.Debugf(`Start discovery %s service on %d port`, dns.UDP, port)
		if err := d.Start(dns.UDP, port, nil); err != nil {
			log.Errorf(`Start discovery %s service on %d port error: %s`, dns.TCP, port, err)
			return
		}
	}()

	return &d, nil
}
