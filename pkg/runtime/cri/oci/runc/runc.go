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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/docker/go-units"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/intelrdt"
	"github.com/opencontainers/runc/libcontainer/specconv"
	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/runc/libcontainer/user"
	"github.com/opencontainers/runc/libcontainer/utils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

const DefaultRootPath = "/run/runc"

type Runc interface {
	CreateContainer(containerID string, opts CreateOptions) (libcontainer.Container, error)
	StartContainer(containerID string) error
	RunContainer(containerID string, bundle string, opts RunOptions) error
	PauseContainer(containerID string) error
	ResumeContainer(containerID string) error
	ListContainer() ([]ContainerState, error)
	StateContainer(containerID string) (*ContainerState, error)
	DeleteContainer(containerID string, opts DeleteOptions) error
	RestoreContainer(ctx context.Context, containerID string, opts RestoreOptions) error
	KillContainer(containerID, sigstr string, opts KillOptions) error
	SpecContainer(opts SpecOptions) error
	ExecContainer(containerID string, args []string, opts ExecOptions) error
	UpdateContainer(containerID string, opts UpdateOptions) error
	CheckpointContainer(containerID string, opts CheckpointOptions) error
	EventsContainer(containerID string, opts EventsOptions) (*EventWatcher, error)
}

// ContainerState represents the platform agnostic pieces relating to a running container's status and state
type ContainerState struct {
	// Version is the OCI version for the container
	Version string `json:"ociVersion"`
	// ID is the container ID
	ID string `json:"id"`
	// InitProcessPid is the init process id in the parent namespace
	InitProcessPid int `json:"pid"`
	// Status is the current status of the container, running, paused, ...
	Status string `json:"status"`
	// Bundle is the path on the filesystem to the Bundle
	Bundle string `json:"Bundle"`
	// Rootfs is a path to a directory containing the container's root filesystem.
	Rootfs string `json:"rootfs"`
	// Created is the unix timestamp for the creation time of the container in UTC
	Created time.Time `json:"created"`
	// Annotations is the user defined annotations added to the config.
	Annotations map[string]string `json:"annotations,omitempty"`
	// The owner of the state directory (the owner of the container).
	Owner string `json:"owner"`
}

type Config struct {
	// Enable debug output for logging
	Debug bool
	// Ignore cgroup permission errors ('true', 'false', or 'auto')
	Rootless string
	// Enable systemd cgroup support, expects cgroupsPath to be of form "slice:prefix:name" for e.g. "system.slice:runc:123321"
	SystemdCgroup bool
	// Path to the criu binary used for checkpoint and restore
	CriuPath string
	// Root directory for storage of container state (this should be located in tmpfs)
	RootPath string
}

type runc struct {
	debug            bool
	useSystemdCgroup bool
	rootless         string
	criuPath         string
	rootPath         string
	notifySocket     string
	listenFDS        string
}

// Open Container Initiative runtime for running applications packaged according to
// the Open Container Initiative (OCI) format and is a compliant
// implementation of the Open Container Initiative specification.
func New(cfg Config) (Runc, error) {
	root := DefaultRootPath

	rootPath, err := filepath.Abs(cfg.RootPath)
	if err != nil {
		return nil, err
	}

	xdgRuntimeDir := ""
	if shouldHonorXDGRuntimeDir() {
		if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
			root = runtimeDir + "/runc"
			xdgRuntimeDir = root
		}
	}

	if rootPath == "" && xdgRuntimeDir != "" {
		// According to the XDG specification, we need to set anything in
		// XDG_RUNTIME_DIR to have a sticky bit if we don't want it to get
		// auto-pruned.
		if err := os.MkdirAll(root, 0700); err != nil {
			fmt.Fprintln(os.Stderr, "the path in $XDG_RUNTIME_DIR must be writable by the user")
			return nil, err
		}
		if err := os.Chmod(root, 0700|os.ModeSticky); err != nil {
			fmt.Fprintln(os.Stderr, "you should check permission of the path in $XDG_RUNTIME_DIR")
			return nil, err
		}
	}

	r := new(runc)
	r.debug = cfg.Debug
	r.useSystemdCgroup = cfg.SystemdCgroup
	r.rootless = cfg.Rootless
	r.criuPath = cfg.CriuPath
	r.rootPath = rootPath
	r.notifySocket = os.Getenv("NOTIFY_SOCKET")
	r.listenFDS = os.Getenv("LISTEN_FDS")

	return r, nil
}

type CreateOptions struct {
	// Path to the root of the bundle directory, defaults to the current directory
	Bundle string
	// Do not create a new session keyring for the container.
	// This will cause the container to inherit the calling processes session key
	NoNewKeyring bool
	// Do not use pivot root to jail process inside rootfs.
	// This should be used whenever the rootfs is on top of a ramdisk
	NoPivot bool
	// Path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
	ConsoleSocket string
	// detach from the container's process
	Detach bool
	// Specify the file to write the process id to
	PidFile string
	// Disable the use of the subreaper used to reap reparented processes
	NoSubreaper bool
	// Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
	PreserveFds int
}

// CreateContainer - creates an instance of a container for a bundle.
// Where 'containerID' is your name for the instance of the container.
// The name you provide for the container instance must be unique on your host.
// CreateOptions:
//	 Bundle <string>: path to the root of the Bundle directory, defaults to the current directory
//   NoNewKeyring <bool>: do not create a new session keyring for the container.
//   NoPivot <bool>: do not use pivot root to jail process inside rootfs. This should be used whenever the rootfs is on top of a ramdisk
//   ConsoleSocket <string>: path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
//   Detach <bool>: detach from the container's process
//   PidFile <string>: specify the file to write the process id to
//   NoSubreaper <bool>: disable the use of the subreaper used to reap reparented processes
//   PreserveFds <int>: pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
func (r *runc) CreateContainer(containerID string, opts CreateOptions) (libcontainer.Container, error) {
	spec, err := setupSpec(opts.Bundle)
	if err != nil {
		return nil, err
	}
	container, err := r.createContainer(containerID, spec, opts)
	if err != nil {
		return nil, err
	}

	notifySocket := newNotifySocket(containerID, r.rootPath, r.notifySocket)
	if notifySocket != nil {
		if err := notifySocket.setupSpec(spec); err != nil {
			return nil, err
		}
	}

	if notifySocket != nil {
		err := notifySocket.setupSocketDirectory()
		if err != nil {
			return nil, err
		}
	}

	// Support on-demand socket activation by passing file descriptors into the container init process.
	listenFDs := make([]*os.File, 0)

	if r.listenFDS != "" {
		listenFDs = activation.Files(false)
	}

	logLevel := "info"
	if r.debug {
		logLevel = "debug"
	}

	rnr := &runner{
		init:            true,
		shouldDestroy:   true,
		container:       container,
		listenFDs:       listenFDs,
		notifySocket:    notifySocket,
		action:          CT_ACT_CREATE,
		criuOpts:        nil,
		logLevel:        logLevel,
		detach:          opts.Detach,
		pidFile:         opts.PidFile,
		preserveFDs:     opts.PreserveFds,
		enableSubreaper: !opts.NoSubreaper,
		consoleSocket:   opts.ConsoleSocket,
	}
	_, err = rnr.run(spec.Process)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// StartContainer - executes the user defined process in a created container.
// Where 'containerID' is your name for the instance of the container.
func (r *runc) StartContainer(containerID string) error {
	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}
	status, err := container.Status()
	if err != nil {
		return err
	}

	switch status {
	case libcontainer.Created:
		notifySocket, err := notifySocketStart(container.ID(), r.rootPath, r.notifySocket)
		if err != nil {
			return err
		}
		if err := container.Exec(); err != nil {
			return err
		}
		if notifySocket != nil {
			return notifySocket.waitForContainer(container)
		}
		return nil
	case libcontainer.Stopped:
		return errors.New("cannot start a container that has stopped")
	case libcontainer.Running:
		return errors.New("cannot start an already running container")
	default:
		return fmt.Errorf("cannot start a container in the %s state\n", status)
	}
}

type RunOptions struct {
	// Path to the root of the bundle directory, defaults to the current directory
	Bundle string
	// Path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
	ConsoleSocket string
	// detach from the container's process
	Detach bool
	// Specify the file to write the process id to
	PidFile string
	// Disable the use of the subreaper used to reap reparented processes
	NoSubreaper bool
	// Do not use pivot root to jail process inside rootfs.
	// This should be used whenever the rootfs is on top of a ramdisk
	NoPivot bool
	// Do not create a new session keyring for the container.
	// This will cause the container to inherit the calling processes session key
	NoNewKeyring bool
	// Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
	PreserveFds int
}

// RunContainer - create and run a container
// Where 'containerID' is your name for the instance of the container.
// The name you provide for the container instance must be unique on your host.
// RunOptions:
//	 Bundle <string>: path to the root of the bundle directory, defaults to the current directory
//   ConsoleSocket <string>: path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
//   Detach <bool>: detach from the container's process
//   PidFile <string>: specify the file to write the process id to
//   NoSubreaper <bool>: disable the use of the subreaper used to reap reparented processes
//   NoPivot <bool>: do not use pivot root to jail process inside rootfs. This should be used whenever the rootfs is on top of a ramdisk
//   NoNewKeyring <bool>: do not create a new session keyring for the container. This will cause the container to inherit the calling processes session key
//   PreserveFds <int>: pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
func (r *runc) RunContainer(containerID string, bundle string, opts RunOptions) error {
	spec, err := setupSpec(bundle)
	if err != nil {
		return err
	}
	return r.runContainer(containerID, spec, CT_ACT_RUN, nil, opts)
}

// PauseContainer - suspends all processes in the instance of the container.
// Where 'containerID' is your name for the instance of the container.
func (r *runc) PauseContainer(containerID string) error {
	rcg, err := shouldUseRootlessCgroupManager(r.rootless, r.useSystemdCgroup)
	if err != nil {
		return err
	}
	if rcg {
		fmt.Println("Warn: runc pause may fail if you don't have the full access to cgroups")
	}
	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}
	return container.Pause()
}

// ResumeContainer - resumes all processes in the instance of the container.
// Where 'containerID' is your name for the instance of the container.
func (r *runc) ResumeContainer(containerID string) error {
	rcg, err := shouldUseRootlessCgroupManager(r.rootless, r.useSystemdCgroup)
	if err != nil {
		return err
	}
	if rcg {
		fmt.Println("Warn: runc resume may fail if you don't have the full access to cgroups")
	}
	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}
	return container.Resume()
}

// ListContainer - lists containers started by run with the given root (default: '/run/runc').
func (r *runc) ListContainer() ([]ContainerState, error) {
	list, err := ioutil.ReadDir(r.rootPath)
	if err != nil {
		return nil, err
	}

	var s = make([]ContainerState, 0)

	for _, item := range list {
		if item.IsDir() {
			// This cast is safe on Linux.
			stat := item.Sys().(*syscall.Stat_t)
			owner, err := user.LookupUid(int(stat.Uid))
			if err != nil {
				owner.Name = fmt.Sprintf("#%d", stat.Uid)
			}

			factory, err := loadFactory(r.rootPath, r.criuPath, r.useSystemdCgroup)
			if err != nil {
				return nil, err
			}

			container, err := factory.Load(item.Name())
			if err != nil {
				fmt.Println(fmt.Fprintf(os.Stderr, "load container %s: %v\n", item.Name(), err))
				continue
			}
			containerStatus, err := container.Status()
			if err != nil {
				fmt.Println(fmt.Fprintf(os.Stderr, "status for %s: %v\n", item.Name(), err))
				continue
			}
			state, err := container.State()
			if err != nil {
				fmt.Println(fmt.Fprintf(os.Stderr, "state for %s: %v\n", item.Name(), err))
				continue
			}
			pid := state.BaseState.InitProcessPid
			if containerStatus == libcontainer.Stopped {
				pid = 0
			}
			bundle, annotations := utils.Annotations(state.Config.Labels)
			s = append(s, ContainerState{
				Version:        state.BaseState.Config.Version,
				ID:             state.BaseState.ID,
				InitProcessPid: pid,
				Status:         containerStatus.String(),
				Bundle:         bundle,
				Rootfs:         state.BaseState.Config.Rootfs,
				Created:        state.BaseState.Created,
				Annotations:    annotations,
				Owner:          owner.Name,
			})
		}
	}

	return s, nil
}

// StateContainer - get current state information for the instance of a container.
// Where 'containerID' is your name for the instance of the container.
func (r *runc) StateContainer(containerID string) (*ContainerState, error) {
	container, err := r.getContainer(containerID)
	if err != nil {
		return nil, err
	}
	status, err := container.Status()
	if err != nil {
		return nil, err
	}
	state, err := container.State()
	if err != nil {
		return nil, err
	}
	pid := state.BaseState.InitProcessPid
	if status == libcontainer.Stopped {
		pid = 0
	}
	bundle, annotations := utils.Annotations(state.Config.Labels)
	return &ContainerState{
		Version:        state.BaseState.Config.Version,
		ID:             state.BaseState.ID,
		InitProcessPid: pid,
		Status:         status.String(),
		Bundle:         bundle,
		Rootfs:         state.BaseState.Config.Rootfs,
		Created:        state.BaseState.Created,
		Annotations:    annotations,
	}, nil
}

type DeleteOptions struct {
	// Forcibly deletes the container if it is still running (uses SIGKILL)
	Force bool
}

// DeleteContainer - delete any resources held by the container often used with detached container.
// Where 'containerID' is your name for the instance of the container.
// DeleteOptions:
//	 Force <bool>: forcibly deletes the container if it is still running (uses SIGKILL)
func (r *runc) DeleteContainer(containerID string, opts DeleteOptions) error {
	container, err := r.getContainer(containerID)
	if err != nil {
		if lerr, ok := err.(libcontainer.Error); ok && lerr.Code() == libcontainer.ContainerNotExists {
			// if there was an aborted start or something of the sort then the container's directory could exist but
			// libcontainer does not see it because the state.json file inside that directory was never created.
			path := filepath.Join(r.rootPath, containerID)
			if e := os.RemoveAll(path); e != nil {
				fmt.Println(fmt.Sprintf("remove %s: %v", path, e))
			}
			if opts.Force {
				return nil
			}
		}
		return err
	}

	status, err := container.Status()
	if err != nil {
		return err
	}
	switch status {
	case libcontainer.Stopped:
		destroy(container)
	case libcontainer.Created:
		return r.killContainer(container)
	default:
		if opts.Force {
			return r.killContainer(container)
		}
		return fmt.Errorf("cannot delete container %s that is not stopped: %s", containerID, status)
	}

	return nil
}

type RestoreOptions struct {
	// Path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
	ConsoleSocket string
	// Path to criu image files for restoring
	ImagePath string
	// Path for saving work files and logs
	WorkPath string
	// Allow open tcp connections
	TcpEstablished bool
	// Allow external unix sockets
	ExtUnixSk bool
	// Allow shell jobs
	ShellJob bool
	// Jandle file locks, for safety
	FileLocks bool
	// Cgroups mode: 'soft' (default), 'full' and 'strict'
	ManageCgroupsMode string
	// Path to the root of the bundle directory
	Bundle string
	// Detach from the container's process
	Detach bool
	// Specify the file to write the process id to
	PidFile string
	// Disable the use of the subreaper used to reap reparented processes
	NoSubreaper bool
	// Do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk
	NoPivot bool
	// Create a namespace, but don't restore its properties
	EmptyNs []string
	// Enable auto deduplication of memory images
	AutoDedup bool
	// Use userfaultfd to lazily restore memory pages
	LazyPages bool
}

// RestoreContainer - restores the saved state of the container instance that was previously saved
// using the runc checkpoint command.
// Where 'containerID' is your name for the instance of the container.
// RestoreOptions:
//   ConsoleSocket <string>: path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
//   ImagePath <string>: path to criu image files for restoring
//   WorkPath <string>: path for saving work files and logs
//   TcpEstablished <bool>: allow open tcp connections
//   ExtUnixSk <bool> allow external unix sockets
//   ShellJob <bool>: allow shell jobs
//   FileLocks <bool>: handle file locks, for safety
//   ManageCgroupsMode <string>: cgroups mode: 'soft' (default), 'full' and 'strict'
//   Bundle <string>: path to the root of the bundle directory
//   Detach <bool>: detach from the container's process
//   PidFile <string>: specify the file to write the process id to
//   NoSubreaper <bool>: disable the use of the subreaper used to reap reparented processes
//   NoPivot <bool>: do not use pivot root to jail process inside rootfs. This should be used whenever the rootfs is on top of a ramdisk
//   EmptyNs <[]string>: create a namespace, but don't restore its properties
//   AutoDedup <bool>: enable auto deduplication of memory images
//   LazyPages <bool>: use userfaultfd to lazily restore memory pages
func (r *runc) RestoreContainer(ctx context.Context, containerID string, opts RestoreOptions) error {
	// Currently this is untested with rootless containers.
	if os.Geteuid() != 0 || system.RunningInUserNS() {
		fmt.Println("Warn: runc checkpoint is untested with rootless containers")
	}

	spec, err := setupSpec(opts.Bundle)
	if err != nil {
		return err
	}

	imagePath := opts.ImagePath
	if imagePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		imagePath = filepath.Join(cwd, "checkpoint")
	}

	if err := os.MkdirAll(imagePath, 0755); err != nil {
		return err
	}

	criOPts := &libcontainer.CriuOpts{
		ImagesDirectory:         imagePath,
		WorkDirectory:           opts.WorkPath,
		ParentImage:             opts.ImagePath,
		TcpEstablished:          opts.TcpEstablished,
		ExternalUnixConnections: opts.ExtUnixSk,
		ShellJob:                opts.ShellJob,
		AutoDedup:               opts.AutoDedup,
		LazyPages:               opts.LazyPages,
	}

	if err := setEmptyNsMask(opts.EmptyNs, criOPts); err != nil {
		return err
	}

	return r.runContainer(containerID, spec, CT_ACT_RESTORE, criOPts, RunOptions{
		Bundle:        opts.Bundle,
		ConsoleSocket: opts.ConsoleSocket,
		Detach:        opts.Detach,
		PidFile:       opts.PidFile,
		NoSubreaper:   opts.NoSubreaper,
		NoPivot:       opts.NoPivot,
	})
}

type KillOptions struct {
	// Send the specified signal to all processes inside the container
	All bool
}

// KillContainer - kill sends the specified signal (default: SIGTERM) to the container's init process
// Where 'containerID' is your name for the instance of the container.
// RestoreOptions:
//   All <bool>: send the specified signal to all processes inside the container
func (r *runc) KillContainer(containerID, sigstr string, opts KillOptions) error {
	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}
	if sigstr == "" {
		sigstr = "SIGTERM"
	}
	signal, err := parseSignal(sigstr)
	if err != nil {
		return err
	}
	return container.Signal(signal, opts.All)
}

type SpecOptions struct {
	// Path to the root of the bundle directory
	Bundle string
}

// SpecContainer - creates the new specification file for the bundle.
// The spec generated is just a starter file. Editing of the spec is required to achieve desired results.
func (r *runc) SpecContainer(opts SpecOptions) error {

	rcg, err := shouldUseRootlessCgroupManager(r.rootless, r.useSystemdCgroup)
	if err != nil {
		return err
	}
	if rcg {
		fmt.Println("Warn: runc pause may fail if you don't have the full access to cgroups")
	}

	spec := specconv.Example()
	if rcg {
		specconv.ToRootless(spec)
	}

	checkNoFile := func(name string) error {
		_, err := os.Stat(name)
		if err == nil {
			return fmt.Errorf("File %s exists. Remove it first", name)
		}
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	if opts.Bundle != "" {
		if err := os.Chdir(opts.Bundle); err != nil {
			return err
		}
	}

	if err := checkNoFile(specConfig); err != nil {
		return err
	}

	data, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(specConfig, data, 0666)
}

type ExecOptions struct {
	// Path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
	ConsoleSocket string
	// Current working directory in the container
	Cwd string
	// Set environment variables
	Env []string
	// Allocate a pseudo-TTY
	Tty bool
	// UID (format: <uid>[:<gid>])
	User string
	// Additional gids
	AdditionalGids []int64
	// Path to the process.json
	Process string
	// Detach from the container's process
	Detach bool
	// Specify the file to write the process id to
	PidFile string
	// Set the asm process label for the process commonly used with selinux
	ProcessLabel string
	// Set the apparmor profile for the process
	Apparmor string
	// Set the no new privileges value for the process
	NoNewPrivs bool
	// Add a capability to the bounding set for the process
	Cap []string
	// Disable the use of the subreaper used to reap reparented processes
	NoSubreaper bool
	// Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
	PreserveFds int
}

// ExecContainer - execute new process inside the container
// Where 'containerID' is your name for the instance of the container and 'args' command arguments
// CreateOptions:
//    ConsoleSocket  <string>: path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal
//    Cwd            <string>: current working directory in the container
//    Env            <[]string>: set environment variables
//    Tty            <bool>: allocate a pseudo-TTY
//    User           <string>: UID (format: <uid>[:<gid>])
//    AdditionalGids <[]int64>: additional gids
//    Process        <string>: path to the process.json
//    Detach         <bool>: detach from the container's process
//    PidFile        <string>: specify the file to write the process id to
//    ProcessLabel   <string>: set the asm process label for the process commonly used with selinux
//    Apparmor       <string>: set the apparmor profile for the process
//    NoNewPrivs     <bool>: set the no new privileges value for the process
//    Cap            <[]string>: add a capability to the bounding set for the process
//    NoSubreaper    <bool>: disable the use of the subreaper used to reap reparented processes
//    PreserveFds    <int>: pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)
func (r *runc) ExecContainer(containerID string, args []string, opts ExecOptions) error {
	pidFile, err := revisePidFile(opts.PidFile)
	if err != nil {
		return err
	}

	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}

	status, err := container.Status()
	if err != nil {
		return err
	}

	if status == libcontainer.Stopped {
		return fmt.Errorf("cannot exec a container that has stopped")
	}

	state, err := container.State()
	if err != nil {
		return err
	}

	bundle := utils.SearchLabels(state.Config.Labels, "bundle")
	p, err := r.getProcess(bundle, args, opts)
	if err != nil {
		return err
	}

	logLevel := "info"
	if r.debug {
		logLevel = "debug"
	}

	rnr := &runner{
		shouldDestroy:   false,
		init:            false,
		container:       container,
		pidFile:         pidFile,
		action:          CT_ACT_RUN,
		logLevel:        logLevel,
		enableSubreaper: !opts.NoSubreaper,
		detach:          opts.Detach,
		preserveFDs:     opts.PreserveFds,
		consoleSocket:   opts.ConsoleSocket,
	}
	_, err = rnr.run(p)
	if err != nil {
		return err
	}

	return nil
}

type UpdateOptions struct {
	// Path to the file containing the resources to update or '-' to read from the standard input
	// The accepted format is as follow (unchanged values can be omitted):
	// {
	//  "memory": {
	// 	   "limit": 0,
	// 	   "reservation": 0,
	// 	   "swap": 0,
	// 	   "kernel": 0,
	// 	   "kernelTCP": 0
	// 	 },
	//   "cpu": {
	// 	   "shares": 0,
	// 	   "quota": 0,
	// 	   "period": 0,
	// 	   "realtimeRuntime": 0,
	// 	   "realtimePeriod": 0,
	// 	   "cpus": "",
	// 	   "mems": ""
	// 	 },
	// 	 "blockIO": {
	// 	   "weight": 0
	// 	 }
	// }
	// Note: if data is to be read from a file or the standard input, all other options are ignored.
	Resources string
	// Specifies per cgroup weight, range is from 10 to 1000
	BlkioWeight int
	// CPU CFS period to be used for hardcapping (in usecs). 0 to use system default
	CpuPeriod string
	// CPU CFS hardcap limit (in usecs). Allowed cpu time in a given period
	CpuQuota string
	// CPU shares (relative weight vs. other containers)
	CpuShare string
	// CPU realtime period to be used for hardcapping (in usecs). 0 to use system default
	CpuRtPeriod string
	// CPU realtime hardcap limit (in usecs). Allowed cpu time in a given period
	CpuRtRuntime string
	// CPU(s) to use
	CpusetCpus string
	// Memory node(s) to use
	CpusetMems string
	// Kernel memory limit (in bytes)
	KernelMemory string
	// Kernel memory limit (in bytes) for tcp buffer
	KernelMemoryTcp string
	// Memory limit (in bytes)
	Memory string
	// Memory reservation or soft_limit (in bytes)
	MemoryReservation string
	// Total memory usage (memory + swap); set '-1' to enable unlimited swap
	MemorySwap string
	// Maximum number of pids allowed in the container
	PidsLimit int
	// The string of Intel RDT/CAT L3 cache schema
	L3CacheSchema string
	// The string of Intel RDT/MBA memory bandwidth schema
	MemBwSchema string
}

// UpdateContainer - update container resource constraints
// Where 'containerID' is your name for the instance of the container
// UpdateOptions:
// Specifies per cgroup weight, range is from 10 to 1000
//    Resources <string>: path to the file containing the resources to update or '-' to read from the standard input
// 				The accepted format is as follow (unchanged values can be omitted):
// 				{
// 				 "memory": {
// 					   "limit": 0,
// 					   "reservation": 0,
// 					   "swap": 0,
// 					   "kernel": 0,
// 					   "kernelTCP": 0
// 					 },
// 				  "cpu": {
// 					   "shares": 0,
// 					   "quota": 0,
// 					   "period": 0,
// 					   "realtimeRuntime": 0,
// 					   "realtimePeriod": 0,
// 					   "cpus": "",
// 					   "mems": ""
// 					 },
// 					 "blockIO": {
// 					   "weight": 0
// 					 }
// 				}
// 				Note: if data is to be read from a file or the standard input, all other options are ignored.
//    BlkioWeight <int>: specifies per cgroup weight, range is from 10 to 1000
//    CpuPeriod <string>: CPU CFS period to be used for hardcapping (in usecs). 0 to use system default
//    CpuQuota <string>: CPU CFS hardcap limit (in usecs). Allowed cpu time in a given period
//    CpuShare <string>: CPU realtime period to be used for hardcapping (in usecs). 0 to use system default
//    CpuRtPeriod <string>: CPU realtime period to be used for hardcapping (in usecs). 0 to use system default
//    CpuRtRuntime <string>: CPU realtime hardcap limit (in usecs). Allowed cpu time in a given period
//    CpusetCpus <string>: CPU(s) to use
//    CpusetMems <string>: memory node(s) to use
//    KernelMemory <string>: kernel memory limit (in bytes)
//    KernelMemoryTcp <string>: kernel memory limit (in bytes) for tcp buffer
//    Memory <string>: memory limit (in bytes)
//    MemoryReservation <string>: memory reservation or soft_limit (in bytes)
//    MemorySwap <string>: total memory usage (memory + swap); set '-1' to enable unlimited swap
//    PidsLimit <int>: maximum number of pids allowed in the container
//    L3CacheSchema <string>: the string of Intel RDT/CAT L3 cache schema
//    MemBwSchema <string>: the string of Intel RDT/MBA memory bandwidth schema
func (r *runc) UpdateContainer(containerID string, opts UpdateOptions) error {
	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}

	resource := specs.LinuxResources{
		Memory: &specs.LinuxMemory{
			Limit:       i64Ptr(0),
			Reservation: i64Ptr(0),
			Swap:        i64Ptr(0),
			Kernel:      i64Ptr(0),
			KernelTCP:   i64Ptr(0),
		},
		CPU: &specs.LinuxCPU{
			Shares:          u64Ptr(0),
			Quota:           i64Ptr(0),
			Period:          u64Ptr(0),
			RealtimeRuntime: i64Ptr(0),
			RealtimePeriod:  u64Ptr(0),
			Cpus:            "",
			Mems:            "",
		},
		BlockIO: &specs.LinuxBlockIO{
			Weight: u16Ptr(0),
		},
		Pids: &specs.LinuxPids{
			Limit: 0,
		},
	}

	config := container.Config()

	if in := opts.Resources; in != "" {
		var (
			f   *os.File
			err error
		)
		switch in {
		case "-":
			f = os.Stdin
		default:
			f, err = os.Open(in)
			if err != nil {
				return err
			}
		}
		err = json.NewDecoder(f).Decode(&r)
		if err != nil {
			return err
		}
	} else {
		if opts.BlkioWeight != 0 {
			resource.BlockIO.Weight = u16Ptr(uint16(opts.BlkioWeight))
		}
		if opts.CpusetCpus != "" {
			resource.CPU.Cpus = opts.CpusetCpus
		}
		if opts.CpusetMems != "" {
			resource.CPU.Mems = opts.CpusetMems
		}

		for _, pair := range []struct {
			opt  string
			val  string
			dest *uint64
		}{
			{"cpu-period", opts.CpuPeriod, resource.CPU.Period},
			{"cpu-rt-period", opts.CpuRtPeriod, resource.CPU.RealtimePeriod},
			{"cpu-share", opts.CpuShare, resource.CPU.Shares},
		} {
			if pair.val != "" {
				var err error
				*pair.dest, err = strconv.ParseUint(pair.val, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid value for %s: %s", pair.opt, err)
				}
			}
		}

		for _, pair := range []struct {
			opt  string
			val  string
			dest *int64
		}{
			{"cpu-quota", opts.CpuQuota, resource.CPU.Quota},
			{"cpu-rt-runtime", opts.CpuRtRuntime, resource.CPU.RealtimeRuntime},
		} {
			if pair.val != "" {
				var err error
				*pair.dest, err = strconv.ParseInt(pair.val, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid value for %s: %s", pair.opt, err)
				}
			}
		}

		for _, pair := range []struct {
			opt  string
			val  string
			dest *int64
		}{
			{"memory", opts.Memory, resource.Memory.Limit},
			{"memory-swap", opts.MemorySwap, resource.Memory.Swap},
			{"kernel-memory", opts.KernelMemory, resource.Memory.Kernel},
			{"kernel-memory-tcp", opts.KernelMemoryTcp, resource.Memory.KernelTCP},
			{"memory-reservation", opts.MemoryReservation, resource.Memory.Reservation},
		} {
			if pair.val != "" {
				var v int64
				if pair.val != "-1" {
					v, err = units.RAMInBytes(pair.val)
					if err != nil {
						return fmt.Errorf("invalid value for %s: %s", pair.opt, err)
					}
				} else {
					v = -1
				}
				*pair.dest = v
			}
		}

		resource.Pids.Limit = int64(opts.PidsLimit)
	}

	// Update the value
	config.Cgroups.Resources.BlkioWeight = *resource.BlockIO.Weight
	config.Cgroups.Resources.CpuPeriod = *resource.CPU.Period
	config.Cgroups.Resources.CpuQuota = *resource.CPU.Quota
	config.Cgroups.Resources.CpuShares = *resource.CPU.Shares
	//CpuWeight is used for cgroupv2 and should be converted
	config.Cgroups.Resources.CpuWeight = cgroups.ConvertCPUSharesToCgroupV2Value(*resource.CPU.Shares)
	//CpuMax is used for cgroupv2 and should be converted
	config.Cgroups.Resources.CpuMax = cgroups.ConvertCPUQuotaCPUPeriodToCgroupV2Value(*resource.CPU.Quota, *resource.CPU.Period)
	config.Cgroups.Resources.CpuRtPeriod = *resource.CPU.RealtimePeriod
	config.Cgroups.Resources.CpuRtRuntime = *resource.CPU.RealtimeRuntime
	config.Cgroups.Resources.CpusetCpus = resource.CPU.Cpus
	config.Cgroups.Resources.CpusetMems = resource.CPU.Mems
	config.Cgroups.Resources.KernelMemory = *resource.Memory.Kernel
	config.Cgroups.Resources.KernelMemoryTCP = *resource.Memory.KernelTCP
	config.Cgroups.Resources.Memory = *resource.Memory.Limit
	config.Cgroups.Resources.MemoryReservation = *resource.Memory.Reservation
	config.Cgroups.Resources.MemorySwap = *resource.Memory.Swap
	config.Cgroups.Resources.PidsLimit = resource.Pids.Limit

	// Update Intel RDT
	l3CacheSchema := opts.L3CacheSchema
	memBwSchema := opts.MemBwSchema
	if l3CacheSchema != "" && !intelrdt.IsCatEnabled() {
		return fmt.Errorf("Intel RDT/CAT: l3 cache schema is not enabled")
	}

	if memBwSchema != "" && !intelrdt.IsMbaEnabled() {
		return fmt.Errorf("Intel RDT/MBA: memory bandwidth schema is not enabled")
	}

	if l3CacheSchema != "" || memBwSchema != "" {
		// If intelRdt is not specified in original configuration, we just don't
		// Apply() to create intelRdt group or attach tasks for this container.
		// In update command, we could re-enable through IntelRdtManager.Apply()
		// and then update intelrdt constraint.
		if config.IntelRdt == nil {
			state, err := container.State()
			if err != nil {
				return err
			}
			config.IntelRdt = &configs.IntelRdt{}
			intelRdtManager := intelrdt.IntelRdtManager{
				Config: &config,
				Id:     container.ID(),
				Path:   state.IntelRdtPath,
			}
			if err := intelRdtManager.Apply(state.InitProcessPid); err != nil {
				return err
			}
		}
		config.IntelRdt.L3CacheSchema = l3CacheSchema
		config.IntelRdt.MemBwSchema = memBwSchema
	}

	return container.Set(config)
}

type CheckpointOptions struct {
	// Path for saving criu image files
	ImagePath string
	// Path for saving work files and logs
	WorkPath string
	// Path for previous criu image files in pre-dump
	ParentPath string
	// Leave the process running after checkpointing
	LeaveRunning bool
	// Allow open tcp connections
	TcpEstablished bool
	// Allow external unix sockets
	ExtUnixSk bool
	// Allow shell jobs
	ShellJob bool
	// Use userfaultfd to lazily restore memory pages
	LazyPages bool
	// Criu writes \\0 to this FD once lazy-pages is ready
	StatusFd string
	// ADDRESS:PORT of the page server
	PageServer string
	// Handle file locks, for safety
	FileLocks bool
	// Dump container's memory information only, leave the container running after this
	PreDump bool
	// Cgroups mode: 'soft' (default), 'full' and 'strict'
	ManageCgroupsMode string
	// Create a namespace, but don't restore its properties
	EmptyNs []string
	// Enable auto deduplication of memory images
	AutoDedup bool
}

// CheckpointContainer - saves the state of the container instance.
// Where 'containerID' is your name for the instance of the container and 'args' command arguments
// CheckpointOptions:
//    ImagePath         <string>: path for saving criu image files
//    WorkPath          <string>: path for saving work files and logs
//    ParentPath        <string>: path for previous criu image files in pre-dump
//    LeaveRunning      <bool>: leave the process running after checkpointing
//    TcpEstablished    <bool>: allow open tcp connections
//    ExtUnixSk         <bool>: allow external unix sockets
//    ShellJob          <bool>: allow shell jobs
//    LazyPages         <bool>: use userfaultfd to lazily restore memory pages
//    StatusFd          <string>: criu writes \\0 to this FD once lazy-pages is ready
//    PageServer        <string>: ADDRESS:PORT of the page server
//    FileLocks         <bool>: handle file locks, for safety
//    PreDump           <bool>: dump container's memory information only, leave the container running after this
//    ManageCgroupsMode <string>: cgroups mode: 'soft' (default), 'full' and 'strict'
//    EmptyNs           <[]string>: create a namespace, but don't restore its properties
//    AutoDedup         <bool>: enable auto deduplication of memory images
func (r *runc) CheckpointContainer(containerID string, opts CheckpointOptions) error {
	// Currently this is untested with rootless containers.
	if os.Geteuid() != 0 || system.RunningInUserNS() {
		fmt.Println("Warn: runc checkpoint is untested with rootless containers")
	}

	container, err := r.getContainer(containerID)
	if err != nil {
		return err
	}

	status, err := container.Status()
	if err != nil {
		return err
	}

	if status == libcontainer.Created || status == libcontainer.Stopped {
		return fmt.Errorf("Container cannot be checkpointed in %s state", status.String())
	}

	imagePath := opts.ImagePath
	if imagePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		imagePath = filepath.Join(cwd, "checkpoint")
	}
	if err := os.MkdirAll(imagePath, 0755); err != nil {
		return err
	}

	criuOpts := &libcontainer.CriuOpts{
		ImagesDirectory:         imagePath,
		WorkDirectory:           opts.WorkPath,
		ParentImage:             opts.ParentPath,
		LeaveRunning:            opts.LeaveRunning,
		TcpEstablished:          opts.TcpEstablished,
		ExternalUnixConnections: opts.ExtUnixSk,
		ShellJob:                opts.ShellJob,
		FileLocks:               opts.FileLocks,
		PreDump:                 opts.PreDump,
		AutoDedup:               opts.AutoDedup,
		LazyPages:               opts.LazyPages,
		StatusFd:                opts.StatusFd,
	}

	if !(criuOpts.LeaveRunning || criuOpts.PreDump) {
		// destroy container unless we tell CRIU to keep it
		defer destroy(container)
	}

	// these are the mandatory criu options for a container
	if err := setPageServer(opts.PageServer, criuOpts); err != nil {
		return err
	}
	if err := setCgroupsMode(opts.ManageCgroupsMode, criuOpts); err != nil {
		return err
	}
	if err := setEmptyNsMask(opts.EmptyNs, criuOpts); err != nil {
		return err
	}

	return container.Checkpoint(criuOpts)
}

func (r *runc) getContainer(containerID string) (libcontainer.Container, error) {
	factory, err := loadFactory(r.rootPath, r.criuPath, r.useSystemdCgroup)
	if err != nil {
		return nil, err
	}
	return factory.Load(containerID)
}

func (r *runc) createContainer(id string, spec *specs.Spec, opts CreateOptions) (libcontainer.Container, error) {
	rootlessCg, err := shouldUseRootlessCgroupManager(r.rootless, r.useSystemdCgroup)
	if err != nil {
		return nil, err
	}
	config, err := specconv.CreateLibcontainerConfig(&specconv.CreateOpts{
		CgroupName:       id,
		Spec:             spec,
		RootlessEUID:     os.Geteuid() != 0,
		RootlessCgroups:  rootlessCg,
		UseSystemdCgroup: r.useSystemdCgroup,
		NoPivotRoot:      opts.NoPivot,
		NoNewKeyring:     opts.NoNewKeyring,
	})
	if err != nil {
		return nil, err
	}
	factory, err := loadFactory(r.rootPath, r.criuPath, r.useSystemdCgroup)
	if err != nil {
		return nil, err
	}
	return factory.Create(id, config)
}

func (r *runc) runContainer(containerID string, spec *specs.Spec, action CtAct, criuOpts *libcontainer.CriuOpts, opts RunOptions) error {

	pidFile, err := revisePidFile(opts.PidFile)
	if err != nil {
		return err
	}

	container, err := r.createContainer(containerID, spec, CreateOptions{
		Bundle:       opts.Bundle,
		NoNewKeyring: opts.NoNewKeyring,
		NoPivot:      opts.NoPivot,
	})
	if err != nil {
		return err
	}

	notifySocket := newNotifySocket(containerID, r.rootPath, r.notifySocket)
	if notifySocket != nil {
		if err := notifySocket.setupSpec(spec); err != nil {
			return err
		}
	}

	if notifySocket != nil {
		if err := notifySocket.setupSocketDirectory(); err != nil {
			return err
		}
		if err := notifySocket.bindSocket(); err != nil {
			return err
		}
	}

	// Support on-demand socket activation by passing file descriptors into the container init process.
	listenFDs := make([]*os.File, 0)

	if r.listenFDS != "" {
		listenFDs = activation.Files(false)
	}

	logLevel := "info"
	if r.debug {
		logLevel = "debug"
	}

	rnr := &runner{
		init:            true,
		shouldDestroy:   true,
		container:       container,
		listenFDs:       listenFDs,
		notifySocket:    notifySocket,
		pidFile:         pidFile,
		logLevel:        logLevel,
		action:          action,
		detach:          opts.Detach,
		criuOpts:        criuOpts,
		enableSubreaper: !opts.NoSubreaper,
		preserveFDs:     opts.PreserveFds,
		consoleSocket:   opts.ConsoleSocket,
	}

	_, err = rnr.run(spec.Process)
	if err != nil {
		return err
	}

	return nil
}

func (r *runc) killContainer(container libcontainer.Container) error {
	_ = container.Signal(unix.SIGKILL, false)
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		if err := container.Signal(unix.Signal(0), false); err != nil {
			destroy(container)
			return nil
		}
	}
	return fmt.Errorf("container init still running")
}

func (r *runc) getProcess(bundle string, args []string, opts ExecOptions) (*specs.Process, error) {
	if opts.Process != "" {
		f, err := os.Open(opts.Process)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		var p specs.Process
		if err := json.NewDecoder(f).Decode(&p); err != nil {
			return nil, err
		}
		return &p, validateProcessSpec(&p)
	}

	// process via cli flags
	if err := os.Chdir(bundle); err != nil {
		return nil, err
	}

	spec, err := loadSpec(specConfig)
	if err != nil {
		return nil, err
	}

	p := spec.Process
	p.Args = args[1:]
	// append the passed env variables
	p.Env = append(p.Env, opts.Env...)
	// set the tty
	p.Terminal = false

	// override the cwd, if passed
	if opts.Cwd != "" {
		p.Cwd = opts.Cwd
	}
	if opts.Apparmor != "" {
		p.ApparmorProfile = opts.Apparmor
	}
	if opts.ProcessLabel != "" {
		p.SelinuxLabel = opts.ProcessLabel
	}
	if opts.Tty {
		p.Terminal = opts.Tty
	}
	if opts.NoNewPrivs {
		p.NoNewPrivileges = opts.NoNewPrivs
	}

	if caps := opts.Cap; len(caps) > 0 {
		for _, c := range caps {
			p.Capabilities.Bounding = append(p.Capabilities.Bounding, c)
			p.Capabilities.Inheritable = append(p.Capabilities.Inheritable, c)
			p.Capabilities.Effective = append(p.Capabilities.Effective, c)
			p.Capabilities.Permitted = append(p.Capabilities.Permitted, c)
			p.Capabilities.Ambient = append(p.Capabilities.Ambient, c)
		}
	}

	// override the user, if passed
	if opts.User != "" {
		u := strings.SplitN(opts.User, ":", 2)
		if len(u) > 1 {
			gid, err := strconv.Atoi(u[1])
			if err != nil {
				return nil, fmt.Errorf("parsing %s as int for gid failed: %v", u[1], err)
			}
			p.User.GID = uint32(gid)
		}
		uid, err := strconv.Atoi(u[0])
		if err != nil {
			return nil, fmt.Errorf("parsing %s as int for uid failed: %v", u[0], err)
		}
		p.User.UID = uint32(uid)
	}

	for _, gid := range opts.AdditionalGids {
		if gid < 0 {
			return nil, fmt.Errorf("additional-gids must be a positive number %d", gid)
		}
		p.User.AdditionalGids = append(p.User.AdditionalGids, uint32(gid))
	}

	return p, validateProcessSpec(p)
}

func setPageServer(pageServer string, options *libcontainer.CriuOpts) error {
	// The dump image can be sent to a criu page server (optional)
	if pageServer != "" {
		addressPort := strings.Split(pageServer, ":")
		if len(addressPort) != 2 {
			return fmt.Errorf("Use --page-server ADDRESS:PORT to specify page server")
		}
		portInt, err := strconv.Atoi(addressPort[1])
		if err != nil {
			return fmt.Errorf("Invalid port number")
		}
		options.PageServer = libcontainer.CriuPageServerInfo{
			Address: addressPort[0],
			Port:    int32(portInt),
		}
	}
	return nil
}

func setCgroupsMode(cgroupsMode string, options *libcontainer.CriuOpts) error {
	if cgroupsMode != "" {
		switch cgroupsMode {
		case "soft":
			options.ManageCgroupsMode = libcontainer.CRIU_CG_MODE_SOFT
			return nil
		case "full":
			options.ManageCgroupsMode = libcontainer.CRIU_CG_MODE_FULL
			return nil
		case "strict":
			options.ManageCgroupsMode = libcontainer.CRIU_CG_MODE_STRICT
			return nil
		default:
			return fmt.Errorf("Invalid manage cgroups mode")
		}
	}
	return nil
}

func setEmptyNsMask(emptyNs []string, options *libcontainer.CriuOpts) error {
	var namespaceMapping = map[specs.LinuxNamespaceType]int{
		specs.NetworkNamespace: unix.CLONE_NEWNET,
	}

	/* Runc doesn't manage network devices and their configuration */
	nsmask := unix.CLONE_NEWNET
	for _, ns := range emptyNs {
		f, exists := namespaceMapping[specs.LinuxNamespaceType(ns)]
		if !exists {
			return fmt.Errorf("namespace %q is not supported", ns)
		}
		nsmask |= f
	}
	options.EmptyNs = uint32(nsmask)
	return nil
}
