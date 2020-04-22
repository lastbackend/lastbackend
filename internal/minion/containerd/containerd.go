package containerd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/lastbackend/lastbackend/internal/minion/containerd/templates"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"github.com/opencontainers/runc/libcontainer/system"
	"errors"
	"gopkg.in/yaml.v3"
)

const (
	defaultNamespace = "lstbknd.net"
)

type Containerd struct {
	Address  string
	Log      string
	Root     string
	State    string
	Config   string
	Template string
	Opt      string
}

type Config struct {
	Containerd      Containerd
	PrivateRegistry string
	DisableSELinux  bool
	Images          string
}

func Run(ctx context.Context, cfg *Config) error {

	log := logger.WithContext(ctx)

	if err := setupContainerdConfig(ctx, cfg); err != nil {
		return err
	}

	go runContainerdServer(ctx, cfg)

	for {
		cli, err := containerd.New("unix://" + cfg.Containerd.Address)
		if err != nil {
			break
		}

		_, err = cli.Version(ctx)
		if err == nil {
			cli.Close()
			break
		}
		cli.Close()

		log.Infof("Waiting for containerd server startup: %v", err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return preloadImages(ctx, cfg)
}

func preloadImages(ctx context.Context, cfg *Config) error {
	log := logger.WithContext(ctx)

	fileInfo, err := os.Stat(cfg.Images)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		log.Errorf("Unable to find images in %s: %v", cfg.Images, err)
		return nil
	}

	if !fileInfo.IsDir() {
		return nil
	}

	fileInfos, err := ioutil.ReadDir(cfg.Images)
	if err != nil {
		log.Errorf("Unable to read images in %s: %v", cfg.Images, err)
		return nil
	}

	client, err := containerd.New(cfg.Containerd.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		filePath := filepath.Join(cfg.Images, fileInfo.Name())

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

func setupContainerdConfig(ctx context.Context, cfg *Config) error {
	log := logger.WithContext(ctx)

	var containerdTemplate string

	registryConfig, err := preparePrivateRegistry(ctx, cfg)
	if err != nil {
		return err
	}

	containerdConfig := templates.ContainerdConfig{
		Opt:               cfg.Containerd.Opt,
		IsRunningInUserNS: system.RunningInUserNS(),
		RegistryConfig:    registryConfig,
	}

	selEnabled, selConfigured, err := selinuxStatus()
	if err != nil {
		return errors.New(fmt.Sprintf( "%v: failed to detect selinux", err))
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

	containerdTemplateBytes, err := ioutil.ReadFile(cfg.Containerd.Template)
	if err == nil {
		log.Infof("Using containerd template at %s", cfg.Containerd.Template)
		containerdTemplate = string(containerdTemplateBytes)
	} else if os.IsNotExist(err) {
		containerdTemplate = templates.ContainerdConfigTemplate
	} else {
		return err
	}
	parsedTemplate, err := templates.ParseTemplateFromConfig(containerdTemplate, containerdConfig)
	if err != nil {
		return err
	}

	return filesystem.WriteFile(cfg.Containerd.Config, parsedTemplate)
}

func preparePrivateRegistry(ctx context.Context, cfg *Config) (*templates.Registry, error) {
	log := logger.WithContext(ctx)

	regTpl := new(templates.Registry)
	regFile, err := ioutil.ReadFile(cfg.PrivateRegistry)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	log.Infof("Using registry config file at %s", cfg.PrivateRegistry)
	if err := yaml.Unmarshal(regFile, &regTpl); err != nil {
		return nil, err
	}

	return regTpl, nil
}

func setLoggerProvider(ctx context.Context, fileeName string) io.WriteCloser {
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

func runContainerdServer(ctx context.Context, cfg *Config) {
	log := logger.WithContext(ctx)

	args := []string{
		"containerd",
		"-c", cfg.Containerd.Config,
		"-a", cfg.Containerd.Address,
		"--state", cfg.Containerd.State,
		"--root", cfg.Containerd.Root,
	}

	if len(os.Getenv("CONTAINERD_LOG_LEVEL")) != 0 {
		args = append(args, "-l", os.Getenv("CONTAINERD_LOG_LEVEL"))
	}

	stdOut := io.Writer(os.Stdout)
	stdErr := io.Writer(os.Stderr)

	if cfg.Containerd.Log != "" {
		stdOut = setLoggerProvider(ctx, cfg.Containerd.Log)
		stdErr = stdOut
	}

	log.Infof("Running containerd server")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
	}

	os.Exit(1)
}
