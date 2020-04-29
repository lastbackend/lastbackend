// +build linux

package rootless

import (
	"context"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/rootless-containers/rootlesskit/pkg/child"
	"github.com/rootless-containers/rootlesskit/pkg/copyup/tmpfssymlink"
	"github.com/rootless-containers/rootlesskit/pkg/network/slirp4netns"
	"github.com/rootless-containers/rootlesskit/pkg/parent"
	"github.com/rootless-containers/rootlesskit/pkg/port/builtin"
)

var (
	pipeFD   = "_LB_ROOTLESS_FD"
	childEnv = "_LB_ROOTLESS_SOCK"
	Sock     = ""
)

func Rootless(stateDir string) error {

	log := logger.WithContext(context.Background())

	defer func() {
		os.Unsetenv(pipeFD)
		os.Unsetenv(childEnv)
	}()

	hasFD := os.Getenv(pipeFD) != ""
	hasChildEnv := os.Getenv(childEnv) != ""

	if hasFD {
		log.Debug("Running rootless child")
		childOpt, err := createChildOpt()
		if err != nil {
			log.Fatal(err)
		}
		if err := child.Child(*childOpt); err != nil {
			log.Fatal("child died", err)
		}
	}

	if hasChildEnv {
		Sock = os.Getenv(childEnv)
		log.Debug("Running rootless process")
		return setupMounts(stateDir)
	}

	log.Debug("Running rootless parent")
	parentOpt, err := createParentOpt(filepath.Join(stateDir, "rootless"))
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv(childEnv, filepath.Join(parentOpt.StateDir, parent.StateFileAPISock))
	if err := parent.Parent(*parentOpt); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

	return nil
}

func parseCIDR(s string) (*net.IPNet, error) {
	if s == "" {
		return nil, nil
	}
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	if !ip.Equal(ipnet.IP) {
		return nil, fmt.Errorf("cidr must be like 10.0.2.0/24, not like 10.0.2.100/24")
	}
	return ipnet, nil
}

func createParentOpt(stateDir string) (*parent.Opt, error) {
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to mkdir %s err: %v", stateDir, err)
	}

	stateDir, err := ioutil.TempDir("", "rootless")
	if err != nil {
		return nil, err
	}

	opt := &parent.Opt{
		StateDir:    stateDir,
		CreatePIDNS: true,
	}

	mtu := 0
	ipnet, err := parseCIDR("10.41.0.0/16")
	if err != nil {
		return nil, err
	}
	disableHostLoopback := true
	binary := "slirp4netns"
	if _, err := exec.LookPath(binary); err != nil {
		return nil, err
	}
	opt.NetworkDriver = slirp4netns.NewParentDriver(binary, mtu, ipnet, disableHostLoopback, "", false, false)
	opt.PortDriver, err = builtin.NewParentDriver(&debugWriter{}, stateDir)
	if err != nil {
		return nil, err
	}

	opt.PipeFDEnvKey = pipeFD

	return opt, nil
}

type debugWriter struct {
}

func (w *debugWriter) Write(p []byte) (int, error) {
	log := logger.WithContext(context.Background())
	s := strings.TrimSuffix(string(p), "\n")
	log.Debug(s)
	return len(p), nil
}

func createChildOpt() (*child.Opt, error) {
	opt := &child.Opt{}
	opt.TargetCmd = os.Args
	opt.PipeFDEnvKey = pipeFD
	opt.NetworkDriver = slirp4netns.NewChildDriver()
	opt.PortDriver = builtin.NewChildDriver(&debugWriter{})
	opt.CopyUpDirs = []string{"/etc", "/run", "/var/lib"}
	opt.CopyUpDriver = tmpfssymlink.NewChildDriver()
	opt.MountProcfs = true
	opt.Reaper = true
	return opt, nil
}

func setupMounts(stateDir string) error {
	mountMap := [][]string{
		{"/run", ""},
		{"/var/run", ""},
		{"/var/log", filepath.Join(stateDir, "logs")},
		{"/var/lib/lastbackend", filepath.Join(stateDir, "lastbackend")},
	}

	for _, v := range mountMap {
		if err := setupMount(v[0], v[1]); err != nil {
			return fmt.Errorf("%v: failed to setup mount %s => %s", err, v[0], v[1])
		}
	}

	return nil
}

func setupMount(target, dir string) error {
	log := logger.WithContext(context.Background())
	log.Infof("Init minion service")

	toCreate := target
	for {
		if toCreate == "/" {
			return fmt.Errorf("missing /%s on the root filesystem", strings.Split(target, "/")[0])
		}

		if err := os.MkdirAll(toCreate, 0700); err == nil {
			break
		}

		toCreate = filepath.Base(toCreate)
	}

	if err := os.MkdirAll(toCreate, 0700); err != nil {
		return fmt.Errorf("%v: failed to create directory %s", err, toCreate)
	}

	log.Debug("Mounting none ", toCreate, " tmpfs")
	if err := unix.Mount("none", toCreate, "tmpfs", 0, ""); err != nil {
		return fmt.Errorf("%v: failed to mount tmpfs to %s", err, toCreate)
	}

	if err := os.MkdirAll(target, 0700); err != nil {
		return fmt.Errorf("%v: failed to create directory %s", err, target)
	}

	if dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("%v: failed to create directory %s", err, dir)
	}

	log.Debug("Mounting ", dir, target, " none bind")
	return unix.Mount(dir, target, "none", unix.MS_BIND, "")
}
