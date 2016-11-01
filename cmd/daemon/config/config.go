package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	r "gopkg.in/dancannon/gorethink.v2"
	"io/ioutil"
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

// Get Rethink DB options used for creating session
func GetRethinkDB() r.ConnectOpts {

	options := r.ConnectOpts{
		MaxOpen:    config.RethinkDB.MaxOpen,
		InitialCap: config.RethinkDB.InitialCap,
		Database:   config.RethinkDB.Database,
		AuthKey:    config.RethinkDB.AuthKey,
	}

	if len(config.RethinkDB.Addresses) > 0 {
		options.Addresses = config.RethinkDB.Addresses
	} else {
		options.Address = config.RethinkDB.Address
	}

	if config.RethinkDB.SSL.CA != "" {
		roots := x509.NewCertPool()
		cert, err := ioutil.ReadFile(config.RethinkDB.SSL.CA)

		if err != nil {
			fmt.Errorf("SSL read error: %s", err.Error())
		}

		roots.AppendCertsFromPEM(cert)

		options.TLSConfig = &tls.Config{
			RootCAs: roots,
		}
	}

	return options
}
