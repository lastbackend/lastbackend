package stream

import (
	"io"
	"net/http"
)

type Stream struct {
	buffer int       // Buffer size
	w      io.Writer // Underlying writer to send data to
	c      int       // Number of bytes written since last call to Count()
}

func New(w io.Writer) (s *Stream) {
	s = new(Stream)
	s.w = w
	s.buffer = 1024
	return
}

func (s *Stream) SetBuffer(size int) {
	s.buffer = size
}

// Write: standard io.Writer interface.  To use this package call
// Write continually.  This will both count the bytes written and
// write to the underlying writer.
func (s *Stream) Write(p []byte) (n int, err error) {
	f, _ := s.w.(http.Flusher)

	n, err = s.w.Write(p)
	if err != nil {
		return n, err
	}

	f.Flush()

	return n, nil
}

// Pipe io.ReadCloser to io.Writer
func (s *Stream) Pipe(reader *io.ReadCloser) {

	buffer := make([]byte, s.buffer)

	for {
		n, err := (*reader).Read(buffer)
		if err != nil {
			(*reader).Close()
			break
		}

		data := buffer[0:n]

		s.Write(data)

		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}
