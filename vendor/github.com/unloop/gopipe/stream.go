package stream

import (
	"io"
	"net/http"
)

type Stream struct {
	buffer int       // Buffer size
	w      io.Writer // Underlying writer to send data to
	c      int       // Number of bytes written since last call to Count()
	done   chan bool // Done pipe
}

func New(w io.Writer) *Stream {
	s := new(Stream)
	s.w = w
	s.buffer = 1024
	s.done = make(chan bool, 1)
	return s
}

func (s *Stream) SetBuffer(size int) *Stream {
	s.buffer = size
	return s
}

// Write: standard io.Writer interface.  To use this package call
// Write continually.  This will both count the bytes written and
// write to the underlying writer.
func (s *Stream) Write(p []byte) (n int, err error) {
	n, err = s.w.Write(p)
	if err != nil {
		return n, err
	}

	if f, ok := s.w.(http.Flusher); ok {
		f.Flush()
	}

	return n, nil
}

// Pipe io.ReadCloser to io.Writer
func (s *Stream) Pipe(reader *io.ReadCloser) {

	var (
		buffer = make([]byte, s.buffer)
	)

	for {
		select {
		case <-s.done:
			(*reader).Close()
			return
		default:
			n, err := (*reader).Read(buffer)
			if err != nil {
				(*reader).Close()
				break
			}

			s.Write(buffer[0:n])

			for i := 0; i < n; i++ {
				buffer[i] = 0
			}
		}
	}
}

// Close pipe streaming
func (s *Stream) Close() {
	s.done <- true
	return
}
