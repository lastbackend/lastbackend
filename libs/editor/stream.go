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

package editor

import (
	"io"
	"strings"
)

type Stream struct {
	c      []string
	reader io.Reader
}

func (c Stream) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

func (c Stream) String() string {
	return strings.Join(c.c, "\n")
}

func (c Stream) Bytes() []byte {
	return []byte(strings.Join(c.c, "\n"))
}

func (c Stream) Length() int {
	return len(c.c)
}
