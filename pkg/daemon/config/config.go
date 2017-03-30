//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package config

import (
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"github.com/lastbackend/lastbackend/pkg/util/serializer/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	cfg = new(Config)
)

func (Config) Configure(path string) error {

	// Parsing config file
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cfg.HttpServer.Port = 3000
	cfg.HttpServer.Host = "0.0.0.0"
	return yaml.Unmarshal(buf, &cfg)
}

func Get() *Config {
	return cfg
}

// Get Etcd DB options used for creating session
func (c *Config) GetEtcdDB() store.Config {
	return store.Config{
		Prefix:    "lastbackend",
		Endpoints: c.Etcd.Endpoints,
		KeyFile:   c.Etcd.TLS.Key,
		CertFile:  c.Etcd.TLS.Cert,
		CAFile:    c.Etcd.TLS.CA,
		Quorum:    c.Etcd.Quorum,
		Codec:     serializer.NewSerializer(json.Encoder{}, json.Decoder{}),
	}
}

func (c *Config) GetVendorConfig(vendor string) (string, string, string) {

	var clientID, clientSecretID, redirectURI string

	switch vendor {
	case "github":
		clientID = c.VCS.Github.Client.ID
		clientSecretID = c.VCS.Github.Client.SecretID
	case "gitlab":
		clientID = c.VCS.Gitlab.Client.ID
		clientSecretID = c.VCS.Gitlab.Client.SecretID
		redirectURI = c.VCS.Gitlab.RedirectUri
	case "bitbucket":
		clientID = c.VCS.Bitbucket.Client.ID
		clientSecretID = c.VCS.Bitbucket.Client.SecretID
		redirectURI = c.VCS.Bitbucket.RedirectUri
	}

	return clientID, clientSecretID, redirectURI
}
