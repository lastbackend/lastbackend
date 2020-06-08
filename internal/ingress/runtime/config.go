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

package runtime

import (
	"bytes"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/ingress/envs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"

	"io/ioutil"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	ConfigName      = "haproxy.cfg"
	logConfigPrefix = "runtime:config"
)

type conf struct {
	Stats struct {
		Port     uint16
		Username string
		Password string
	}
	Resolvers map[string]uint16
	Frontend  map[uint16]*confFrontend
	Backend   map[string]*confBackend
}

type confFrontend struct {
	Type  string
	Rules map[string]map[string]string
}

type confBackend struct {
	Domain   string
	Type     string
	Upstream string
	Port     uint16
}

func NewHAProxyConfig(port uint16, username, password string) *conf {
	c := new(conf)
	c.Stats.Port = port
	c.Stats.Username = username
	c.Stats.Password = password
	c.Resolvers = make(map[string]uint16, 0)
	c.Frontend = make(map[uint16]*confFrontend, 0)
	c.Backend = make(map[string]*confBackend, 0)
	return c
}

func (c conf) Check() error {

	log.Debug("config check")
	var (
		_, path, _ = envs.Get().GetTemplate()
	)

	cfgPath := filepath.Join(path, ConfigName)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Debug("config not found: create new")
		return c.Sync()
	}

	return nil
}

func (c conf) Sync() error {

	log.Debug("config sync")

	var (
		routes       = envs.Get().GetState().Routes().GetRouteManifests()
		tpl, path, _ = envs.Get().GetTemplate()
	)

	log.Debugf("Update routes: %d", len(routes))

	var cfg = conf{}
	cfg.Stats.Username = c.Stats.Username
	cfg.Stats.Password = c.Stats.Password
	cfg.Stats.Port = c.Stats.Port

	if cfg.Stats.Username != models.EmptyString && cfg.Stats.Password != models.EmptyString {
		if cfg.Stats.Port == 0 {
			cfg.Stats.Port = 9000
		}
	}

	cfg.Resolvers = envs.Get().GetResolvers()
	cfg.Frontend = make(map[uint16]*confFrontend, 0)
	cfg.Backend = make(map[string]*confBackend, 0)

	cfg.Frontend[80] = new(confFrontend)
	cfg.Frontend[80].Type = "http"
	cfg.Frontend[80].Rules = make(map[string]map[string]string, 0)

	for n, r := range routes {

		log.Debugf("route configure: %s", n)

		var tp string
		switch r.Port {
		case 80:
			tp = "http"
			break
		case 443:
			tp = "https"
			break
		default:
			tp = "tcp"
		}

		if r.Port == 0 {
			continue
		}

		var frontend *confFrontend

		if _, ok := cfg.Frontend[r.Port]; ok {
			frontend = cfg.Frontend[r.Port]
		} else {
			frontend = new(confFrontend)
			frontend.Type = tp
			frontend.Rules = make(map[string]map[string]string, 0)
			cfg.Frontend[r.Port] = frontend
		}

		if _, ok := frontend.Rules[r.Endpoint]; !ok {
			frontend.Rules[r.Endpoint] = make(map[string]string, 0)
		}

		for _, b := range r.Rules {

			name := fmt.Sprintf("%s_%d", strings.Replace(n, ":", "_", -1), b.Port)
			log.Debugf("create new backend: %s", name)

			backend := new(confBackend)
			backend.Type = tp
			backend.Port = uint16(b.Port)
			backend.Upstream = b.Upstream
			backend.Domain = r.Endpoint

			cfg.Backend[name] = backend
			frontend.Rules[r.Endpoint][b.Path] = name
		}

	}

	buf := &bytes.Buffer{}
	tpl.Execute(buf, cfg)
	log.Debugf("config path: %s", path)

	var (
		f   *os.File
		err error
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debugf("config directory does not exists: %s", path)
		if err := os.MkdirAll(path, 0644); err != nil {
			log.Errorf("can not be create config dir: %s", err.Error())
			return err
		}
	}

	cfgPath := filepath.Join(path, ConfigName)
	testPath := fmt.Sprintf("%s.test", cfgPath)

	f, err = os.Open(testPath)
	if os.IsNotExist(err) {
		log.Debugf("config file not exists: %s", testPath)
		f, err = os.Create(testPath)
		if err != nil {
			log.Errorf("can not be create config file: %s", err.Error())
		}
	}
	f.Close()

	if err := ioutil.WriteFile(testPath, buf.Bytes(), 0644); err != nil {
		log.Errorf("can no write test config: %s", err.Error())
		return err
	}

	if err := c.Validate(testPath); err != nil {
		log.Errorf("config is not working (%s)", err.Error())
		return err
	}

	f, err = os.Open(cfgPath)
	if os.IsNotExist(err) {
		log.Debugf("config file not exists: %s", cfgPath)
		f, err = os.Create(cfgPath)
		if err != nil {
			log.Errorf("can not be create config file: %s", err.Error())
		}
	}
	f.Close()

	return ioutil.WriteFile(cfgPath, buf.Bytes(), 0644)
}

func (conf) Validate(path string) error {

	log.Debugf("%s:> config validate", logConfigPrefix)

	var hpbin = envs.Get().GetHaproxy()

	cmd := exec.Command(hpbin, "-c", "-V", "-f", path)
	err := cmd.Start()

	if err != nil {
		log.Errorf("can not be check config: %s", err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 0 {
					return errors.New(string(exiterr.Stderr))
				}
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	return nil
}
