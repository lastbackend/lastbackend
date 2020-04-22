// +build !windows

package rootless

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"errors"
	"golang.org/x/sys/unix"
	"github.com/lastbackend/lastbackend/tools/logger"
)

func setupMounts(stateDir string) error {
	mountMap := [][]string{
		{"/run", ""},
		{"/var/run", ""},
		{"/var/log", filepath.Join(stateDir, "logs")},
		{"/var/lib/cni", filepath.Join(stateDir, "cni")},
		{"/var/lib/kubelet", filepath.Join(stateDir, "kubelet")},
		{"/etc/rancher", filepath.Join(stateDir, "etc", "rancher")},
	}

	for _, v := range mountMap {
		if err := setupMount(v[0], v[1]); err != nil {
			return errors.New(fmt.Sprintf("%v: failed to setup mount %s => %s", err, v[0], v[1]))
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
		return errors.New(fmt.Sprintf("%v: failed to create directory %s", err, toCreate))
	}

	log.Debug("Mounting none ", toCreate, " tmpfs")
	if err := unix.Mount("none", toCreate, "tmpfs", 0, ""); err != nil {
		return errors.New(fmt.Sprintf("%v: failed to mount tmpfs to %s", err, toCreate))
	}

	if err := os.MkdirAll(target, 0700); err != nil {
		return errors.New(fmt.Sprintf("%v: failed to create directory %s", err, target))
	}

	if dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.New(fmt.Sprintf("%v: failed to create directory %s", err, dir))
	}

	log.Debug("Mounting ", dir, target, " none bind")
	return unix.Mount(dir, target, "none", unix.MS_BIND, "")
}
