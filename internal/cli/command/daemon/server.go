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
	"github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func SetServerConfigFromFile(configPath string, cfg *config.Config) error {
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

func SetServerConfigFromEnvs(cfg *config.Config) error {

	token := os.Getenv("LB_ACCESS_TOKEN")
	if token != "" {
		cfg.Security.Token = token
	}
	host := os.Getenv("LB_BIND_ADDRESS")
	if token != "" {
		cfg.Server.Host = host
	}
	port := os.Getenv("LB_BIND_PORT")
	if port != "" {
		cfg.Server.Port = converter.StringToUint(port)
	}
	clusterName := os.Getenv("LB_CLUSTER_NAME")
	if clusterName != "" {
		cfg.ClusterName = clusterName
	}
	clusterDesc := os.Getenv("LB_CLUSTER_DESCRIPTION")
	if clusterDesc != "" {
		cfg.ClusterDescription = clusterDesc
	}
	tlsVerify := os.Getenv("LB_TLS_VERIFY")
	if tlsVerify != "" {
		cfg.Server.TLS.Verify = converter.StringToBool(tlsVerify)
	}
	tlsCa := os.Getenv("LB_TLS_CA_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileCA = tlsCa
	}
	tlsCert := os.Getenv("LB_TLS_CERT_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileCert = tlsCert
	}
	tlsKey := os.Getenv("LB_TLS_PRIVATE_KEY_FILE")
	if tlsCa != "" {
		cfg.Server.TLS.FileKey = tlsKey
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
	workdir := os.Getenv("LB_WORKDIR")
	if workdir != "" {
		cfg.WorkDir = workdir
	}
	rootless := os.Getenv("LB_ROOTLESS")
	if rootless != "" {
		cfg.Rootless = converter.StringToBool(rootless)
	}

	return nil
}

func SetServerConfigFromFlagsAndEnvs(flags *pflag.FlagSet, cfg *config.Config) error {

	token, err := flags.GetString("access-token")
	if err != nil {
		return fmt.Errorf(`"access-token" flag is non-string, programmer error, please correct`)
	}

	bindAddress, err := flags.GetString("bind-address")
	if err != nil {
		return fmt.Errorf(`"bind-address" flag is non-string, programmer error, please correct`)
	}

	bindPort, err := flags.GetUint("bind-port")
	if err != nil {
		return fmt.Errorf(`"bind-port" flag is non-unit, programmer error, please correct`)
	}

	clusterName, err := flags.GetString("cluster-name")
	if err != nil {
		return fmt.Errorf(`"cluster-name" flag is non-string, programmer error, please correct`)
	}

	clusterDesc, err := flags.GetString("cluster-desc")
	if err != nil {
		return fmt.Errorf(`"cluster-desc" flag is non-string, programmer error, please correct`)
	}

	tlsVerify, err := flags.GetBool("tls-verify")
	if err != nil {
		return fmt.Errorf(`"tls-verify" flag is non-bool, programmer error, please correct`)
	}

	tlsCaFile, err := flags.GetString("tls-ca-file")
	if err != nil {
		return fmt.Errorf(`"tls-ca-file" flag is non-string, programmer error, please correct`)
	}

	tlsCertFile, err := flags.GetString("tls-cert-file")
	if err != nil {
		return fmt.Errorf(`"tls-cert-file" flag is non-string, programmer error, please correct`)
	}

	tlsKeyFile, err := flags.GetString("tls-private-key-file")
	if err != nil {
		return fmt.Errorf(`"tls-private-key-file" flag is non-string, programmer error, please correct`)
	}

	vaultToken, err := flags.GetString("vault-token")
	if err != nil {
		return fmt.Errorf(`"vault-token" flag is non-string, programmer error, please correct`)
	}

	vaultEndpoint, err := flags.GetString("vault-endpoint")
	if err != nil {
		return fmt.Errorf(`"vault-endpoint" flag is non-string, programmer error, please correct`)
	}

	internalDomain, err := flags.GetString("domain-internal")
	if err != nil {
		return fmt.Errorf(`"domain-internal" flag is non-string, programmer error, please correct`)
	}

	externalDomain, err := flags.GetString("domain-external")
	if err != nil {
		return fmt.Errorf(`"domain-external" flag is non-string, programmer error, please correct`)
	}

	workdir, err := flags.GetString("workdir")
	if err != nil {
		return fmt.Errorf(`"workdir" flag is non-string, programmer error, please correct`)
	}

	rootless, err := flags.GetBool("rootless")
	if err != nil {
		return fmt.Errorf(`"rootless" flag is non-bool, programmer error, please correct`)
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
	if clusterName != "" {
		cfg.ClusterName = clusterName
	}
	if clusterDesc != "" {
		cfg.ClusterDescription = clusterDesc
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
	if workdir != "" && (cfg.WorkDir != "" && workdir != config.DefaultWorkDir || cfg.WorkDir == "") {
		cfg.WorkDir = workdir
	}
	if rootless {
		cfg.Rootless = rootless
	}

	return nil
}
