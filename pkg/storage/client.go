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

package storage

import (
	"crypto/tls"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd3"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"time"
)

func createEtcd3Storage(c store.Config, log logger.ILogger) (store.IStore, store.DestroyFunc, error) {

	tlsConfig, err := getTLSConfig(c.CertFile, c.KeyFile, c.CAFile)
	if err != nil {
		return nil, nil, err
	}

	cfg := clientv3.Config{
		Endpoints:   c.Endpoints,
		TLS:         tlsConfig,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		client.Close()
	}

	return etcd3.New(client, c.Codec, c.Prefix, log), destroyFunc, nil
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
