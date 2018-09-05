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

package resource

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// MB - MegaByte size
	MB = 1000 * 1000
	// MIB - MegaByte size
	MIB = 1024 * 1024
	// GB - GigaByte size
	GB = 1000 * 1000 * 1000
	// GIB - GibiByte size
	GIB = 1024 * 1024 * 1024
)

// parseResource - parse resource size string
// mb,mib,gb,gib,kb,kib,
func DecodeResource(res string) (int64, error) {

	var i = int64(1024)

	rq := regexp.MustCompile("([0-9]+)\\w*.*")
	rt := regexp.MustCompile("\\d*(\\w+)?.*")
	mq := rq.FindStringSubmatch(res)
	mt := rt.FindStringSubmatch(res)

	if len(mt) == 2 {
		switch strings.ToLower(mt[1]) {
		case "mb":
			i = MB
			break
		case "mib":
			i = MIB
			break
		case "gb":
			i = GB
			break
		case "gib":
			i = GIB
			break
		}
	}

	if len(mq) == 2 {

		q, err := strconv.ParseInt(mq[1], 10, 64)
		if err != nil {
			return i, err
		}
		i*=q
	}

	return i, nil
}

func EncodeResource(res int64) string {

	var r string

	if res < MB {
		return fmt.Sprintf("%d", res/1024)
	}

	if res > MB && res < GB {
		return fmt.Sprintf("%dMB", res/MB)
	}

	if res > GB  {
		return fmt.Sprintf("%dGB", res/GB)
	}

	return r
}
