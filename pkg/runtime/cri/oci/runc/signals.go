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
	"os"
	"os/signal"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/runc/libcontainer/utils"
	"golang.org/x/sys/unix"
)

const signalBufferSize = 2048

// newSignalHandler returns a signal handler for processing SIGCHLD and SIGWINCH signals
// while still forwarding all other signals to the process.
// If notifySocket is present, use it to read systemd notifications from the container and
// forward them to notifySocketHost.
func newSignalHandler(enableSubreaper bool, notifySocket *notifySocket) *signalHandler {
	if enableSubreaper {
		// set us as the subreaper before registering the signal handler for the container
		if err := system.SetSubreaper(1); err != nil {
			fmt.Println(err)
		}
	}
	// ensure that we have a large buffer size so that we do not miss any signals
	// in case we are not processing them fast enough.
	s := make(chan os.Signal, signalBufferSize)
	// handle all signals for the process.
	signal.Notify(s)
	return &signalHandler{
		signals:      s,
		notifySocket: notifySocket,
	}
}

// exit models a process exit status with the pid and exit status.
type exit struct {
	pid    int
	status int
}

type signalHandler struct {
	signals      chan os.Signal
	notifySocket *notifySocket
}

// forward handles the main signal event loop forwarding, resizing, or reaping depending on the signal received.
func (h *signalHandler) forward(process *libcontainer.Process, tty *tty, detach bool) (int, error) {
	// make sure we know the pid of our main process so that we can return
	// after it dies.
	if detach && h.notifySocket == nil {
		return 0, nil
	}

	pid1, err := process.Pid()
	if err != nil {
		return -1, err
	}

	if h.notifySocket != nil {
		if detach {
			h.notifySocket.run(pid1)
			return 0, nil
		}
		h.notifySocket.run(os.Getpid())
		go h.notifySocket.run(0)
	}

	// Perform the initial tty resize. Always ignore errors resizing because
	// stdout might have disappeared (due to races with when SIGHUP is sent).
	_ = tty.resize()
	// Handle and forward signals.
	for s := range h.signals {
		switch s {
		case unix.SIGWINCH:
			// Ignore errors resizing, as above.
			_ = tty.resize()
		case unix.SIGCHLD:
			exits, err := h.reap()
			if err != nil {
				fmt.Println(err)
			}
			for _, e := range exits {
				fmt.Println(fmt.Sprintf("process exited (pid: %s, status: %d)", e.pid, e.status))
				if e.pid == pid1 {
					// call Wait() on the process even though we already have the exit
					// status because we must ensure that any of the go specific process
					// fun such as flushing pipes are complete before we return.
					process.Wait()
					return e.status, nil
				}
			}
		default:
			fmt.Println(fmt.Sprintf("sending signal to process %s", s))
			if err := unix.Kill(pid1, s.(unix.Signal)); err != nil {
				fmt.Println(err)
			}
		}
	}
	return -1, nil
}

// reap runs wait4 in a loop until we have finished processing any existing exits
// then returns all exits to the main event loop for further processing.
func (h *signalHandler) reap() (exits []exit, err error) {
	var (
		ws  unix.WaitStatus
		rus unix.Rusage
	)
	for {
		pid, err := unix.Wait4(-1, &ws, unix.WNOHANG, &rus)
		if err != nil {
			if err == unix.ECHILD {
				return exits, nil
			}
			return nil, err
		}
		if pid <= 0 {
			return exits, nil
		}
		exits = append(exits, exit{
			pid:    pid,
			status: utils.ExitStatus(ws),
		})
	}
}
