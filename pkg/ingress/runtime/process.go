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

package runtime

import (
	"bytes"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	hPid         string = "/var/run/haproxy.pid"
	IPTablesExec string = "iptables"
)

type Process struct {
	process *os.Process
}

func (hp *Process) manage() error {

	var (
		process *os.Process
	)

	pid := hp.getPid()

	process, err := os.FindProcess(int(pid))
	if err != nil {
		fmt.Printf("Failed to find process: %s\n", err)
		return err
	}

	if err := process.Signal(syscall.Signal(0)); err == nil {
		hp.process = process
	}

	if hp.process == nil {
		log.Debug("running process not found: start new")
		if process, err = hp.start(); err != nil {
			fmt.Printf("Failed to start process: %s", err)
			return err
		}
		hp.process = process
	}

	go func() {
		for {
			if err := hp.process.Signal(syscall.Signal(0)); err != nil {
				log.Debug("process exited")
				if process, err = hp.start(); err != nil {
					fmt.Printf("Failed to start process: %s", err)
				}
				hp.process = process
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return nil
}

func (hp *Process) start() (*os.Process, error) {

	log.Debug("start new process")

	bin := envs.Get().GetHaproxy()
	_, path, pid := envs.Get().GetTemplate()

	if pid == types.EmptyString {
		pid = hPid
	}

	cmd := exec.Command(bin, "-f", filepath.Join(path, ConfigName), "-D", "-p", pid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Errorf("failed to start haproxy: %s", err.Error())
		return nil, err
	}
	time.Sleep(1 * time.Second)

	return cmd.Process, nil
}

func (hp *Process) reload() error {

	log.Debug("reload haproxy process")
	var ports = make(map[int]bool, 0)
	routes := envs.Get().GetState().Routes().GetRouteManifests()

	for _, r := range routes {

		for _, rule := range r.Rules {
			if _, ok := ports[rule.Port]; !ok {
				ports[rule.Port] = true
			}
		}
	}

	if hp.process == nil {
		log.Error("process is not running")
		return nil
	}

	defer func() {
		//for port := range ports {
		//	c := exec.Command(IPTablesExec, "-D", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port), "--syn", "-j", "DROP")
		//	c.Start()
		//}
	}()

	bin := envs.Get().GetHaproxy()
	_, path, pidpath := envs.Get().GetTemplate()
	if pidpath == types.EmptyString {
		pidpath = hPid
	}

	pid := hp.getPid()
	cmd := exec.Command(bin, "-f", filepath.Join(path, ConfigName), "-p", pidpath, "-sf", fmt.Sprintf("%d", pid))
	cmd.Stdout = os.Stdout

	//for port := range ports {
	//	c := exec.Command(IPTablesExec, "-I", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port), "--syn", "-j", "DROP")
	//	c.Start()
	//}

	err := cmd.Start()
	if err != nil {
		log.Errorf("failed to start haproxy: %s", err.Error())
		return err
	}

	hp.process = cmd.Process
	return nil
}

func (hp *Process) getPid() int {

	_, _, pidpath := envs.Get().GetTemplate()

	if pidpath == types.EmptyString {
		pidpath = hPid
	}

	pf, err := os.Open(pidpath)
	if err != nil && !os.IsNotExist(err) {
		log.Errorf("can not open pid file: %s", err.Error())
		return 0
	}
	pf.Close()

	if os.IsNotExist(err) {
		return 0
	}

	d, err := ioutil.ReadFile(pidpath)
	if err != nil {
		return 0
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		log.Errorf("error parsing pid from %s: %s", pidpath, err)
		return 0
	}

	return pid
}
