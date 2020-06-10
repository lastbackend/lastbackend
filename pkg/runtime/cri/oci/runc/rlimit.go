// +build linux
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

package runc

import (
	"fmt"

	"golang.org/x/sys/unix"
)

var rlimitMap = map[string]int{
	"RLIMIT_CPU":        unix.RLIMIT_CPU,
	"RLIMIT_FSIZE":      unix.RLIMIT_FSIZE,
	"RLIMIT_DATA":       unix.RLIMIT_DATA,
	"RLIMIT_STACK":      unix.RLIMIT_STACK,
	"RLIMIT_CORE":       unix.RLIMIT_CORE,
	"RLIMIT_RSS":        unix.RLIMIT_RSS,
	"RLIMIT_NPROC":      unix.RLIMIT_NPROC,
	"RLIMIT_NOFILE":     unix.RLIMIT_NOFILE,
	"RLIMIT_MEMLOCK":    unix.RLIMIT_MEMLOCK,
	"RLIMIT_AS":         unix.RLIMIT_AS,
	"RLIMIT_LOCKS":      unix.RLIMIT_LOCKS,
	"RLIMIT_SIGPENDING": unix.RLIMIT_SIGPENDING,
	"RLIMIT_MSGQUEUE":   unix.RLIMIT_MSGQUEUE,
	"RLIMIT_NICE":       unix.RLIMIT_NICE,
	"RLIMIT_RTPRIO":     unix.RLIMIT_RTPRIO,
	"RLIMIT_RTTIME":     unix.RLIMIT_RTTIME,
}

func strToRlimit(key string) (int, error) {
	rl, ok := rlimitMap[key]
	if !ok {
		return 0, fmt.Errorf("wrong rlimit value: %s", key)
	}
	return rl, nil
}
