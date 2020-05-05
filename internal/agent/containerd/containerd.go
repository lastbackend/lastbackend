package containerd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/services/server"
	srvconfig "github.com/containerd/containerd/services/server/config"
	"github.com/containerd/containerd/sys"
	"github.com/lastbackend/lastbackend/internal/agent/containerd/templates"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/prometheus/common/version"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const (
	defaultNamespace = "lstbknd.net"
	// DefaultRootDir is the default location used by containerd to store
	// persistent data
	DefaultRootDir = "/var/lib/containerd"
	// DefaultStateDir is the default location used by containerd to store
	// transient data
	DefaultStateDir = "/run/containerd"
	// DefaultAddress is the default unix socket address
	DefaultAddress = "/run/containerd/containerd.sock"
	// DefaultDebugAddress is the default unix socket address for pprof data
	DefaultDebugAddress = "/run/containerd/debug.sock"
	// DefaultFIFODir is the default location used by client-side cio library
	// to store FIFOs.
	DefaultFIFODir = "/run/containerd/fifo"
	// DefaultRuntime is the default linux runtime
	DefaultRuntime = "io.containerd.runc.v2"
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

func defaultConfig() *srvconfig.Config {
	return &srvconfig.Config{
		Version: 1,
		Root:    DefaultRootDir,
		State:   DefaultStateDir,
		GRPC: srvconfig.GRPCConfig{
			Address: DefaultAddress,
		},
		Debug: srvconfig.Debug{
			Level:   "info",
			Address: DefaultDebugAddress,
		},
		DisabledPlugins: []string{},
		RequiredPlugins: []string{},
	}
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

func (c *Containerd) runServer() error {
	log := logger.WithContext(c.ctx)

	var (
		start   = time.Now()
		serverC = make(chan *server.Server, 1)
		done    = handleSignals(c.ctx, serverC)
	)

	config := defaultConfig()
	
	if err := srvconfig.LoadConfig(c.config, config); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	if c.root != "" {
		config.Root = c.root
	}
	
	if c.state != "" {
		config.State = c.state
	}
	
	if c.address != "" {
		config.GRPC.Address = c.address
	}
	
	// Make sure top-level directories are created early.
	if err := server.CreateTopLevelDirectories(config); err != nil {
		return err
	}
	
	if config.GRPC.Address == "" {
		return errors.New("grpc address cannot be empty")
	}
	if config.TTRPC.Address == "" {
		// If TTRPC was not explicitly configured, use defaults based on GRPC.
		config.TTRPC.Address = fmt.Sprintf("%s.ttrpc", config.GRPC.Address)
		config.TTRPC.UID = config.GRPC.UID
		config.TTRPC.GID = config.GRPC.GID
	}
	
	log.WithFields(logger.Fields{
		"version":  version.Version,
		"revision": version.Revision,
	}).Info("starting containerd")
	
	server, err := server.New(c.ctx, config)
	if err != nil {
		return err
	}
	
	serverC <- server
	
	if config.Debug.Address != "" {
		var l net.Listener
		if filepath.IsAbs(config.Debug.Address) {
			if l, err = sys.GetLocalListener(config.Debug.Address, config.Debug.UID, config.Debug.GID); err != nil {
				return fmt.Errorf("err %v: failed to get listener for debug endpoint", err)
			}
		} else {
			if l, err = net.Listen("tcp", config.Debug.Address); err != nil {
				return fmt.Errorf("err %v: failed to get listener for debug endpoint", err)
			}
		}
		serve(c.ctx, l, server.ServeDebug)
	}

	if config.Metrics.Address != "" {
		l, err := net.Listen("tcp", config.Metrics.Address)
		if err != nil {
			return fmt.Errorf("err %v: failed to get listener for metrics endpoint", err)
		}
		serve(c.ctx, l, server.ServeMetrics)
	}

	// setup the ttrpc endpoint
	tl, err := sys.GetLocalListener(config.TTRPC.Address, config.TTRPC.UID, config.TTRPC.GID)
	if err != nil {
		return fmt.Errorf("err %v: failed to get listener for main ttrpc endpoint", err)
	}
	serve(c.ctx, tl, server.ServeTTRPC)

	if config.GRPC.TCPAddress != "" {
		l, err := net.Listen("tcp", config.GRPC.TCPAddress)
		if err != nil {
			return fmt.Errorf("err %v: failed to get listener for TCP grpc endpoint", err)
		}
		serve(c.ctx, l, server.ServeTCP)
	}

	// setup the main grpc endpoint
	l, err := sys.GetLocalListener(config.GRPC.Address, config.GRPC.UID, config.GRPC.GID)
	if err != nil {
		return fmt.Errorf("err %v: failed to get listener for main endpoint", err)
	}
	
	serve(c.ctx, l, server.ServeGRPC)

	log.Infof("containerd successfully booted in %fs", time.Since(start).Seconds())


	<-done

	return nil
}

func serve(ctx context.Context, l net.Listener, serveFunc func(net.Listener) error) {
	log := logger.WithContext(ctx)

	path := l.Addr().String()
	log.WithFields(logger.Fields{"address": path}).Info("serving...")

	go func() {
		defer l.Close()
		if err := serveFunc(l); err != nil {
			log.WithFields(logger.Fields{"address": path}).Fatalf("err %v:serve failure")
		}
	}()
}

func handleSignals(ctx context.Context, serverC chan *server.Server) chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		var server *server.Server
		for {
			select {
			case s := <-serverC:
				server = s
			case <-ctx.Done():
				if server == nil {
					close(done)
					return
				}
				server.Stop()
				close(done)
				return
			}
		}
	}()
	return done
}
