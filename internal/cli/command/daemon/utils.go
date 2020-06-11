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
	"io/ioutil"
	"os"

	"github.com/lastbackend/lastbackend/internal/daemon/config"
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func SetConfigFromFile(configPath string, cfg *config.Config) error {
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

func SetConfigFromEnvs(cfg *config.Config) error {

	// **********************************************************
	// General config *******************************************
	// **********************************************************
	token := os.Getenv("LB_ACCESS_TOKEN")
	if token != "" {
		cfg.Security.Token = token
	}
	noSchedule := os.Getenv("LB_NO_SCHEDULE")
	if noSchedule != "" {
		cfg.DisableSchedule = converter.StringToBool(noSchedule)
	}
	disableServer := os.Getenv("LB_AGENT")
	if disableServer != "" {
		cfg.DisableServer = converter.StringToBool(disableServer)
	}
	rootDir := os.Getenv("LB_ROOT_DIR")
	if rootDir != "" {
		cfg.RootDir = rootDir
	}
	debug := os.Getenv("LB_DEBUG")
	if debug != "" {
		cfg.Debug = converter.StringToBool(debug)
	}

	// **********************************************************
	// Server config ********************************************
	// **********************************************************

	host := os.Getenv("LB_BIND_ADDRESS")
	if host != "" {
		cfg.APIServer.Host = host
	}
	port := os.Getenv("LB_BIND_PORT")
	if port != "" {
		cfg.APIServer.Port = converter.StringToUint(port)
	}
	tlsVerify := os.Getenv("LB_TLS_VERIFY")
	if tlsVerify != "" {
		cfg.APIServer.TLS.Verify = converter.StringToBool(tlsVerify)
	}
	tlsCa := os.Getenv("LB_TLS_CA_FILE")
	if tlsCa != "" {
		cfg.APIServer.TLS.FileCA = tlsCa
	}
	tlsCert := os.Getenv("LB_TLS_CERT_FILE")
	if tlsCert != "" {
		cfg.APIServer.TLS.FileCert = tlsCert
	}
	tlsKey := os.Getenv("LB_TLS_PRIVATE_KEY_FILE")
	if tlsKey != "" {
		cfg.APIServer.TLS.FileKey = tlsKey
	}
	vaultToken := os.Getenv("LB_VAULT_TOKEN")
	if vaultToken != "" {
		cfg.Vault.Token = vaultToken
	}
	vaultEndpoint := os.Getenv("LB_VAULT_ENDPOINT")
	if vaultEndpoint != "" {
		cfg.Vault.Token = vaultEndpoint
	}
	domainInternal := os.Getenv("LB_DOMAIN_INTERNAL")
	if domainInternal != "" {
		cfg.Vault.Token = domainInternal
	}
	domainExternal := os.Getenv("LB_DOMAIN_EXTERNAL")
	if domainExternal != "" {
		cfg.Domain.External = domainExternal
	}

	// **********************************************************
	// Agent config *********************************************
	// **********************************************************
	storageDriver := os.Getenv("LB_STORAGE_DRIVER")
	if storageDriver != "" {
		cfg.StorageDriver = storageDriver
	}
	nodeHost := os.Getenv("LB_NODE_BIND_ADDRESS")
	if nodeHost != "" {
		cfg.NodeServer.Host = nodeHost
	}
	nodePort := os.Getenv("LB_NODE_BIND_PORT")
	if nodePort != "" {
		cfg.NodeServer.Port = converter.StringToUint(nodePort)
	}
	nodeTLSVerify := os.Getenv("LB_NODE_TLS_VERIFY")
	if nodeTLSVerify != "" {
		cfg.NodeServer.TLS.Verify = converter.StringToBool(nodeTLSVerify)
	}
	nodeTLSCa := os.Getenv("LB_NODE_TLS_CA_FILE")
	if nodeTLSCa != "" {
		cfg.NodeServer.TLS.FileCA = nodeTLSCa
	}
	nodeTLSCert := os.Getenv("LB_NODE_TLS_CERT_FILE")
	if nodeTLSCert != "" {
		cfg.NodeServer.TLS.FileCert = nodeTLSCert
	}
	nodeTLSKey := os.Getenv("LB_NODE_TLS_PRIVATE_KEY_FILE")
	if nodeTLSKey != "" {
		cfg.NodeServer.TLS.FileKey = nodeTLSKey
	}
	apiAddress := os.Getenv("LB_API_ADDRESS")
	if apiAddress != "" {
		cfg.NodeClient.Address = apiAddress
	}
	apiTLSVerify := os.Getenv("LB_API_TLS_VERIFY")
	if apiTLSVerify != "" {
		cfg.NodeClient.TLS.Verify = converter.StringToBool(apiTLSVerify)
	}
	apiTLSCa := os.Getenv("LB_API_TLS_CA_FILE")
	if apiTLSCa != "" {
		cfg.NodeClient.TLS.FileCA = apiTLSCa
	}
	apiTLSCert := os.Getenv("LB_API_TLS_CERT_FILE")
	if apiTLSCa != "" {
		cfg.NodeClient.TLS.FileCert = apiTLSCert
	}
	apiTLSKey := os.Getenv("LB_API_TLS_PRIVATE_KEY_FILE")
	if apiTLSCa != "" {
		cfg.NodeClient.TLS.FileKey = apiTLSKey
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

func SetConfigFromFlags(flags *pflag.FlagSet, cfg *config.Config) error {

	debug, err := flags.GetBool("debug")
	if err != nil {
		return errors.Wrapf(err, `"debugr" flag is non-bool, programmer error, please correct`)
	}

	disableServer, err := flags.GetBool("agent")
	if err != nil {
		return errors.Wrapf(err, "\"agent\" flag is non-bool, programmer error, please correct")
	}

	disableSchedule, err := flags.GetBool("no-schedule")
	if err != nil {
		return errors.Wrapf(err, "\"no-schedule\" flag is non-bool, programmer error, please correct")
	}

	rootDir, err := flags.GetString("root-dir")
	if err != nil {
		return errors.Wrapf(err, `"root-dir" flag is non-string, programmer error, please correct`)
	}

	token, err := flags.GetString("access-token")
	if err != nil {
		return errors.Wrapf(err, `"access-token" flag is non-string, programmer error, please correct`)
	}

	bindAddress, err := flags.GetString("bind-address")
	if err != nil {
		return errors.Wrapf(err, `"bind-address" flag is non-string, programmer error, please correct`)
	}

	bindPort, err := flags.GetUint("bind-port")
	if err != nil {
		return errors.Wrapf(err, `"bind-port" flag is non-unit, programmer error, please correct`)
	}

	tlsVerify, err := flags.GetBool("tls-verify")
	if err != nil {
		return errors.Wrapf(err, `"tls-verify" flag is non-bool, programmer error, please correct`)
	}

	tlsCaFile, err := flags.GetString("tls-ca-file")
	if err != nil {
		return errors.Wrapf(err, `"tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	tlsCertFile, err := flags.GetString("tls-cert-file")
	if err != nil {
		return errors.Wrapf(err, `"tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	tlsKeyFile, err := flags.GetString("tls-private-key-file")
	if err != nil {
		return errors.Wrapf(err, `"tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	vaultToken, err := flags.GetString("vault-token")
	if err != nil {
		return errors.Wrapf(err, `"vault-token" flag is non-string, programmer error, please correct`)
	}

	vaultEndpoint, err := flags.GetString("vault-endpoint")
	if err != nil {
		return errors.Wrapf(err, `"vault-endpoint" flag is non-string, programmer error, please correct`)
	}

	internalDomain, err := flags.GetString("domain-internal")
	if err != nil {
		return errors.Wrapf(err, `"domain-internal" flag is non-string, programmer error, please correct`)
	}

	externalDomain, err := flags.GetString("domain-external")
	if err != nil {
		return errors.Wrapf(err, `"domain-external" flag is non-string, programmer error, please correct`)
	}

	nodeBindAddress, err := flags.GetString("node-bind-address")
	if err != nil {
		return errors.Wrapf(err, `"bind-address" flag is non-string, programmer error, please correct`)
	}

	nodeBindPort, err := flags.GetUint("node-bind-port")
	if err != nil {
		return errors.Wrapf(err, `"bind-port" flag is non-uint, programmer error, please correct`)
	}

	nodeTLSVerify, err := flags.GetBool("node-tls-verify")
	if err != nil {
		return errors.Wrapf(err, `"node-tls-verify" flag is non-bool, programmer error, please correct`)
	}

	nodeTLSCaFile, err := flags.GetString("node-tls-ca-file")
	if err != nil {
		return errors.Wrapf(err, `"node-tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	nodeTLSCertFile, err := flags.GetString("node-tls-cert-file")
	if err != nil {
		return errors.Wrapf(err, `"node-tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	nodeTLSKeyFile, err := flags.GetString("node-tls-private-key-file")
	if err != nil {
		return errors.Wrapf(err, `"node-tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	apiAddress, err := flags.GetString("api-address")
	if err != nil {
		return errors.Wrapf(err, `"api-address" flag is non-string, programmer error, please correct`)
	}

	apiTLSVerify, err := flags.GetBool("api-tls-verify")
	if err != nil {
		return errors.Wrapf(err, `"api-tls-verify" flag is non-bool, programmer error, please correct`)
	}

	apiTLSCaFile, err := flags.GetString("api-tls-ca-file")
	if err != nil {
		return errors.Wrapf(err, `"api-tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	apiTLSCertFile, err := flags.GetString("api-tls-cert-file")
	if err != nil {
		return errors.Wrapf(err, `"api-tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	apiTLSKeyFile, err := flags.GetString("api-tls-private-key-file")
	if err != nil {
		return errors.Wrapf(err, `"api-tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	storageDriver, err := flags.GetString("storage-driver")
	if err != nil {
		return errors.Wrapf(err, `"storage-driver" flag is non-string, programmer error, please correct`)
	}

	manifestDir, err := flags.GetString("manifest-dir")
	if err != nil {
		return errors.Wrapf(err, `"manifest-dir" flag is non-string, programmer error, please correct`)
	}

	cidr, err := flags.GetString("services-cidr")
	if err != nil {
		return errors.Wrapf(err, `"services-cidr" flag is non-string, programmer error, please correct`)
	}

	if debug {
		cfg.Debug = debug
	}
	if disableSchedule {
		cfg.DisableSchedule = disableSchedule
	}
	if disableServer {
		cfg.DisableServer = disableServer
	}
	if rootDir != "" {
		cfg.RootDir = rootDir
	}
	if token != "" {
		cfg.Security.Token = token
	}
	if bindAddress != "" && (cfg.APIServer.Host != "" && bindAddress != config.DefaultBindServerAddress || cfg.APIServer.Host == "") {
		cfg.APIServer.Host = bindAddress
	}
	if bindPort != 0 && (cfg.APIServer.Port != 0 && bindPort != config.DefaultBindServerPort || cfg.APIServer.Port == 0) {
		cfg.APIServer.Port = bindPort
	}
	if tlsVerify {
		cfg.APIServer.TLS.Verify = tlsVerify
	}
	if tlsCaFile != "" {
		cfg.APIServer.TLS.FileCA = tlsCaFile
	}
	if tlsCertFile != "" {
		cfg.APIServer.TLS.FileCert = tlsCertFile
	}
	if tlsKeyFile != "" {
		cfg.APIServer.TLS.FileKey = tlsKeyFile
	}
	if vaultToken != "" {
		cfg.Vault.Token = vaultToken
	}
	if vaultEndpoint != "" {
		cfg.Vault.Endpoint = vaultEndpoint
	}
	if internalDomain != "" && (cfg.Domain.Internal != "" && internalDomain != config.DefaultInternalDomain || cfg.Domain.Internal == "") {
		cfg.Domain.Internal = internalDomain
	}
	if externalDomain != "" {
		cfg.Domain.External = externalDomain
	}
	if nodeBindAddress != "" && (cfg.NodeServer.Host != "" && nodeBindAddress != config.DefaultBindServerAddress || cfg.NodeServer.Host == "") {
		cfg.NodeServer.Host = nodeBindAddress
	}
	if nodeBindPort != 0 && (cfg.NodeServer.Port != 0 && nodeBindPort != config.DefaultBindServerPort || cfg.NodeServer.Port == 0) {
		cfg.NodeServer.Port = nodeBindPort
	}
	if nodeTLSVerify {
		cfg.NodeServer.TLS.Verify = nodeTLSVerify
	}
	if nodeTLSCaFile != "" {
		cfg.NodeServer.TLS.FileCA = nodeTLSCaFile
	}
	if nodeTLSCertFile != "" {
		cfg.NodeServer.TLS.FileCert = nodeTLSCertFile
	}
	if nodeTLSKeyFile != "" {
		cfg.NodeServer.TLS.FileKey = nodeTLSKeyFile
	}
	if apiAddress != "" {
		cfg.NodeClient.Address = apiAddress
	}
	if apiTLSVerify {
		cfg.NodeClient.TLS.Verify = apiTLSVerify
	}
	if apiTLSCaFile != "" {
		cfg.NodeClient.TLS.FileCA = apiTLSCaFile
	}
	if apiTLSCertFile != "" {
		cfg.NodeClient.TLS.FileCert = apiTLSCertFile
	}
	if apiTLSKeyFile != "" {
		cfg.NodeClient.TLS.FileKey = apiTLSKeyFile
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
