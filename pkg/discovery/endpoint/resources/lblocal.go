//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package resources

import (
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/discovery/endpoint"
	"github.com/lastbackend/lastbackend/pkg/util"
	"github.com/miekg/dns"
	"time"
)

const logLevel = 3

func LbLocalR(w dns.ResponseWriter, r *dns.Msg) {

	var (
		log = context.Get().GetLogger()
	)

	log.V(logLevel).Debug("Resource: dns request lblocal.")

	var (
		v4 bool
		rr dns.RR
	)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		log.V(logLevel).Debug("Resource: dns.OpcodeQuery")

		for _, q := range m.Question {

			switch r.Question[0].Qtype {
			case dns.TypeTXT:
				log.V(logLevel).Debug("Resource: get txt type query")
				t := new(dns.TXT)
				t.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}
				m.Answer = append(m.Answer, t)
				m.Extra = append(m.Extra, rr)
			default:
				log.V(logLevel).Debug("Resource: get unknown query type")
				fallthrough
			case dns.TypeAAAA, dns.TypeA:
				log.V(logLevel).Debug("Resource: get A or AAAA type query")

				if q.Name[len(q.Name)-1:] != "." {
					q.Name += "."
				}

				log.V(logLevel).Debugf("Resource: find ip addresses for domain: ", q.Name)

				// Generate A and AAAA records
				ips, err := endpoint.Get(util.Trim(q.Name, `.`))
				if err != nil {
					log.V(logLevel).Errorf("Resource: get endpoint `%s` err: %s", q.Name, err.Error())
					w.WriteMsg(m)
					return
				}

				if ips == nil {
					w.WriteMsg(m)
					return
				}

				log.V(logLevel).Debugf("Resource: ips list: %#v for %s", ips, q.Name)

				for _, ip := range ips {
					v4 = ip.To4() != nil

					if v4 {
						rr = new(dns.A)
						rr.(*dns.A).Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
						rr.(*dns.A).A = ip.To4()
					} else {
						rr = new(dns.AAAA)
						rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0}
						rr.(*dns.AAAA).AAAA = ip
					}

					m.Answer = append(m.Answer, rr)
				}
			}
		}
	case dns.OpcodeUpdate:
		log.Debugf("Resource: dns.OpcodeUpdate")
	}

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacMD5, 300, time.Now().Unix())
		} else {
			log.V(logLevel).Errorf("Resource: tsig status err: %s", w.TsigStatus().Error())
		}
	}

	log.Debugf("Resource: send message info  %#v", m)

	w.WriteMsg(m)
}
