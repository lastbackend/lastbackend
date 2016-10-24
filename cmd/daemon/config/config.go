package config

import (
	"k8s.io/client-go/1.5/rest"
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
