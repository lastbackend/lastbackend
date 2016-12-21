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
