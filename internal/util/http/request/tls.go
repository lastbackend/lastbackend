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

package request

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func NewTLSConfig(c *Config) (*tls.Config, error) {

	if !(c.HasCA() || c.HasCertAuth() || c.TLS.Insecure || len(c.TLS.ServerName) > 0) {
		return nil, nil
	}
	if c.HasCA() && c.TLS.Insecure {
		return nil, fmt.Errorf("certificates file with the insecure flag is not allowed")
	}
	if err := loadTLSFiles(c); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		InsecureSkipVerify:       c.TLS.Insecure,
		ServerName:               c.TLS.ServerName,
		PreferServerCipherSuites: true,
	}

	if c.HasCA() {
		// Load CA cert
		tlsConfig.RootCAs = rootCertPool(c.TLS.CAData)
	}

	var staticCert *tls.Certificate

	if c.HasCertAuth() {
		cert, err := tls.X509KeyPair(c.TLS.CertData, c.TLS.KeyData)
		if err != nil {
			return nil, err
		}
		staticCert = &cert
	}

	if c.HasCertAuth() {
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if staticCert != nil {
				return staticCert, nil
			}
			return new(tls.Certificate), nil
		}
	}

	return tlsConfig, nil
}

func loadTLSFiles(c *Config) (err error) {

	c.TLS.CAData, err = getDataFromSliceOrFile(c.TLS.CAData, c.TLS.CAFile)
	if err != nil {
		return err
	}

	c.TLS.CertData, err = getDataFromSliceOrFile(c.TLS.CertData, c.TLS.CertFile)
	if err != nil {
		return err
	}

	c.TLS.KeyData, err = getDataFromSliceOrFile(c.TLS.KeyData, c.TLS.KeyFile)
	if err != nil {
		return err
	}

	return nil
}

func getDataFromSliceOrFile(data []byte, file string) ([]byte, error) {

	if len(data) > 0 {
		return data, nil
	}

	if len(file) == 0 {
		return nil, nil
	}

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	return fileData, nil
}

func rootCertPool(caData []byte) *x509.CertPool {
	if len(caData) == 0 {
		return nil
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(caData)
	if !ok {
		panic("failed to parse root certificate")
	}
	return certPool
}
