package editor

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Context struct {
	Name   string
	input  *Stream
	output *Stream
}

func (c *Context) Write(content io.Reader, split bufio.SplitFunc) error {

	var (
		err error
	)

	c.input, err = c.readerToStream(content, split)
	if err != nil {
		return err
	}

	c.Name, err = c.write(c.input.reader)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) Read(filename string, split bufio.SplitFunc) error {

	var (
		err error
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	c.output, err = c.readerToStream(file, split)
	if err != nil {
		return err
	}

	return nil
}

func (Context) write(content io.Reader) (string, error) {

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	io.Copy(f, content)
	f.Close()

	return f.Name(), nil
}

func (c *Context) readerToStream(content io.Reader, split bufio.SplitFunc) (*Stream, error) {

	var stream = new(Stream)

	scanner := bufio.NewScanner(content)
	scanner.Split(split)

	for scanner.Scan() {
		stream.c = append(stream.c, scanner.Text())
	}

	stream.reader = strings.NewReader(stream.String())

	return stream, scanner.Err()
}
