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

package util

import (
	"net"
	"strings"
)

func RemoveDuplicates(data []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range data {
		if encountered[data[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[data[v]] = true
			// Append to result slice.
			result = append(result, data[v])
		}
	}

	// Return the new slice.
	return result
}

func Trim(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func ConvertStringIPToNetIP(data []string) ([]net.IP, error) {
	var ips = []net.IP{}
	for i := range data {
		ips = append(ips, net.ParseIP(data[i]))
	}
	return ips, nil
}
