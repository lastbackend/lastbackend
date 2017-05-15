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
	"github.com/lastbackend/lastbackend/pkg/common/config"
)

var _cfg = new(Config)

func Set(cfg *config.Config) *Config{
	_cfg.LogLevel = cfg.LogLevel
	_cfg.Token = cfg.Token
	_cfg.Etcd = cfg.Etcd
	_cfg.Registry = cfg.Registry
	_cfg.APIServer = cfg.APIServer
	return _cfg
}

func Get() *Config {
	return _cfg
}

// Get Etcd DB options used for creating session
func (c *Config) GetEtcdDB() store.Config {
	return store.Config{
		Prefix:    "lastbackend",
		Endpoints: *c.Etcd.Endpoints,
		KeyFile:   *c.Etcd.TLS.Key,
		CertFile:  *c.Etcd.TLS.Cert,
		CAFile:    *c.Etcd.TLS.CA,
		Codec:     serializer.NewSerializer(json.Encoder{}, json.Decoder{}),
	}
}
