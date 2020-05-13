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
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type notifySocket struct {
	socket     *net.UnixConn
	host       string
	socketPath string
}

func newNotifySocket(id, rootPath, notifySocketHost string) *notifySocket {
	if notifySocketHost == "" {
		return nil
	}
	notifySocket := &notifySocket{
		socket:     nil,
		host:       notifySocketHost,
		socketPath: filepath.Join(filepath.Join(rootPath, id), "notify", "notify.sock"),
	}
	return notifySocket
}

func notifySocketStart(id, rootPath, notifySocketHost string) (*notifySocket, error) {
	notifySocket := newNotifySocket(id, rootPath, notifySocketHost)
	if notifySocket == nil {
		return nil, nil
	}

	if err := notifySocket.bindSocket(); err != nil {
		return nil, err
	}
	return notifySocket, nil
}

func (s *notifySocket) Close() error {
	return s.socket.Close()
}

func (s *notifySocket) setupSpec(spec *specs.Spec) error {
	pathInContainer := filepath.Join("/run/notify", path.Base(s.socketPath))
	mount := specs.Mount{
		Destination: path.Dir(pathInContainer),
		Source:      path.Dir(s.socketPath),
		Options:     []string{"bind", "nosuid", "noexec", "nodev", "ro"},
	}
	spec.Mounts = append(spec.Mounts, mount)
	spec.Process.Env = append(spec.Process.Env, fmt.Sprintf("NOTIFY_SOCKET=%s", pathInContainer))
	return nil
}

func (s *notifySocket) bindSocket() error {
	addr := net.UnixAddr{
		Name: s.socketPath,
		Net:  "unixgram",
	}

	socket, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		return err
	}

	err = os.Chmod(s.socketPath, 0777)
	if err != nil {
		socket.Close()
		return err
	}

	s.socket = socket
	return nil
}

func (s *notifySocket) setupSocketDirectory() error {
	return os.Mkdir(path.Dir(s.socketPath), 0755)
}

func (n *notifySocket) waitForContainer(container libcontainer.Container) error {
	s, err := container.State()
	if err != nil {
		return err
	}
	return n.run(s.InitProcessPid)
}

func (n *notifySocket) run(pid1 int) error {
	if n.socket == nil {
		return nil
	}
	notifySocketHostAddr := net.UnixAddr{Name: n.host, Net: "unixgram"}
	client, err := net.DialUnix("unixgram", nil, &notifySocketHostAddr)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	fileChan := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 4096)
			r, err := n.socket.Read(buf)
			if err != nil {
				return
			}
			got := buf[0:r]
			// systemd-ready sends a single datagram with the state string as payload,
			// so we don't need to worry about partial messages.
			for _, line := range bytes.Split(got, []byte{'\n'}) {
				if bytes.HasPrefix(got, []byte("READY=")) {
					fileChan <- line
					return
				}
			}

		}
	}()

	for {
		select {
		case <-ticker.C:
			_, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid1)))
			if err != nil {
				return nil
			}
		case b := <-fileChan:
			var out bytes.Buffer
			_, err = out.Write(b)
			if err != nil {
				return err
			}

			_, err = out.Write([]byte{'\n'})
			if err != nil {
				return err
			}

			_, err = client.Write(out.Bytes())
			if err != nil {
				return err
			}

			// now we can inform systemd to use pid1 as the pid to monitor
			newPid := fmt.Sprintf("MAINPID=%d\n", pid1)
			client.Write([]byte(newPid))
			return nil
		}
	}
}
