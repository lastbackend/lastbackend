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

package v3

import (
	"crypto/tls"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	store2 "github.com/lastbackend/lastbackend/internal/pkg/storage/etcd/store"
	"github.com/lastbackend/lastbackend/internal/util/serializer"
	"github.com/lastbackend/lastbackend/internal/util/serializer/json"
	"path"
	"time"
)

func GetClient(conf *Config) (store2.Store, store2.DestroyFunc, error) {

	if conf == nil {
		panic("need to set the storage configuration")
	}

	tlsConfig, err := getTLSConfig(conf.TLS.Cert, conf.TLS.Key, conf.TLS.CA)
	if err != nil {
		return nil, nil, err
	}

	cfg := clientv3.Config{
		Endpoints:   conf.Endpoints,
		TLS:         tlsConfig,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		//client.Close()
	}

	codec := serializer.NewSerializer(json.Encoder{}, json.Decoder{})
	var st = &dbstore{
		client:     client,
		codec:      codec,
		pathPrefix: path.Join("/", conf.Prefix),
		watcher:    newWatcher(client),
	}

	st.watcher.prefix = path.Join("/", conf.Prefix)
	st.opts = append(st.opts, clientv3.WithSerializable())

	return st, destroyFunc, nil
}

func getTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {

	if len(certFile) == 0 || len(keyFile) == 0 || len(caFile) == 0 {
		return nil, nil
	}

	tlsInfo := transport.TLSInfo{
		CertFile:      certFile,
		KeyFile:       keyFile,
		TrustedCAFile: caFile,
	}

	return tlsInfo.ClientConfig()
}
