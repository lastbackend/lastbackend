package config

import (
	etcd "github.com/coreos/etcd/client"
	"k8s.io/client-go/1.5/rest"
	"time"
)

var config Config

func Get() *Config {
	return &config
}

func GetK8S() *rest.Config {
	return &rest.Config{
		Host: config.K8S.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile:   config.K8S.SSL.CA,
			KeyFile:  config.K8S.SSL.Key,
			CertFile: config.K8S.SSL.Cert,
		},
	}
}

func GetEtcd() (etcd.Client, error) {
	db_config := etcd.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	return etcd.New(db_config)
}
