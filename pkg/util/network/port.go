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

package network

import (
	"strings"
	"strconv"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func ParsePortMap(s string) (int, string, error) {

	var (
		port int
		proto string
		err error
	)

	pm := strings.Split(s, "/")
	switch len(pm) {
	case 0:
		break
	case 1:
		port, err = strconv.Atoi(pm[0])
		if err != nil {
			break
		}
		proto = "tcp"
		break
	case 2:
		port, err = strconv.Atoi(pm[0])
		if err != nil {
			return port, proto, err
		}
		proto = strings.ToLower(pm[1])
		break
	default:
		err = errors.New("Invalid port map declaration")
		return port, proto, err
	}

	return port, proto, nil
}
