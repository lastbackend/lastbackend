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

package local

import (
	"net"
	"strings"
)

func getInterface() *net.Interface {

	ifas, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for _, ifa := range ifas {
		ips, err := ifa.Addrs()
		if err != nil {
			return nil
		}

		for _, ip := range ips {
			if ip.String() == "127.0.0.1" || strings.HasPrefix(ifa.Name, "lo") {
				return &ifa
			}
		}
	}

	return nil
}
