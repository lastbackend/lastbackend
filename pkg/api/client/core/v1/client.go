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

package v1

import (
	"crypto/tls"

	"github.com/coreos/etcd/pkg/transport"
)

type Client struct {
	endpoint string
	tls      *tls.Config
}

func Get(conf Config) (*Client, error) {

	tlsConfig, err := getTLSConfig(conf.TLS.Cert, conf.TLS.Key, conf.TLS.CA)
	if err != nil {
		return nil, err
	}

	c := &Client{
		endpoint: conf.Endpoint,
		tls:      tlsConfig,
	}

	return c, nil
}

func getTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {

	if len(certFile) == 0 || len(keyFile) == 0 || len(caFile) == 0 {
		return nil, nil
	}

	tlsInfo := transport.TLSInfo{
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}

	return tlsInfo.ClientConfig()
}
