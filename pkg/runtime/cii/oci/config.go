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

package oci

type ConfigDocker struct {
	// Root is a path to the "root directory" where data not
	// explicitly handled by other options will be stored.
	Root string `toml:"root"`
	// RunRoot is a path to the "run directory" where state information not
	// explicitly handled by other options will be stored.
	RunRoot string `toml:"runroot"`
	// Storage is the name of the storage driver which handles actually
	// storing the contents of containers.
	StorageDriver string `toml:"storage_driver"`
	// StorageOption is a list of storage driver specific options.
	StorageOptions []string `toml:"storage_option"`
	// LogDir is the default log directory where all logs will go unless kubelet
	// tells us to put them somewhere else.
	LogDir string `toml:"log_dir"`
	// DefaultTransport is a value we prefix to image names that fail to
	// validate source references.
	DefaultTransport string `toml:"default_transport"`
	// GlobalAuthFile is a path to a file like /var/lib/kubelet/config.json
	// containing credentials necessary for pulling images from secure
	// registries.
	GlobalAuthFile string `toml:"global_auth_file"`
	// PauseImage is the name of an image which we use to instantiate infra
	// containers.
	PauseImage string `toml:"pause_image"`
	// PauseImageAuthFile, if not empty, is a path to a file like
	// /var/lib/kubelet/config.json containing credentials necessary
	// for pulling PauseImage
	PauseImageAuthFile string `toml:"pause_image_auth_file"`
	// PauseCommand is the path of the binary we run in an infra
	// container that's been instantiated using PauseImage.
	PauseCommand string `toml:"pause_command"`
	// SignaturePolicyPath is the name of the file which decides what sort
	// of policy we use when deciding whether or not to trust an image that
	// we've pulled.  Outside of testing situations, it is strongly advised
	// that this be left unspecified so that the default system-wide policy
	// will be used.
	SignaturePolicyPath string `toml:"signature_policy"`
	// InsecureRegistries is a list of registries that must be contacted w/o
	// TLS verification.
	InsecureRegistries []string `toml:"insecure_registries"`
	// ImageVolumes controls how volumes specified in image config are handled
	ImageVolumes ImageVolumesType `toml:"image_volumes"`
	// Registries holds a list of registries used to pull unqualified images
	Registries []string `toml:"registries"`
}

// ImageVolumesType describes image volume handling strategies
type ImageVolumesType string
