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

func OtherR(w dns.ResponseWriter, r *dns.Msg) {

	var log = context.Get().GetLogger()

	log.Debug(`Dns request OTHER`)

	var (
		v4 bool
		rr dns.RR
	)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		log.Debug(`dns.OpcodeQuery`)

		for _, q := range m.Question {

			switch r.Question[0].Qtype {
			case dns.TypeTXT:
				t := new(dns.TXT)
				t.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}
				m.Answer = append(m.Answer, t)
				m.Extra = append(m.Extra, rr)
			default:
				fallthrough
			case dns.TypeAAAA, dns.TypeA:

				if q.Name[len(q.Name)-1:] != "." {
					q.Name += "."
				}

				log.Debugf("Find ip addresses for domain: ", q.Name)

				// Generate A and AAAA records
				ips, err := endpoint.Get(util.Trim(q.Name, `.`))
				if err != nil {
					log.Error(err)
					w.WriteMsg(m)
					return
				}

				if ips == nil {
					w.WriteMsg(m)
					return
				}

				log.Debugf("Ips list: %#v for %s", ips, q.Name)

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
		log.Debug(`dns.OpcodeUpdate`)
	}

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacMD5, 300, time.Now().Unix())
		} else {
			log.Error(`Status `, w.TsigStatus().Error())
		}
	}

	log.Info(`Send info `, m)
	w.WriteMsg(m)
}
