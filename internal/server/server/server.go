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

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	gw "github.com/lastbackend/lastbackend/internal/server/api/v1"
	"github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	v1 "github.com/lastbackend/lastbackend/internal/server/server/v1"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	defaultHTTPPort = 2967
	defaultGRPCPort = 2968
	defaultPROMPort = 2662
)

type Server struct {
	config  config.Config
	storage storage.IStorage
}

func NewServer(stg storage.IStorage, cfg config.Config) *Server {

	hs := new(Server)
	hs.storage = stg
	hs.config = cfg

	return hs
}

func (hs *Server) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	log := logger.WithContext(ctx)

	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", hs.config.Server.Host, defaultGRPCPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcprometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpcprometheus.UnaryServerInterceptor),
	)

	gw.RegisterV1Server(grpcServer, v1.NewV1Server())
	grpcprometheus.Register(grpcServer)

	group.Go(func() error {
		log.Info("grpc")
		return grpcServer.Serve(sock)
	})

	mw := middleware.New(hs.storage, hs.config.Security.Token)
	mw.Add(mw.Logger)
	mw.Add(mw.RequestID)
	mw.Add(mw.Authenticate)

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	runtime.SetHTTPBodyMarshaler(mux)

	var srv *http.Server
	if hs.config.Server.TLS.Verify {
		srv = &http.Server{
			Addr: fmt.Sprintf("%s:%d", hs.config.Server.Host, defaultHTTPPort),
			// add handler with middleware
			Handler:   mw.Apply(mux, hs.config),
			TLSConfig: configTLS(hs.config.Server.TLS.FileCA),
		}

	} else {
		srv = &http.Server{
			Addr: fmt.Sprintf("%s:%d", hs.config.Server.Host, defaultHTTPPort),
			// add handler with middleware
			Handler: mw.Apply(mux, hs.config),
		}
	}

	opts := []grpc.DialOption{

		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}

	if hs.config.Server.TLS.Verify {
		// enable security options
		creds, _ := credentials.NewServerTLSFromFile(hs.config.Server.TLS.FileCert, hs.config.Server.TLS.FileKey)
		opts = append(opts, grpc.WithTransportCredentials(creds))

	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	group.Go(func() error {
		log.Info("handler")
		return gw.RegisterV1HandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", hs.config.Server.Host, defaultGRPCPort), opts)
	})

	group.Go(func() error {
		log.Info("http")
		if hs.config.Server.TLS.Verify {
			return srv.ListenAndServeTLS(hs.config.Server.TLS.FileCert, hs.config.Server.TLS.FileKey)
		} else {
			return srv.ListenAndServe()
		}
	})

	group.Go(func() error {
		log.Info("listen")
		return http.ListenAndServe(fmt.Sprintf("%s:%d", hs.config.Server.Host, defaultPROMPort), promhttp.Handler())
	})

	return group.Wait()
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
		// Ensure that we only use our "CA" to validate certificates
		ClientCAs: caCertPool,
		// Force it server side
		PreferServerCipherSuites: true,
		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}

	return TLSConfig
}
