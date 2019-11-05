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

package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/util/http/cors"
	"io/ioutil"
	"net/http"
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

type NotFoundHandler struct {
	http.Handler
}

func (NotFoundHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"code": 404, "status": "Not Found", "message": "Not Found"}`))
}

type MethodNotAllowedHandler struct {
	http.Handler
}

func (MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"code": 405, "status": "Method Not Allowed", "message": "Method Not Allowed"}`))
}

func Handle(ctx context.Context, h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	headers := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cors.Headers(w, r)
			h.ServeHTTP(w, r)
		}
	}

	h = headers(h)
	for _, m := range middleware {
		h = m(ctx, h)
	}

	return h
}

func Listen(host string, port int, router http.Handler) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router)
}

func ListenWithTLS(host string, port int, caFile, certFile, keyFile string, router http.Handler) error {

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	server.TLSConfig = configTLS(caFile)

	return server.ListenAndServeTLS(certFile, keyFile)
}

func configTLS(caFile string) *tls.Config {

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		panic("failed to parse root certificate")
	}

	TLSConfig := &tls.Config{
		// Reject any TLS certificate that cannot be validated
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Ensure that we only use our "FileCA" to validate certificates
		ClientCAs: caCertPool,
		// Force it server side
		PreferServerCipherSuites: true,
		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}
	return TLSConfig
}
