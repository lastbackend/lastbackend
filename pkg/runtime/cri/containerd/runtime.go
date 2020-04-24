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

package containerd

import "github.com/containerd/containerd"

const defaultNamespace = "lstbknd"
const containerdDefaultAddress = "/run/containerd/containerd.sock"

type Runtime struct {
	client *containerd.Client
}

type Config struct {
	Address string
	TLS     *TLSConfig
}

type TLSConfig struct {
	CAPath   string
	CertPath string
	KeyPath  string
}

func New(cfg Config) (*Runtime, error) {
	r := new(Runtime)

	address := cfg.Address
	if len(cfg.Address) == 0 {
		address = containerdDefaultAddress
	}

	client, err := containerd.New(address, containerd.WithDefaultNamespace(defaultNamespace))
	if err != nil {
		return nil, err
	}

	r.client = client

	return r, nil
}
