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
