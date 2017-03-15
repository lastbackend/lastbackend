package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	r "gopkg.in/dancannon/gorethink.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"reflect"
	"time"
)

var ExternalConfig interface{}
var config = new(Config)

func (Config) Configure(path string) error {

	// Parsing config file
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if !isNil(ExternalConfig) {
		err = yaml.Unmarshal(buf, ExternalConfig)
		if err != nil {
			return err
		}
	}

	return yaml.Unmarshal(buf, &config)
}

func Get() *Config {
	return config
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

// Get Etcd DB options used for creating session
func GetEtcdDB() *etcd.Client {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: config.Etcd.TimeOut * time.Second,
	})
	if err != nil {
		_ = fmt.Errorf("Etcd error: %c", err.Error())
	}

	return cli
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
			_ = fmt.Errorf("SSL read error: %c", err.Error())
		}

		roots.AppendCertsFromPEM(cert)

		options.TLSConfig = &tls.Config{
			RootCAs: roots,
		}
	}

	return options
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
