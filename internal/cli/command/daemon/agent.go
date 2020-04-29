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

package daemon

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func SetAgentConfigFromFile(configPath string, cfg *config.Config) error {
	_, err := os.Stat(configPath)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, &cfg)
}

func SetAgentConfigFromEnvs(cfg *config.Config) error {

	token := os.Getenv("LB_ACCESS_TOKEN")
	if token != "" {
		cfg.Security.Token = token
	}
	host := os.Getenv("LB_NODE_BIND_ADDRESS")
	if token != "" {
		cfg.Server.Host = host
	}
	port := os.Getenv("LB_NODE_BIND_PORT")
	if port != "" {
		cfg.Server.Port = converter.StringToUint(port)
	}
	tlsVerify := os.Getenv("LB_NODE_TLS_VERIFY")
	if tlsVerify != "" {
		cfg.Server.TLS.Verify = converter.StringToBool(tlsVerify)
	}
	tlsCa := os.Getenv("LB_NODE_TLS_CA_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileCA = tlsCa
	}
	tlsCert := os.Getenv("LB_NODE_TLS_CERT_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileCert = tlsCert
	}
	tlsKey := os.Getenv("LB_NODE_TLS_PRIVATE_KEY_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileKey = tlsKey
	}
	apiAddress := os.Getenv("LB_API_ADDRESS")
	if apiAddress != "" {
		cfg.API.Address = apiAddress
	}
	apiTlsVerify := os.Getenv("LB_API_TLS_VERIFY")
	if apiTlsVerify != "" {
		cfg.API.TLS.Verify = converter.StringToBool(apiTlsVerify)
	}
	apiTlsCa := os.Getenv("LB_API_TLS_CA_FILE")
	if apiTlsCa != "" {
		cfg.API.TLS.FileCA = apiTlsCa
	}
	apiTlsCert := os.Getenv("LB_API_TLS_CERT_FILE")
	if apiTlsCa != "" {
		cfg.API.TLS.FileCert = apiTlsCert
	}
	apiTlsKey := os.Getenv("LB_API_TLS_PRIVATE_KEY_FILE")
	if apiTlsCa != "" {
		cfg.API.TLS.FileKey = apiTlsKey
	}
	workdir := os.Getenv("LB_WORKDIR")
	if workdir != "" {
		cfg.WorkDir = workdir
	}
	manifefstdir := os.Getenv("LB_MANIFESTDIR")
	if manifefstdir != "" {
		cfg.ManifestDir = manifefstdir
	}
	registryConfig := os.Getenv("LB_REGISTRY_CONFIG")
	if registryConfig != "" {
		cfg.Registry.Config = registryConfig
	}
	rootless := os.Getenv("LB_ROOTLESS")
	if rootless != "" {
		cfg.Rootless = converter.StringToBool(rootless)
	}
	disableSeLinux := os.Getenv("LB_DISABLE_SELINUX")
	if disableSeLinux != "" {
		cfg.DisableSeLinux = converter.StringToBool(disableSeLinux)
	}
	cidr := os.Getenv("LB_SERVICES_CIDR")
	if cidr != "" {
		cfg.CIDR = cidr
	}
	return nil
}

func SetAgentConfigFromFlags(flags *pflag.FlagSet, cfg *config.Config) error {

	token, err := flags.GetString("access-token")
	if err != nil {
		return fmt.Errorf(`"access-token" flag is non-string, programmer error, please correct`)
	}

	bindAddress, err := flags.GetString("node-bind-address")
	if err != nil {
		return fmt.Errorf(`"bind-address" flag is non-string, programmer error, please correct`)
	}

	bindPort, err := flags.GetUint("node-bind-port")
	if err != nil {
		return fmt.Errorf(`"bind-port" flag is non-uint, programmer error, please correct`)
	}

	tlsVerify, err := flags.GetBool("node-tls-verify")
	if err != nil {
		return fmt.Errorf(`"node-tls-verify" flag is non-bool, programmer error, please correct`)
	}

	tlsCaFile, err := flags.GetString("node-tls-ca-file")
	if err != nil {
		return fmt.Errorf(`"node-tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	tlsCertFile, err := flags.GetString("node-tls-cert-file")
	if err != nil {
		return fmt.Errorf(`"node-tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	tlsKeyFile, err := flags.GetString("node-tls-private-key-file")
	if err != nil {
		return fmt.Errorf(`"node-tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	apiAddress, err := flags.GetString("api-address")
	if err != nil {
		return fmt.Errorf(`"api-address" flag is non-string, programmer error, please correct`)
	}

	apiTlsVerify, err := flags.GetBool("api-tls-verify")
	if err != nil {
		return fmt.Errorf(`"api-tls-verify" flag is non-bool, programmer error, please correct`)
	}

	apiTlsCaFile, err := flags.GetString("api-tls-ca-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	apiTlsCertFile, err := flags.GetString("api-tls-cert-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	apiTlsKeyFile, err := flags.GetString("api-tls-private-key-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	workdir, err := flags.GetString("workdir")
	if err != nil {
		return fmt.Errorf(`"workdir" flag is non-string, programmer error, please correct`)
	}

	manifestdir, err := flags.GetString("manifestdir")
	if err != nil {
		return fmt.Errorf(`"manifestdir" flag is non-string, programmer error, please correct`)
	}

	registryConfig, err := flags.GetString("registry-config-path")
	if err != nil {
		return fmt.Errorf(`"registry-config-path" flag is non-string, programmer error, please correct`)
	}

	rootless, err := flags.GetBool("rootless")
	if err != nil {
		return fmt.Errorf(`"rootless" flag is non-bool, programmer error, please correct`)
	}

	disableSeLinux, err := flags.GetBool("disable-selinux")
	if err != nil {
		return fmt.Errorf(`"disable-selinux" flag is non-bool, programmer error, please correct`)
	}

	cidr, err := flags.GetString("services-cidr")
	if err != nil {
		return fmt.Errorf(`"services-cidr" flag is non-string, programmer error, please correct`)
	}

	if token != "" {
		cfg.Security.Token = token
	}
	if bindAddress != "" && (cfg.Server.Host != "" && bindAddress != config.DefaultBindServerAddress || cfg.Server.Host == "") {
		cfg.Server.Host = bindAddress
	}
	if bindPort != 0 && (cfg.Server.Port != 0 && bindPort != config.DefaultBindServerPort || cfg.Server.Port == 0) {
		cfg.Server.Port = bindPort
	}
	if tlsVerify {
		cfg.Server.TLS.Verify = tlsVerify
	}
	if tlsCaFile != "" {
		cfg.Server.TLS.FileCA = tlsCaFile
	}
	if tlsCertFile != "" {
		cfg.Server.TLS.FileCert = tlsCertFile
	}
	if tlsKeyFile != "" {
		cfg.Server.TLS.FileKey = tlsKeyFile
	}
	if apiAddress != "" {
		cfg.API.Address = apiAddress
	}
	if apiTlsVerify {
		cfg.API.TLS.Verify = apiTlsVerify
	}
	if apiTlsCaFile != "" {
		cfg.API.TLS.FileCA = apiTlsCaFile
	}
	if apiTlsCertFile != "" {
		cfg.API.TLS.FileCert = apiTlsCertFile
	}
	if apiTlsKeyFile != "" {
		cfg.API.TLS.FileKey = apiTlsKeyFile
	}
	if workdir != "" && (cfg.WorkDir != "" && workdir != config.DefaultWorkDir || cfg.WorkDir == "") {
		cfg.WorkDir = workdir
	}
	if manifestdir != "" {
		cfg.ManifestDir = manifestdir
	}
	if registryConfig != "" {
		cfg.Registry.Config = registryConfig
	}
	if rootless {
		cfg.Rootless = rootless
	}
	if disableSeLinux {
		cfg.DisableSeLinux = disableSeLinux
	}
	if cidr != "" && (cfg.CIDR != "" && cidr != config.DefaultCIDR || cfg.CIDR == "") {
		cfg.CIDR = cidr
	}

	return nil
}
