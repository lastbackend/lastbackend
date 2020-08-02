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

package resources
//
//import (
//	"time"
//
//	"github.com/lastbackend/lastbackend/tools/log"
//	"github.com/miekg/dns"
//	"net"
//)
//
//func other(w dns.ResponseWriter, r *dns.Msg) {
//
//	log.Debugf("%s:other:> dns request `.`", logPrefix)
//
//	var (
//		v4 bool
//		rr dns.RR
//	)
//
//	m := new(dns.Msg)
//	m.SetReply(r)
//	m.Compress = false
//
//	switch r.Opcode {
//	case dns.OpcodeQuery:
//		log.Debugf("%s:other:> dns.OpcodeQuery", logPrefix)
//
//		for _, q := range m.Question {
//
//			switch r.Question[0].Qtype {
//			case dns.TypeTXT:
//				log.Debugf("%s:other:> get txt type query", logPrefix)
//				t := new(dns.TXT)
//				t.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}
//				m.Answer = append(m.Answer, t)
//				m.Extra = append(m.Extra, rr)
//			default:
//				log.Debugf("%s:other:> get unknown query type", logPrefix)
//				fallthrough
//			case dns.TypeAAAA, dns.TypeA:
//				log.Debugf("%s:other:> get A or AAAA type query", logPrefix)
//
//				if q.Name[len(q.Name)-1:] != "." {
//					q.Name += "."
//				}
//
//				// GenerateConfig A and AAAA records
//				ips := make([]net.IP, 0)
//
//				log.Debugf("%s:other:> ips list: %#v for %s", logPrefix, ips, q.Name)
//
//				for _, ip := range ips {
//					v4 = ip.To4() != nil
//
//					if v4 {
//						rr = new(dns.A)
//						rr.(*dns.A).Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
//						rr.(*dns.A).A = ip.To4()
//					} else {
//						rr = new(dns.AAAA)
//						rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0}
//						rr.(*dns.AAAA).AAAA = ip
//					}
//
//					m.Answer = append(m.Answer, rr)
//				}
//			}
//		}
//	case dns.OpcodeUpdate:
//		log.Debugf("%s:other:> dns.OpcodeUpdate", logPrefix)
//	}
//
//	if r.IsTsig() != nil {
//		if w.TsigStatus() == nil {
//			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacMD5, 300, time.Now().Unix())
//		} else {
//			log.Errorf("%s:other:> tsig status err: %s", logPrefix, w.TsigStatus())
//		}
//	}
//
//	log.Debugf("%s:other:> send message info  %#v", logPrefix, m)
//
//	if err := w.WriteMsg(m); err != nil {
//		log.Errorf("%s:other:> write message err: %v", logPrefix, err)
//	}
//}
