package container

import (
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"k8s.io/client-go/1.5/pkg/api"
)

const kind = "container"

type ContainerList struct {
	ListMeta   common.ListMeta `json:"meta"`
	Containers []Container     `json:"containers"`
}

type ContainerStatus struct {
	ContainerStates []api.ContainerState `json:"container_states"`
}

type Container struct {
	TypeMeta   common.TypeMeta `json:"spec"`
	Name       string          `json:"name"`
	Image      string          `json:"image"`
	Command    []string        `json:"command,omitempty"`
	Args       []string        `json:"args,omitempty"`
	WorkingDir string          `json:"workdir,omitempty"`
	Ports      []Port          `json:"ports,omitempty"`
	Env        []EnvVar        `json:"env,omitempty"`
	Volumes    []Volume        `json:"volumes,omitempty"`
}

// Port represents a network port in a single container
type Port struct {
	// Optional: If specified, this must be an IANA_SVC_NAME  Each named port
	// in a pod must have a unique name.
	Name string `json:"name,omitempty"`
	// Optional: If specified, this must be a valid port number, 0 < x < 65536.
	// If Host network is specified, this must match Container port.
	HostPort int32 `json:"host,omitempty"`
	// Required: This must be a valid port number, 0 < x < 65536.
	ContainerPort int32 `json:"container"`
	// Required: Supports "TCP" and "UDP".
	Protocol string `json:"protocol,omitempty"`
	// Optional: What host IP to bind the external port to.
	HostIP string `json:"ip,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// VolumeMount describes a mounting of a Volume within a container.
type Volume struct {
	// Required: This must match the Name of a Volume [above].
	Name string `json:"name"`
	// Optional: Defaults to false (read-write).
	ReadOnly bool `json:"readonly,omitempty"`
	// Required. Must not contain ':'.
	MountPath string `json:"mountpath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	SubPath string `json:"subpath,omitempty"`
}

func CreateContainerList(containers []api.Container) *ContainerList {

	containerList := ContainerList{
		ListMeta:   common.ListMeta{Total: len(containers)},
		Containers: make([]Container, 0),
	}

	for _, c := range containers {

		var p = Container{
			TypeMeta:   common.NewTypeMeta(kind),
			Name:       c.Name,
			Image:      c.Image,
			Command:    c.Command,
			Args:       c.Args,
			WorkingDir: c.WorkingDir,
		}

		for _, port := range c.Ports {
			p.Ports = append(p.Ports, Port{
				Name:          port.Name,
				HostPort:      port.HostPort,
				ContainerPort: port.ContainerPort,
				Protocol:      string(port.Protocol),
				HostIP:        port.HostIP,
			})
		}

		for _, env := range c.Env {
			p.Env = append(p.Env, EnvVar{
				Name:  env.Name,
				Value: env.Value,
			})
		}

		for _, volume := range c.VolumeMounts {
			p.Volumes = append(p.Volumes, Volume{
				Name:      volume.Name,
				ReadOnly:  volume.ReadOnly,
				MountPath: volume.MountPath,
				SubPath:   volume.SubPath,
			})
		}

		containerList.Containers = append(containerList.Containers, p)
	}

	return &containerList
}
