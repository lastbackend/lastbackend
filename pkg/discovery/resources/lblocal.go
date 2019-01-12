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

package resources

import (
	"context"
	"net"
	"regexp"
	"time"

	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

func lbLocal(w dns.ResponseWriter, r *dns.Msg) {

	log.V(logLevel).Debugf("%s:lb.local:> dns request `lb.local.`", logPrefix)

	var (
		err error
		v4  bool
		rr  dns.RR
	)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())

	switch r.Opcode {
	case dns.OpcodeQuery:
		log.V(logLevel).Debugf("%s:lb.local:> dns.OpcodeQuery", logPrefix)

		for _, q := range m.Question {

			switch r.Question[0].Qtype {
			case dns.TypeTXT:
				log.V(logLevel).Debugf("%s:lb.local:> get txt type query", logPrefix)
				t := new(dns.TXT)
				t.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}
				m.Authoritative = true
				m.Answer = append(m.Answer, t)
				m.Extra = append(m.Extra, rr)
			default:
				log.V(logLevel).Debugf("%s:lb.local:> get unknown query type", logPrefix)
				fallthrough
			case dns.TypeAAAA, dns.TypeA:
				log.V(logLevel).Debugf("%s:lb.local:> get A or AAAA type query", logPrefix)

				if q.Name[len(q.Name)-1:] != "." {
					q.Name += "."
				}

				log.V(logLevel).Debugf("%s:lb.local:> find ip addresses for domain: %s", logPrefix, q.Name)

				// GenerateConfig A and AAAA records
				ips := make([]net.IP, 0)

				endpoint := util.Trim(q.Name, `.`)
				item := envs.Get().GetCache().Endpoint().Get(endpoint)

				if item != nil {
					data := util.RemoveDuplicates(item)
					ips, err = util.ConvertStringIPToNetIP(data)
					if err != nil {
						log.Error(err)
						return
					}
				} else {

					rg, _ := regexp.Compile("^(.+)\\.(.+)\\.lb\\.local$")
					match := rg.FindStringSubmatch(endpoint)

					if len(match) != 0 {

						log.V(logLevel).Debugf("%s:lb.local:> find endpoint %s:%s", logPrefix, match[2], match[1])

						e, err := em.Get(match[2], match[1])
						if err != nil {
							log.V(logLevel).Errorf("%s:lb.local:> get endpoint `%s` err: %v", endpoint, logPrefix, err)
						}

						if e != nil {
							envs.Get().GetCache().Endpoint().Set(endpoint, []string{e.Spec.IP})

							ips, err = util.ConvertStringIPToNetIP([]string{e.Spec.IP})
							if err != nil {
								log.Error(err)
								return
							}
						}

						if e == nil {
							break
						}
					}

				}

				if len(ips) == 0 {
					defaultIPs := viper.GetStringSlice("discovery.default_ips")
					ips, err = util.ConvertStringIPToNetIP(defaultIPs)
					if err != nil {
						log.Error(err)
						return
					}
				}

				log.V(logLevel).Debugf("%s:lb.local:> ips list: %s for %s", logPrefix, ips, q.Name)

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

					m.Authoritative = true
					m.RecursionAvailable = true
					m.RecursionDesired = true
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	case dns.OpcodeUpdate:
		log.V(logLevel).Debugf("%s:lb.local:> dns.OpcodeUpdate", logPrefix)
	}

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacMD5, 300, time.Now().Unix())
		} else {
			log.V(logLevel).Errorf("%s:lb.local:> tsig status err: %s", logPrefix, w.TsigStatus())
		}
	}

	log.V(logLevel).Debugf("%s:lb.local:> send message info  %#v", logPrefix, m)

	w.WriteMsg(m)
}
