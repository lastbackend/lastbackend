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
	apiTLSVerify := os.Getenv("LB_API_TLS_VERIFY")
	if apiTLSVerify != "" {
		cfg.API.TLS.Verify = converter.StringToBool(apiTLSVerify)
	}
	apiTLSCa := os.Getenv("LB_API_TLS_CA_FILE")
	if apiTLSCa != "" {
		cfg.API.TLS.FileCA = apiTLSCa
	}
	apiTLSCert := os.Getenv("LB_API_TLS_CERT_FILE")
	if apiTLSCa != "" {
		cfg.API.TLS.FileCert = apiTLSCert
	}
	apiTLSKey := os.Getenv("LB_API_TLS_PRIVATE_KEY_FILE")
	if apiTLSCa != "" {
		cfg.API.TLS.FileKey = apiTLSKey
	}
	rootDir := os.Getenv("LB_ROOT_DIR")
	if rootDir != "" {
		cfg.RootDir = rootDir
	}
	storageDriver := os.Getenv("LB_STORAGE_DRIVER")
	if rootDir != "" {
		cfg.StorageDriver = storageDriver
	}
	manifefstDir := os.Getenv("LB_MANIFEST_DIR")
	if manifefstDir != "" {
		cfg.ManifestDir = manifefstDir
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

	apiTLSVerify, err := flags.GetBool("api-tls-verify")
	if err != nil {
		return fmt.Errorf(`"api-tls-verify" flag is non-bool, programmer error, please correct`)
	}

	apiTLSCaFile, err := flags.GetString("api-tls-ca-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	apiTLSCertFile, err := flags.GetString("api-tls-cert-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	apiTLSKeyFile, err := flags.GetString("api-tls-private-key-file")
	if err != nil {
		return fmt.Errorf(`"api-tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	rootDir, err := flags.GetString("root-dir")
	if err != nil {
		return fmt.Errorf(`"root-dir" flag is non-string, programmer error, please correct`)
	}

	storageDriver, err := flags.GetString("storage-driver")
	if err != nil {
		return fmt.Errorf(`"storage-driver" flag is non-string, programmer error, please correct`)
	}

	manifestDir, err := flags.GetString("manifest-dir")
	if err != nil {
		return fmt.Errorf(`"manifest-dir" flag is non-string, programmer error, please correct`)
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
	if apiTLSVerify {
		cfg.API.TLS.Verify = apiTLSVerify
	}
	if apiTLSCaFile != "" {
		cfg.API.TLS.FileCA = apiTLSCaFile
	}
	if apiTLSCertFile != "" {
		cfg.API.TLS.FileCert = apiTLSCertFile
	}
	if apiTLSKeyFile != "" {
		cfg.API.TLS.FileKey = apiTLSKeyFile
	}
	if rootDir != "" {
		cfg.RootDir = rootDir
	}
	if storageDriver != "" {
		cfg.StorageDriver = storageDriver
	}
	if manifestDir != "" {
		cfg.ManifestDir = manifestDir
	}
	if cidr != "" && (cfg.CIDR != "" && cidr != config.DefaultCIDR || cfg.CIDR == "") {
		cfg.CIDR = cidr
	}

	return nil
}
