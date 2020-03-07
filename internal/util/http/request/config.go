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

package request

import (
	"time"
)

// config holds the common attributes that can be passed to a Last.Backend client on
// initialization.
type Config struct {
	BearerToken string
	Timeout     time.Duration
	TLS         *TLSConfig
	Headers     map[string]string
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

// HasTokenAuth returns whether the configuration has token authentication or not.
func (c *Config) HasTokenAuth() bool {
	return len(c.BearerToken) != 0
}

// HasCA returns whether the configuration has a certificate authority or not.
func (c *Config) HasCA() bool {
	return len(c.TLS.CAData) > 0 || len(c.TLS.CAFile) > 0
}

// HasCertAuth returns whether the configuration has certificate authentication or not.
func (c *Config) HasCertAuth() bool {
	return (len(c.TLS.CertData) != 0 || len(c.TLS.CertFile) != 0) && (len(c.TLS.KeyData) != 0 || len(c.TLS.KeyFile) != 0)
}
