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

package docker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/containerd/containerd/errdefs"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
)

const (
	// GrantPull represents the capability to fetch manifests and blobs by digest
	GrantPull RegistryGrant = 1 << iota
	// GrantResolve represents the capability to fetch manifests by name
	GrantResolve
	// GrantPush represents the capability to push blobs and manifests
	GrantPush
)

// authorizer is used to authorize HTTP requests based on 401 HTTP responses.
// This responsible for caching tokens or credentials used by requests.
type Authorizer interface {

	// Authorize sets the Authorization header for this request:
	//   "Bearer <some bearer token>"
	//   "Basic <base64 encoded credentials>"
	// The request remains unmodified if no authorization is found for the request.
	Authorize(context.Context, *http.Request) error

	// AddResponses adds a 401 response for the authorizer, which must be considered
	// when authorizing requests. The last answer must be unauthorized,
	// and previous requests are used to consider redirects and retries that could lead to 401.
	AddResponses(context.Context, []*http.Response) error
}

// HostOptions is used to configure registry hosts
type HostOptions struct {
	HostDir       func(string) (string, error)
	Credentials   func(host string) (string, string, error)
	DefaultTLS    *tls.Config
	DefaultScheme string
}

type hostConfig struct {
	scheme      string
	host        string
	path        string
	grants      RegistryGrant
	caCerts     []string
	clientPairs [][2]string
	skipVerify  *bool
}

type RegistryHosts func(string) ([]RegistryHost, error)

// RegistryGrant represents the set of operations for which the
// RegistryHost host may be trusted to perform.
//
// Public RegistryHost  : pull, push, resolve
// Private RegistryHost : pull, push, resolve
// Public Mirror    : pull
// Private Mirror   : pull, resolve
type RegistryGrant uint8

func (rg RegistryGrant) Has(t RegistryGrant) bool {
	return rg&t == t
}

func newRegistry(client *http.Client, url url.URL, grant RegistryGrant) RegistryHost {
	return RegistryHost{
		Client: client,
		Host:   url.Host,
		Scheme: url.Scheme,
		Path:   url.Path,
		Grant:  grant,
	}
}

// Registries fetches the registry hosts for a given namespace,
// provided by the host component of an distribution image reference.
type Registries func(string) ([]RegistryHost, error)

// RegistryHost represents a complete configuration for a RegistryHost node
// that provides connection configuration and location.
type RegistryHost struct {
	Client     *http.Client
	Authorizer Authorizer
	Host       string
	Scheme     string
	Path       string
	Grant      RegistryGrant
}

func ConfigureHosts(ctx context.Context, options HostOptions) RegistryHosts {
	return func(host string) ([]RegistryHost, error) {
		var hosts []hostConfig
		if options.HostDir != nil {
			dir, err := options.HostDir(host)
			if err != nil && !errdefs.IsNotFound(err) {
				return nil, err
			}
			if dir != "" {
				log.G(ctx).WithField("dir", dir).Debug("loading host directory")
				hosts, err = loadHostDir(ctx, dir)
				if err != nil {
					return nil, err
				}
			}

		}

		// If hosts was not set, add a default host
		// NOTE: Check nil here and not empty, the host may be
		// intentionally configured to not have any endpoints
		if hosts == nil {
			hosts = make([]hostConfig, 1)
		}
		if len(hosts) > 0 && hosts[len(hosts)-1].host == "" {
			if host == "docker.io" {
				hosts[len(hosts)-1].scheme = "https"
				hosts[len(hosts)-1].host = "registry-1.docker.io"
			} else {
				hosts[len(hosts)-1].host = host
				if options.DefaultScheme != "" {
					hosts[len(hosts)-1].scheme = options.DefaultScheme
				} else {
					hosts[len(hosts)-1].scheme = "https"
				}
			}
			hosts[len(hosts)-1].path = "/v2"
			hosts[len(hosts)-1].grants = GrantPull | GrantResolve | GrantPush
		}

		defaultTransport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			TLSClientConfig:       options.DefaultTLS,
			ExpectContinueTimeout: 5 * time.Second,
		}

		client := &http.Client{
			Transport: defaultTransport,
		}

		authOpts := []AuthorizerOpt{WithAuthClient(client)}
		if options.Credentials != nil {
			authOpts = append(authOpts, WithAuthCreds(options.Credentials))
		}
		authorizer := NewDockerAuthorizer(authOpts...)

		rhosts := make([]RegistryHost, len(hosts))
		for i, host := range hosts {

			rhosts[i].Scheme = host.scheme
			rhosts[i].Host = host.host
			rhosts[i].Path = host.path
			rhosts[i].Grant = host.grants

			if host.caCerts != nil || host.clientPairs != nil || host.skipVerify != nil {
				var tlsConfig *tls.Config
				if options.DefaultTLS != nil {
					tlsConfig = options.DefaultTLS.Clone()
				} else {
					tlsConfig = &tls.Config{}
				}
				if host.skipVerify != nil {
					tlsConfig.InsecureSkipVerify = *host.skipVerify
				}
				if host.caCerts != nil {
					if tlsConfig.RootCAs == nil {
						rootPool, err := rootSystemPool()
						if err != nil {
							return nil, errors.Wrap(err, "unable to initialize cert pool")
						}
						tlsConfig.RootCAs = rootPool
					}
					for _, f := range host.caCerts {
						data, err := ioutil.ReadFile(f)
						if err != nil {
							return nil, errors.Wrapf(err, "unable to read CA cert %q", f)
						}
						if !tlsConfig.RootCAs.AppendCertsFromPEM(data) {
							return nil, errors.Errorf("unable to load CA cert %q", f)
						}
					}
				}

				if host.clientPairs != nil {
					for _, pair := range host.clientPairs {
						certPEMBlock, err := ioutil.ReadFile(pair[0])
						if err != nil {
							return nil, errors.Wrapf(err, "unable to read CERT file %q", pair[0])
						}
						var keyPEMBlock []byte
						if pair[1] != "" {
							keyPEMBlock, err = ioutil.ReadFile(pair[1])
							if err != nil {
								return nil, errors.Wrapf(err, "unable to read CERT file %q", pair[1])
							}
						} else {
							// Load key block from same PEM file
							keyPEMBlock = certPEMBlock
						}
						cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
						if err != nil {
							return nil, errors.Wrap(err, "failed to load X509 key pair")
						}

						tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
					}
				}
				tr := defaultTransport.Clone()
				tr.TLSClientConfig = tlsConfig
				c := *client
				c.Transport = tr

				rhosts[i].Client = &c
				rhosts[i].Authorizer = NewDockerAuthorizer(append(authOpts, WithAuthClient(&c))...)
			} else {
				rhosts[i].Client = client
				rhosts[i].Authorizer = authorizer
			}

		}

		return rhosts, nil
	}

}

func HostDirFromRoot(root string) func(string) (string, error) {
	return func(host string) (string, error) {
		for _, p := range hostPaths(root, host) {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			} else if !os.IsNotExist(err) {
				return "", err
			}
		}
		return "", ErrNotFound
	}
}

func hostPaths(root, host string) []string {
	ch := hostDirectory(host)
	if ch == host {
		return []string{filepath.Join(root, host)}
	}

	return []string{
		filepath.Join(root, ch),
		filepath.Join(root, host),
	}
}

func rootSystemPool() (*x509.CertPool, error) {
	return x509.SystemCertPool()
}

// hostDirectory converts ":port" to "_port_" in directory names
func hostDirectory(host string) string {
	idx := strings.LastIndex(host, ":")
	if idx > 0 {
		return host[:idx] + "_" + host[idx+1:] + "_"
	}
	return host
}

func loadHostDir(ctx context.Context, hostsDir string) ([]hostConfig, error) {
	b, err := ioutil.ReadFile(filepath.Join(hostsDir, "hosts.toml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(b) == 0 {
		// If hosts.toml does not exist, fallback to checking for
		// certificate files based on Docker's certificate file
		// pattern (".crt", ".cert", ".key" files)
		return loadCertFiles(ctx, hostsDir)
	}

	hosts, err := parseHostsFile(ctx, hostsDir, b)
	if err != nil {
		log.G(ctx).WithError(err).Error("failed to decode hosts.toml")
		// Fallback to checking certificate files
		return loadCertFiles(ctx, hostsDir)
	}

	return hosts, nil
}

type hostFileConfig struct {
	Grants []string `toml:"grants"`
	// CACert can be a string or an array of strings
	CACert     toml.Primitive `toml:"ca"`
	Client     toml.Primitive `toml:"client"`
	SkipVerify bool           `toml:"skip_verify"`
}

type configFile struct {
	// hostConfig holds defaults for all hosts as well as
	// for the default server
	hostFileConfig

	// Server specifies the default server. When `host` is
	// also specified, those hosts are tried first.
	Server string `toml:"server"`

	// HostConfigs store the per-host configuration
	HostConfigs map[string]hostFileConfig `toml:"host"`
}

func parseHostsFile(ctx context.Context, baseDir string, b []byte) ([]hostConfig, error) {
	var c configFile
	md, err := toml.Decode(string(b), &c)
	if err != nil {
		return nil, err
	}

	var orderedHosts []string
	for _, key := range md.Keys() {
		if len(key) >= 2 {
			if key[0] == "host" && (len(orderedHosts) == 0 || orderedHosts[len(orderedHosts)-1] != key[1]) {
				orderedHosts = append(orderedHosts, key[1])
			}
		}
	}

	if c.Server != "" {
		c.HostConfigs[c.Server] = c.hostFileConfig
		orderedHosts = append(orderedHosts, c.Server)
	} else if len(orderedHosts) == 0 {
		c.HostConfigs[""] = c.hostFileConfig
		orderedHosts = append(orderedHosts, "")
	}
	hosts := make([]hostConfig, len(orderedHosts))
	for i, server := range orderedHosts {
		hostConfig := c.HostConfigs[server]

		if !strings.HasPrefix(server, "http") {
			server = "https://" + server
		}
		u, err := url.Parse(server)
		if err != nil {
			return nil, errors.Errorf("unable to parse server %v", server)
		}
		hosts[i].scheme = u.Scheme
		hosts[i].host = u.Host

		// TODO: Handle path based on registry protocol
		// Define a registry protocol type
		//   OCI v1    - Always use given path as is
		//   Docker v2 - Always ensure ends with /v2/
		if len(u.Path) > 0 {
			u.Path = path.Clean(u.Path)
			if !strings.HasSuffix(u.Path, "/v2") {
				u.Path = u.Path + "/v2"
			}
		} else {
			u.Path = "/v2"
		}
		hosts[i].path = u.Path

		if hosts[i].scheme == "https" {
			hosts[i].skipVerify = &hostConfig.SkipVerify
		}

		if len(hostConfig.Grants) > 0 {
			for _, c := range hostConfig.Grants {
				switch strings.ToLower(c) {
				case "pull":
					hosts[i].grants |= GrantPull
				case "resolve":
					hosts[i].grants |= GrantResolve
				case "push":
					hosts[i].grants |= GrantPush
				default:
					return nil, errors.Errorf("unknown capability %v", c)
				}
			}
		} else {
			hosts[i].grants = GrantPull | GrantResolve | GrantPush
		}

		baseKey := []string{}
		if server != "" {
			baseKey = append(baseKey, "host", server)
		}
		caKey := append(baseKey, "ca")
		if md.IsDefined(caKey...) {
			switch t := md.Type(caKey...); t {
			case "String":
				var caCert string
				if err := md.PrimitiveDecode(hostConfig.CACert, &caCert); err != nil {
					return nil, errors.Wrap(err, "failed to decode \"ca\"")
				}
				hosts[i].caCerts = []string{makeAbsPath(caCert, baseDir)}
			case "Array":
				var caCerts []string
				if err := md.PrimitiveDecode(hostConfig.CACert, &caCerts); err != nil {
					return nil, errors.Wrap(err, "failed to decode \"ca\"")
				}
				for i, p := range caCerts {
					caCerts[i] = makeAbsPath(p, baseDir)
				}

				hosts[i].caCerts = caCerts
			default:
				return nil, errors.Errorf("invalid type %v for \"ca\"", t)
			}
		}

		clientKey := append(baseKey, "client")
		if md.IsDefined(clientKey...) {
			switch t := md.Type(clientKey...); t {
			case "String":
				var clientCert string
				if err := md.PrimitiveDecode(hostConfig.Client, &clientCert); err != nil {
					return nil, errors.Wrap(err, "failed to decode \"ca\"")
				}
				hosts[i].clientPairs = [][2]string{{makeAbsPath(clientCert, baseDir), ""}}
			case "Array":
				var clientCerts []interface{}
				if err := md.PrimitiveDecode(hostConfig.Client, &clientCerts); err != nil {
					return nil, errors.Wrap(err, "failed to decode \"ca\"")
				}
				for _, pairs := range clientCerts {
					switch p := pairs.(type) {
					case string:
						hosts[i].clientPairs = append(hosts[i].clientPairs, [2]string{makeAbsPath(p, baseDir), ""})
					case []interface{}:
						var pair [2]string
						if len(p) > 2 {
							return nil, errors.Errorf("invalid pair %v for \"client\"", p)
						}
						for pi, cp := range p {
							s, ok := cp.(string)
							if !ok {
								return nil, errors.Errorf("invalid type %T for \"client\"", cp)
							}
							pair[pi] = makeAbsPath(s, baseDir)
						}
						hosts[i].clientPairs = append(hosts[i].clientPairs, pair)
					default:
						return nil, errors.Errorf("invalid type %T for \"client\"", p)
					}
				}
			default:
				return nil, errors.Errorf("invalid type %v for \"client\"", t)
			}
		}
	}

	return hosts, nil
}

func makeAbsPath(p string, base string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(base, p)
}

func loadCertFiles(ctx context.Context, certsDir string) ([]hostConfig, error) {
	fs, err := ioutil.ReadDir(certsDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	hosts := make([]hostConfig, 1)
	for _, f := range fs {
		if !f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".crt") {
			hosts[0].caCerts = append(hosts[0].caCerts, filepath.Join(certsDir, f.Name()))
		}
		if strings.HasSuffix(f.Name(), ".cert") {
			var pair [2]string
			certFile := f.Name()
			pair[0] = filepath.Join(certsDir, certFile)
			// Check if key also exists
			keyFile := certFile[:len(certFile)-5] + ".key"
			if _, err := os.Stat(keyFile); err == nil {
				pair[1] = filepath.Join(certsDir, keyFile)
			} else if !os.IsNotExist(err) {
				return nil, err
			}
			hosts[0].clientPairs = append(hosts[0].clientPairs, pair)
		}
	}
	return hosts, nil
}
