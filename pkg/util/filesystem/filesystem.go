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

package filesystem

import (
	"bytes"
	"io"
	"os"
)

const (
	LineSeparator = '\n'
)

func LineSeek(lines int, f *os.File) (int64, error) {

	count := 0
	pos, err := f.Seek(0, io.SeekEnd)

	if err != nil {
		return 0, err
	}

	chunk := 4096

	b1 := make([]byte, 1)
	if _, err := f.ReadAt(b1, pos-1); err != nil {
		return 0, err
	}

	if '\n' == b1[0] {
		pos = pos - 1
	}

	for {
		rf := pos - int64(chunk)
		ids := make([]int64, 0)

		if rf <= 0 {
			chunk += int(rf)
			rf = 0
			ids = append(ids, 0)
			count++
		}

		b := make([]byte, chunk)

		_, err := f.ReadAt(b, rf)

		if err != nil {
			return 0, err
		}

		i := 0

		for {
			pos := bytes.IndexByte(b[i:], LineSeparator)
			if pos == -1 {
				break
			}
			i = i + pos + 1
			ids = append(ids, int64(i))
			count++
		}

		var lpos int64

		if len(ids) == 0 {
			lpos = 0
		} else {
			lpos = ids[0]
		}

		if count == lines {
			return lpos + rf, nil
		}

		if count > lines {
			left := count - lines
			pos := ids[left]
			return pos + rf, nil

		}

		if rf == 0 {
			return 0, nil
		}
		pos = lpos + rf - 1
	}

}
