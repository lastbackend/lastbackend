//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

type ServerFlags struct {
	AccessToken        string
	ClusterName        string
	ClusterDescription string
	Server             struct {
		Host string
		Port uint
		TLS  struct {
			FileCert string
			FileKey  string
			FileCA   string
		}
	}
	Vault struct {
		Token    string
		Endpoint string
	}
	Domain struct {
		Internal string
		External string
	}
	Storage struct {
		Driver string
		Etcd   struct {
			FileCert  string
			FileKey   string
			FileCA    string
			Endpoints []string
		}
	}
	CIDR string
}

func NewServerFlags() *ServerFlags {
	s := new(ServerFlags)
	return s
}

func (f *ServerFlags) AddFlags(mainfs *pflag.FlagSet) {

	fs := pflag.NewFlagSet("", pflag.ExitOnError)

	defer func() {
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()

	fs.StringVarP(&f.AccessToken, "access-token", "", "", "Access token to API server")
	fs.StringVarP(&f.ClusterName, "cluster-name", "", "", "Cluster name info")
	fs.StringVarP(&f.ClusterDescription, "cluster-desc", "", "", "Cluster description")
	fs.StringVarP(&f.Server.Host, "bind-address", "", "0.0.0.0", "Bind address for listening")
	fs.UintVarP(&f.Server.Port, "bind-port", "", 2967, "Bind address for listening")
	fs.StringVarP(&f.Server.TLS.FileCert, "tls-cert-file", "", "", "TLS cert file path")
	fs.StringVarP(&f.Server.TLS.FileKey, "tls-private-key-file", "", "", "TLS private key file path")
	fs.StringVarP(&f.Server.TLS.FileCA, "tls-ca-file", "", "", "TLS certificate authority file path")
	fs.StringVarP(&f.Vault.Token, "vault-token", "", "", "Vault access token")
	fs.StringVarP(&f.Vault.Endpoint, "vault-endpoint", "", "", "Vault access endpoint")
	fs.StringVarP(&f.Domain.Internal, "domain-internal", "", "lb.local", "Default external domain for cluster")
	fs.StringVarP(&f.Domain.External, "domain-external", "", "", "Internal domain name for cluster")
	fs.StringVarP(&f.Storage.Driver, "storage", "", "etcd", "Set storage driver (Allow: etcd, mock)")
	fs.StringVarP(&f.Storage.Etcd.FileCert, "etcd-cert-file", "", "", "ETCD database cert file path")
	fs.StringVarP(&f.Storage.Etcd.FileKey, "etcd-private-key-file", "", "", "ETCD database private key file path")
	fs.StringVarP(&f.Storage.Etcd.FileCA, "etcd-ca-file", "", "", "ETCD database certificate authority file")
	fs.StringSliceVarP(&f.Storage.Etcd.Endpoints, "etcd-endpoints", "", []string{"127.0.0.1:2379"}, "ETCD database endpoints list")
	fs.StringVarP(&f.CIDR, "services-cidr", "", "172.0.0.0/24", "Services IP CIDR for internal IPAM service")

}

func AddGlobalFlags(fs *pflag.FlagSet) {

	// lookup flags in global flag set and re-register the values with our flagset
	global := pflag.CommandLine
	local := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	pflagRegister(global, local, "verbose", "v")

	fs.AddFlagSet(local)
}

func pflagRegister(global, local *pflag.FlagSet, globalName string, shaortName string) {
	if f := global.Lookup(globalName); f != nil {
		f.Name = normalize(f.Name)
		f.Shorthand = shaortName
		local.AddFlag(f)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagset (pflag): %s", globalName))
	}
}

func normalize(s string) string {
	return strings.Replace(s, "_", "-", -1)
}

