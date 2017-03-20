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

package store

import (
	"github.com/lastbackend/lastbackend/pkg/serializer"
)

// Config is configuration for creating a storage backend.
type Config struct {
	// Prefix is the prefix to all keys passed to storage.Interface methods.
	Prefix string
	// Enpoints is the list of storage servers to connect with.
	Endpoints []string
	// TLS credentials
	KeyFile  string
	CertFile string
	CAFile   string
	// Quorum indicates that whether read operations should be quorum-level consistent.
	Quorum bool
	Codec  serializer.Codec
}

func NewDefaultConfig(prefix string, codec serializer.Codec) *Config {
	return &Config{
		Prefix: prefix,
		Codec:  codec,
	}
}