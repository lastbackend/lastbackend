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

package config

import (
	"time"
)

const (
	defaultTimeout = 10
)

// config holds the common attributes that can be passed to a Last.Backend client on
// initialization.
type Config struct {
	BearerToken string
	TLS         *TLSConfig
	Timeout     time.Duration
	Headers     map[string]string
}

func (c Config) getDefault() *Config {
	cfg := new(Config)
	cfg.Timeout = defaultTimeout
	return cfg
}

// TLSConfig contains settings to enable transport layer security
type TLSConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool

	// Override for the server name passed to the server for SNI and used to verify certificates..
	ServerName string

	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string

	// Bytes of the PEM-encoded server trusted root certificates. Supersedes CAFile.
	CAData []byte
	// Bytes of the PEM-encoded client certificate. Supersedes CertFile.
	CertData []byte
	// Bytes of the PEM-encoded client key. Supersedes KeyFile.
	KeyData []byte
}
