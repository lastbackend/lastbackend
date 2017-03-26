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
	"github.com/lastbackend/lastbackend/pkg/serializer"
	"github.com/lastbackend/lastbackend/pkg/serializer/json"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	cfg            = new(Config)
	ExternalConfig interface{}
)

func (Config) Configure(path string) error {

	// Parsing config file
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if !validator.IsNil(ExternalConfig) {
		err = yaml.Unmarshal(buf, ExternalConfig)
		if err != nil {
			return err
		}
	}

	return yaml.Unmarshal(buf, &cfg)
}

func Get() *Config {
	return cfg
}

// Get Etcd DB options used for creating session
func GetEtcdDB() store.Config {
	return store.Config{
		Prefix:    "lastbackend",
		Endpoints: cfg.Etcd.Endpoints,
		KeyFile:   cfg.Etcd.TLS.Key,
		CertFile:  cfg.Etcd.TLS.Cert,
		CAFile:    cfg.Etcd.TLS.CA,
		Quorum:    cfg.Etcd.Quorum,
		Codec:     serializer.NewSerializer(json.Encoder{}, json.Decoder{}),
	}
}

func (c *Config) GetVendorConfig(vendor string) (string, string, string) {

	var clientID, clientSecretID, redirectURI string

	switch vendor {
	case "github":
		clientID = c.VCS.Github.User.Platform.Client.ID
		clientSecretID = c.VCS.Github.User.Platform.Client.SecretID
	case "gitlab":
		clientID = c.VCS.Gitlab.User.Platform.Client.ID
		clientSecretID = c.VCS.Gitlab.User.Platform.Client.SecretID
		redirectURI = c.VCS.Gitlab.User.Platform.RedirectUri
	case "bitbucket":
		clientID = c.VCS.Bitbucket.User.Platform.Client.ID
		clientSecretID = c.VCS.Bitbucket.User.Platform.Client.SecretID
		redirectURI = c.VCS.Bitbucket.User.Platform.RedirectUri
	}

	return clientID, clientSecretID, redirectURI
}
