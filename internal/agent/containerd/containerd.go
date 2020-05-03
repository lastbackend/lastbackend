package containerd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/lastbackend/lastbackend/internal/agent/containerd/templates"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/opencontainers/runc/libcontainer/system"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"github.com/containerd/containerd/services/server"
	srvconfig "github.com/containerd/containerd/services/server/config"
)

const (
	defaultNamespace = "lstbknd.net"
)

type Containerd struct {
	ctx    context.Context
	cancel context.CancelFunc

	address        string
	opt            string
	config         string
	template       string
	registry       string
	state          string
	root           string
	disableSELinux bool
	images         string
	log            string
}

type Config struct {
	// The containerd managed opt directory provides a way for users to install containerd
	// dependencies using the existing distribution infrastructure.
	Opt            string
	Address        string
	ConfigPath     string
	Template       string
	Registry       string
	State          string
	Root           string
	DisableSELinux bool
	Images         string
	Log            string
}

func New(cfg Config) (*Containerd, error) {

	log := logger.WithContext(context.Background())

	c := new(Containerd)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.address = cfg.Address
	c.opt = cfg.Opt
	c.config = cfg.ConfigPath
	c.template = cfg.Template
	c.registry = cfg.Registry
	c.state = cfg.State
	c.root = cfg.Root
	c.disableSELinux = cfg.DisableSELinux
	c.images = cfg.Images
	c.log = cfg.Log

	var containerdTemplate string

	registryConfig, err := c.preparePrivateRegistry()
	if err != nil {
		return nil, err
	}

	containerdConfig := templates.ContainerdConfig{
		Opt:               cfg.Opt,
		IsRunningInUserNS: system.RunningInUserNS(),
		RegistryConfig:    registryConfig,
	}

	selEnabled, selConfigured, err := selinuxStatus()
	if err != nil {
		return nil, fmt.Errorf("%v: failed to detect selinux", err)
	}

	containerdConfig.SELinuxEnabled = selEnabled

	if cfg.DisableSELinux {
		containerdConfig.SELinuxEnabled = false
		if selEnabled {
			log.Warn("SELinux is enabled for system but has been disabled for containerd by override")
		}
	}

	if containerdConfig.SELinuxEnabled && !selConfigured {
		log.Warnf("SELinux is enabled for last.backend but process is not running in context '%s', last.backend-selinux policy may need to be applied", SELinuxContextType)
	}

	containerdTemplateBytes, err := ioutil.ReadFile(cfg.Template)
	if err == nil {
		containerdTemplate = string(containerdTemplateBytes)
	} else if os.IsNotExist(err) {
		containerdTemplate = templates.ContainerdConfigTemplate
	} else {
		return nil, err
	}
	parsedTemplate, err := templates.ParseTemplateFromConfig(containerdTemplate, containerdConfig)
	if err != nil {
		return nil, err
	}

	if err := filesystem.WriteFile(cfg.ConfigPath, parsedTemplate); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Containerd) Run() error {

	log := logger.WithContext(context.Background())

	go c.runServer()

	for {
		cli, err := containerd.New("unix://" + c.address)
		if err != nil {
			break
		}

		_, err = cli.Version(c.ctx)
		if err == nil {
			cli.Close()
			break
		}
		cli.Close()

		log.Infof("Waiting for containerd server startup: %v", err)

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return c.preloadImages()
}

func (c *Containerd) Stop() {
	c.cancel()
}

func (c Containerd) preloadImages() error {
	log := logger.WithContext(c.ctx)

	fileInfo, err := os.Stat(c.images)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		log.Errorf("Unable to find images in %s: %v", c.images, err)
		return nil
	}

	if !fileInfo.IsDir() {
		return nil
	}

	fileInfos, err := ioutil.ReadDir(c.images)
	if err != nil {
		log.Errorf("Unable to read images in %s: %v", c.images, err)
		return nil
	}

	client, err := containerd.New(c.address)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		filePath := filepath.Join(c.images, fileInfo.Name())

		file, err := os.Open(filePath)
		if err != nil {
			log.Errorf("Unable to read %s: %v", filePath, err)
			continue
		}

		log.Debugf("Import %s", filePath)
		_, err = client.Import(namespaces.WithNamespace(context.Background(), defaultNamespace), file)
		if err != nil {
			log.Errorf("Unable to import %s: %v", filePath, err)
		}
	}

	return nil
}

func (c *Containerd) preparePrivateRegistry() (*templates.Registry, error) {
	log := logger.WithContext(c.ctx)

	regTpl := new(templates.Registry)
	regFile, err := ioutil.ReadFile(c.registry)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	log.Infof("Using registry config file at %s", c.registry)
	if err := yaml.Unmarshal(regFile, &regTpl); err != nil {
		return nil, err
	}

	return regTpl, nil
}

func (c *Containerd) setLoggerProvider(ctx context.Context, fileeName string) io.WriteCloser {
	log := logger.WithContext(ctx)
	log.Infof("Logging containerd to %s", fileeName)

	return &lumberjack.Logger{
		Filename:   fileeName,
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

func (c *Containerd) runServer() {

	args := []string{
		"-c", c.config,
		"-a", c.address,
		"--state", c.state,
		"--root", c.root,
	}

	app:= command.App()
	if err := app.Run(args); err != nil {
		fmt.Println(fmt.Sprintf("containerd err: %v", err))
		os.Exit(1)
	}

	//
	//cfg := new(srvconfig.Config)
	//cfg.Root = c.root
	//cfg.State = c.state
	//
	//scd, err :=server.New(c.ctx, cfg)
	//if err != nil {
	//	os.Exit(1)
	//}

	os.Exit(1)
}
